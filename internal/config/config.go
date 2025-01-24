package config

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/frk/gosql/internal/typesutil"

	"golang.org/x/tools/go/packages"

	_ "github.com/lib/pq"
)

// The Config struct is used to configure the gosql tool.
type Config struct {
	// The directory in which the tool will search for files to process.
	// If not provided, the current working directory will be used by default.
	WorkingDirectory String `json:"working_directory"`
	// If set to true, the tool will search the hierarchy of the working
	// directory for files to process.
	Recursive Bool `json:"recursive"`
	// List of files to be used as input for the tool.
	// The files must be located in the working directory.
	InputFiles StringSlice `json:"input_files"`
	// List of regular expressions to match input files that the tool should
	// process. The regular expressions must match files that are located in
	// the working directory.
	InputFileRegexps StringSlice `json:"input_file_regexps"`
	// The format used for generating the name of the output files.
	//
	// The format can contain one (and only one) "%s" placeholder which the
	// tool will replace with the input file's base name, if no placeholder is
	// present then the input file's base name will be prefixed to the format.
	//
	// If not provided, the format "%s_gosql.go" will be used by default.
	OutputFileNameFormat String `json:"output_file_name_format"`
	// The connection string of the database that will be used for type checking.
	// This value is required.
	DatabaseDSN String `json:"database_dsn"`
	// If set to true, the generator will quote postgres identifiers like
	// column names, table names, etc.
	QuoteIdentifiers Bool `json:"quote_identifiers"`
	// If set to a non-empty string, it specifies the struct tag to be used
	// for constructing the column keys of a FilterXxx type. A valid tag must
	// begin with a letter (A-z) or an underscore (_), subsequent characters
	// in the tag can be letters, underscores, and digits (0-9).
	// If set to "" (empty string), the generator will default to use the
	// field names instead of struct tags to construct the column keys.
	//
	// If not provided, the tag "json" will be used by default.
	FilterColumnKeyTag String `json:"filter_column_key_tag"`
	// If set, instructs the generator to use only the base of a tag/field
	// chain to construct the column keys of a FilterXxx type.
	//
	// If not provided, `false` will be used by default.
	FilterColumnKeyBase Bool `json:"filter_column_key_base"`
	// The separator to be used to join a chain of tag/field values for
	// constructing the column keys of a FilterXxx type. The separator can
	// be at most one byte long.
	//
	// If not provided, the separator "." will be used by default.
	FilterColumnKeySeparator String `json:"filter_column_key_separator"`
	// The name to be used for the generated method.
	//
	// If not provided, the name "Exec" will be used by default.
	MethodName String `json:"method_name"`
	// If set, the generator will produce methods that take context.Context
	// as its first argument and they will pass that argument to the XxxContext
	// query executing methods of the Conn type.
	//
	// If not provided, `false` will be used by default.
	MethodWithContext Bool `json:"method_with_context"`
	// The Go type to be used as the "querier" argument for the generated methods.
	//
	// The string value must be of the format "[*]package/path.TypeName" and
	// the type represented by it must implement the following interface:
	//	{
	// 	     Exec(query string, args ...interface{}) (sql.Result, error)
	// 	     ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	// 	     Query(query string, args ...interface{}) (*sql.Rows, error)
	// 	     QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	// 	     QueryRow(query string, args ...interface{}) *sql.Row
	// 	     QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	// 	}
	//
	// If not provided, the type github.com/frk/gosql.Conn will be used by default.
	MethodArgumentType GoType `json:"method_argument_type"`

	// A list of filter value converters that will be used to map
	// filter value converting functions to specific types.
	FilterValueConverters []FilterValueConverter `json:"filter_value_converters"`

	// holds the compiled expressions of the InputFileRegexps slice.
	compiledInputFileRegexps []*regexp.Regexp
}

type FilterValueConverter struct {
	// Go type for which to use the filter value converter.
	Type ObjectIdent `json:"type"`
	// Go func which should be used as the filter value converter.
	Func ObjectIdent `json:"func"`
}

var DefaultConfig = Config{
	WorkingDirectory:         String{Value: "."},
	Recursive:                Bool{Value: false},
	InputFiles:               StringSlice{},
	InputFileRegexps:         StringSlice{},
	OutputFileNameFormat:     String{Value: "%s_gosql.go"},
	DatabaseDSN:              String{Value: ""},
	QuoteIdentifiers:         Bool{Value: false},
	FilterColumnKeyTag:       String{Value: "json"},
	FilterColumnKeyBase:      Bool{Value: false},
	FilterColumnKeySeparator: String{Value: "."},
	MethodName:               String{Value: "Exec"},
	MethodWithContext:        Bool{Value: false},
	MethodArgumentType: GoType{
		Name:    "Conn",
		PkgPath: "github.com/frk/gosql",
		PkgName: "gosql",
	},
}

// ParseFlags unmarshals the cli flags into the receiver.
func (c *Config) ParseFlags(printUsage func()) {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.Usage = printUsage
	fs.Var(&c.WorkingDirectory, "wd", "")
	fs.Var(&c.Recursive, "r", "")
	fs.Var(&c.InputFiles, "f", "")
	fs.Var(&c.InputFileRegexps, "rx", "")
	fs.Var(&c.OutputFileNameFormat, "o", "")
	fs.Var(&c.DatabaseDSN, "db", "")
	fs.Var(&c.QuoteIdentifiers, "qid", "")
	fs.Var(&c.FilterColumnKeyTag, "fcktag", "")
	fs.Var(&c.FilterColumnKeyBase, "fckbase", "")
	fs.Var(&c.FilterColumnKeySeparator, "fcksep", "")
	fs.Var(&c.MethodWithContext, "with-ctx", "")
	fs.Var(&c.MethodArgumentType, "argtype", "")
	fs.StringVar(&ConfigFile, "config", "", "The filepath to a specific configuration file")
	_ = fs.Parse(os.Args[1:])
}

// The filepath to the config file which will be used by ParseFile to load the configuration.
// If not provided ParseFile will look for a config file in the caller's git project root.
//
// NOTE This value is set by the ParseFlags method.
var ConfigFile string

// ParseFile looks for a gosql config file in the git project's root of the receiver's
// working directory, if it finds such a file it will then unmarshal it into the receiver.
func (c *Config) ParseFile() error {
	if ConfigFile == "" {
		dir, err := filepath.Abs(c.WorkingDirectory.Value)
		if err != nil {
			return err
		}

		var isRoot bool
		var configName string
		for len(dir) > 1 && dir[0] == '/' {
			isRoot, configName, err = examineDir(dir)
			if err != nil {
				return err
			}
			if isRoot {
				break
			}
			dir = filepath.Dir(dir) // parent dir will be examined next
		}
		if configName != "" {
			ConfigFile = filepath.Join(dir, configName)
		}
	}

	// if explicitly set or found in project's root, unamrshal the config file
	if ConfigFile != "" {
		f, err := os.Open(ConfigFile)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := json.NewDecoder(f).Decode(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) FileFilterFunc() (filter func(filePath string) bool) {
	if len(c.InputFiles.Value) == 0 && len(c.InputFileRegexps.Value) == 0 {
		return nil
	}

	// copy file paths for matching
	allowFilePaths := make([]string, len(c.InputFiles.Value))
	copy(allowFilePaths, c.InputFiles.Value)

	// copy regular expressions for matching
	allowRegexps := make([]*regexp.Regexp, len(c.compiledInputFileRegexps))
	copy(allowRegexps, c.compiledInputFileRegexps)

	return func(filePath string) bool {
		for _, fp := range allowFilePaths {
			if filePath == fp {
				return true
			}
		}
		for _, rx := range allowRegexps {
			if rx.MatchString(filePath) {
				return true
			}
		}
		return false
	}
}

// Validate checks the config for errors and updates some of the values to a more "normalized" format.
func (c *Config) Validate() (err error) {
	// check that the working directory can be openned
	f, err := os.Open(c.WorkingDirectory.Value)
	if err != nil {
		return fmt.Errorf("failed to open working directory: %q -- %v", c.WorkingDirectory.Value, err)
	}
	f.Close()

	// check that the dsn can be used to initiallize connections
	if len(c.DatabaseDSN.Value) == 0 {
		return fmt.Errorf("missing database connection string")
	}
	db, err := sql.Open("postgres", c.DatabaseDSN.Value)
	if err != nil {
		return fmt.Errorf("error opening database: %q -- %v", c.DatabaseDSN.Value, err)
	} else {
		defer db.Close()
		if err := db.Ping(); err != nil {
			return fmt.Errorf("error connecting to database: %q -- %v", c.DatabaseDSN.Value, err)
		}
	}

	// update file paths to absolutes
	for i, fp := range c.InputFiles.Value {
		abs, err := filepath.Abs(fp)
		if err != nil {
			return fmt.Errorf("error resolving absolute path of file: %q -- %v", fp, err)
		}
		c.InputFiles.Value[i] = abs
	}

	// compile the input file regexeps
	c.compiledInputFileRegexps = make([]*regexp.Regexp, len(c.InputFileRegexps.Value))
	for i, expr := range c.InputFileRegexps.Value {
		rx, err := regexp.Compile(expr)
		if err != nil {
			return fmt.Errorf("error compiling regular expression: %q -- %v", expr, err)
		}
		c.compiledInputFileRegexps[i] = rx
	}

	// check that the output filename format contains at most one "%" and
	// that it is followed by an "s" to form the "%s" verb
	if n := strings.Count(c.OutputFileNameFormat.Value, "%"); n == 0 {
		// modify the output filename format
		c.OutputFileNameFormat.Value = "%s" + c.OutputFileNameFormat.Value
	} else if n > 1 || (n == 1 && !strings.Contains(c.OutputFileNameFormat.Value, "%s")) {
		return fmt.Errorf("bad output filename format: %q", c.OutputFileNameFormat.Value)
	}

	// check the filter column key configuration for errors
	rxFCKTag := regexp.MustCompile(`^(?:[A-Za-z_]\w*)?$`)
	if !rxFCKTag.MatchString(c.FilterColumnKeyTag.Value) {
		return fmt.Errorf("bad filter column key tag: %q", c.FilterColumnKeyTag.Value)
	}
	if len(c.FilterColumnKeySeparator.Value) > 1 {
		return fmt.Errorf("bad filter column key separator: %q", c.FilterColumnKeySeparator.Value)
	}

	// check that the method argument type implements the required interface
	if c.MethodArgumentType.String() != "github.com/frk/gosql.Conn" {
		if err := checkMethodArgumentType(&c.MethodArgumentType, c); err != nil {
			return err
		}
	}

	//
	// TODO: needs to confirm the validity of the types & functions
	//
	// for i, fvc := range c.FilterValueConverters {
	// 	if err := checkFilterValueConverter(fvc, c); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

// examineDir reports if the directory at the given path is the root directory
// of a git project and if it is, it will also report the name of the gosql config
// file (either ".gosql" or ".gosql.json") if such a file exists in that root
// directory, otherwise confName will be empty.
func examineDir(path string) (isRoot bool, confName string, err error) {
	d, err := os.Open(path)
	if err != nil {
		return false, "", err
	}
	defer d.Close()

	infoList, err := d.Readdir(-1)
	if err != nil {
		return false, "", err
	}

	for _, info := range infoList {
		name := info.Name()
		if name == ".git" && info.IsDir() {
			isRoot = true
		}
		if (name == ".gosql" || name == ".gosql.json") && !info.IsDir() {
			confName = name
		}
	}

	// NOTE(mkopriva): currently we don't care about .gosql files that live outside
	// of the git project root directory, if, in the future, the rules are expanded
	// then this will need to be either removed or accordingly updated.
	if !isRoot {
		confName = ""
	}
	return isRoot, confName, nil
}

func checkMethodArgumentType(t *GoType, c *Config) (err error) {
	cfg := &packages.Config{Mode: packages.NeedName |
		packages.NeedFiles |
		packages.NeedCompiledGoFiles |
		packages.NeedSyntax |
		packages.NeedTypes |
		packages.NeedImports |
		packages.NeedDeps |
		packages.NeedTypesInfo}
	pkgs, err := packages.Load(cfg, t.PkgPath)
	if err != nil {
		return fmt.Errorf("failed to load package of method argument type: %q -- %v", t, err)
	}

	for _, syn := range pkgs[0].Syntax {
		for _, dec := range syn.Decls {
			gd, ok := dec.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}

			for _, spec := range gd.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok || typeSpec.Name.Name != t.Name {
					continue
				}

				obj, ok := pkgs[0].TypesInfo.Defs[typeSpec.Name]
				if !ok {
					return fmt.Errorf("bad method argument type: %q", t)
				}

				typeName, ok := obj.(*types.TypeName)
				if !ok {
					return fmt.Errorf("bad method argument type: %q", t)
				}

				named, ok := typeName.Type().(*types.Named)
				if !ok {
					return fmt.Errorf("bad method argument type: %q", t)
				}

				if !typesutil.ImplementsGosqlConn(named) {
					return fmt.Errorf("bad method argument type: %q --"+
						" does not implement the \"github.com/frk/gosql.Conn\" interface.", t)
				}

				// Use the package name from the AST since it may be
				// different from the last segment of the package's path.
				t.PkgName = pkgs[0].Name
				return nil
			}
		}
	}

	return fmt.Errorf("could not find method argument type: %q", t)
}

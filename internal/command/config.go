package command

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
	"strconv"
	"strings"

	"github.com/frk/gosql/internal/typesutil"

	"golang.org/x/tools/go/packages"

	_ "github.com/lib/pq"
)

// The Config struct is used to configure the gosql tool.
type Config struct {
	// The connection string of the database that will be used for type checking.
	// This value is required.
	DatabaseDSN String `json:"database_dsn"`
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
	// The Go type to be used as the argument for the generated methods.
	//
	// The string value must be of the format "[*]package/path.TypeName" and
	// the type represented by it must implement the following interface:
	//	{
	// 	     Exec(query string, args ...interface{}) (sql.Result, error)
	// 	     Query(query string, args ...interface{}) (*sql.Rows, error)
	// 	     QueryRow(query string, args ...interface{}) *sql.Row
	// 	}
	//
	// If not provided, the type gosql.Conn will be used by default.
	MethodArgumentType String `json:"method_argument_type"`

	// holds the compiled expressions of the InputFileRegexps slice.
	compiledInputFileRegexps []*regexp.Regexp
}

var DefaultConfig = Config{
	DatabaseDSN:              String{Value: ""},
	WorkingDirectory:         String{Value: "."},
	Recursive:                Bool{Value: false},
	InputFiles:               StringSlice{},
	InputFileRegexps:         StringSlice{},
	OutputFileNameFormat:     String{Value: "%s_gosql.go"},
	QuoteIdentifiers:         Bool{Value: false},
	FilterColumnKeyTag:       String{Value: "json"},
	FilterColumnKeyBase:      Bool{Value: false},
	FilterColumnKeySeparator: String{Value: "."},
	MethodArgumentType:       String{Value: "gosql.Conn"},
}

// ParseFlags unmarshals the cli flags into the receiver.
func (c *Config) ParseFlags() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.Usage = printUsage
	fs.Var(&c.DatabaseDSN, "db", "")
	fs.Var(&c.WorkingDirectory, "wd", "")
	fs.Var(&c.Recursive, "r", "")
	fs.Var(&c.InputFiles, "f", "")
	fs.Var(&c.InputFileRegexps, "rx", "")
	fs.Var(&c.OutputFileNameFormat, "o", "")
	fs.Var(&c.QuoteIdentifiers, "qid", "")
	fs.Var(&c.FilterColumnKeyTag, "fcktag", "")
	fs.Var(&c.FilterColumnKeyBase, "fckbase", "")
	fs.Var(&c.FilterColumnKeySeparator, "fcksep", "")
	fs.Var(&c.MethodArgumentType, "argtype", "")
	_ = fs.Parse(os.Args[1:])
}

// ParseFile looks for a gosql config file in the git project's root of the receiver's
// working directory, if it finds such a file it will then unmarshal it into the receiver.
func (c *Config) ParseFile() error {
	dir, err := filepath.Abs(c.WorkingDirectory.Value)
	if err != nil {
		return err
	}

	var isRoot bool
	var confName string
	for len(dir) > 1 && dir[0] == '/' {
		isRoot, confName, err = examineDir(dir)
		if err != nil {
			return err
		}
		if isRoot {
			break
		}
		dir = filepath.Dir(dir) // parent dir will be examined next
	}

	// if found, unamrshal the config file
	if confName != "" {
		confpath := filepath.Join(dir, confName)
		f, err := os.Open(confpath)
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

// validate checks the config for errors and updates some of the values to a more "normalized" format.
func (c *Config) validate() (err error) {
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
	if c.MethodArgumentType.Value != "gosql.Conn" {
		if err := checkMethodArgumentType(c.MethodArgumentType.Value); err != nil {
			return err
		}
	}
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

func checkMethodArgumentType(argtype string) (err error) {
	// split into import path and type name
	var importPath, typeIdent string
	if i := strings.LastIndex(argtype, "."); i < 0 {
		return fmt.Errorf("bad method argument type: %q", argtype)
	} else {
		importPath, typeIdent = argtype[:i], argtype[i+1:]
	}

	cfg := &packages.Config{Mode: packages.NeedTypesInfo}
	pkgs, err := packages.Load(cfg, importPath)
	if err != nil {
		return fmt.Errorf("failed to load package of method argument type: %q -- %v", argtype, err)
	}

	for _, syn := range pkgs[0].Syntax {
		for _, dec := range syn.Decls {
			gd, ok := dec.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}

			for _, spec := range gd.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok || typeSpec.Name.Name != typeIdent {
					continue
				}

				obj, ok := pkgs[0].TypesInfo.Defs[typeSpec.Name]
				if !ok {
					return fmt.Errorf("bad method argument type: %q", argtype)
				}

				typeName, ok := obj.(*types.TypeName)
				if !ok {
					return fmt.Errorf("bad method argument type: %q", argtype)
				}

				named, ok := typeName.Type().(*types.Named)
				if !ok {
					return fmt.Errorf("bad method argument type: %q", argtype)
				}

				if !typesutil.ImplementsGosqlConn(named) {
					return fmt.Errorf("bad method argument type: %q", argtype)
				}

				// all good
				return nil
			}
		}
	}

	return fmt.Errorf("could not find method argument type: %q", argtype)
}

// String implements both the flag.Value and the json.Unmarshal interfaces
// enforcing priority of flags over json, meaning that json.Unmarshal will
// not override the value if it was previously set by flag.Var.
type String struct {
	Value string
	IsSet bool
}

// Get implements the flag.Getter interface.
func (s String) Get() interface{} {
	return s.Value
}

// String implements the flag.Value interface.
func (s String) String() string {
	return s.Value
}

// Set implements the flag.Value interface.
func (s *String) Set(value string) error {
	s.Value = value
	s.IsSet = true
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *String) UnmarshalJSON(data []byte) error {
	if !s.IsSet {
		if len(data) == 0 || string(data) == `null` {
			return nil
		}

		var value string
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		s.Value = value
		s.IsSet = true
	}
	return nil
}

// Bool implements both the flag.Value and the json.Unmarshal interfaces
// enforcing priority of flags over json, meaning that json.Unmarshal will
// not override the value if it was previously set by flag.Var.
type Bool struct {
	Value bool
	IsSet bool
}

// IsBoolFlag indicates that the Bool type can be used as a boolean flag.
func (b Bool) IsBoolFlag() bool {
	return true
}

// Get implements the flag.Getter interface.
func (b Bool) Get() interface{} {
	return b.String()
}

// String implements the flag.Value interface.
func (b Bool) String() string {
	return strconv.FormatBool(b.Value)
}

// Set implements the flag.Value interface.
func (b *Bool) Set(value string) error {
	if len(value) > 0 {
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		b.Value = v
		b.IsSet = true
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (b *Bool) UnmarshalJSON(data []byte) error {
	if !b.IsSet {
		if len(data) == 0 || string(data) == `null` {
			return nil
		}

		var value bool
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		b.Value = value
		b.IsSet = true
	}
	return nil
}

// StringSlice implements both the flag.Value and the json.Unmarshal interfaces
// enforcing priority of flags over json, meaning that json.Unmarshal will
// not override the value if it was previously set by flag.Var.
type StringSlice struct {
	Value []string
	IsSet bool
}

// Get implements the flag.Getter interface.
func (ss StringSlice) Get() interface{} {
	return ss.String()
}

// String implements the flag.Value interface.
func (ss StringSlice) String() string {
	return strings.Join(ss.Value, ",")
}

// Set implements the flag.Value interface.
func (ss *StringSlice) Set(value string) error {
	if len(value) > 0 {
		ss.Value = append(ss.Value, value)
		ss.IsSet = true
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ss *StringSlice) UnmarshalJSON(data []byte) error {
	if !ss.IsSet {
		if len(data) == 0 || string(data) == `null` {
			return nil
		}

		var value []string
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		if len(value) > 0 {
			ss.Value = value
			ss.IsSet = true
		}
	}
	return nil
}

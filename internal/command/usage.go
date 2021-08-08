package command

import (
	"fmt"
	"os"
)

func PrintUsage() {
	fmt.Fprint(os.Stderr, usage)
}

const usage = `usage: gosql [-db] [-wd] [-r] [-f] [-rx] [-o] [-qid] [-argtype]
	[-fcktag] [-fckbase] [-fcksep] [-with-ctx] [--config]

gosql generates SQL queries based .... (todo: write doc)


The -db flag specifies the connection string of the database that the tool should
used for type checking. This value is required.


The -wd flag specifies the directory whose files the tool will process. When used
together with the -f or -rx flags the tool will process only those files that match
the -f and -rx values. If left unespecified, the current working directory will be
used by default.


The -r flag instructs the tool to process the files in the whole hierarchy of the
working directory. When used together with the -f or -rx flags the tool will process
only those files that match the -f and -rx values.


The -f flag specifies a file to be used as input for the tool. The file must be
located in the working directory. The flag can be used more than once to specify
multiple files.


The -rx flag specifies a regular expressions to match input files that the tool should
process. The regular expressions must match files located in the working directory.
The flag can be used more than once to specify multiple regular expressions.


The -o flag specifies the format to be used for generating the name of the output files.
The format can contain one (and only one) "%s" placeholder which the tool will replace
with the input file's base name, if no placeholder is present then the input file's base
name will be prefixed to the format.
If left unspecified, the format "%s_gosql.go" will be used by default.


The -qid flag instructs the generator to quote postgres identifiers like
column names, table names, etc.


The -fcktag flag if set to a non-empty string, specifies the struct tag to be used
for constructing the column keys of a FilterXxx type. A valid tag must begin with a
letter (A-z) or an underscore (_), subsequent characters in the tag can be letters,
underscores, and digits (0-9). If set to "" (empty string), the generator will default
to use the field names instead of struct tags to construct the column keys.
If left unspecified, the tag "json" will be used by default.


The -fckbase flag if set, instructs the generator to use only the base of a tag/field
chain to construct the column keys of a FilterXxx type.
If left unspecified, the value false will be used by default.


The -fcksep flag specifies the separator to be used to join a chain of tag/field values
for constructing the column keys of a FilterXxx type. The separator can be at most one byte long.
If left unspecified, the separator "." will be used by default.


The -argtype flag specifies the Go type to be used as the argument for the generated
methods. The string value must be of the format "[*]package/path.TypeName" and the type
represented by it must implement the following interface:
    {
        Exec(query string, args ...interface{}) (sql.Result, error)
        Query(query string, args ...interface{}) (*sql.Rows, error)
        QueryRow(query string, args ...interface{}) *sql.Row
    }
If left unspecified, the type gosql.Conn will be used by default.

` //`

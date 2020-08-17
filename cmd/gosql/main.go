package main

type Config struct {
	// Connection string to the postgres database that will be used for type checking.
	// TODO
	// - accept as argument
	// - default to env var
	// - allow individual target structs to somehow specify the database
	//   connection string...
	DatabaseURL string
	// The source file, or directory, that contains the targets for the generator.
	//
	// NOTE(mkopriva):
	// 1. Any files in SourcePath whose filename ends with _test.go will be
	//    ignored by the generator.
	// 2. Any files in SourcePath that have a matching suffix to the OutputFileSuffix
	//    value will be ignored by the generator.
	SourcePath string
	// A regular expression, which, if not empty, will be used to filter
	// out files whose filepaths DO NOT match the pattern.
	FilePattern string
	// If set and SourcePath is a directory, traverse all sub directories of SourcePath.
	Recursive bool
	// Used to create the file name of the file that will hold the output
	// of the generator by appending the suffix to the source file's name.
	OutputFileSuffix string
	// Indicates how the key in a filter's column map should be constructed.
	FilterKeyType string
}

func main() {
	//
}

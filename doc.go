// WIP
//
// Package gosql implements a generator for SQL queries and the Go code that
// executes those queries.
//
// The generator works off of user-declared struct types that conform, in name
// and in structure, to a specific format.
//
// -----------------------------------------------------------------------------
//
// To generate the query code the user has to first declare a struct type
//
// The generator produces a method on the target type that
//
// Query types
//
//	InsertXxx
//	UpdateXxx
//	SelectXxx
//	DeleteXxx
//
// Filter type
//	FilterXxx
//
// Tags
//
// There are two types of tags recognized by the package: "rel" and "sql".
//
// The "rel" tag is used in the spec types to identify the target relation
// and to link that relation to the tagged field.
//
// The "sql" tag:
//
// Options:
//
//	pk	  - indicates the field's corresponding column to be the primary key
//	ro	  - specifies the field to be read only
//	wo	  - specifies the field to be write only
//	default   - blah blah
//	nullempty - blah blah
//	json      - blah blah
//	xml       - blah blah
//	+         - blah blah
//	cast      - blah blah
//	coalesce  - blah blah
//
// Blocks
//
// blocks are blah blah blah
//
//	Where      - blah blah
//	Join       - blah blah
//	From       - blah blah
//	Using      - blah blah
//	OnConflict - blah blah
//
// Directives
//
// directives are blah blah blahs
//
//	All
//	Default
//	Force
//	Return
//	Return
package gosql

// NOTE(mkopriva): This is used by default for UPDATEs which don't specify
// a WHERE clause, if multiple fields are tagged as pkeys then a composite
// primary key is assumed.

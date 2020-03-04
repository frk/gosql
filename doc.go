// Package gosql implements the generation of SQL queries and the Go code that
// executes those queries. The code is generated from struct types that are
// predeclared by the user of the package, these struct types act as specifictions
// for the code to be generated.
//
// Specs
//
// specifications are a thing blah blah blah
//
//	InsertXxx
//	UpdateXxx
//	SelectXxx
//	DeleteXxx
//	FilterXxx
//
// Tags
//
// This is a test blah blah blahs
//
// Options:
//
//	pk	indicates the field's corresponding column to be the primary key
//	ro	specifies the field to be read only
//	wo	specifies the field to be write only
//	auto
//	nullempty
//	json
//	xml
//	+
//	cast
//	coalesce
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

package main

// Vim script to execute the generator from the editor.
//------------------------------------------------------------------------------
// Handling of enums: postgres enums are text based and if the Go app declares
// an enum using integers rather than strings one has to map those integers onto
// text to be able to read/write enum values from/to the database. If possible
// it may be of great value to generate the mapping somehow.
//------------------------------------------------------------------------------

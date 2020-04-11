package convert

// [ ( x1 , y1 ) , ... , ( xn , yn ) ]
// ( ( x1 , y1 ) , ... , ( xn , yn ) )
//
// Square brackets ([]) indicate an open path, while parentheses (()) indicate
// a closed path. When the outermost parentheses are omitted, as in the third
// through fifth syntaxes, a closed path is assumed.
//
// - will require struct tag option to specify whether the path should be closed
//   or open, this unfortunately will make it a "static" path type
//
// [][2]float64

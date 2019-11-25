# wip

-----

## FQL

Package fql implements a tokenizer for FQL.

FQL stands for *Filtering Query Language*. The design of the query language is heavily inspired by [Ghost's](https://github.com/TryGhost/Ghost) GQL spec outlined [here](https://github.com/TryGhost/Ghost/issues/5604).

Example: `status:active;createdAt:>d1483228800`

-----

#### Syntax:

- A filter **expression** consists of one or more filter **rules** or **groups** separated by **logical operators**.
- A filter **group** is a filter **expression** wrapped in parentheses `( ... )` which can be used to define the order of evaluation.
- The **logical operators** can be one of the following:
	- `;` represents the logical **AND** operator.
	- `,` represents the logical **OR** operator.
- A filter **rule** consists of a **key**, the key-value **separator**, an optional **relational operator**, and a **value** in that order.
- The key-value **separator** is the colon character (`:`).
- The rule **key** represents a resource's field on which the filter rule should be applied. Keys can contain only alphanumeric characters and the underscore character `[_0-9a-Z]`.
- The rule **value** represents the value that will be compared against a resource's field indicated by the rule **key** using the **relational operator**. The value can be one of the following:
	- `null`
	- a **boolean**: either `true` or `false`.
	- a **number**: integer or float, can optionally be preceded by a unary operator `-` or `+`.
	- a **string**: delimited by double quotes `"`, can also contain escaped double quotes `\"`.
	- a **timestamp**: denoted by a preceding `d` (date) and represented by the Unix timestamp; that is, the number of seconds since the Unix Epoch on January 1st, 1970 UTC. For example the value `d1483228800` represents `2017-01-01 00:00:00 +0000 UTC`.
- The **relational operator** is used to compare the rule **value** against a resource's field. If no operator is provided the default `=` (is equal) is assumed. The **relational operator** can be one of the following, depending on the **value**:
	- `>`
	- `<`
	- `>=`
	- `<=`
	- `!`

NOTE: the operators `>`, `<`, `>=`, and `<=` are NOT applicable with `null`, **boolean**, or **string** values. Attempting to use these operators with values of the aforementioned types will result in an error being returned by the tokenizer.

-------------

#### Grammar:

The FQL grammar is specified below using the Extended Backus-Naur Form (EBNF):

```ebnf
filter_expr  = ( filter_rule | filter_group ) { seq_op filter_expr } ;
filter_group = '(' filter_expr ')' ;
		
filter_rule = filter_key ':' [ rel_op ] filter_val ;
filter_key  = identifier ;
filter_val  = bool_val | num_val | time_val | text_val | null_val ;
		
bool_val  = 'true' | 'false' ;
num_val   = [ unary_op ] ( int_val | float_val ) ;
time_val  = 'd' [ unary_op ] int_val ;
text_val  = interpreted_string_lit ;
null_val  = 'null' ;
int_val   = '0' | (( digit - '0' ) { digit } ) ;
float_val = float_lit ;
		
rel_op   = '>' | '>=' | '<' | '<=' | '!' ;
seq_op   = ',' | ';' ;
unary_op = '+' | '-' ;
		
digit = '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' ;
```

- The definitions of the following rules can be found in Go's language spec.

	```
	identifier = .
	interpreted_string_lit = .
	float_lit = .
	```


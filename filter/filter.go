package filter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/frk/fql"
	"github.com/frk/gosql"
)

type sqlWriter struct {
	// The underlying string builder used by sqlWriter.
	strings.Builder
	// The current position of parameters. (for figuring out "N" in "$N")
	p int
	// List of params to be passed to the query executing code.
	params []interface{}
}

type sqlNode interface {
	// Indicates whether or not the node can be followed by "AND" or "OR".
	canAndOr() bool
	// Writes the node with the given sqlWriter.
	write(w *sqlWriter)
}

// An sqlNode representing a binary expression.
type sqlBinary struct {
	col   string      // the LHS column
	op    string      // the binary comparison operator
	val   interface{} // the RHS value
	isany bool
}

func (sqlBinary) canAndOr() bool { return true }

func (n sqlBinary) write(w *sqlWriter) {
	w.WriteString(n.col)
	w.WriteString(" ")
	w.WriteString(n.op)
	w.WriteString(" ")
	w.WriteString(gosql.OrdinalParameters[w.p])

	w.params = append(w.params, n.val)
	w.p += 1
}

// An sqlNode representing a unary expression.
type sqlUnary struct {
	col  string // the LHS column
	pred string // the unary predicate
}

func (sqlUnary) canAndOr() bool { return true }

func (n sqlUnary) write(w *sqlWriter) {
	w.WriteString(n.col)
	w.WriteString(" ")
	w.WriteString(n.pred)
}

// An sqlNode representing a PostgreSQL specific text-search expression.
type sqlTextSearch struct {
	col string // the LHS ts_vector column
	val string // the search text
}

func (sqlTextSearch) canAndOr() bool { return true }

func (n sqlTextSearch) write(w *sqlWriter) {
	w.WriteString(n.col)
	w.WriteString(" @@ to_tsquery('simple', ")
	w.WriteString(gosql.OrdinalParameters[w.p])
	w.WriteString(")")

	w.params = append(w.params, toTSQuery(n.val))
	w.p += 1
}

// An sqlNode that represents the "AND" logical operator.
type sqlAnd struct{}

func (sqlAnd) canAndOr() bool { return false }

func (sqlAnd) write(w *sqlWriter) { w.WriteString(" AND ") }

// An sqlNode that represents the "OR" logical operator.
type sqlOr struct{}

func (sqlOr) canAndOr() bool { return false }

func (sqlOr) write(w *sqlWriter) { w.WriteString(" OR ") }

// An sqlNode representing the left/opening parenthesis "(".
type sqlLParen struct{}

func (sqlLParen) canAndOr() bool { return false }

func (sqlLParen) write(w *sqlWriter) { w.WriteByte('(') }

// An sqlNode representing the right/closing parenthesis ")".
type sqlRParen struct{}

func (sqlRParen) canAndOr() bool { return true }

func (sqlRParen) write(w *sqlWriter) { w.WriteByte(')') }

// The Constructor type can be used to dynamically construct a "filter" for an SQL query.
//
// The Constructor type implements the gosql.Filter and the gosql.FilterConstructor interface.
type Constructor struct {
	filter filter
	// The colmap field maps "public facing keys" to valid column identifiers
	// of the relation with which the Filter instance is associated.
	// Note that colmap, once set by Init, is NOT to be modified, it is read-only.
	colmap map[string]Column
	// The identifier of the ts_vector column that can be used for full text search.
	tscol string
	// Indicates whether or not an error should be returned if any of the
	// Constructor's methods encounter a column that has no entry in the colmap.
	strict bool
}

// The Column type XXX
type Column struct {
	// Name is the name of the column.
	Name string
	// ConvertValue, when set, will be used convert
	// the column's FQL-parsed value.
	ConvertValue func(any) (any, error)
	// IsNULLable indicates whether or not the column's NULLable.
	IsNULLable bool
}

var (
	// make sure Constructor implements gosql.FilterConstructor
	_ gosql.FilterConstructor = (*Constructor)(nil)
)

// canAndOr reports whether or not the sqlAnd/sqlOr nodes can be used
// given the current state of the sqlNode list.
func (c *Constructor) canAndOr() bool {
	return len(c.filter.where) > 0 && c.filter.where[len(c.filter.where)-1].canAndOr()
}

// StrictSwitch switches the strict mode on and off.
func (c *Constructor) StrictSwitch() { c.strict = !c.strict }

// Init initializes the Constructor's colmap and tscol fields using the given values.
//
// The Init method implements part of the gosql.FilterConstructor interface.
func (c *Constructor) Init(colmap map[string]string, tscol string) {
	c.tscol = tscol
	c.colmap = make(map[string]Column, len(colmap))
	for k, v := range colmap {
		c.colmap[k] = Column{Name: v}
	}
}

// InitV2 initializes the Constructor's colmap and tscol fields using the given values.
//
// The Init method implements part of the gosql.FilterConstructor interface.
func (c *Constructor) InitV2(colmap map[string]Column, tscol string) {
	c.tscol = tscol
	c.colmap = colmap
}

// Col prepares a new, column-specific predicate for the WHERE clause.
// The column is assumed to already be vetted.
//
// The Col method implements part of the gosql.FilterConstructor interface.
func (c *Constructor) Col(column string, op string, value interface{}) {
	switch strings.ToUpper(op) {
	case "ANY":
		switch v := value.(type) {
		case []int64:
			c.ColAnyInt64s(column, v)
		case []int:
			c.ColAnyInts(column, v)
		case []int16:
			c.ColAnyInt16s(column, v)
		case []string:
			c.ColAnyStrings(column, v)
		default:
			// XXX ignore if unsupported type
		}
	case "IN":
		switch v := value.(type) {
		case []int64:
			c.ColInInt64s(column, v)
		case []int:
			c.ColInInts(column, v)
		case []int16:
			c.ColInInt16s(column, v)
		case []string:
			c.ColInStrings(column, v)
		default:
			// XXX ignore if unsupported type
		}

	default:
		if c.canAndOr() {
			c.filter.where = append(c.filter.where, sqlAnd{})
		}
		switch value {
		case nil:
			switch op {
			case "=":
				c.filter.where = append(c.filter.where, sqlUnary{col: column, pred: "IS NULL"})
			case "<>":
				c.filter.where = append(c.filter.where, sqlUnary{col: column, pred: "IS NOT NULL"})
			}
		case true:
			switch op {
			case "=":
				c.filter.where = append(c.filter.where, sqlUnary{col: column, pred: "IS TRUE"})
			case "<>":
				c.filter.where = append(c.filter.where, sqlUnary{col: column, pred: "IS NOT TRUE"})
			}
		case false:
			switch op {
			case "=":
				c.filter.where = append(c.filter.where, sqlUnary{col: column, pred: "IS FALSE"})
			case "<>":
				c.filter.where = append(c.filter.where, sqlUnary{col: column, pred: "IS NOT FALSE"})
			}
		default:
			c.filter.where = append(c.filter.where, sqlBinary{
				col: column,
				op:  op,
				val: value,
			})
		}
	}
}

// TextSearch prepares a new, text-search predicate for the WHERE clause.
func (c *Constructor) TextSearch(value string) {
	if len(c.tscol) > 0 {
		if c.canAndOr() {
			c.filter.where = append(c.filter.where, sqlAnd{})
		}
		c.filter.where = append(c.filter.where, sqlTextSearch{col: c.tscol, val: value})
	}
}

// And adds the AND logical operator argument for the WHERE clause. If the
// given nest function is not nil it will be executed and its result will be
// wrapped in parentheses.
//
// The And method implements part of the gosql.FilterConstructor interface.
func (c *Constructor) And(nest func()) {
	if c.canAndOr() {
		c.filter.where = append(c.filter.where, sqlAnd{})
	}

	if nest != nil {
		c.filter.where = append(c.filter.where, sqlLParen{})
		nest()
		c.filter.where = append(c.filter.where, sqlRParen{})
	}
}

// Or adds the OR logical operator argument for the WHERE clause. If the
// given nest function is not nil it will be executed and its result will be
// wrapped in parentheses.
//
// The Or method implements part of the gosql.FilterConstructor interface.
func (c *Constructor) Or(nest func()) {
	if c.canAndOr() {
		c.filter.where = append(c.filter.where, sqlOr{})
	}

	if nest != nil {
		c.filter.where = append(c.filter.where, sqlLParen{})
		nest()
		c.filter.where = append(c.filter.where, sqlRParen{})
	}
}

// OrderBy adds a new argument to the ORDER BY clause.
//
// The column is assumed to already be vetted.
func (c *Constructor) OrderBy(column string, desc, nullsfirst bool) {
	if len(c.filter.orderby) > 0 {
		c.filter.orderby += ", "
	}
	c.filter.orderby += column

	if desc {
		c.filter.orderby += " DESC"
	} else {
		c.filter.orderby += " ASC"
	}

	if nullsfirst {
		c.filter.orderby += " NULLS FIRST"
	} else {
		c.filter.orderby += " NULLS LAST"
	}
}

// OrderByV2 adds a new argument to the ORDER BY clause.
//
// The column is assumed to already be vetted.
func (c *Constructor) OrderByV2(column string, desc, nullable, nullsfirst bool) {
	if len(c.filter.orderby) > 0 {
		c.filter.orderby += ", "
	}
	c.filter.orderby += column

	if desc {
		c.filter.orderby += " DESC"
	} else {
		c.filter.orderby += " ASC"
	}

	if nullable {
		if nullsfirst {
			c.filter.orderby += " NULLS FIRST"
		} else {
			c.filter.orderby += " NULLS LAST"
		}
	}
}

// Limit sets the given value as the argument for the LIMIT clause.
func (c *Constructor) Limit(count int64) {
	if count > 0 {
		c.filter.limit = count
	}
}

// Offset sets the given value as the argument for the OFFSET clause.
func (c *Constructor) Offset(start int64) {
	if start > 0 {
		c.filter.offset = start
	}
}

// Page, using the given value, calculates the argument for the OFFSET clause.
func (c *Constructor) Page(page int64) {
	if page > 0 && c.filter.limit >= 0 {
		c.filter.offset = c.filter.limit * (page - 1)
	}
}

// UnmarshalFQL unmarshals the given string using the github.com/frk/fql package.
// If the given string is not valid FQL an error will be returned. If the Constructor
// is in strict mode and the column keys in the FQL do not have an entry in the
// Constructor's colmap, the UnknownColumnKeyError will be returned.
func (c *Constructor) UnmarshalFQL(str string) error {
	str = strings.Trim(str, ";,") // sanitize

	z := fql.NewTokenizer(str)
	for {
		tok, err := z.Next()
		if err != nil {
			if err != fql.EOF {
				return err
			}
			break
		}

		switch tok {
		case fql.LPAREN:
			if c.canAndOr() {
				c.filter.where = append(c.filter.where, sqlAnd{})
			}
			c.filter.where = append(c.filter.where, sqlLParen{})
		case fql.RPAREN:
			c.filter.where = append(c.filter.where, sqlRParen{})
		case fql.AND:
			c.filter.where = append(c.filter.where, sqlAnd{})
		case fql.OR:
			c.filter.where = append(c.filter.where, sqlOr{})
		case fql.RULE:
			rule := z.Rule()
			if col, ok := c.colmap[rule.Key]; ok {
				val := rule.Val
				if val == fql.Null {
					val = nil
				}

				if col.ConvertValue != nil {
					v, err := col.ConvertValue(val)
					if err != nil {
						return err
					}
					val = v
				}

				c.Col(col.Name, cmpop2string[rule.Cmp], val)
			} else if c.strict {
				return UnknownColumnKeyError{Key: rule.Key}
			}
		}
	}
	return nil
}

// UnmarshalSort parses the given string as a comma separated list of column
// keys that can optionally be preceded by a hyphen to indicate the descending
// sort order. The keys are then used to build the ORDER BY clause of the
// filter. Empty items between commas in the str are ignored. If the Constructor
// is in strict mode and the column keys in the input do not have an entry in
// the Constructor's colmap, the UnknownColumnKeyError will be returned.
//
// Example value:
//
//	"-created_at,label"
func (c *Constructor) UnmarshalSort(str string) error {
	start, end := 0, len(str)
	for start < end {
		pos := start
		for pos < end && str[pos] != ',' {
			pos += 1
		}

		desc := (str[start] == '-')
		if desc {
			start += 1
		}

		if key := str[start:pos]; len(key) > 0 {
			if col, ok := c.colmap[key]; ok {
				c.OrderByV2(col.Name, desc, col.IsNULLable, false)
			} else if c.strict {
				return UnknownColumnKeyError{Key: key}
			}
		}
		start = pos + 1
	}
	return nil
}

// Filter returns the constructed gosql.Filter instance.
func (c *Constructor) Filter() gosql.Filter {
	return &c.filter
}

// CountFilter returns an alternative implementation of the constructed gosql.Filter,
// one that that omits clauses that are unnecessary in the context of a SELECT COUNT
// query, like, for example, the LIMIT clause.
func (c *Constructor) CountFilter() gosql.Filter {
	return (*countFilter)(&c.filter)
}

// filter is the result of the Constructor. The currently supported SQL clauses
// that the filter type can produce are:
//   - WHERE
//   - ORDER BY
//   - LIMIT
//   - OFFSET
type filter struct {
	// The arguments for the WHERE clause as a list of sqlNodes.
	where []sqlNode
	// The arguments for the ORDER BY clause as plain, valid SQL string.
	orderby string
	// The argument for the LIMIT clause.
	limit int64
	// The argument for the OFFSET clause.
	offset int64
}

var (

	// make sure filter implements gosql.Filter
	_ gosql.Filter = (*filter)(nil)
)

// ToSQL returns the SQL representation of the filter together with a list of
// parameters to be used with the SQL string. The given ppos argument specifies
// the current parameter position of the caller.
//
// The ToSQL method implements the gosql.Filter interface.
func (f *filter) ToSQL(ppos int) (filterString string, params []interface{}) {
	w := sqlWriter{p: ppos}
	if len(f.where) > 0 {
		w.WriteString(" WHERE ")
		for _, node := range f.where {
			node.write(&w)
		}
	}

	if len(f.orderby) > 0 {
		w.WriteString(" ORDER BY ")
		w.WriteString(f.orderby)
	}
	if f.limit > 0 {
		w.WriteString(" LIMIT ")
		w.WriteString(strconv.FormatInt(f.limit, 10))
	}
	if f.offset > 0 {
		w.WriteString(" OFFSET ")
		w.WriteString(strconv.FormatInt(f.offset, 10))
	}

	return w.String(), w.params
}

// countFilter is an alternative implementation of gosql.Filter that is
// intended to be used with SELECT COUNT queries.
type countFilter filter

var (
	// make sure countFilter implements gosql.Filter
	_ gosql.Filter = (*countFilter)(nil)
)

// ToSQL implements the gosql.Filter interface.
func (f *countFilter) ToSQL(ppos int) (filterString string, params []interface{}) {
	w := sqlWriter{p: ppos}
	if len(f.where) > 0 {
		w.WriteString(" WHERE ")
		for _, node := range f.where {
			node.write(&w)
		}
	}

	return w.String(), w.params
}

// The UnknownColumnKeyError indicates that a provided key has no matching
// entry in the colmap of the Filter from which the error was retruned.
type UnknownColumnKeyError struct {
	Key string
}

// Error implements the error interface.
func (e UnknownColumnKeyError) Error() string {
	return fmt.Sprintf("filter: unknown column key %q.", e.Key)
}

// toTSQuery re-formats the given string as a postgresql ts_query. Each lexeme in
// the returned ts_query is labeled with `:*` to allow for prefix matching. The input
// is expected to be very basic, containg one or more search terms. For input that's
// more complex, the output is not guaranteed to be a valid ts_query.
func toTSQuery(str string) (out string) {
	list := strings.Split(strings.TrimSpace(str), " ")
	if len(list) == 0 {
		return
	}

	if str := list[0]; str != "" {
		out += str + ":*"
	}
	for _, str := range list[1:] {
		if str != "" {
			out += " & " + str + ":*"
		}
	}
	return out
}

// maps FQL comparison operators to SQL ones.
var cmpop2string = map[fql.CmpOp]string{
	fql.CmpEq: "=",
	fql.CmpNe: "<>",
	fql.CmpGt: ">",
	fql.CmpLt: "<",
	fql.CmpGe: ">=",
	fql.CmpLe: "<=",
}

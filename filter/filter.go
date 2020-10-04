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
	col string      // the LHS column
	op  string      // the binary comparison operator
	val interface{} // the RHS value
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

// The Filter type can be used to dynamically construct a "filter" for an SQL query.
// The supported SQL clauses that the Filter type can produce are:
//     - WHERE
//     - ORDER BY
//     - LIMIT
//     - OFFSET
//
// The Filter type implements the gosql.Filter and the gosql.FilterConstructor interface.
type Filter struct {
	// The colmap field maps "public facing keys" to valid column identifiers
	// of the relation with which the Filter instance is associated.
	// Note that colmap, once set by Init, is NOT to be modified, it is read-only.
	colmap map[string]string
	// The identifier of the ts_vector column that can be used for full text search.
	tscol string
	// The arguments for the WHERE clause as a list of sqlNodes.
	where []sqlNode
	// The arguments for the ORDER BY clause as plain, valid SQL string.
	orderby string
	// The argument for the LIMIT clause.
	limit int64
	// The argument for the OFFSET clause.
	offset int64
	// Indicates whether or not an error should be returned if any of the
	// Filter's methods encounter a column that has no entry in the colmap.
	strict bool
}

var (
	// make sure Filter implements gosql.Filter
	_ gosql.Filter = (*Filter)(nil)
	// make sure Filter implements gosql.FilterConstructor
	_ gosql.FilterConstructor = (*Filter)(nil)
)

// canAndOr reports whether or not the sqlAnd/sqlOr nodes can be used
// given the current state of the sqlNode list.
func (f *Filter) canAndOr() bool {
	return len(f.where) > 0 && f.where[len(f.where)-1].canAndOr()
}

// Strict switches the strict mode on and off.
func (f *Filter) Strict() { f.strict = !f.strict }

// ToSQL returns the SQL representation of the Filter together with a list of
// parameters to be used with the SQL string. The given ppos argument specifies
// the current parameter position of the caller.
//
// The ToSQL method implements the gosql.Filter interface.
func (f *Filter) ToSQL(ppos int) (filterString string, params []interface{}) {
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

// Init initializes the Filter's colmap and tscol fields using the given values.
//
// The Init method implements part of the gosql.FilterConstructor interface.
func (f *Filter) Init(colmap map[string]string, tscol string) {
	f.colmap = colmap
	f.tscol = tscol
}

// Col prepares a new, column-specific predicate for the WHERE clause.
// The column is assumed to already be vetted.
//
// The Col method implements part of the gosql.FilterConstructor interface.
func (f *Filter) Col(column string, op string, value interface{}) {
	if f.canAndOr() {
		f.where = append(f.where, sqlAnd{})
	}
	switch value {
	case nil:
		switch op {
		case "=":
			f.where = append(f.where, sqlUnary{col: column, pred: "IS NULL"})
		case "<>":
			f.where = append(f.where, sqlUnary{col: column, pred: "IS NOT NULL"})
		}
	case true:
		switch op {
		case "=":
			f.where = append(f.where, sqlUnary{col: column, pred: "IS TRUE"})
		case "<>":
			f.where = append(f.where, sqlUnary{col: column, pred: "IS NOT TRUE"})
		}
	case false:
		switch op {
		case "=":
			f.where = append(f.where, sqlUnary{col: column, pred: "IS FALSE"})
		case "<>":
			f.where = append(f.where, sqlUnary{col: column, pred: "IS NOT FALSE"})
		}
	default:
		f.where = append(f.where, sqlBinary{col: column, op: op, val: value})
	}
}

// TextSearch prepares a new, text-search predicate for the WHERE clause.
func (f *Filter) TextSearch(value string) *Filter {
	if len(f.tscol) > 0 {
		if f.canAndOr() {
			f.where = append(f.where, sqlAnd{})
		}
		f.where = append(f.where, sqlTextSearch{col: f.tscol, val: value})
	}
	return f
}

// And adds the AND logical operator argument for the WHERE clause. If the
// given nest function is not nil it will be executed and its result will be
// wrapped in parentheses.
//
// The And method implements part of the gosql.FilterConstructor interface.
func (f *Filter) And(nest func()) {
	if f.canAndOr() {
		f.where = append(f.where, sqlAnd{})
	}

	if nest != nil {
		f.where = append(f.where, sqlLParen{})
		nest()
		f.where = append(f.where, sqlRParen{})
	}
}

// Or adds the OR logical operator argument for the WHERE clause. If the
// given nest function is not nil it will be executed and its result will be
// wrapped in parentheses.
//
// The Or method implements part of the gosql.FilterConstructor interface.
func (f *Filter) Or(nest func()) {
	if f.canAndOr() {
		f.where = append(f.where, sqlOr{})
	}

	if nest != nil {
		f.where = append(f.where, sqlLParen{})
		nest()
		f.where = append(f.where, sqlRParen{})
	}
}

// OrderBy adds a new argument to the ORDER BY clause.
// The column is assumed to already be vetted.
func (f *Filter) OrderBy(column string, desc, nullsfirst bool) {
	if len(f.orderby) > 0 {
		f.orderby += ", "
	}
	f.orderby += column

	if desc {
		f.orderby += " DESC"
	} else {
		f.orderby += " ASC"
	}

	if nullsfirst {
		f.orderby += " NULLS FIRST"
	} else {
		f.orderby += " NULLS LAST"
	}
}

// Limit sets the given value as the argument for the LIMIT clause.
func (f *Filter) Limit(count int64) {
	if count > 0 {
		f.limit = count
	}
}

// Offset sets the given value as the argument for the OFFSET clause.
func (f *Filter) Offset(start int64) {
	if start > 0 {
		f.offset = start
	}
}

// Page, using the given value, calculates the argument for the OFFSET clause.
func (f *Filter) Page(page int64) {
	if page > 0 && f.limit >= 0 {
		f.offset = f.limit * (page - 1)
	}
}

// UnmarshalFQL unmarshals the given string using the github.com/frk/fql package.
// If the given string is not valid FQL an error will be returned. If the Filter
// is in strict mode and the column keys in the FQL do not have an entry in the
// Filter's colmap, the UnknownColumnKeyError will be returned.
func (f *Filter) UnmarshalFQL(str string) error {
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
			f.where = append(f.where, sqlLParen{})
		case fql.RPAREN:
			f.where = append(f.where, sqlRParen{})
		case fql.AND:
			f.where = append(f.where, sqlAnd{})
		case fql.OR:
			f.where = append(f.where, sqlOr{})
		case fql.RULE:
			rule := z.Rule()
			if col, ok := f.colmap[rule.Key]; ok {
				val := rule.Val
				if val == fql.Null {
					val = nil
				}

				f.Col(col, cmpop2string[rule.Cmp], val)
			} else if f.strict {
				return UnknownColumnKeyError{Key: rule.Key}
			}
		}
	}
	return nil
}

// UnmarshalSort parses the given string as a comma separated list of column
// keys that can optionally be preceded by a hyphen to indicate the descending
// sort order. The keys are then used to build the ORDER BY clause of the
// Filter. Empty items between commas in the str are ignored. If the Filter is in
// strict mode and the column keys in the input do not have an entry in the Filter's
// colmap, the UnknownColumnKeyError will be returned.
//
// Example value:
//
//	"-created_at,label"
//
func (f *Filter) UnmarshalSort(str string) error {
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
			if col, ok := f.colmap[key]; ok {
				f.OrderBy(col, desc, false)
			} else if f.strict {
				return UnknownColumnKeyError{Key: key}
			}
		}
		start = pos + 1
	}
	return nil
}

// filterForCount is an alternative implementation of gosql.Filter that is
// intended to be used with SELECT COUNT queries.
type filterForCount Filter

var (
	// make sure filterForCount implements gosql.Filter
	_ gosql.Filter = (*filterForCount)(nil)
)

// ForCount returns a gosql.Filter implementation that produces an output that
// omits clauses that are unnecessary in the context of a SELECT COUNT query,
// like, for example, the LIMIT clause.
func (f *Filter) ForCount() gosql.Filter {
	return (*filterForCount)(f)
}

// ToSQL implements the gosql.Filter interface.
func (f *filterForCount) ToSQL(ppos int) (filterString string, params []interface{}) {
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

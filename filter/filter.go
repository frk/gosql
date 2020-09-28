package filter

import (
	"strconv"
	"strings"

	"github.com/frk/fql"
)

type Filter struct {
	where    string
	orderby  string
	limit    int64
	offset   int64
	params   []interface{}
	canAndor bool

	colmap map[string]string
	tscol  string
}

func (f *Filter) Init(colmap map[string]string, tscol string) {
	f.colmap = colmap
	f.tscol = tscol
}

func (f *Filter) Col(column string, op string, value interface{}) *Filter {
	if f.canAndor {
		f.where += ` AND `
	}

	var comparison string
	switch value {
	case nil:
		switch op {
		case "=":
			comparison = `IS NULL`
		case "<>":
			comparison = `IS NOT NULL`
		}
	case true:
		switch op {
		case "=":
			comparison = `IS TRUE`
		case "<>":
			comparison = `IS NOT TRUE`
		}
	case false:
		switch op {
		case "=":
			comparison = `IS FALSE`
		case "<>":
			comparison = `IS NOT FALSE`
		}
	default:
		f.params = append(f.params, value)
		comparison = op + ` $` + strconv.Itoa(len(f.params))
	}

	f.where += column + ` ` + comparison
	f.canAndor = true
	return f
}

func (f *Filter) And(nest func()) {
	if !f.canAndor {
		return
	}

	f.where += ` AND `
	f.canAndor = false
	if nest != nil {
		f.where += `(`
		nest()
		f.where += `)`
		f.canAndor = true
	}
}

func (f *Filter) Or(nest func()) {
	if !f.canAndor {
		return
	}

	f.where += ` OR `
	f.canAndor = false
	if nest != nil {
		f.where += `(`
		nest()
		f.where += `)`
		f.canAndor = true
	}
}

func (f *Filter) Limit(count int64) {
	if count > 0 {
		f.limit = count
	}
}

func (f *Filter) Offset(start int64) {
	if start > 0 {
		f.offset = start
	}
}

func (f *Filter) Page(page int64) {
	if page > 0 && f.limit >= 0 {
		f.offset = f.limit * (page - 1)
	}
}

// UnmarshalFQL
func (f *Filter) UnmarshalFQL(fqlString string, colmap map[string]string, strict bool) error {
	fqlString = strings.Trim(fqlString, ";,")
	z := fql.NewTokenizer(fqlString)
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
			f.GroupStart()
		case fql.RPAREN:
			f.GroupEnd()
		case fql.AND:
			f.And(nil)
		case fql.OR:
			f.Or(nil)
		case fql.RULE:
			rule := z.Rule()
			if col, ok := colmap[rule.Key]; ok {
				val := rule.Val
				if val == fql.Null {
					val = nil
				}

				f.Col(col, cmpop2string[rule.Cmp], val)
			} else if strict {
				// TODO return unknown column error
			}
		}
	}
	return nil
}

// UnmarshalSort parses sortString as a comma separated list of column keys that can
// optionally be preceded by a hyphen to indicate the descending sort order. The
// keys are then used to build the Order By clause of the filter.
// Empty items between commas in the sortString are ignored.
func (f *Filter) UnmarshalSort(sortString string, colmap map[string]string, strict bool) error {
	start, end := 0, len(sortString)
	for start < end {
		pos := start
		for pos < end && sortString[pos] != ',' {
			pos += 1
		}

		desc := (sortString[start] == '-')
		if desc {
			start += 1
		}

		if key := sortString[start:pos]; len(key) > 0 {
			if col, ok := colmap[key]; ok {
				f.OrderBy(col, desc, false)
			} else if strict {
				// TODO retrun unknown column error
			}
		}
		start = pos + 1
	}
	return nil
}

func (f *Filter) TextSearch(document, value string) *Filter {
	if f.canAndor {
		f.where += ` AND `
	}

	f.params = append(f.params, value)
	pos := `$` + strconv.Itoa(len(f.params))

	f.where += document + ` @@ to_tsquery('simple', ` + pos + `)`
	f.canAndor = true
	return f
}

func (f *Filter) OrderBy(column string, desc, nullsfirst bool) {
	if len(f.orderby) > 0 {
		f.orderby += `, `
	}
	f.orderby += column

	if desc {
		f.orderby += ` DESC`
	} else {
		f.orderby += ` ASC`
	}

	if nullsfirst {
		f.orderby += ` NULLS FIRST`
	} else {
		f.orderby += ` NULLS LAST`
	}
}

// Params returns a slice of the collected params which should be passed
// directly to the corresponding query.
func (f *Filter) Params() []interface{} {
	return f.params
}

func (f *Filter) ToSQL() (out string) {
	if len(f.where) > 0 {
		out += ` WHERE ` + f.where
	}
	if len(f.orderby) > 0 {
		out += ` ORDER BY ` + f.orderby
	}
	if f.limit > 0 {
		out += ` LIMIT ` + strconv.FormatInt(f.limit, 10)
	}
	if f.offset > 0 {
		out += ` OFFSET ` + strconv.FormatInt(f.offset, 10)
	}
	return out
}

func (f *Filter) ToSQLWhereClause() (out string) {
	if len(f.where) > 0 {
		out += ` WHERE ` + f.where
	}
	return out
}

func (f *Filter) GroupStart() *Filter {
	f.where += `(`
	return f
}

func (f *Filter) GroupEnd() *Filter {
	f.where += `)`
	return f
}

// formatTSQuery returns a valid postgresql tsquery from the given string s. Each
// lexeme in the returned query is labeled with `:*` to allow for prefix matching.
func formatTSQuery(s string) (out string) {
	list := strings.Split(strings.TrimSpace(s), " ")
	if len(list) == 0 {
		return
	}

	if s := list[0]; s != "" {
		out += s + ":*"
	}
	for _, s := range list[1:] {
		if s != "" {
			out += " & " + s + ":*"
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

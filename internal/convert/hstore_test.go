package convert

import (
	"database/sql"
	"testing"
)

func TestHStore_ValuerAndScanner(t *testing.T) {
	A := strptr

	test_valuer_with_scanner{{
		valuer: func() interface{} {
			return new(HStoreFromStringMap)
		},
		scanner: func() (interface{}, interface{}) {
			v := &HStoreToStringMap{Val: new(map[string]string)}
			return v, v.Val
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "hstore", input: nil, output: new(map[string]string)},
			{typ: "hstore", input: map[string]string{}, output: &map[string]string{}},
			{
				typ:    "hstore",
				input:  map[string]string{"a": "1", "b": "2"},
				output: &map[string]string{"a": "1", "b": "2"}},
			{
				typ:    "hstore",
				input:  map[string]string{"a": `'`},
				output: &map[string]string{"a": `'`}},
			{
				typ:    "hstore",
				input:  map[string]string{"a": `"`, "b": `\`},
				output: &map[string]string{"a": `"`, "b": `\`}},
			{
				typ:    "hstore",
				input:  map[string]string{"a": `\"`, "b": `\\`},
				output: &map[string]string{"a": `\"`, "b": `\\`}},
			{
				typ:    "hstore",
				input:  map[string]string{"a": "\a", "b": "\b", "c": "\t", "d": "\n"},
				output: &map[string]string{"a": "\a", "b": "\b", "c": "\t", "d": "\n"}},
			{
				typ:    "hstore",
				input:  map[string]string{"text": `foo' "bar", baz \quux`},
				output: &map[string]string{"text": `foo' "bar", baz \quux`}},
		},
	}, {
		valuer: func() interface{} {
			return new(HStoreFromStringPtrMap)
		},
		scanner: func() (interface{}, interface{}) {
			v := &HStoreToStringPtrMap{Val: new(map[string]*string)}
			return v, v.Val
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "hstore", input: nil, output: new(map[string]*string)},
			{typ: "hstore", input: map[string]*string{}, output: &map[string]*string{}},
			{
				typ:    "hstore",
				input:  map[string]*string{"a": A("1"), "b": A("2"), "c": nil},
				output: &map[string]*string{"a": A("1"), "b": A("2"), "c": nil}},
			{
				typ:    "hstore",
				input:  map[string]*string{"a": A(`'`)},
				output: &map[string]*string{"a": A(`'`)}},
			{
				typ:    "hstore",
				input:  map[string]*string{"a": A(`"`), "b": A(`\`)},
				output: &map[string]*string{"a": A(`"`), "b": A(`\`)}},
			{
				typ:    "hstore",
				input:  map[string]*string{"a": A(`\"`), "b": A(`\\`)},
				output: &map[string]*string{"a": A(`\"`), "b": A(`\\`)}},
			{
				typ:    "hstore",
				input:  map[string]*string{"a": A("\a"), "b": A("\b"), "c": A("\t"), "d": A("\n")},
				output: &map[string]*string{"a": A("\a"), "b": A("\b"), "c": A("\t"), "d": A("\n")}},
			{
				typ:    "hstore",
				input:  map[string]*string{"text": A(`foo' "bar", baz \quux`)},
				output: &map[string]*string{"text": A(`foo' "bar", baz \quux`)}},
		},
	}, {
		valuer: func() interface{} {
			return new(HStoreFromNullStringMap)
		},
		scanner: func() (interface{}, interface{}) {
			v := &HStoreToNullStringMap{Val: new(map[string]sql.NullString)}
			return v, v.Val
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "hstore", input: nil, output: new(map[string]sql.NullString)},
			{typ: "hstore", input: map[string]sql.NullString{}, output: &map[string]sql.NullString{}},
			{
				typ:    "hstore",
				input:  map[string]sql.NullString{"a": {"1", true}, "b": {"2", true}, "c": {"", false}},
				output: &map[string]sql.NullString{"a": {"1", true}, "b": {"2", true}, "c": {"", false}}},
			{
				typ:    "hstore",
				input:  map[string]sql.NullString{"a": {`'`, true}},
				output: &map[string]sql.NullString{"a": {`'`, true}}},
			{
				typ:    "hstore",
				input:  map[string]sql.NullString{"a": {`"`, true}, "b": {`\`, true}},
				output: &map[string]sql.NullString{"a": {`"`, true}, "b": {`\`, true}}},
			{
				typ:    "hstore",
				input:  map[string]sql.NullString{"a": {`\"`, true}, "b": {`\\`, true}},
				output: &map[string]sql.NullString{"a": {`\"`, true}, "b": {`\\`, true}}},
			{
				typ: "hstore",
				input: map[string]sql.NullString{"a": {"\a", true}, "b": {"\b", true},
					"c": {"\t", true}, "d": {"\n", true}},
				output: &map[string]sql.NullString{"a": {"\a", true}, "b": {"\b", true},
					"c": {"\t", true}, "d": {"\n", true}}},
			{
				typ:    "hstore",
				input:  map[string]sql.NullString{"text": {`foo' "bar", baz \quux`, true}},
				output: &map[string]sql.NullString{"text": {`foo' "bar", baz \quux`, true}}},
		},
	}}.execute(t)
}

func TestHStore_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil // string
		},
		rows: []test_valuer_row{
			{typ: "hstore", in: nil, want: nil},
			{typ: "hstore", in: "", want: strptr("")},
			{
				typ:  "hstore",
				in:   `"a"=>"1", "b"=>"2"`,
				want: strptr(`"a"=>"1", "b"=>"2"`)},
			{
				typ:  "hstore",
				in:   `"text"=>"foo' \"bar\", baz \\quux"`,
				want: strptr(`"text"=>"foo' \"bar\", baz \\quux"`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		rows: []test_valuer_row{
			{typ: "hstore", in: nil, want: nil},
			{typ: "hstore", in: "", want: strptr("")},
			{
				typ:  "hstore",
				in:   []byte(`"a"=>"1", "b"=>"2"`),
				want: strptr(`"a"=>"1", "b"=>"2"`)},
			{
				typ:  "hstore",
				in:   []byte(`"text"=>"foo' \"bar\", baz \\quux"`),
				want: strptr(`"text"=>"foo' \"bar\", baz \\quux"`)},
		},
	}}.execute(t)
}

func TestHStore_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			v := new(string)
			return v, v
		},
		rows: []test_scanner_row{
			{typ: "hstore", in: ``, want: strptr(``)},
			{
				typ:  "hstore",
				in:   `"a"=>"1"`,
				want: strptr(`"a"=>"1"`)},
			{
				typ:  "hstore",
				in:   `"a"=>"1", "b"=>"2","c"=>"3"`,
				want: strptr(`"a"=>"1", "b"=>"2", "c"=>"3"`)},
			{
				typ:  "hstore",
				in:   `"a"=>"1", "b"=>NULL,"c"=>NULL`,
				want: strptr(`"a"=>"1", "b"=>NULL, "c"=>NULL`)},
			{
				typ:  "hstore",
				in:   `"a"=>"1", "b"=>"\\\"","c"=>"\\\\\\\\"`,
				want: strptr(`"a"=>"1", "b"=>"\\\"", "c"=>"\\\\\\\\"`)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			v := new([]byte)
			return v, v
		},
		rows: []test_scanner_row{
			{typ: "hstore", in: ``, want: bytesptr(``)},
			{
				typ:  "hstore",
				in:   `"a"=>"1"`,
				want: bytesptr(`"a"=>"1"`)},
			{
				typ:  "hstore",
				in:   `"a"=>"1", "b"=>"2","c"=>"3"`,
				want: bytesptr(`"a"=>"1", "b"=>"2", "c"=>"3"`)},
			{
				typ:  "hstore",
				in:   `"a"=>"1", "b"=>NULL,"c"=>NULL`,
				want: bytesptr(`"a"=>"1", "b"=>NULL, "c"=>NULL`)},
			{
				typ:  "hstore",
				in:   `"a"=>"1", "b"=>"\\\"","c"=>"\\\\\\\\"`,
				want: bytesptr(`"a"=>"1", "b"=>"\\\"", "c"=>"\\\\\\\\"`)},
		},
	}}.execute(t)
}

package pg2go

import (
	"database/sql"
	"testing"
)

func TestHStore(t *testing.T) {
	A := strptr

	testlist{{
		valuer: func() interface{} {
			return new(HStoreFromStringMap)
		},
		scanner: func() (interface{}, interface{}) {
			v := HStoreToStringMap{Val: new(map[string]string)}
			return v, v.Val
		},
		data: []testdata{
			{input: nil, output: map[string]string(nil)},
			{input: map[string]string{}, output: map[string]string{}},
			{
				input:  map[string]string{"a": "1", "b": "2"},
				output: map[string]string{"a": "1", "b": "2"}},
			{
				input:  map[string]string{"a": `'`},
				output: map[string]string{"a": `'`}},
			{
				input:  map[string]string{"a": `"`, "b": `\`},
				output: map[string]string{"a": `"`, "b": `\`}},
			{
				input:  map[string]string{"a": `\"`, "b": `\\`},
				output: map[string]string{"a": `\"`, "b": `\\`}},
			{
				input:  map[string]string{"a": "\a", "b": "\b", "c": "\t", "d": "\n"},
				output: map[string]string{"a": "\a", "b": "\b", "c": "\t", "d": "\n"}},
			{
				input:  map[string]string{"text": `foo' "bar", baz \quux`},
				output: map[string]string{"text": `foo' "bar", baz \quux`}},
		},
	}, {
		valuer: func() interface{} {
			return new(HStoreFromStringPtrMap)
		},
		scanner: func() (interface{}, interface{}) {
			v := HStoreToStringPtrMap{Val: new(map[string]*string)}
			return v, v.Val
		},
		data: []testdata{
			{input: nil, output: map[string]*string(nil)},
			{input: map[string]*string{}, output: map[string]*string{}},
			{
				input:  map[string]*string{"a": A("1"), "b": A("2"), "c": nil},
				output: map[string]*string{"a": A("1"), "b": A("2"), "c": nil}},
			{
				input:  map[string]*string{"a": A(`'`)},
				output: map[string]*string{"a": A(`'`)}},
			{
				input:  map[string]*string{"a": A(`"`), "b": A(`\`)},
				output: map[string]*string{"a": A(`"`), "b": A(`\`)}},
			{
				input:  map[string]*string{"a": A(`\"`), "b": A(`\\`)},
				output: map[string]*string{"a": A(`\"`), "b": A(`\\`)}},
			{
				input:  map[string]*string{"a": A("\a"), "b": A("\b"), "c": A("\t"), "d": A("\n")},
				output: map[string]*string{"a": A("\a"), "b": A("\b"), "c": A("\t"), "d": A("\n")}},
			{
				input:  map[string]*string{"text": A(`foo' "bar", baz \quux`)},
				output: map[string]*string{"text": A(`foo' "bar", baz \quux`)}},
		},
	}, {
		valuer: func() interface{} {
			return new(HStoreFromNullStringMap)
		},
		scanner: func() (interface{}, interface{}) {
			v := HStoreToNullStringMap{Val: new(map[string]sql.NullString)}
			return v, v.Val
		},
		data: []testdata{
			{input: nil, output: map[string]sql.NullString(nil)},
			{input: map[string]sql.NullString{}, output: map[string]sql.NullString{}},
			{
				input:  map[string]sql.NullString{"a": {"1", true}, "b": {"2", true}, "c": {"", false}},
				output: map[string]sql.NullString{"a": {"1", true}, "b": {"2", true}, "c": {"", false}}},
			{
				input:  map[string]sql.NullString{"a": {`'`, true}},
				output: map[string]sql.NullString{"a": {`'`, true}}},
			{
				input:  map[string]sql.NullString{"a": {`"`, true}, "b": {`\`, true}},
				output: map[string]sql.NullString{"a": {`"`, true}, "b": {`\`, true}}},
			{
				input:  map[string]sql.NullString{"a": {`\"`, true}, "b": {`\\`, true}},
				output: map[string]sql.NullString{"a": {`\"`, true}, "b": {`\\`, true}}},
			{
				input: map[string]sql.NullString{"a": {"\a", true}, "b": {"\b", true},
					"c": {"\t", true}, "d": {"\n", true}},
				output: map[string]sql.NullString{"a": {"\a", true}, "b": {"\b", true},
					"c": {"\t", true}, "d": {"\n", true}}},
			{
				input:  map[string]sql.NullString{"text": {`foo' "bar", baz \quux`, true}},
				output: map[string]sql.NullString{"text": {`foo' "bar", baz \quux`, true}}},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			v := new(string)
			return v, v
		},
		data: []testdata{
			{input: string(""), output: string("")},
			{
				input:  string(`"a"=>"1"`),
				output: string(`"a"=>"1"`)},
			{
				input:  string(`"a"=>"1", "b"=>"2"`),
				output: string(`"a"=>"1", "b"=>"2"`)},
			{
				input:  string(`"text"=>"foo' \"bar\", baz \\quux"`),
				output: string(`"text"=>"foo' \"bar\", baz \\quux"`)},
			{
				input:  string(`"a"=>"1", "b"=>NULL,"c"=>NULL`),
				output: string(`"a"=>"1", "b"=>NULL, "c"=>NULL`)},
			{
				input:  string(`"a"=>"1", "b"=>"\\\"","c"=>"\\\\\\\\"`),
				output: string(`"a"=>"1", "b"=>"\\\"", "c"=>"\\\\\\\\"`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			s := new([]byte)
			return s, s
		},
		data: []testdata{
			{input: []byte(""), output: []byte("")},
			{
				input:  []byte(`"a"=>"1"`),
				output: []byte(`"a"=>"1"`)},
			{
				input:  []byte(`"a"=>"1", "b"=>"2","c"=>"3"`),
				output: []byte(`"a"=>"1", "b"=>"2", "c"=>"3"`)},
			{
				input:  []byte(`"a"=>"1", "b"=>NULL,"c"=>NULL`),
				output: []byte(`"a"=>"1", "b"=>NULL, "c"=>NULL`)},
			{
				input:  []byte(`"a"=>"1", "b"=>"\\\"","c"=>"\\\\\\\\"`),
				output: []byte(`"a"=>"1", "b"=>"\\\"", "c"=>"\\\\\\\\"`)},
			{
				input:  []byte(`"text"=>"foo' \"bar\", baz \\quux"`),
				output: []byte(`"text"=>"foo' \"bar\", baz \\quux"`)},
		},
	}}.execute(t, "hstore")
}

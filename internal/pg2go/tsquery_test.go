package pg2go

import (
	"testing"
)

func TestTSQuery(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			v := new(string)
			return v, v
		},
		data: []testdata{
			{input: string(`'fat'`), output: string(`'fat'`)},
			{input: string(`'fat' & 'rat'`), output: string(`'fat' & 'rat'`)},
			{
				input:  string(`'fat':AB & ( 'rat' | 'cat' ) & !'bat':*`),
				output: string(`'fat':AB & ( 'rat' | 'cat' ) & !'bat':*`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			v := new([]byte)
			return v, v
		},
		data: []testdata{
			{input: []byte(`'fat'`), output: []byte(`'fat'`)},
			{input: []byte(`'fat' & 'rat'`), output: []byte(`'fat' & 'rat'`)},
			{
				input:  []byte(`'fat':AB & ( 'rat' | 'cat' ) & !'bat':*`),
				output: []byte(`'fat':AB & ( 'rat' | 'cat' ) & !'bat':*`)},
		},
	}}.execute(t, "tsquery")
}

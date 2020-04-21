package pg2go

import (
	"testing"
)

func TestMoney(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(MoneyFromInt64)
		},
		scanner: func() (interface{}, interface{}) {
			v := MoneyToInt64{Val: new(int64)}
			return v, v.Val
		},
		data: []testdata{
			{input: int64(0), output: int64(0)},
			{input: int64(99), output: int64(99)},
			{input: int64(120), output: int64(120)},
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
			{input: string(`$0.00`), output: string(`$0.00`)},
			{input: string(`$0.99`), output: string(`$0.99`)},
			{input: string(`$1.20`), output: string(`$1.20`)},
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
			{input: []byte(`$0.00`), output: []byte(`$0.00`)},
			{input: []byte(`$0.99`), output: []byte(`$0.99`)},
			{input: []byte(`$1.20`), output: []byte(`$1.20`)},
		},
	}}.execute(t, "money")
}

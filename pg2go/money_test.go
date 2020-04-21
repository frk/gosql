package pg2go

import (
	"testing"
)

func TestMoney(t *testing.T) {
	testlist2{{
		valuer:  MoneyFromInt64,
		scanner: MoneyToInt64,
		data: []testdata{
			{input: int64(0), output: int64(0)},
			{input: int64(99), output: int64(99)},
			{input: int64(120), output: int64(120)},
		},
	}, {
		data: []testdata{
			{input: string(`$0.00`), output: string(`$0.00`)},
			{input: string(`$0.99`), output: string(`$0.99`)},
			{input: string(`$1.20`), output: string(`$1.20`)},
		},
	}, {
		data: []testdata{
			{input: []byte(`$0.00`), output: []byte(`$0.00`)},
			{input: []byte(`$0.99`), output: []byte(`$0.99`)},
			{input: []byte(`$1.20`), output: []byte(`$1.20`)},
		},
	}}.execute(t, "money")
}

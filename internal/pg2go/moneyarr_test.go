package pg2go

import (
	"testing"
)

func TestMoneyArray(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(MoneyArrayFromInt64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := MoneyArrayToInt64Slice{Val: new([]int64)}
			return v, v.Val
		},
		data: []testdata{
			{input: []int64(nil), output: []int64(nil)},
			{input: []int64{}, output: []int64{}},
			{input: []int64{0}, output: []int64{0}},
			{input: []int64{0, 99}, output: []int64{0, 99}},
			{input: []int64{120, 0, 99}, output: []int64{120, 0, 99}},
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
			{input: string(`{}`), output: string(`{}`)},
			{input: string(`{$0.00}`), output: string(`{$0.00}`)},
			{
				input:  string(`{$1.20,$0.99}`),
				output: string(`{$1.20,$0.99}`)},
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
			{input: []byte(`{}`), output: []byte(`{}`)},
			{input: []byte(`{$0.00}`), output: []byte(`{$0.00}`)},
			{
				input:  []byte(`{$1.20,$0.99}`),
				output: []byte(`{$1.20,$0.99}`)},
		},
	}}.execute(t, "moneyarr")
}

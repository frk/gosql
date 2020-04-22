package pgsql

import (
	"testing"
)

func TestMoneyArray(t *testing.T) {
	testlist2{{
		valuer:  MoneyArrayFromInt64Slice,
		scanner: MoneyArrayToInt64Slice,
		data: []testdata{
			{input: []int64(nil), output: []int64(nil)},
			{input: []int64{}, output: []int64{}},
			{input: []int64{0}, output: []int64{0}},
			{input: []int64{0, 99}, output: []int64{0, 99}},
			{input: []int64{120, 0, 99}, output: []int64{120, 0, 99}},
		},
	}, {
		data: []testdata{
			{input: string(`{}`), output: string(`{}`)},
			{input: string(`{$0.00}`), output: string(`{$0.00}`)},
			{
				input:  string(`{$1.20,$0.99}`),
				output: string(`{$1.20,$0.99}`)},
		},
	}, {
		data: []testdata{
			{input: []byte(`{}`), output: []byte(`{}`)},
			{input: []byte(`{$0.00}`), output: []byte(`{$0.00}`)},
			{
				input:  []byte(`{$1.20,$0.99}`),
				output: []byte(`{$1.20,$0.99}`)},
		},
	}}.execute(t, "moneyarr")
}

package pgsql

import (
	"testing"
)

func TestVarBitArray(t *testing.T) {
	testlist2{{
		valuer:  VarBitArrayFromBoolSliceSlice,
		scanner: VarBitArrayToBoolSliceSlice,
		data: []testdata{
			{input: [][]bool(nil), output: [][]bool(nil)},
			{input: [][]bool{}, output: [][]bool{}},
			{
				input:  [][]bool{{true}},
				output: [][]bool{{true}}},
			{
				input:  [][]bool{{false}},
				output: [][]bool{{false}}},
			{
				input:  [][]bool{{false, true}, {true, true}, {true, false, false}},
				output: [][]bool{{false, true}, {true, true}, {true, false, false}}},
			{
				input:  [][]bool{{false, true}, nil, {}},
				output: [][]bool{{false, true}, nil, {}}},
		},
	}, {
		valuer:  VarBitArrayFromUint8SliceSlice,
		scanner: VarBitArrayToUint8SliceSlice,
		data: []testdata{
			{input: [][]uint8(nil), output: [][]uint8(nil)},
			{input: [][]uint8{}, output: [][]uint8{}},
			{
				input:  [][]uint8{{1}},
				output: [][]uint8{{1}}},
			{
				input:  [][]uint8{{0}},
				output: [][]uint8{{0}}},
			{
				input:  [][]uint8{{0, 1}, {1, 1}, {1, 0, 0}},
				output: [][]uint8{{0, 1}, {1, 1}, {1, 0, 0}}},
			{
				input:  [][]uint8{{0, 1}, nil, {}},
				output: [][]uint8{{0, 1}, nil, {}}},
		},
	}, {
		valuer:  VarBitArrayFromStringSlice,
		scanner: VarBitArrayToStringSlice,
		data: []testdata{
			{input: []string(nil), output: []string(nil)},
			{input: []string{}, output: []string{}},
			{
				input:  []string{"1"},
				output: []string{"1"}},
			{
				input:  []string{"0"},
				output: []string{"0"}},
			{
				input:  []string{"101010", "01", "10", "1111100"},
				output: []string{"101010", "01", "10", "1111100"}},
			{
				input:  []string{"1011", "", "010", ""},
				output: []string{"1011", "", "010", ""}},
		},
	}, {
		valuer:  VarBitArrayFromInt64Slice,
		scanner: VarBitArrayToInt64Slice,
		data: []testdata{
			{input: []int64(nil), output: []int64(nil)},
			{input: []int64{}, output: []int64{}},
			{
				input:  []int64{1},
				output: []int64{1}},
			{
				input:  []int64{0},
				output: []int64{0}},
			{
				input:  []int64{42, 1, 2, 124},
				output: []int64{42, 1, 2, 124}},
		},
	}, {
		data: []testdata{
			{input: string(`{}`), output: string(`{}`)},
			{
				input:  string(`{1}`),
				output: string(`{1}`)},
			{
				input:  string(`{0}`),
				output: string(`{0}`)},
			{
				input:  string(`{101010,01,10,1111100}`),
				output: string(`{101010,01,10,1111100}`)},
			{
				input:  string(`{0,NULL,""}`),
				output: string(`{0,NULL,""}`)},
		},
	}, {
		data: []testdata{
			{input: []byte(`{}`), output: []byte(`{}`)},
			{
				input:  []byte(`{1}`),
				output: []byte(`{1}`)},
			{
				input:  []byte(`{0}`),
				output: []byte(`{0}`)},
			{
				input:  []byte(`{101010,01,10,1111100}`),
				output: []byte(`{101010,01,10,1111100}`)},
			{
				input:  []byte(`{0,NULL,""}`),
				output: []byte(`{0,NULL,""}`)},
		},
	}}.execute(t, "varbitarr")
}

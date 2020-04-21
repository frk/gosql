package pg2go

import (
	"testing"
)

func TestVarBitArray(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(VarBitArrayFromBoolSliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := VarBitArrayToBoolSliceSlice{Val: new([][]bool)}
			return s, s.Val
		},
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
		valuer: func() interface{} {
			return new(VarBitArrayFromUint8SliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := VarBitArrayToUint8SliceSlice{Val: new([][]uint8)}
			return s, s.Val
		},
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
		valuer: func() interface{} {
			return new(VarBitArrayFromStringSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := VarBitArrayToStringSlice{Val: new([]string)}
			return s, s.Val
		},
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
		valuer: func() interface{} {
			return new(VarBitArrayFromInt64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := VarBitArrayToInt64Slice{Val: new([]int64)}
			return s, s.Val
		},
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
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			s := new(string)
			return s, s
		},
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
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			s := new([]byte)
			return s, s
		},
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

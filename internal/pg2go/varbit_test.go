package pg2go

import (
	"testing"
)

func TestVarBit(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(VarBitFromInt64)
		},
		scanner: func() (interface{}, interface{}) {
			s := VarBitToInt64{Val: new(int64)}
			return s, s.Val
		},
		data: []testdata{
			{input: int64(0), output: int64(0)},
			{input: int64(1), output: int64(1)},
			{input: int64(1024), output: int64(1024)},
			// {input: int64(-512), output: int64(-512)},
		},
	}, {
		valuer: func() interface{} {
			return new(VarBitFromBoolSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := VarBitToBoolSlice{Val: new([]bool)}
			return s, s.Val
		},
		data: []testdata{
			{input: []bool(nil), output: []bool(nil)},
			{input: []bool{}, output: []bool{}},
			{input: []bool{false}, output: []bool{false}},
			{input: []bool{true}, output: []bool{true}},
			{
				input:  []bool{true, false, false, true, false},
				output: []bool{true, false, false, true, false}},
		},
	}, {
		valuer: func() interface{} {
			return new(VarBitFromUint8Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := VarBitToUint8Slice{Val: new([]uint8)}
			return s, s.Val
		},
		data: []testdata{
			{input: []uint8(nil), output: []uint8(nil)},
			{input: []uint8{}, output: []uint8{}},
			{input: []uint8{0}, output: []uint8{0}},
			{input: []uint8{1}, output: []uint8{1}},
			{
				input:  []uint8{1, 0, 0, 1, 0},
				output: []uint8{1, 0, 0, 1, 0}},
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
			{input: string(""), output: string("")},
			{input: string("0"), output: string("0")},
			{input: string("1"), output: string("1")},
			{
				input:  string("10010"),
				output: string("10010")},
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
			{input: []byte("0"), output: []byte("0")},
			{input: []byte("1"), output: []byte("1")},
			{
				input:  []byte("10010"),
				output: []byte("10010")},
		},
	}}.execute(t, "varbit")
}

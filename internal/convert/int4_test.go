package convert

import (
	"testing"
)

func TestInt4(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return nil // int
		},
		scanner: func() (interface{}, interface{}) {
			s := new(int)
			return s, s
		},
		data: []testdata{
			{input: int(-2147483648), output: int(-2147483648)},
			{input: int(2147483647), output: int(2147483647)},
		},
	}, {
		valuer: func() interface{} {
			return nil // int8
		},
		scanner: func() (interface{}, interface{}) {
			s := new(int8)
			return s, s
		},
		data: []testdata{
			{input: int8(-128), output: int8(-128)},
			{input: int8(127), output: int8(127)},
		},
	}, {
		valuer: func() interface{} {
			return nil // int16
		},
		scanner: func() (interface{}, interface{}) {
			s := new(int16)
			return s, s
		},
		data: []testdata{
			{input: int16(-32768), output: int16(-32768)},
			{input: int16(32767), output: int16(32767)},
		},
	}, {
		valuer: func() interface{} {
			return nil // int32
		},
		scanner: func() (interface{}, interface{}) {
			s := new(int32)
			return s, s
		},
		data: []testdata{
			{input: int32(-2147483648), output: int32(-2147483648)},
			{input: int32(2147483647), output: int32(2147483647)},
		},
	}, {
		valuer: func() interface{} {
			return nil // int64
		},
		scanner: func() (interface{}, interface{}) {
			s := new(int64)
			return s, s
		},
		data: []testdata{
			{input: int64(-2147483648), output: int64(-2147483648)},
			{input: int64(2147483647), output: int64(2147483647)},
		},
	}, {
		valuer: func() interface{} {
			return nil // uint
		},
		scanner: func() (interface{}, interface{}) {
			s := new(uint)
			return s, s
		},
		data: []testdata{
			{input: uint(0), output: uint(0)},
			{input: uint(2147483647), output: uint(2147483647)},
		},
	}, {
		valuer: func() interface{} {
			return nil // uint8
		},
		scanner: func() (interface{}, interface{}) {
			s := new(uint8)
			return s, s
		},
		data: []testdata{
			{input: uint8(0), output: uint8(0)},
			{input: uint8(255), output: uint8(255)},
		},
	}, {
		valuer: func() interface{} {
			return nil // uint16
		},
		scanner: func() (interface{}, interface{}) {
			s := new(uint16)
			return s, s
		},
		data: []testdata{
			{input: uint16(0), output: uint16(0)},
			{input: uint16(65535), output: uint16(65535)},
		},
	}, {
		valuer: func() interface{} {
			return nil // uint32
		},
		scanner: func() (interface{}, interface{}) {
			s := new(uint32)
			return s, s
		},
		data: []testdata{
			{input: uint32(0), output: uint32(0)},
			{input: uint32(2147483647), output: uint32(2147483647)},
		},
	}, {
		valuer: func() interface{} {
			return nil // uint64
		},
		scanner: func() (interface{}, interface{}) {
			s := new(uint64)
			return s, s
		},
		data: []testdata{
			{input: uint64(0), output: uint64(0)},
			{input: uint64(2147483647), output: uint64(2147483647)},
		},
	}, {
		valuer: func() interface{} {
			return nil // float32
		},
		scanner: func() (interface{}, interface{}) {
			s := new(float32)
			return s, s
		},
		data: []testdata{
			{input: float32(-2147483648.0), output: float32(-2147483648.0)},
			// XXX float32(2147483647.0) gets turned into 2147483648 which is outside int4 range
			{input: float32(214748364.0), output: float32(214748364.0)},
		},
	}, {
		valuer: func() interface{} {
			return nil // float64
		},
		scanner: func() (interface{}, interface{}) {
			s := new(float64)
			return s, s
		},
		data: []testdata{
			{input: float64(-2147483648.0), output: float64(-2147483648.0)},
			{input: float64(2147483647.0), output: float64(2147483647.0)},
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
			{input: string("-2147483648"), output: string(`-2147483648`)},
			{input: string("2147483647"), output: string(`2147483647`)},
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
			{input: []byte("-2147483648"), output: []byte(`-2147483648`)},
			{input: []byte("2147483647"), output: []byte(`2147483647`)},
		},
	}}.execute(t, "int4")
}

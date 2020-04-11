package convert

import (
	"testing"
)

func TestInt2(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return nil // int
		},
		scanner: func() (interface{}, interface{}) {
			s := new(int)
			return s, s
		},
		data: []testdata{
			{input: int(-32768), output: iptr(-32768)},
			{input: int(32767), output: iptr(32767)},
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
			{input: int8(-128), output: i8ptr(-128)},
			{input: int8(127), output: i8ptr(127)},
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
			{input: int16(-32768), output: i16ptr(-32768)},
			{input: int16(32767), output: i16ptr(32767)},
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
			{input: int32(-32768), output: i32ptr(-32768)},
			{input: int32(32767), output: i32ptr(32767)},
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
			{input: int64(-32768), output: i64ptr(-32768)},
			{input: int64(32767), output: i64ptr(32767)},
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
			{input: uint(0), output: uptr(0)},
			{input: uint(32767), output: uptr(32767)},
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
			{input: uint8(0), output: u8ptr(0)},
			{input: uint8(255), output: u8ptr(255)},
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
			{input: uint16(0), output: u16ptr(0)},
			{input: uint16(32767), output: u16ptr(32767)},
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
			{input: uint32(0), output: u32ptr(0)},
			{input: uint32(32767), output: u32ptr(32767)},
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
			{input: uint64(0), output: u64ptr(0)},
			{input: uint64(32767), output: u64ptr(32767)},
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
			{input: float32(-32768.0), output: f32ptr(-32768.0)},
			{input: float32(32767.0), output: f32ptr(32767.0)},
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
			{input: float64(-32768.0), output: f64ptr(-32768.0)},
			{input: float64(32767.0), output: f64ptr(32767.0)},
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
			{input: string("-32768"), output: strptr(`-32768`)},
			{input: string("32767"), output: strptr(`32767`)},
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
			{input: []byte("-32768"), output: bytesptr(`-32768`)},
			{input: []byte("32767"), output: bytesptr(`32767`)},
		},
	}}.execute(t, "int2")
}

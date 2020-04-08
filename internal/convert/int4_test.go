package convert

import (
	"testing"
)

func TestInt4(t *testing.T) {
	test_valuer_with_scanner{{
		valuer: func() interface{} {
			return nil // int
		},
		scanner: func() (interface{}, interface{}) {
			s := new(int)
			return s, s
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "int4", input: int(-2147483648), output: iptr(-2147483648)},
			{typ: "int4", input: int(2147483647), output: iptr(2147483647)},
		},
	}, {
		valuer: func() interface{} {
			return nil // int8
		},
		scanner: func() (interface{}, interface{}) {
			s := new(int8)
			return s, s
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "int4", input: int8(-128), output: i8ptr(-128)},
			{typ: "int4", input: int8(127), output: i8ptr(127)},
		},
	}, {
		valuer: func() interface{} {
			return nil // int16
		},
		scanner: func() (interface{}, interface{}) {
			s := new(int16)
			return s, s
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "int4", input: int16(-32768), output: i16ptr(-32768)},
			{typ: "int4", input: int16(32767), output: i16ptr(32767)},
		},
	}, {
		valuer: func() interface{} {
			return nil // int32
		},
		scanner: func() (interface{}, interface{}) {
			s := new(int32)
			return s, s
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "int4", input: int32(-2147483648), output: i32ptr(-2147483648)},
			{typ: "int4", input: int32(2147483647), output: i32ptr(2147483647)},
		},
	}, {
		valuer: func() interface{} {
			return nil // int64
		},
		scanner: func() (interface{}, interface{}) {
			s := new(int64)
			return s, s
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "int4", input: int64(-2147483648), output: i64ptr(-2147483648)},
			{typ: "int4", input: int64(2147483647), output: i64ptr(2147483647)},
		},
	}, {
		valuer: func() interface{} {
			return nil // uint
		},
		scanner: func() (interface{}, interface{}) {
			s := new(uint)
			return s, s
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "int4", input: uint(0), output: uptr(0)},
			{typ: "int4", input: uint(2147483647), output: uptr(2147483647)},
		},
	}, {
		valuer: func() interface{} {
			return nil // uint8
		},
		scanner: func() (interface{}, interface{}) {
			s := new(uint8)
			return s, s
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "int4", input: uint8(0), output: u8ptr(0)},
			{typ: "int4", input: uint8(255), output: u8ptr(255)},
		},
	}, {
		valuer: func() interface{} {
			return nil // uint16
		},
		scanner: func() (interface{}, interface{}) {
			s := new(uint16)
			return s, s
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "int4", input: uint16(0), output: u16ptr(0)},
			{typ: "int4", input: uint16(65535), output: u16ptr(65535)},
		},
	}, {
		valuer: func() interface{} {
			return nil // uint32
		},
		scanner: func() (interface{}, interface{}) {
			s := new(uint32)
			return s, s
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "int4", input: uint32(0), output: u32ptr(0)},
			{typ: "int4", input: uint32(2147483647), output: u32ptr(2147483647)},
		},
	}, {
		valuer: func() interface{} {
			return nil // uint64
		},
		scanner: func() (interface{}, interface{}) {
			s := new(uint64)
			return s, s
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "int4", input: uint64(0), output: u64ptr(0)},
			{typ: "int4", input: uint64(2147483647), output: u64ptr(2147483647)},
		},
	}, {
		valuer: func() interface{} {
			return nil // float32
		},
		scanner: func() (interface{}, interface{}) {
			s := new(float32)
			return s, s
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "int4", input: float32(-2147483648.0), output: f32ptr(-2147483648.0)},
			// XXX float32(2147483647.0) gets turned into 2147483648 which is outside int4 range
			{typ: "int4", input: float32(214748364.0), output: f32ptr(214748364.0)},
		},
	}, {
		valuer: func() interface{} {
			return nil // float64
		},
		scanner: func() (interface{}, interface{}) {
			s := new(float64)
			return s, s
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "int4", input: float64(-2147483648.0), output: f64ptr(-2147483648.0)},
			{typ: "int4", input: float64(2147483647.0), output: f64ptr(2147483647.0)},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			s := new(string)
			return s, s
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "int4", input: "-2147483648", output: strptr(`-2147483648`)},
			{typ: "int4", input: "2147483647", output: strptr(`2147483647`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			s := new([]byte)
			return s, s
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "int4", input: []byte("-2147483648"), output: bytesptr(`-2147483648`)},
			{typ: "int4", input: []byte("2147483647"), output: bytesptr(`2147483647`)},
		},
	}}.execute(t)
}

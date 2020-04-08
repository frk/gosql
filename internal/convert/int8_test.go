package convert

import (
	"testing"
)

func TestInt8(t *testing.T) {
	test_valuer_with_scanner{{
		valuer: func() interface{} {
			return nil // int
		},
		scanner: func() (interface{}, interface{}) {
			s := new(int)
			return s, s
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "int8", input: int(-2147483648), output: iptr(-2147483648)},
			{typ: "int8", input: int(9223372036854775807), output: iptr(9223372036854775807)},
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
			{typ: "int8", input: int8(-128), output: i8ptr(-128)},
			{typ: "int8", input: int8(127), output: i8ptr(127)},
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
			{typ: "int8", input: int16(-32768), output: i16ptr(-32768)},
			{typ: "int8", input: int16(32767), output: i16ptr(32767)},
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
			{typ: "int8", input: int32(-2147483648), output: i32ptr(-2147483648)},
			{typ: "int8", input: int32(2147483647), output: i32ptr(2147483647)},
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
			{typ: "int8", input: int64(-9223372036854775808), output: i64ptr(-9223372036854775808)},
			{typ: "int8", input: int64(9223372036854775807), output: i64ptr(9223372036854775807)},
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
			{typ: "int8", input: uint(0), output: uptr(0)},
			{typ: "int8", input: uint(2147483647), output: uptr(2147483647)},
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
			{typ: "int8", input: uint8(0), output: u8ptr(0)},
			{typ: "int8", input: uint8(255), output: u8ptr(255)},
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
			{typ: "int8", input: uint16(0), output: u16ptr(0)},
			{typ: "int8", input: uint16(65535), output: u16ptr(65535)},
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
			{typ: "int8", input: uint32(0), output: u32ptr(0)},
			{typ: "int8", input: uint32(2147483647), output: u32ptr(2147483647)},
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
			{typ: "int8", input: uint64(0), output: u64ptr(0)},
			{typ: "int8", input: uint64(9223372036854775807), output: u64ptr(9223372036854775807)},
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
			{typ: "int8", input: float32(-2147483648.0), output: f32ptr(-2147483648.0)},
			{typ: "int8", input: float32(2147483647.0), output: f32ptr(2147483647.0)},
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
			// XXX float64(-9223372036854775808.0) is outside int8 range
			{typ: "int8", input: float64(-922337203685477580.0), output: f64ptr(-922337203685477580.0)},
			// XXX float64(9223372036854775807.0) is outside int8 range
			{typ: "int8", input: float64(922337203685477580.0), output: f64ptr(922337203685477580.0)},
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
			{typ: "int8", input: "-9223372036854775808", output: strptr(`-9223372036854775808`)},
			{typ: "int8", input: "9223372036854775807", output: strptr(`9223372036854775807`)},
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
			{typ: "int8", input: []byte("-9223372036854775808"), output: bytesptr(`-9223372036854775808`)},
			{typ: "int8", input: []byte("9223372036854775807"), output: bytesptr(`9223372036854775807`)},
		},
	}}.execute(t)
}

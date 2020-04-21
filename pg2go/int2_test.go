package pg2go

import (
	"testing"
)

func TestInt2(t *testing.T) {
	testlist2{{
		data: []testdata{
			{input: int(-32768), output: int(-32768)},
			{input: int(32767), output: int(32767)},
		},
	}, {
		data: []testdata{
			{input: int8(-128), output: int8(-128)},
			{input: int8(127), output: int8(127)},
		},
	}, {
		data: []testdata{
			{input: int16(-32768), output: int16(-32768)},
			{input: int16(32767), output: int16(32767)},
		},
	}, {
		data: []testdata{
			{input: int32(-32768), output: int32(-32768)},
			{input: int32(32767), output: int32(32767)},
		},
	}, {
		data: []testdata{
			{input: int64(-32768), output: int64(-32768)},
			{input: int64(32767), output: int64(32767)},
		},
	}, {
		data: []testdata{
			{input: uint(0), output: uint(0)},
			{input: uint(32767), output: uint(32767)},
		},
	}, {
		data: []testdata{
			{input: uint8(0), output: uint8(0)},
			{input: uint8(255), output: uint8(255)},
		},
	}, {
		data: []testdata{
			{input: uint16(0), output: uint16(0)},
			{input: uint16(32767), output: uint16(32767)},
		},
	}, {
		data: []testdata{
			{input: uint32(0), output: uint32(0)},
			{input: uint32(32767), output: uint32(32767)},
		},
	}, {
		data: []testdata{
			{input: uint64(0), output: uint64(0)},
			{input: uint64(32767), output: uint64(32767)},
		},
	}, {
		data: []testdata{
			{input: float32(-32768.0), output: float32(-32768.0)},
			{input: float32(32767.0), output: float32(32767.0)},
		},
	}, {
		data: []testdata{
			{input: float64(-32768.0), output: float64(-32768.0)},
			{input: float64(32767.0), output: float64(32767.0)},
		},
	}, {
		data: []testdata{
			{input: string("-32768"), output: string(`-32768`)},
			{input: string("32767"), output: string(`32767`)},
		},
	}, {
		data: []testdata{
			{input: []byte("-32768"), output: []byte(`-32768`)},
			{input: []byte("32767"), output: []byte(`32767`)},
		},
	}}.execute(t, "int2")
}

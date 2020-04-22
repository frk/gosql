package pgsql

import (
	"testing"
)

func TestNumeric(t *testing.T) {
	testlist2{{
		data: []testdata{
			{
				input:  int(-9223372036854775808),
				output: int(-9223372036854775808)},
			{
				input:  int(9223372036854775807),
				output: int(9223372036854775807)},
		},
	}, {
		data: []testdata{
			{
				input:  int8(-128),
				output: int8(-128)},
			{
				input:  int8(127),
				output: int8(127)},
		},
	}, {
		data: []testdata{
			{
				input:  int16(-32768),
				output: int16(-32768)},
			{
				input:  int16(32767),
				output: int16(32767)},
		},
	}, {
		data: []testdata{
			{
				input:  int32(-2147483648),
				output: int32(-2147483648)},
			{
				input:  int32(2147483647),
				output: int32(2147483647)},
		},
	}, {
		data: []testdata{
			{
				input:  int64(-9223372036854775808),
				output: int64(-9223372036854775808)},
			{
				input:  int64(9223372036854775807),
				output: int64(9223372036854775807)},
		},
	}, {
		data: []testdata{
			{
				input:  uint(0),
				output: uint(0)},
			{
				input:  uint(9223372036854775807),
				output: uint(9223372036854775807)},
		},
	}, {
		data: []testdata{
			{
				input:  uint8(0),
				output: uint8(0)},
			{
				input:  uint8(255),
				output: uint8(255)},
		},
	}, {
		data: []testdata{
			{
				input:  uint16(0),
				output: uint16(0)},
			{
				input:  uint16(65535),
				output: uint16(65535)},
		},
	}, {
		data: []testdata{
			{
				input:  uint32(0),
				output: uint32(0)},
			{
				input:  uint32(4294967295),
				output: uint32(4294967295)},
		},
	}, {
		data: []testdata{
			{
				input:  uint64(0),
				output: uint64(0)},
			{
				input:  uint64(9223372036854775807),
				output: uint64(9223372036854775807)},
		},
	}, {
		data: []testdata{
			{
				input:  float32(-21474.83648),
				output: float32(-21474.83648)},
			{
				input:  float32(21474.83647),
				output: float32(21474.83647)},
		},
	}, {
		data: []testdata{
			{
				input:  float64(-92233720.3685477580),
				output: float64(-92233720.3685477580)},
			{
				input:  float64(922337203.685477580),
				output: float64(922337203.685477580)},
		},
	}, {
		data: []testdata{
			{
				input:  string(`-92233720.36854775808`),
				output: string(`-92233720.36854775808`)},
			{
				input:  string(`922337203.6854775807`),
				output: string(`922337203.6854775807`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte(`-92233720.36854775808`),
				output: []byte(`-92233720.36854775808`)},
			{
				input:  []byte(`922337203.6854775807`),
				output: []byte(`922337203.6854775807`)},
		},
	}}.execute(t, "numeric")
}

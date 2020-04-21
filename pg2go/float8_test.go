package pg2go

import (
	"testing"
)

func TestFloat8(t *testing.T) {
	testlist2{{
		data: []testdata{
			{input: float32(0), output: float32(0)},
			{input: float32(1), output: float32(1)},
			{input: float32(3.14), output: float32(3.14)},
			{input: float32(0.15), output: float32(0.15)},
		},
	}, {
		data: []testdata{
			{input: float64(0), output: float64(0)},
			{input: float64(1), output: float64(1)},
			{input: float64(3.14), output: float64(3.14)},
			{input: float64(0.15), output: float64(0.15)},
		},
	}, {
		data: []testdata{
			{input: string("0"), output: string(`0`)},
			{input: string("1"), output: string(`1`)},
			{input: string("3.14"), output: string(`3.14`)},
			{input: string("0.15"), output: string(`0.15`)},
		},
	}, {
		data: []testdata{
			{input: []byte("0"), output: []byte(`0`)},
			{input: []byte("1"), output: []byte(`1`)},
			{input: []byte("3.14"), output: []byte(`3.14`)},
			{input: []byte("0.15"), output: []byte(`0.15`)},
		},
	}}.execute(t, "float8")
}

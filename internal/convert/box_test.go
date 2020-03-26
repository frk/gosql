package convert

import (
	"testing"
)

func TestBox_Valuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return new(BoxFromFloat64Array2Array2)
		},
		rows: []test_valuer_row{
			{typ: "box", in: nil, want: strptr(`(0,0),(0,0)`)},
			{typ: "box", in: [2][2]float64{{1, 1}, {0, 0}}, want: strptr(`(1,1),(0,0)`)},
			{typ: "box", in: [2][2]float64{{0, 0}, {1, 1}}, want: strptr(`(1,1),(0,0)`)},
			{typ: "box",
				in: [2][2]float64{{4.5203, 0.79322}, {3.2, 5.63333}},
				// TODO(mkopriva) figure out whether postgres can be made to return
				// a string exactly matching the input Go input ...
				want: strptr(`(4.5202999999999998,5.6333299999999999),(3.2000000000000002,0.79322000000000004)`),
			},
		},
	}}.execute(t)
}

func TestBox_Scanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			s := BoxToFloat64Array2Array2{A: new([2][2]float64)}
			return s, s.A
		},
		rows: []test_scanner_row{
			{typ: "box", in: nil, want: new([2][2]float64)},
			{typ: "box", in: `(1,1),(0,0)`, want: &[2][2]float64{{1, 1}, {0, 0}}},
			{typ: "box", in: `(0,0),(1,1)`, want: &[2][2]float64{{1, 1}, {0, 0}}},
			{typ: "box",
				in:   `(4.5203,0.79322),(3.2,5.63333)`,
				want: &[2][2]float64{{4.5203, 5.63333}, {3.2, 0.79322}},
			},
		},
	}}.execute(t)
}

func TestBox_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil // string
		},
		rows: []test_valuer_row{
			{typ: "box", in: nil, want: nil},
			{typ: "box", in: "(0,0),(0,0)", want: strptr(`(0,0),(0,0)`)},
			{typ: "box", in: "(1,1),(0,0)", want: strptr(`(1,1),(0,0)`)},
			{typ: "box", in: "(0,0),(1,1)", want: strptr(`(1,1),(0,0)`)},
			{
				typ: "box",
				in:  "(4.5203,0.79322),(3.2,5.63333)",
				want: strptr(`(4.5202999999999998,5.6333299999999999),` +
					`(3.2000000000000002,0.79322000000000004)`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		rows: []test_valuer_row{
			{typ: "box", in: nil, want: nil},
			{typ: "box", in: []byte("(0,0),(0,0)"), want: strptr(`(0,0),(0,0)`)},
			{typ: "box", in: []byte("(1,1),(0,0)"), want: strptr(`(1,1),(0,0)`)},
			{typ: "box", in: []byte("(0,0),(1,1)"), want: strptr(`(1,1),(0,0)`)},
			{
				typ: "box",
				in:  []byte("(4.5203,0.79322),(3.2,5.63333)"),
				want: strptr(`(4.5202999999999998,5.6333299999999999),` +
					`(3.2000000000000002,0.79322000000000004)`)},
		},
	}}.execute(t)
}

func TestBox_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "box", in: `(0,0),(0,0)`, want: strptr("(0,0),(0,0)")},
			{typ: "box", in: `(1,1),(0,0)`, want: strptr("(1,1),(0,0)")},
			{typ: "box", in: `(0,0),(1,1)`, want: strptr("(1,1),(0,0)")},
			{
				typ: "box",
				in:  `(4.5203,0.79322),(3.2,5.63333)`,
				want: strptr(`(4.5202999999999998,5.6333299999999999),` +
					`(3.2000000000000002,0.79322000000000004)`)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "box", in: `(0,0),(0,0)`, want: bytesptr("(0,0),(0,0)")},
			{typ: "box", in: `(1,1),(0,0)`, want: bytesptr("(1,1),(0,0)")},
			{typ: "box", in: `(0,0),(1,1)`, want: bytesptr("(1,1),(0,0)")},
			{
				typ: "box",
				in:  `(4.5203,0.79322),(3.2,5.63333)`,
				want: bytesptr(`(4.5202999999999998,5.6333299999999999),` +
					`(3.2000000000000002,0.79322000000000004)`)},
		},
	}}.execute(t)
}

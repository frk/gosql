package convert

import (
	"testing"
)

func TestBox2Float64a2a2(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			s := Box2Float64a2a2{Ptr: new([2][2]float64)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "box", in: nil, want: new([2][2]float64)},
			{typ: "box", in: `(1,1),(0,0)`, want: &[2][2]float64{{1, 1}, {0, 0}}},
			{typ: "box", in: `(0,0),(1,1)`, want: &[2][2]float64{{1, 1}, {0, 0}}},
			{typ: "box", in: `(4.5203,0.79322),(3.2,5.63333)`, want: &[2][2]float64{{4.5203, 5.63333}, {3.2, 0.79322}}},
		},
	}}.execute(t)
}

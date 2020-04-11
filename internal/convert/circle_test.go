package convert

import (
	"testing"
)

func TestCircle(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		data: []testdata{
			{
				input:  string("<(0,0),3.5>"),
				output: strptr("<(0,0),3.5>")},
			{
				input:  string("<(0.5,1),5>"),
				output: strptr("<(0.5,1),5>")},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		data: []testdata{
			{
				input:  nil,
				output: new([]byte)},
			{
				input:  []byte("<(0,0),3.5>"),
				output: bytesptr("<(0,0),3.5>")},
			{
				input:  []byte("<(0.5,1),5>"),
				output: bytesptr("<(0.5,1),5>")},
		},
	}}.execute(t, "circle")
}

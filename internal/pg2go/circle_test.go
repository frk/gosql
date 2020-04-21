package pg2go

import (
	"testing"
)

func TestCircle(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			s := new(string)
			return s, s
		},
		data: []testdata{
			{
				input:  string("<(0,0),3.5>"),
				output: string("<(0,0),3.5>")},
			{
				input:  string("<(0.5,1),5>"),
				output: string("<(0.5,1),5>")},
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
			{
				input:  []byte("<(0,0),3.5>"),
				output: []byte("<(0,0),3.5>")},
			{
				input:  []byte("<(0.5,1),5>"),
				output: []byte("<(0.5,1),5>")},
		},
	}}.execute(t, "circle")
}

package pg2go

import (
	"testing"
)

func TestCircle(t *testing.T) {
	testlist2{{
		data: []testdata{
			{
				input:  string("<(0,0),3.5>"),
				output: string("<(0,0),3.5>")},
			{
				input:  string("<(0.5,1),5>"),
				output: string("<(0.5,1),5>")},
		},
	}, {
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

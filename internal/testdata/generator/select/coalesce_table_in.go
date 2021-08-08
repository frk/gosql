package testdata

import (
	"time"
)

type SelectCoalesceTableQuery struct {
	Columns struct {
		A  string    `sql:"col_a"`
		B  string    `sql:"col_b"`
		C  int       `sql:"col_c"`
		D  time.Time `sql:"col_d"`
		E  string    `sql:"col_e"`
		A2 string    `sql:"col_a,coalesce('foo')"`
		B2 string    `sql:"col_b,coalesce"`
		E2 string    `sql:"col_e,coalesce('red')"`
	} `rel:"test_coalesce:c"`
}

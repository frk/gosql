// WARNING(mkopriva): IF NOT NECESSARY TRY NOT TO DO ANYTHING THAT WILL REORDER
// LINES OR ADD NEW LINES IN THE MIDDLE OF THE FILE, ONLY ADD NEW CODE AT THE TAIL END.
package testdata

import "time"

//OK: simple select
type SelectPostgresTestOK_Simple struct {
	Columns struct {
		A int       `sql:"col_a"`
		B string    `sql:"col_b"`
		C bool      `sql:"col_c"`
		D float64   `sql:"col_d"`
		E time.Time `sql:"col_e"`
	} `rel:"column_tests_1"`
}

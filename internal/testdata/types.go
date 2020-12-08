// WARNING(mkopriva): IF NOT NECESSARY TRY NOT TO DO ANYTHING THAT WILL REORDER
// LINES OR ADD NEW LINES IN THE MIDDLE OF THE FILE, ONLY ADD NEW CODE AT THE TAIL END.
package testdata

import (
	"time"

	"github.com/frk/gosql/internal/testdata/common"
)

// column_tests_1 record
type CT1 struct {
	A int       `sql:"col_a"`
	B string    `sql:"col_b"`
	C bool      `sql:"col_c"`
	D float64   `sql:"col_d"`
	E time.Time `sql:"col_e"`
}

type CT1_part struct {
	A int    `sql:"col_a"`
	B string `sql:"col_b"`
	C bool   `sql:"col_c"`

	// omitted fields
	// D float64   `sql:"col_d"`
	// E time.Time `sql:"col_e"`
}

type CT1_bad struct {
	XYZ string `sql:"col_xyz"` // not in table
}

// column_tests_2 record (partial)
type CT2 struct {
	Foo int    `sql:"col_foo"`
	Bar string `sql:"col_bar"`
	Baz bool   `sql:"col_baz"`
}

type T2 struct {
	Foo int    `sql:"foo"`
	Bar string `sql:"bar"`
	Baz bool   `sql:"baz"`
}

type COLOR_ENUM string

type CT3 struct {
	ColorText string        `sql:"color_text"`
	ColorEnum COLOR_ENUM    `sql:"color_enum"`
	SomeTime  common.MyTime `sql:"some_time"`
}

type CT3b struct {
	ColorText string         `sql:"color_text"`
	ColorEnum COLOR_ENUM     `sql:"color_enum"`
	SomeTime  *common.MyTime `sql:"some_time"`
}

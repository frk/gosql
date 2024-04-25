// WARNING(mkopriva): IF NOT NECESSARY TRY NOT TO DO ANYTHING THAT WILL REORDER
// LINES OR ADD NEW LINES IN THE MIDDLE OF THE FILE, ONLY ADD NEW CODE AT THE TAIL END.
package testdata

import (
	"database/sql"
	"database/sql/driver"

	"github.com/frk/gosql/internal/testdata/common"
)

type CT4 struct {
	ColorText string        `sql:"color_text"`
	ColorEnum COLOR_ENUM    `sql:"color_enum"`
	H1Styles  cssStyles     `sql:"h1_styles"`
	SomeTime  common.MyTime `sql:"some_time"`
}

// custom composite type test
type cssStyles struct {
	// ...
}

var _ driver.Valuer = cssStyles{}
var _ sql.Scanner = (*cssStyles)(nil)

func (s cssStyles) Value() (driver.Value, error) {
	return nil, nil
}

func (s *cssStyles) Scan(src any) error {
	return nil
}

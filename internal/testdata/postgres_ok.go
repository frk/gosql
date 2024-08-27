// WARNING(mkopriva): IF NOT NECESSARY TRY NOT TO DO ANYTHING THAT WILL REORDER
// LINES OR ADD NEW LINES IN THE MIDDLE OF THE FILE, ONLY ADD NEW CODE AT THE TAIL END.
package testdata

import (
	"time"

	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

// OK: simple select
type SelectPostgresTestOK_Simple struct {
	Columns struct {
		A int       `sql:"col_a"`
		B string    `sql:"col_b"`
		C bool      `sql:"col_c"`
		D float64   `sql:"col_d"`
		E time.Time `sql:"col_e"`
	} `rel:"column_tests_1"`
}

type SelectPostgresTestOK_Enums struct {
	Rel   CT3 `rel:"column_tests_3"`
	Where struct {
		SomeTime common.MyTime `sql:"some_time <"`
	}
}

type InsertPostgresTestOK_Enums struct {
	Rel []*CT3 `rel:"column_tests_3"`
}

type SelectPostgresTestOK_CustomTypePointer struct {
	Rel   CT3b `rel:"column_tests_3"`
	Where struct {
		SomeTime common.MyTime `sql:"some_time <"`
	}
}

type InsertPostgresTestOK_CustomTypePointer struct {
	Rel []*CT3b `rel:"column_tests_3"`
}

type SelectPostgresTestOK_CompositeType struct {
	Rel CT4 `rel:"column_tests_4"`
}

type SelectPostgresTestOK_CompositeTypePointer struct {
	Rel *CT4 `rel:"column_tests_4"`
}

type InsertPostgresTestOK_CompositeTypeSlice struct {
	Rel []*CT4 `rel:"column_tests_4"`
}

type SelectPostgresTestOK_WhereJoinedUnaryNullColumn struct {
	Columns CT1 `rel:"column_tests_1:a"`
	Join    struct {
		_ gosql.LeftJoin `sql:"column_tests_2:b"`
	}
	Where struct {
		_ gosql.Column `sql:"b.col_baz isnull"`
	}
}

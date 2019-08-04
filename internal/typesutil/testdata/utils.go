package testdata

import (
	"time"

	"github.com/frk/gosql"
)

type GetDirectiveNameTest struct {
	_ time.Time
	_ Column
	_ Relation
	_ struct{}
	_ string

	_ gosql.Column
	_ gosql.Relation
	_ gosql.RightJoin
}

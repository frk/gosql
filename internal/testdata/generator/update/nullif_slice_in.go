package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type UpdateNullIfSliceQuery struct {
	Data []*common.ConflictData `rel:"test_onconflict:k"`
}

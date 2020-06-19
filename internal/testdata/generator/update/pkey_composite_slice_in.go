package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type UpdatePKeyCompositeSliceQuery struct {
	Data []*common.ConflictData `rel:"test_composite_pkey:p"`
}

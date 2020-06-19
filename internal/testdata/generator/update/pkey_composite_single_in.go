package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type UpdatePKeyCompositeSingleQuery struct {
	Data *common.ConflictData `rel:"test_composite_pkey:p"`
}

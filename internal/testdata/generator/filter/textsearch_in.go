package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type FilterTextSearchRecords struct {
	User *common.User2    `rel:"test_user:u"`
	_    gosql.TextSearch `sql:"u._search_document"`
	common.FilterMaker
}

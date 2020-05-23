// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"time"

	"github.com/frk/gosql"
)

var _FilterBasicAliasRecords_colmap = map[string]string{
	"Id":        `u."id"`,
	"Email":     `u."email"`,
	"FullName":  `u."full_name"`,
	"CreatedAt": `u."created_at"`,
}

func (f *FilterBasicAliasRecords) TextSearch(v string) {
	// No search document specified.
}

func (f *FilterBasicAliasRecords) UnmarshalFQL(fqlString string) error {
	return f.Filter.UnmarshalFQL(fqlString, _FilterBasicAliasRecords_colmap, false)
}

func (f *FilterBasicAliasRecords) UnmarshalSort(sortString string) error {
	return f.Filter.UnmarshalSort(sortString, _FilterBasicAliasRecords_colmap, false)
}

func (f *FilterBasicAliasRecords) Id(op string, val int) *FilterBasicAliasRecords {
	f.Filter.Col(`u."id"`, op, val)
	return f
}

func (f *FilterBasicAliasRecords) Email(op string, val string) *FilterBasicAliasRecords {
	f.Filter.Col(`u."email"`, op, val)
	return f
}

func (f *FilterBasicAliasRecords) FullName(op string, val string) *FilterBasicAliasRecords {
	f.Filter.Col(`u."full_name"`, op, val)
	return f
}

func (f *FilterBasicAliasRecords) CreatedAt(op string, val time.Time) *FilterBasicAliasRecords {
	f.Filter.Col(`u."created_at"`, op, val)
	return f
}

func (f *FilterBasicAliasRecords) AND(nest ...func(*FilterBasicAliasRecords)) *FilterBasicAliasRecords {
	if len(nest) == 0 {
		f.Filter.AND()
		return f
	}
	f.Filter.AND(func(_ *gosql.Filter) {
		nest[0](f)
	})
	return f
}

func (f *FilterBasicAliasRecords) OR(nest ...func(*FilterBasicAliasRecords)) *FilterBasicAliasRecords {
	if len(nest) == 0 {
		f.Filter.OR()
		return f
	}
	f.Filter.OR(func(_ *gosql.Filter) {
		nest[0](f)
	})
	return f
}

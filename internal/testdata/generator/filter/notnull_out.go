// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"time"

	"github.com/frk/gosql/filter"
)

var _FilterNotNULLRecords_colmap = map[string]filter.Column{
	"A": {Name: `v."col_a"`, IsNULLable: true},
	"B": {Name: `v."col_b"`},
	"C": {Name: `v."col_c"`, IsNULLable: true},
	"D": {Name: `v."col_d"`},
	"E": {Name: `v."col_e"`, IsNULLable: true},
}

func (f *FilterNotNULLRecords) Init() {
	f.FilterMaker.InitV2(_FilterNotNULLRecords_colmap, "")
}

func (f *FilterNotNULLRecords) A(op string, val int) *FilterNotNULLRecords {
	f.FilterMaker.Col(`v."col_a"`, op, val)
	return f
}

func (f *FilterNotNULLRecords) B(op string, val string) *FilterNotNULLRecords {
	f.FilterMaker.Col(`v."col_b"`, op, val)
	return f
}

func (f *FilterNotNULLRecords) C(op string, val bool) *FilterNotNULLRecords {
	f.FilterMaker.Col(`v."col_c"`, op, val)
	return f
}

func (f *FilterNotNULLRecords) D(op string, val float64) *FilterNotNULLRecords {
	f.FilterMaker.Col(`v."col_d"`, op, val)
	return f
}

func (f *FilterNotNULLRecords) E(op string, val time.Time) *FilterNotNULLRecords {
	f.FilterMaker.Col(`v."col_e"`, op, val)
	return f
}

func (f *FilterNotNULLRecords) And(nest func(*FilterNotNULLRecords)) *FilterNotNULLRecords {
	if nest == nil {
		f.FilterMaker.And(nil)
		return f
	}
	f.FilterMaker.And(func() {
		nest(f)
	})
	return f
}

func (f *FilterNotNULLRecords) Or(nest func(*FilterNotNULLRecords)) *FilterNotNULLRecords {
	if nest == nil {
		f.FilterMaker.Or(nil)
		return f
	}
	f.FilterMaker.Or(func() {
		nest(f)
	})
	return f
}

package common

import (
	"github.com/frk/gosql/filter"
)

type FilterMaker struct {
	// ...
}

func (fm *FilterMaker) Init(colmap map[string]string, tscol string) {
	// ...
}

func (fm *FilterMaker) InitV2(colmap map[string]filter.Column, tscol string) {
	// ...
}

func (fm *FilterMaker) Col(col string, op string, val interface{}) {
	// ...
}

func (fm *FilterMaker) And(nest func()) {
	// ...
}

func (fm *FilterMaker) Or(nest func()) {
	// ...
}

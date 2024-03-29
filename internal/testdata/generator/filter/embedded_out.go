// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

var _FilterEmbeddedRecords_colmap = map[string]string{
	"fooField.barField.value": `n."foo_bar_baz_val"`,
	"fooField.bazField.value": `n."foo_baz_val"`,
	"barField.value":          `n."foo2_bar_baz_val"`,
	"bazField.value":          `n."foo2_baz_val"`,
}

func (f *FilterEmbeddedRecords) Init() {
	f.FilterMaker.Init(_FilterEmbeddedRecords_colmap, "")
}

func (f *FilterEmbeddedRecords) FOOBarEBazVal(op string, val string) *FilterEmbeddedRecords {
	f.FilterMaker.Col(`n."foo_bar_baz_val"`, op, val)
	return f
}

func (f *FilterEmbeddedRecords) FOOBazVal(op string, val string) *FilterEmbeddedRecords {
	f.FilterMaker.Col(`n."foo_baz_val"`, op, val)
	return f
}

func (f *FilterEmbeddedRecords) EFooBarEBazVal(op string, val string) *FilterEmbeddedRecords {
	f.FilterMaker.Col(`n."foo2_bar_baz_val"`, op, val)
	return f
}

func (f *FilterEmbeddedRecords) EFooBazVal(op string, val string) *FilterEmbeddedRecords {
	f.FilterMaker.Col(`n."foo2_baz_val"`, op, val)
	return f
}

func (f *FilterEmbeddedRecords) And(nest func(*FilterEmbeddedRecords)) *FilterEmbeddedRecords {
	if nest == nil {
		f.FilterMaker.And(nil)
		return f
	}
	f.FilterMaker.And(func() {
		nest(f)
	})
	return f
}

func (f *FilterEmbeddedRecords) Or(nest func(*FilterEmbeddedRecords)) *FilterEmbeddedRecords {
	if nest == nil {
		f.FilterMaker.Or(nil)
		return f
	}
	f.FilterMaker.Or(func() {
		nest(f)
	})
	return f
}

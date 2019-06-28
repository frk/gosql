package gosql

import (
	"go/types"
	"strings"

	"github.com/frk/gosql/internal/typesutil"
	"github.com/frk/tagutil"
)

func analyze(named *types.Named) (*command, error) {
	a := new(analyzer)
	a.pkg = named.Obj().Pkg().Path()
	a.cmd = &command{name: named.Obj().Name()}

	var ok bool
	if a.ctyp, ok = named.Underlying().(*types.Struct); !ok {
		return nil, &analysisError{code: badCmdTypeError, args: args{a.cmd.name}}
	}

	key := strings.ToLower(a.cmd.name)
	if len(key) > 5 {
		key = key[:6]
	}

	switch key {
	case "insert":
		a.cmd.typ = cmdtypeInsert
	case "update":
		a.cmd.typ = cmdtypeUpdate
	case "select":
		a.cmd.typ = cmdtypeSelect
	case "delete":
		a.cmd.typ = cmdtypeDelete
	case "filter":
		a.cmd.typ = cmdtypeFilter
	default:
		return nil, &analysisError{code: badCmdNameError, args: args{a.cmd.name}}
	}

	if err := a.run(); err != nil {
		return nil, err
	}
	return a.cmd, nil
}

// analyzer holds the state of the analysis
type analyzer struct {
	pkg  string        // the package path of the command under analysis
	ctyp *types.Struct // the struct type of the command under analysis
	rtyp *types.Struct // the struct type of the record under analysis
	cmd  *command      // the result of the analysis
}

func (a *analyzer) run() error {
	for i := 0; i < a.ctyp.NumFields(); i++ {
		f := a.ctyp.Field(i)

		tag := tagutil.New(a.ctyp.Tag(i))
		if rel := tag.First("rel"); len(rel) > 0 {
			if a.cmd.rec != nil {
				return &analysisError{code: manyRecordError, args: args{a.cmd.name}}
			}

			rec, err := a.analyzeRecord(f)
			if err != nil {
				return err
			}
			rec.rel.ident = a.analyzeIdent(rel)

			_ = a.rtyp // TODO analyze fields

			a.cmd.rec = rec
			continue
		}
	}

	if a.cmd.rec == nil {
		return &analysisError{code: noRecordError, args: args{a.cmd.name}}
	}

	return nil
}

// // analyze kind
// switch x := ftyp.(type) {
// case *types.Basic:
// 	rec.typ.kind = basickindmap[x.Kind()]
// case *types.Array:
// 	rec.typ.kind = gokindArray
// case *types.Slice:
// 	rec.typ.kind = gokindSlice
// case *types.Map:
// 	rec.typ.kind = gokindMap
// case *types.Pointer:
// 	rec.typ.kind = gokindPtr
// case *types.Struct:
// 	rec.typ.kind = gokindStruct
// default:
// 	// *types.Chan, *types.Signature, *types.Interface ...
// 	return &badRecordTypeError{fieldName: rec.field, cmdName: a.cmd.name}
// }

func (a *analyzer) analyzeRecord(field *types.Var) (*record, error) {
	var (
		rec   = &record{field: field.Name()}
		ftyp  = field.Type()
		named *types.Named
		err   error
		ok    bool
	)
	if named, ok = ftyp.(*types.Named); ok {
		ftyp = named.Underlying()
	}
	if iface, ok := ftyp.(*types.Interface); ok { // check for iterator interface
		if named, err = a.analyzeIterator(iface, rec); err != nil {
			return nil, err
		}
	} else if sig, ok := ftyp.(*types.Signature); ok { // check for iterator func
		if named, err = a.analyzeIteratorFunc(sig, rec); err != nil {
			return nil, err
		}
	}

	// if unnamed and not an iterator, check for slices and pointers
	if named == nil {
		if slice, ok := ftyp.(*types.Slice); ok { // allows []T
			ftyp = slice.Elem()
			rec.typ.isslice = true
		}
		if ptr, ok := ftyp.(*types.Pointer); ok { // allows *T
			ftyp = ptr.Elem()
			rec.typ.ispointer = true
		}
		if named, ok = ftyp.(*types.Named); ok {
			ftyp = named.Underlying()
		}

		// fail if still unnamed but is slice or pointer
		if named == nil && (rec.typ.isslice || rec.typ.ispointer) {
			return nil, &analysisError{code: badRecordTypeError, args: args{a.cmd.name}}
		}
	}

	if named != nil {
		pkg := named.Obj().Pkg()
		rec.typ.name = named.Obj().Name()
		rec.typ.pkgpath = pkg.Path()
		rec.typ.pkgname = pkg.Name()
		rec.typ.pkglocal = pkg.Name()
		rec.typ.isimported = (pkg.Path() != a.pkg)
		rec.typ.isscanner = false // TODO isscanner(named)
		rec.typ.isvaluer = false  // TODO isvaluer(named)
		rec.typ.istime = false    // TODO istime(named)
	}

	if a.rtyp, ok = ftyp.(*types.Struct); !ok {
		return nil, &analysisError{code: badRecordTypeError, args: args{a.cmd.name}}
	}

	rec.typ.kind = gokindStruct
	return rec, nil
}

func (a *analyzer) analyzeIterator(iface *types.Interface, rec *record) (*types.Named, error) {
	if iface.NumExplicitMethods() != 1 {
		return nil, &analysisError{code: badIteratorTypeError, args: args{a.cmd.name, rec.field}}
	}

	mth := iface.ExplicitMethod(0)
	sig := mth.Type().(*types.Signature)
	named, err := a.analyzeIteratorFunc(sig, rec)
	if err != nil {
		return nil, err
	}

	rec.iter.method = mth.Name()
	return named, nil
}

func (a *analyzer) analyzeIteratorFunc(sig *types.Signature, rec *record) (*types.Named, error) {
	// must take 1 argument and return one value of type error. "func(T) error"
	if sig.Params().Len() != 1 || sig.Results().Len() != 1 || !typesutil.IsError(sig.Results().At(0)) {
		return nil, &analysisError{code: badIteratorTypeError, args: args{a.cmd.name, rec.field}}
	}

	typ := sig.Params().At(0).Type()
	if ptr, ok := typ.(*types.Pointer); ok { // allows *T
		typ = ptr.Elem()
		rec.typ.ispointer = true
	}

	// make sure that the argument type is a named struct type
	named, ok := typ.(*types.Named)
	if !ok {
		return nil, &analysisError{code: badIteratorTypeError, args: args{a.cmd.name, rec.field}}
	} else if _, ok := named.Underlying().(*types.Struct); !ok {
		return nil, &analysisError{code: badIteratorTypeError, args: args{a.cmd.name, rec.field}}
	}

	rec.iter = new(iterator)
	return named, nil
}

// Used to analyze the value of the `rel` tag, the expected format is: "[qualifier.]name[:alias]".
func (a *analyzer) analyzeIdent(val string) (id ident) {
	if i := strings.LastIndexByte(val, '.'); i > -1 {
		id.qualifier = val[:i]
		val = val[i+1:]
	}
	if i := strings.LastIndexByte(val, ':'); i > -1 {
		id.alias = val[i+1:]
		val = val[:i]
	}
	id.name = val
	return id
}

type cmdtype uint

const (
	cmdtypeInsert cmdtype = iota + 1
	cmdtypeUpdate
	cmdtypeSelect
	cmdtypeDelete
	cmdtypeFilter
)

type command struct {
	name string  // name of the target struct type
	typ  cmdtype // the type of the command
	rec  *record
}

// record holds the combined information on a go struct type and on the
// db relation that's associated with that struct type.
type record struct {
	typ   gotype    // info about the record's go type
	rel   relation  // the associated relation
	field string    // name of the field that holds the record in the command's type
	iter  *iterator // if set, indicates that the record is handled by an iterator
}

type relation struct {
	ident  ident // the relation identifier
	isview bool  // indicates that the relation is a table view
}

// attribute holds information about a record's struct field and the corresponding db column.
type attribute struct {
	typ   gotype    // info about the field's type
	col   column    // the corresponding column
	field string    // name of the struct field
	path  fieldpath // a list of names of the field's parent structs, if any. (excluding the top record)
	// identifies the field's corresponding column as a primary key
	//
	// NOTE(mkopriva): This is used by default for UPDATEs which don't specify
	// a WHERE clause, if multiple fields are tagged as pkeys then we should
	// assume a composite primary key
	ispkey bool
	// indicates that the corresponding column's value is set automatically
	// by the database and therefore the column should be omitted
	// from the generated INSERT/UPDATE statements
	auto bool
	// indicates that the DEFAULT marker should be used during INSERT/UPDATE
	usedefault bool
	// indicates that if the field's value is EMPTY then NULL should
	// be stored in the column during INSERT/UPDATE
	nullempty bool
	// indicates that field should only be read from the database and never written
	readonly bool
	// indicates that field should only be written into the database and never read
	writeonly bool
	// indicates that the column value should be marshaled/unmarshaled
	// to/from json before/after being stored/retrieved.
	usejson bool
	// indicates that the column value should be wrapped
	// in a COALESCE call when read from the db.
	usecoalesce   bool
	coalescevalue string
	// for UPDATEs, if set to true, it indicates that the provided field
	// value should be added to the already existing column value.
	binadd bool

	// Tag  Tag
	// // In the context of statements that support the RETURNING clause this
	// // flag indicates that the field's column should be included in the clause
	// // and it's value scanned to the record's field.
	// Return    bool
	// Key           string
}

type fieldpath []string

type column struct {
	ident      ident  // the column identifier
	found      bool   // indicates that the column was found in the associated relation
	typname    string // name of the db type
	typisenum  bool   // indicates that the column's type is an enum type
	isnullable bool   // indicates that the column can be set to NULL
}

type iterator struct {
	method string
}

type ident struct {
	qualifier string
	name      string
	alias     string
}

type gotype struct {
	name       string // the name of a named type or empty string for unnamed types
	kind       gokind // the kind of the go type
	pkgpath    string // the package import path
	pkgname    string // the package's name
	pkglocal   string // the local package name (including ".")
	isimported bool   // indicates whether or not the package is imported
	isslice    bool   // reports whether or not the type's a slice type
	ispointer  bool   // reports whether or not the type's a pointer type
	isscanner  bool   // reports whether or not the type implements the sql.Scanner interface
	isvaluer   bool   // reports whether or not the type implements the driver.Valuer interface
	istime     bool   // reposrts whether or not the type is time.Time
}

type gokind uint

const (
	// basic
	gokindInvalid gokind = iota
	gokindBool
	gokindInt
	gokindInt8
	gokindInt16
	gokindInt32
	gokindInt64
	gokindUint
	gokindUint8
	gokindUint16
	gokindUint32
	gokindUint64
	gokindUintptr
	gokindFloat32
	gokindFloat64
	gokindComplex64
	gokindComplex128
	gokindString
	gokindUnsafeptr

	// non-basic
	gokindArray
	gokindChan
	gokindFunc
	gokindInterface
	gokindMap
	gokindPtr
	gokindSlice
	gokindStruct
)

var basickindmap = map[types.BasicKind]gokind{
	types.Invalid:       gokindInvalid,
	types.Bool:          gokindBool,
	types.Int:           gokindInt,
	types.Int8:          gokindInt8,
	types.Int16:         gokindInt16,
	types.Int32:         gokindInt32,
	types.Int64:         gokindInt64,
	types.Uint:          gokindUint,
	types.Uint8:         gokindUint8,
	types.Uint16:        gokindUint16,
	types.Uint32:        gokindUint32,
	types.Uint64:        gokindUint64,
	types.Uintptr:       gokindUintptr,
	types.Float32:       gokindFloat32,
	types.Float64:       gokindFloat64,
	types.Complex64:     gokindComplex64,
	types.Complex128:    gokindComplex128,
	types.String:        gokindString,
	types.UnsafePointer: gokindUnsafeptr,
}
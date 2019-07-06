package gosql

import (
	"go/types"
	"regexp"
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
	rtyp *types.Struct // the struct type of the relation under analysis
	cmd  *command      // the result of the analysis
}

func (a *analyzer) run() error {
	for i := 0; i < a.ctyp.NumFields(); i++ {
		fld := a.ctyp.Field(i)
		tag := tagutil.New(a.ctyp.Tag(i))

		if reltag := tag.First("rel"); len(reltag) > 0 {
			if a.cmd.rel != nil {
				return &analysisError{code: manyRecordError, args: args{a.cmd.name}}
			}

			a.cmd.rel = new(relinfo)
			a.cmd.rel.field = fld.Name()
			a.cmd.rel.ident = a.analyzeIdent(reltag)
			if err := a.analyzeDatatype(fld); err != nil {
				return err
			}

			_ = a.rtyp // TODO analyze fields
			continue
		}
	}

	if a.cmd.rel == nil {
		return &analysisError{code: noRecordError, args: args{a.cmd.name}}
	}

	return nil
}

func (a *analyzer) analyzeDatatype(field *types.Var) error {
	var (
		rel   = a.cmd.rel
		ftyp  = field.Type()
		named *types.Named
		err   error
		ok    bool
	)
	if named, ok = ftyp.(*types.Named); ok {
		ftyp = named.Underlying()
	}

	// Check whether the relation field's type is an interface or a function,
	// if so, it is then expected to be a "valid" iterator, and it is analyzed as such.
	//
	// Failure of the iterator analysis will cause the whole analysis to exit
	// as there's currently no support for non-iterator interfaces nor functions.
	if iface, ok := ftyp.(*types.Interface); ok {
		if named, err = a.analyzeIterator(iface, rel); err != nil {
			return err
		}
	} else if sig, ok := ftyp.(*types.Signature); ok {
		if named, err = a.analyzeIteratorFunc(sig, rel); err != nil {
			return err
		}
	}

	// if unnamed and not an iterator, check for slices and pointers
	if named == nil {
		if slice, ok := ftyp.(*types.Slice); ok { // allows []T
			ftyp = slice.Elem()
			rel.datatype.isslice = true
		}
		if ptr, ok := ftyp.(*types.Pointer); ok { // allows *T
			ftyp = ptr.Elem()
			rel.datatype.ispointer = true
		}
		if named, ok = ftyp.(*types.Named); !ok {
			// fail if still unnamed but is slice or pointer
			if rel.datatype.isslice || rel.datatype.ispointer {
				return &analysisError{code: badRecordTypeError, args: args{a.cmd.name}}
			}
		}
	}

	if named != nil {
		pkg := named.Obj().Pkg()
		rel.datatype.name = named.Obj().Name()
		rel.datatype.pkgpath = pkg.Path()
		rel.datatype.pkgname = pkg.Name()
		rel.datatype.pkglocal = pkg.Name()
		rel.datatype.isimported = (pkg.Path() != a.pkg)
		rel.datatype.isscanner = typesutil.ImplementsScanner(named)
		rel.datatype.isvaluer = typesutil.ImplementsValuer(named)
		rel.datatype.istime = typesutil.IsTime(named)
		rel.datatype.isafterscanner = typesutil.ImplementsAfterScanner(named)
		ftyp = named.Underlying()
	}

	rel.datatype.kind = a.analyzeKind(ftyp)
	if rel.datatype.kind != kindStruct {
		// NOTE currently only the struct kind is supported as the relation's associated datatype
		return &analysisError{code: badRecordTypeError, args: args{a.cmd.name}}
	}

	styp := ftyp.(*types.Struct)
	return a.analyzeFields(styp)
}

func (a *analyzer) analyzeFields(styp *types.Struct) error {
	type iteration struct {
		styp *types.Struct
		typ  *typeinfo
		idx  int    // keep track of the field index
		pfx  string // column prefix
	}

	// lifo stack
	stack := []*iteration{{styp: styp, typ: &a.cmd.rel.datatype.typeinfo}}

stackloop:
	for len(stack) > 0 {
		i := stack[len(stack)-1]
		for i.idx < i.styp.NumFields() {
			fld := i.styp.Field(i.idx)
			tag := tagutil.New(i.styp.Tag(i.idx))
			sqltag := tag.First("sql")

			// instead of incrementing the index in the for-statement
			// it is done here manually to ensure that it is not skipped
			// when continuing to the outer loop
			i.idx++

			// ignore field if:
			// - no column name or sql tag was provided
			if sqltag == "" ||
				// - explicitly marked to be ignored
				sqltag == "-" ||
				// - has blank name, i.e. it's practically inaccessible
				fld.Name() == "_" ||
				// - it's unexported and the field's struct type is imported
				(!fld.Exported() && i.typ.isimported) {
				continue
			}

			f := new(fieldinfo)
			f.name = fld.Name()
			f.isembedded = fld.Embedded()
			f.isexported = fld.Exported()
			f.tag = tag
			f.auto = tag.HasOption("sql", "auto")
			f.ispkey = tag.HasOption("sql", "pk")
			f.nullempty = tag.HasOption("sql", "nullempty")
			f.readonly = tag.HasOption("sql", "ro")
			f.writeonly = tag.HasOption("sql", "wo")
			f.usejson = tag.HasOption("sql", "json")
			f.binadd = tag.HasOption("sql", "+")
			f.coalesce = a.analyzeCoalesceinfo(tag)

			ftyp := fld.Type()
			if slice, ok := ftyp.(*types.Slice); ok {
				f.typ.isslice = true
				ftyp = slice.Elem()
			}
			if ptr, ok := ftyp.(*types.Pointer); ok {
				f.typ.ispointer = true
				ftyp = ptr.Elem()
			}
			if named, ok := ftyp.(*types.Named); ok {
				pkg := named.Obj().Pkg()
				f.typ.name = named.Obj().Name()
				f.typ.pkgpath = pkg.Path()
				f.typ.pkgname = pkg.Name()
				f.typ.pkglocal = pkg.Name()
				f.typ.isimported = (pkg.Path() != a.pkg)
				f.typ.isscanner = typesutil.ImplementsScanner(named)
				f.typ.isvaluer = typesutil.ImplementsValuer(named)
				f.typ.istime = typesutil.IsTime(named)
				ftyp = named.Underlying()
			}
			f.typ.kind = a.analyzeKind(ftyp)
			i.typ.fields = append(i.typ.fields, f)

			// if the field's type is a struct and the `sql` tag's
			// value starts with the ">" (descend) marker, then
			// analyze its fields as well
			if f.typ.kind == kindStruct && strings.HasPrefix(sqltag, ">") && !f.typ.isslice {
				j := &iteration{styp: ftyp.(*types.Struct), typ: &f.typ}
				j.pfx = i.pfx + strings.TrimPrefix(sqltag, ">")
				stack = append(stack, j)
				continue stackloop
			}

			f.column.ident = a.analyzeIdent(i.pfx + sqltag)
		}
		stack = stack[:len(stack)-1]
	}
	return nil
}

func (a *analyzer) analyzeIterator(iface *types.Interface, rel *relinfo) (*types.Named, error) {
	if iface.NumExplicitMethods() != 1 {
		return nil, &analysisError{code: badIteratorTypeError, args: args{a.cmd.name, rel.field}}
	}

	mth := iface.ExplicitMethod(0)
	sig := mth.Type().(*types.Signature)
	named, err := a.analyzeIteratorFunc(sig, rel)
	if err != nil {
		return nil, err
	}

	rel.datatype.iter.method = mth.Name()
	return named, nil
}

func (a *analyzer) analyzeIteratorFunc(sig *types.Signature, rel *relinfo) (*types.Named, error) {
	// must take 1 argument and return one value of type error. "func(T) error"
	if sig.Params().Len() != 1 || sig.Results().Len() != 1 || !typesutil.IsError(sig.Results().At(0).Type()) {
		return nil, &analysisError{code: badIteratorTypeError, args: args{a.cmd.name, rel.field}}
	}

	typ := sig.Params().At(0).Type()
	if ptr, ok := typ.(*types.Pointer); ok { // allows *T
		typ = ptr.Elem()
		rel.datatype.ispointer = true
	}

	// make sure that the argument type is a named struct type
	named, ok := typ.(*types.Named)
	if !ok {
		return nil, &analysisError{code: badIteratorTypeError, args: args{a.cmd.name, rel.field}}
	} else if _, ok := named.Underlying().(*types.Struct); !ok {
		return nil, &analysisError{code: badIteratorTypeError, args: args{a.cmd.name, rel.field}}
	}

	rel.datatype.iter = new(iterator)
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

func (a *analyzer) analyzeKind(typ types.Type) typekind {
	switch x := typ.(type) {
	case *types.Basic:
		return basickindmap[x.Kind()]
	case *types.Array:
		return kindArray
	case *types.Chan:
		return kindChan
	case *types.Signature:
		return kindFunc
	case *types.Interface:
		return kindInterface
	case *types.Map:
		return kindMap
	case *types.Pointer:
		return kindPtr
	case *types.Slice:
		return kindSlice
	case *types.Struct:
		return kindStruct
	}
	return 0 // unsupported / unknown
}

var reCoalesceValue = regexp.MustCompile(`(?i)^coalesce$|^coalesce\((.*)\)$`)

func (a *analyzer) analyzeCoalesceinfo(tag tagutil.Tag) *coalesceinfo {
	if sqltag := tag["sql"]; len(sqltag) > 0 {
		for _, opt := range sqltag[1:] {
			if strings.HasPrefix(opt, "coalesce") {
				cls := new(coalesceinfo)
				if match := reCoalesceValue.FindStringSubmatch(opt); len(match) > 1 {
					cls.defval = match[1]
				}
				return cls
			}
		}
	}
	return nil
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
	rel  *relinfo
}

// relinfo holds the information on a go struct type and on the
// db relation that's associated with that struct type.
type relinfo struct {
	field    string // name of the field that references the relation in the command
	ident    ident  // the relation identifier
	datatype datatype
	isview   bool // indicates that the relation is a table view
}

// datatype holds information on the type of data a command should read from,
// or write to, the associated database relation.
type datatype struct {
	typeinfo // type info on the datatype
	// if set, indicates that the datatype is handled by an iterator
	iter *iterator
	// reports whether or not the type implements the afterscanner interface
	isafterscanner bool
}

type iterator struct {
	method string
}

type ident struct {
	qualifier string
	name      string
	alias     string
}

type typeinfo struct {
	name       string   // the name of a named type or empty string for unnamed types
	kind       typekind // the kind of the go type
	pkgpath    string   // the package import path
	pkgname    string   // the package's name
	pkglocal   string   // the local package name (including ".")
	isimported bool     // indicates whether or not the package is imported
	isslice    bool     // reports whether or not the type's a slice type
	ispointer  bool     // reports whether or not the type's a pointer type
	isscanner  bool     // reports whether or not the type implements the sql.Scanner interface
	isvaluer   bool     // reports whether or not the type implements the driver.Valuer interface
	istime     bool     // reposrts whether or not the type is time.Time
	// if the typeinfo represents a struct type then this slice will hold
	// the info about the fields of that struct type
	fields []*fieldinfo
}

// fieldinfo holds information about a record's struct field and the corresponding db column.
type fieldinfo struct {
	typ  typeinfo // info about the field's type
	name string   // name of the struct field
	// indicates whether or not the field is embedded
	isembedded bool
	// indicates whether or not the field is exported
	isexported bool
	// the field's parsed tag
	tag tagutil.Tag
	// the corresponding column
	column column
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
	// if set it indicates that the column value should be wrapped
	// in a COALESCE call when read from the db.
	coalesce *coalesceinfo
	// for UPDATEs, if set to true, it indicates that the provided field
	// value should be added to the already existing column value.
	binadd bool
}

type column struct {
	ident      ident  // the column identifier
	found      bool   // indicates that the column was found in the associated relation
	typname    string // name of the db type
	typisenum  bool   // indicates that the column's type is an enum type
	isnullable bool   // indicates that the column can be set to NULL
}

type coalesceinfo struct {
	defval string
}

type typekind uint

const (
	// basic
	kindInvalid typekind = iota
	kindBool
	kindInt
	kindInt8
	kindInt16
	kindInt32
	kindInt64
	kindUint
	kindUint8
	kindUint16
	kindUint32
	kindUint64
	kindUintptr
	kindFloat32
	kindFloat64
	kindComplex64
	kindComplex128
	kindString
	kindUnsafeptr

	// non-basic
	kindArray
	kindChan
	kindFunc
	kindInterface
	kindMap
	kindPtr
	kindSlice
	kindStruct
)

var basickindmap = map[types.BasicKind]typekind{
	types.Invalid:       kindInvalid,
	types.Bool:          kindBool,
	types.Int:           kindInt,
	types.Int8:          kindInt8,
	types.Int16:         kindInt16,
	types.Int32:         kindInt32,
	types.Int64:         kindInt64,
	types.Uint:          kindUint,
	types.Uint8:         kindUint8,
	types.Uint16:        kindUint16,
	types.Uint32:        kindUint32,
	types.Uint64:        kindUint64,
	types.Uintptr:       kindUintptr,
	types.Float32:       kindFloat32,
	types.Float64:       kindFloat64,
	types.Complex64:     kindComplex64,
	types.Complex128:    kindComplex128,
	types.String:        kindString,
	types.UnsafePointer: kindUnsafeptr,
}

var typekind2string = map[typekind]string{
	// builtin basic only
	kindBool:       "bool",
	kindInt:        "int",
	kindInt8:       "int8",
	kindInt16:      "int16",
	kindInt32:      "int32",
	kindInt64:      "int64",
	kindUint:       "uint",
	kindUint8:      "uint8",
	kindUint16:     "uint16",
	kindUint32:     "uint32",
	kindUint64:     "uint64",
	kindUintptr:    "uintptr",
	kindFloat32:    "float32",
	kindFloat64:    "float64",
	kindComplex64:  "complex64",
	kindComplex128: "complex128",
	kindString:     "string",
}
package gosql

import (
	"go/types"
	"regexp"
	"strconv"
	"strings"

	"github.com/frk/gosql/internal/errors"
	"github.com/frk/gosql/internal/typesutil"
	"github.com/frk/tagutil"
)

var (
	// NOTE(mkopriva): Identifiers must begin with a letter (a-z) or an underscore (_).
	// Subsequent characters in an identifier can be letters, underscores, and digits (0-9).

	// Matches a valid identifier.
	reident = regexp.MustCompile(`^[A-Za-z_]\w*$`)

	// Matches a valid identifier. The identifier can optionally be prefixed
	// by another identifier and concatenated to it by a dot. It can also have
	// another optional identifier at the right end concatenated to id by a colon.
	// Expected input format: [schema_name.]relation_name[:alias_name]
	rerelid = regexp.MustCompile(`^(?:[A-Za-z_]\w*\.)?[A-Za-z_]\w*(?:\:[A-Za-z_]\w*)?$`)

	// Matches a valid identifier. The identifier can optionally be prefixed
	// by another identifier and concatenated to it by a dot.
	// Expected input format: [rel_alias.]column_name
	recolid = regexp.MustCompile(`^(?:[A-Za-z_]\w*\.)?[A-Za-z_]\w*$`)

	// Matches a few reserved identifiers.
	rereserved = regexp.MustCompile(`^(?i:true|false|` +
		`current_date|current_time|current_timestamp|` +
		`current_role|current_schema|current_user|` +
		`localtime|localtimestamp|` +
		`session_user)$`)

	// Matches coalesce or coalesce(<value>) where <value> is expected to
	// be a single value literal.
	recoalesce = regexp.MustCompile(`(?i)^coalesce$|^coalesce\((.*)\)$`)
)

// TODO(mkopriva): to provide more detailed error messages either pass in the
// details about the file under analysis, or make sure that the caller has that
// information and appends it to the error.
func analyze(named *types.Named) (*typespec, error) {
	a := new(analyzer)
	a.pkg = named.Obj().Pkg().Path()
	a.named = named
	a.spec = new(typespec)
	a.spec.name = named.Obj().Name()

	var ok bool
	if a.spectyp, ok = named.Underlying().(*types.Struct); !ok {
		panic(a.spec.name + " typespec kind not supported.") // this shouldn't happen
	}

	key := strings.ToLower(a.spec.name)
	if len(key) > 5 {
		key = key[:6]
	}

	switch key {
	case "insert":
		a.spec.kind = speckindInsert
	case "update":
		a.spec.kind = speckindUpdate
	case "select":
		a.spec.kind = speckindSelect
	case "delete":
		a.spec.kind = speckindDelete
	case "filter":
		a.spec.kind = speckindFilter
	default:
		panic(a.spec.name + " typespec kind has unsupported name prefix.") // this shouldn't happen
	}

	if err := a.run(); err != nil {
		return nil, err
	}
	return a.spec, nil
}

// analyzer holds the state of the analysis
type analyzer struct {
	pkg     string        // the package path of the typespec under analysis
	named   *types.Named  // the named type of the typespec under analysis
	spectyp *types.Struct // the struct type of the typespec under analysis
	reltyp  *types.Struct // the struct type of the relation under analysis
	spec    *typespec     // the result of the analysis
}

func (a *analyzer) run() (err error) {
	for i := 0; i < a.spectyp.NumFields(); i++ {
		fld := a.spectyp.Field(i)
		tag := tagutil.New(a.spectyp.Tag(i))

		if reltag := tag.First("rel"); len(reltag) > 0 {
			rid, err := a.relid(reltag, fld)
			if err != nil {
				return err
			}

			a.spec.rel = new(relfield)
			a.spec.rel.name = fld.Name()
			a.spec.rel.relid = rid

			switch fname := strings.ToLower(a.spec.rel.name); {
			case fname == "count" && a.isint(fld.Type()):
				if a.spec.kind != speckindSelect {
					return errors.IllegalCountFieldError

				}
				a.spec.selkind = selectcount
			case fname == "exists" && a.isbool(fld.Type()):
				if a.spec.kind != speckindSelect {
					return errors.IllegalExistsFieldError
				}
				a.spec.selkind = selectexists
			case fname == "notexists" && a.isbool(fld.Type()):
				if a.spec.kind != speckindSelect {
					return errors.IllegalNotExistsFieldError
				}
				a.spec.selkind = selectnotexists
			case fname == "_" && typesutil.IsDirective("Relation", fld.Type()):
				if a.spec.kind != speckindDelete {
					return errors.IllegalRelationDirectiveError
				}
				a.spec.rel.isdir = true
			default:
				if err := a.relrecordtype(a.spec.rel, fld); err != nil {
					return err
				}
			}
			continue
		}

		// TODO(mkopriva): allow for embedding a struct with "common feature fields",
		// and make sure to also allow imported and local-unexported struct types.

		// fields with gosql directive types
		if dirname := typesutil.GetDirectiveName(fld); fld.Name() == "_" && len(dirname) > 0 {
			switch strings.ToLower(dirname) {
			case "all":
				if a.spec.kind != speckindUpdate && a.spec.kind != speckindDelete {
					return errors.IllegalAllDirectiveError
				}
				if a.spec.all || a.spec.where != nil || len(a.spec.filter) > 0 {
					return errors.ConflictWhereProducerError
				}
				a.spec.all = true
			case "default":
				if a.spec.kind != speckindInsert && a.spec.kind != speckindUpdate {
					return errors.IllegalDefaultDirectiveError
				}
				if a.spec.defaults, err = a.collist(tag["sql"], fld); err != nil {
					return err
				}
			case "force":
				if a.spec.kind != speckindInsert && a.spec.kind != speckindUpdate {
					return errors.IllegalForceDirectiveError
				}
				if a.spec.force, err = a.collist(tag["sql"], fld); err != nil {
					return err
				}
			case "return":
				if len(a.spec.rel.rec.fields) == 0 {
					// TODO test
					return errors.ReturnDirectiveWithNoRelfieldError
				}
				if a.spec.kind != speckindInsert && a.spec.kind != speckindUpdate && a.spec.kind != speckindDelete {
					return errors.IllegalReturnDirectiveError
				}
				if a.spec.returning != nil || a.spec.result != nil || len(a.spec.rowsaffected) > 0 {
					return errors.ConflictResultProducerError
				}
				if a.spec.returning, err = a.collist(tag["sql"], fld); err != nil {
					return err
				}
			case "limit":
				if err := a.limitvar(fld, tag.First("sql")); err != nil {
					return err
				}
			case "offset":
				if err := a.offsetvar(fld, tag.First("sql")); err != nil {
					return err
				}
			case "orderby":
				if err := a.orderbydir(tag["sql"], fld); err != nil {
					return err
				}
			case "override":
				if err := a.overridedir(tag.First("sql"), fld); err != nil {
					return err
				}
			case "textsearch":
				if err := a.textsearch(tag.First("sql"), fld); err != nil {
					return err
				}
			default:
				return errors.IllegalCommandDirectiveError
			}
			continue
		}

		// fields with specific names
		switch fname := strings.ToLower(fld.Name()); fname {
		case "where":
			if err := a.whereblock(fld); err != nil {
				return err
			}
		case "join", "from", "using":
			if err := a.joinblock(fld); err != nil {
				return err
			}
		case "onconflict":
			if err := a.onconflictblock(fld); err != nil {
				return err
			}
		case "result":
			if err := a.resultfield(fld); err != nil {
				return err
			}
		case "limit":
			if err := a.limitvar(fld, tag.First("sql")); err != nil {
				return err
			}
		case "offset":
			if err := a.offsetvar(fld, tag.First("sql")); err != nil {
				return err
			}
		case "rowsaffected":
			if err := a.rowsaffected(fld); err != nil {
				return err
			}
		default:
			// if no match by field name, look for specific field types
			if a.isaccessible(fld, a.named) {
				switch {
				case a.isfilter(fld.Type()):
					if a.spec.kind != speckindSelect && a.spec.kind != speckindUpdate && a.spec.kind != speckindDelete {
						return errors.IllegalFilterFieldError
					}
					if a.spec.all || a.spec.where != nil || len(a.spec.filter) > 0 {
						return errors.ConflictWhereProducerError
					}
					a.spec.filter = fld.Name()
				case a.iserrorhandler(fld.Type()):
					if len(a.spec.erh) > 0 {
						return errors.ConflictErrorHandlerFieldError
					}
					a.spec.erh = fld.Name()
				}
			}
		}

	}

	if a.spec.rel == nil {
		return errors.NoRelfieldError
	}
	return nil
}

func (a *analyzer) relrecordtype(rel *relfield, field *types.Var) error {
	var (
		ftyp  = field.Type()
		named *types.Named
		err   error
		ok    bool
	)
	if named, ok = ftyp.(*types.Named); ok {
		ftyp = named.Underlying()
	}

	// Check whether the relation field's type is an interface or a function,
	// if so, it is then expected to be an iterator, and it is analyzed as such.
	//
	// Failure of the iterator analysis will cause the whole analysis to exit
	// as there's currently no support for non-iterator interfaces nor functions.
	if iface, ok := ftyp.(*types.Interface); ok {
		if named, err = a.iterator(iface, named, rel); err != nil {
			return err
		}
	} else if sig, ok := ftyp.(*types.Signature); ok {
		if named, err = a.iteratorfunc(sig, rel); err != nil {
			return err
		}
	} else {
		// If not an iterator, check for slices, arrays, and pointers.
		if slice, ok := ftyp.(*types.Slice); ok { // allows []T / []*T
			ftyp = slice.Elem()
			rel.rec.isslice = true
		} else if array, ok := ftyp.(*types.Array); ok { // allows [N]T / [N]*T
			ftyp = array.Elem()
			rel.rec.isarray = true
			rel.rec.arraylen = array.Len()
		}
		if ptr, ok := ftyp.(*types.Pointer); ok { // allows *T
			ftyp = ptr.Elem()
			rel.rec.ispointer = true
		}

		// Get the name of the base type, if applicable.
		if rel.rec.isslice || rel.rec.isarray || rel.rec.ispointer {
			if named, ok = ftyp.(*types.Named); !ok {
				// Fail if the type is a slice, an array, or a pointer
				// while its base type remains unnamed.
				return errors.BadRelfieldTypeError
			}
		}
	}

	if named != nil {
		pkg := named.Obj().Pkg()
		rel.rec.base.name = named.Obj().Name()
		rel.rec.base.pkgpath = pkg.Path()
		rel.rec.base.pkgname = pkg.Name()
		rel.rec.base.pkglocal = pkg.Name()
		rel.rec.base.isimported = a.isimported(named)
		rel.rec.isafterscanner = typesutil.ImplementsAfterScanner(named)
		ftyp = named.Underlying()
	}

	rel.rec.base.kind = a.typekind(ftyp)
	if rel.rec.base.kind != kindstruct {
		return errors.BadRelfieldTypeError
	}

	styp := ftyp.(*types.Struct)
	return a.relfields(rel, styp)
}

func (a *analyzer) relfields(rel *relfield, styp *types.Struct) error {
	// The loopstate type holds the state of a loop over a struct's fields.
	type loopstate struct {
		styp *types.Struct // the struct type whose fields are being analyzed
		typ  *typeinfo     // info on the struct type; holds the resulting slice of analyzed fieldinfo
		idx  int           // keeps track of the field index
		pfx  string        // column prefix
		path []*fieldelem
	}

	// LIFO stack of states used for depth first traversal of struct fields.
	stack := []*loopstate{{styp: styp, typ: &rel.rec.base}}

stackloop:
	for len(stack) > 0 {
		loop := stack[len(stack)-1]
		for loop.idx < loop.styp.NumFields() {
			fld := loop.styp.Field(loop.idx)
			tag := tagutil.New(loop.styp.Tag(loop.idx))
			sqltag := tag.First("sql")

			// Instead of incrementing the index in the for-statement
			// it is done here manually to ensure that it is not skipped
			// when continuing to the outer loop.
			loop.idx++

			// Ignore the field if:
			// - no column name or sql tag was provided
			if sqltag == "" ||
				// - explicitly marked to be ignored
				sqltag == "-" ||
				// - has blank name, i.e. it's practically inaccessible
				fld.Name() == "_" ||
				// - it's unexported and the field's struct type is imported
				(!fld.Exported() && loop.typ.isimported) {
				continue
			}

			f := new(fieldinfo)
			f.tag = tag
			f.name = fld.Name()
			f.isembedded = fld.Embedded()
			f.isexported = fld.Exported()

			// Analyze the field's type.
			ftyp := fld.Type()
			f.typ, ftyp = a.typeinfo(ftyp)

			// If the field's type is a struct and the `sql` tag's
			// value starts with the ">" (descend) marker, then it is
			// considered to be a "parent" field element whose child
			// fields need to be analyzed as well.
			if f.typ.is(kindstruct) && strings.HasPrefix(sqltag, ">") {
				loop2 := new(loopstate)
				loop2.styp = ftyp.(*types.Struct)
				loop2.typ = &f.typ
				loop2.pfx = loop.pfx + strings.TrimPrefix(sqltag, ">")

				// Allocate path of the appropriate size an copy it.
				loop2.path = make([]*fieldelem, len(loop.path))
				_ = copy(loop2.path, loop.path)

				// If the parent node is a pointer to a struct,
				// get the struct type info.
				typ := f.typ
				if typ.kind == kindptr {
					typ = *typ.elem
				}

				fe := new(fieldelem)
				fe.name = f.name
				fe.tag = f.tag
				fe.isembedded = f.isembedded
				fe.isexported = f.isexported
				fe.typename = typ.name
				fe.typepkgpath = typ.pkgpath
				fe.typepkgname = typ.pkgname
				fe.typepkglocal = typ.pkglocal
				fe.isimported = typ.isimported
				fe.ispointer = (f.typ.kind == kindptr)
				loop2.path = append(loop2.path, fe)

				stack = append(stack, loop2)
				continue stackloop
			}

			// If the field is not a struct to be descended,
			// it is considered to be a "leaf" field and as
			// such the analysis of leaf-specific information
			// needs to be carried out.
			f.path = loop.path
			f.auto = tag.HasOption("sql", "auto")
			f.ispkey = tag.HasOption("sql", "pk")
			f.nullempty = tag.HasOption("sql", "nullempty")
			f.readonly = tag.HasOption("sql", "ro")
			f.writeonly = tag.HasOption("sql", "wo")
			f.usejson = tag.HasOption("sql", "json")
			f.usexml = tag.HasOption("sql", "xml")
			f.binadd = tag.HasOption("sql", "+")
			f.cancast = tag.HasOption("sql", "cast")
			f.usecoalesce, f.coalesceval = a.coalesceinfo(tag)

			// Resolve the column id.
			colid, err := a.colid(loop.pfx+sqltag, fld)
			if err != nil {
				return err
			}
			f.colid = colid

			// Add the field to the list.
			rel.rec.fields = append(rel.rec.fields, f)
		}
		stack = stack[:len(stack)-1]
	}
	return nil
}

// The typeinfo method analyzes the given type and returns the result. The analysis
// looks only for information of "named types" and in case of slice, array, map,
// or pointer types it will analyze the element type of those types. The second
// return value is the base element type of the given type.
func (a *analyzer) typeinfo(tt types.Type) (typ typeinfo, base types.Type) {
	base = tt

	if named, ok := base.(*types.Named); ok {
		pkg := named.Obj().Pkg()
		typ.name = named.Obj().Name()
		typ.pkgpath = pkg.Path()
		typ.pkgname = pkg.Name()
		typ.pkglocal = pkg.Name()
		typ.isimported = a.isimported(named)
		typ.isscanner = typesutil.ImplementsScanner(named)
		typ.isvaluer = typesutil.ImplementsValuer(named)
		typ.istime = typesutil.IsTime(named)
		typ.isjsmarshaler = typesutil.ImplementsJSONMarshaler(named)
		typ.isjsunmarshaler = typesutil.ImplementsJSONUnmarshaler(named)
		typ.isxmlmarshaler = typesutil.ImplementsXMLMarshaler(named)
		typ.isxmlunmarshaler = typesutil.ImplementsXMLUnmarshaler(named)
		base = named.Underlying()
	}

	typ.kind = a.typekind(base)

	var elem typeinfo // element info
	switch T := base.(type) {
	case *types.Basic:
		typ.isrune = T.Name() == "rune"
		typ.isbyte = T.Name() == "byte"
	case *types.Slice:
		elem, base = a.typeinfo(T.Elem())
		typ.elem = &elem
	case *types.Array:
		elem, base = a.typeinfo(T.Elem())
		typ.elem = &elem
		typ.arraylen = T.Len()
	case *types.Map:
		key, _ := a.typeinfo(T.Key())
		elem, base = a.typeinfo(T.Elem())
		typ.key = &key
		typ.elem = &elem
	case *types.Pointer:
		elem, base = a.typeinfo(T.Elem())
		typ.elem = &elem
	case *types.Interface:
		// If base is an unnamed interface type check at least whether
		// or not it declares, or embeds, one of the relevant methods.
		if typ.name == "" {
			typ.isscanner = typesutil.IsScanner(T)
			typ.isvaluer = typesutil.IsValuer(T)
			typ.isjsmarshaler = typesutil.IsJSONMarshaler(T)
			typ.isjsunmarshaler = typesutil.IsJSONUnmarshaler(T)
			typ.isxmlmarshaler = typesutil.IsXMLMarshaler(T)
			typ.isxmlunmarshaler = typesutil.IsXMLUnmarshaler(T)
		}
	}
	return typ, base
}

func (a *analyzer) iterator(iface *types.Interface, named *types.Named, rel *relfield) (*types.Named, error) {
	if iface.NumExplicitMethods() != 1 {
		return nil, errors.BadIteratorTypeError
	}

	mth := iface.ExplicitMethod(0)

	// Make sure that the method is exported or, if it's not, then at least
	// ensure that the receiver type is local, i.e. not imported, otherwise
	// the method will not be accessible.
	if !mth.Exported() && named != nil && (named.Obj().Pkg().Path() != a.pkg) {
		return nil, errors.BadIteratorTypeError
	}

	sig := mth.Type().(*types.Signature)
	named, err := a.iteratorfunc(sig, rel)
	if err != nil {
		return nil, err
	}

	rel.rec.itermethod = mth.Name()
	return named, nil
}

func (a *analyzer) iteratorfunc(sig *types.Signature, rel *relfield) (*types.Named, error) {
	// Must take 1 argument and return one value of type error. "func(T) error"
	if sig.Params().Len() != 1 || sig.Results().Len() != 1 || !typesutil.IsError(sig.Results().At(0).Type()) {
		return nil, errors.BadIteratorTypeError
	}

	typ := sig.Params().At(0).Type()
	if ptr, ok := typ.(*types.Pointer); ok { // allows *T
		typ = ptr.Elem()
		rel.rec.ispointer = true
	}

	// Make sure that the argument type is a named struct type.
	named, ok := typ.(*types.Named)
	if !ok {
		return nil, errors.BadIteratorTypeError
	} else if _, ok := named.Underlying().(*types.Struct); !ok {
		return nil, errors.BadIteratorTypeError
	}

	rel.rec.isiter = true
	return named, nil
}

func (a *analyzer) typekind(typ types.Type) typekind {
	switch x := typ.(type) {
	case *types.Basic:
		return basickind2typekind[x.Kind()]
	case *types.Array:
		return kindarray
	case *types.Chan:
		return kindchan
	case *types.Signature:
		return kindfunc
	case *types.Interface:
		return kindinterface
	case *types.Map:
		return kindmap
	case *types.Pointer:
		return kindptr
	case *types.Slice:
		return kindslice
	case *types.Struct:
		return kindstruct
	}
	return 0 // unsupported / unknown
}

func (a *analyzer) coalesceinfo(tag tagutil.Tag) (use bool, val string) {
	if sqltag := tag["sql"]; len(sqltag) > 0 {
		for _, opt := range sqltag[1:] {
			if strings.HasPrefix(opt, "coalesce") {
				use = true
				if match := recoalesce.FindStringSubmatch(opt); len(match) > 1 {
					val = match[1]
				}
				break
			}
		}
	}
	return use, val
}

func (a *analyzer) whereblock(field *types.Var) (err error) {
	if a.spec.kind != speckindSelect && a.spec.kind != speckindUpdate && a.spec.kind != speckindDelete {
		return errors.IllegalWhereBlockError
	}
	if a.spec.all || a.spec.where != nil || len(a.spec.filter) > 0 {
		return errors.ConflictWhereProducerError
	}
	// The loopstate type holds the state of a loop over a struct's fields.
	type loopstate struct {
		wb  *whereblock
		ns  *typesutil.NamedStruct // the struct type of the whereblock
		idx int                    // keeps track of the field index
	}

	wb := new(whereblock)
	wb.name = field.Name()
	ns, err := typesutil.GetStruct(field)
	if err != nil { // fails only if non struct
		return errors.BadWhereBlockTypeError
	}

	// LIFO stack of states used for depth first traversal of struct fields.
	stack := []*loopstate{{wb: wb, ns: ns}}

stackloop:
	for len(stack) > 0 {
		loop := stack[len(stack)-1]
		for loop.idx < loop.ns.Struct.NumFields() {
			fld := loop.ns.Struct.Field(loop.idx)
			tag := tagutil.New(loop.ns.Struct.Tag(loop.idx))
			sqltag := tag.First("sql")

			// Instead of incrementing the index in the for-statement
			// it is done here manually to ensure that it is not skipped
			// when continuing to the outer loop.
			loop.idx++

			if sqltag == "-" || sqltag == "" {
				continue
			}

			// Skip the field if it's unexported and the ns.Struct's
			// type is imported. Unless it is one of the directive
			// fields that do not require direct access at runtime.
			if fld.Name() != "_" && !a.isaccessible(fld, ns.Named) {
				continue
			}

			item := new(whereitem)
			loop.wb.items = append(loop.wb.items, item)

			// Analyze the bool operation for any but the first
			// item in a whereblock. Fail if a value was provided
			// but it is not "or" nor "and".
			if len(loop.wb.items) > 1 {
				item.op = booland // default to "and"
				if booltag := tag.First("bool"); len(booltag) > 0 {
					v := strings.ToLower(booltag)
					if v == "or" {
						item.op = boolor
					} else if v != "and" {
						return errors.BadBoolTagValueError
					}
				}
			}

			// Nested whereblocks are marked with ">" and should be
			// analyzed before any other fields in the current block.
			if sqltag == ">" {
				ns, err := typesutil.GetStruct(fld)
				if err != nil {
					return errors.BadWhereBlockTypeError
				}

				wb := new(whereblock)
				wb.name = fld.Name()
				item.node = wb

				loop2 := new(loopstate)
				loop2.wb = wb
				loop2.ns = ns
				stack = append(stack, loop2)
				continue stackloop
			}

			lhs, op, op2, rhs := a.splitcmpexpr(sqltag)

			// Analyze directive where item.
			if fld.Name() == "_" {
				if !typesutil.IsDirective("Column", fld.Type()) {
					continue
				}

				// If the expression in a gosql.Column tag's value
				// contains a right-hand-side, it is expected to be
				// either another column or a value-literal to which
				// the main column should be compared.
				if len(rhs) > 0 {
					colid, err := a.colid(lhs, fld)
					if err != nil {
						return err
					}

					wn := new(wherecolumn)
					wn.colid = colid
					wn.cmp = string2cmpop[op]
					wn.sop = string2scalarrop[op2]

					if a.iscolid(rhs) {
						wn.colid2, _ = a.colid(rhs, fld) // ignore error since iscolid returned true
					} else {
						wn.lit = rhs // assume literal expression
					}

					if wn.cmp.isunary() {
						// TODO add test
						return errors.IllegalUnaryComparisonOperatorError
					} else if wn.sop > 0 && !wn.cmp.canusescalar() {
						return errors.BadCmpopComboError
					}

					item.node = wn
					continue
				}

				// Assume column with unary predicate.
				colid, err := a.colid(lhs, fld)
				if err != nil {
					return err
				}

				// If no operator was provided, default to "istrue"
				if len(op) == 0 {
					op = "istrue"
				}
				cmp := string2cmpop[op]
				if !cmp.isunary() {
					return errors.BadUnaryCmpopError
				}
				if len(op2) > 0 {
					return errors.ExtraScalarropError
				}

				wn := new(wherecolumn)
				wn.colid = colid
				wn.cmp = cmp
				item.node = wn
				continue
			}

			// Check whether the field is supposed to be used to
			// produce a [NOT] BETWEEN [SYMMETRIC] predicate clause.
			//
			// A valid "between" field MUST be of type struct with
			// the number of fields equal to 2, where each of the
			// fields is marked with an "x" or a "y" in their `sql`
			// tag to indicate their position in the clause.
			if strings.Contains(op, "between") {
				if len(op2) > 0 {
					// TODO test
					return errors.ExtraScalarropError
				}
				ns, err := typesutil.GetStruct(fld)
				if err != nil {
					return errors.BadBetweenTypeError
				} else if ns.Struct.NumFields() != 2 {
					return errors.BadBetweenTypeError
				}

				var x, y interface{}
				for i := 0; i < 2; i++ {
					fld := ns.Struct.Field(i)
					tag := tagutil.New(ns.Struct.Tag(i))
					sqltag := tag.First("sql")

					if fld.Name() == "_" && typesutil.IsDirective("Column", fld.Type()) {
						sqltag2 := strings.ToLower(tag.Second("sql"))

						colid, err := a.colid(sqltag, fld)
						if err != nil {
							return err
						}
						if sqltag2 == "x" {
							x = colid
						} else if sqltag2 == "y" {
							y = colid
						}
						continue
					}

					if a.isaccessible(fld, ns.Named) {
						v := new(varinfo)
						v.name = fld.Name()
						v.typ, _ = a.typeinfo(fld.Type())

						if sqltag == "x" {
							x = v
						} else if sqltag == "y" {
							y = v
						}
					}
				}

				if x == nil || y == nil {
					return errors.NoBetweenXYArgsError
				}

				colid, err := a.colid(lhs, fld)
				if err != nil {
					return err
				}

				bw := new(wherebetween)
				bw.name = fld.Name()
				bw.colid = colid
				bw.cmp = string2cmpop[op]
				bw.x, bw.y = x, y
				item.node = bw
				continue
			}

			// Analyze field where item.
			colid, err := a.colid(lhs, fld)
			if err != nil {
				return err
			}

			// If no comparison operator was provided default to "="
			if len(op) == 0 {
				op = "="
			}
			cmp := string2cmpop[op]
			if cmp.isunary() {
				// TODO add test
				return errors.IllegalUnaryComparisonOperatorError
			}

			sop := string2scalarrop[op2]
			if sop > 0 && !cmp.canusescalar() {
				return errors.BadCmpopComboError
			}

			wf := new(wherefield)
			wf.name = fld.Name()
			wf.colid = colid
			wf.typ, _ = a.typeinfo(fld.Type())
			wf.cmp = cmp
			wf.sop = sop
			wf.modfunc = a.funcname(tag["sql"][1:])

			if wf.sop > 0 && wf.typ.kind != kindslice && wf.typ.kind != kindarray {
				return errors.BadScalarFieldTypeError
			}

			item.node = wf
		}
		stack = stack[:len(stack)-1]
	}

	a.spec.where = wb
	return nil
}

func (a *analyzer) joinblock(field *types.Var) (err error) {
	joinblockname := strings.ToLower(field.Name())
	if joinblockname == "join" && a.spec.kind != speckindSelect {
		return errors.IllegalJoinBlockError
	} else if joinblockname == "from" && a.spec.kind != speckindUpdate {
		return errors.IllegalFromBlockError
	} else if joinblockname == "using" && a.spec.kind != speckindDelete {
		return errors.IllegalUsingBlockError
	}

	join := new(joinblock)
	ns, err := typesutil.GetStruct(field)
	if err != nil {
		return errors.BadJoinBlockTypeError
	}

	for i := 0; i < ns.Struct.NumFields(); i++ {
		fld := ns.Struct.Field(i)
		tag := tagutil.New(ns.Struct.Tag(i))
		sqltag := tag.First("sql")

		if sqltag == "-" || sqltag == "" {
			continue
		}

		// In a joinblock all fields are expected to be directives
		// with the blank identifier as their name.
		if fld.Name() != "_" {
			continue
		}

		switch dirname := strings.ToLower(typesutil.GetDirectiveName(fld)); dirname {
		case "relation":
			if joinblockname != "from" && joinblockname != "using" {
				return errors.IllegalJoinBlockRelationDirectiveError
			} else if len(join.rel.name) > 0 {
				return errors.ConflictJoinBlockRelationDirectiveError
			}
			id, err := a.relid(sqltag, fld)
			if err != nil {
				return err
			}
			join.rel = id
		case "leftjoin", "rightjoin", "fulljoin", "crossjoin":
			id, err := a.relid(sqltag, fld)
			if err != nil {
				return err
			}

			var conds []*joincond
			for _, val := range tag["sql"][1:] {
				vals := strings.Split(val, ";")
				for i, val := range vals {

					cond := new(joincond)
					if len(conds) > 0 && i == 0 {
						cond.op = booland
					} else if len(conds) > 0 && i > 0 {
						cond.op = boolor
					}

					lhs, op, op2, rhs := a.splitcmpexpr(val)
					if cond.col1, err = a.colid(lhs, fld); err != nil {
						return err
					}

					// optional right-hand side
					if a.iscolid(rhs) {
						cond.col2, _ = a.colid(rhs, fld) // ignore error since iscolid returned true
					} else {
						cond.lit = rhs
					}

					cond.cmp = string2cmpop[op]
					cond.sop = string2scalarrop[op2]

					if len(rhs) > 0 {
						if cond.cmp.isunary() {
							// TODO add test
							return errors.IllegalUnaryComparisonOperatorError
						} else if cond.sop > 0 && !cond.cmp.canusescalar() {
							return errors.BadCmpopComboError
						}
					} else {
						if !cond.cmp.isunary() {
							return errors.BadUnaryCmpopError
						} else if len(op2) > 0 {
							return errors.ExtraScalarropError
						}
					}

					conds = append(conds, cond)
				}
			}

			item := new(joinitem)
			item.typ = string2jointype[dirname]
			item.rel = id
			item.conds = conds
			join.items = append(join.items, item)
		default:
			return errors.IllegalJoinBlockDirectiveError
		}

	}

	a.spec.join = join
	return nil
}

func (a *analyzer) onconflictblock(field *types.Var) (err error) {
	if a.spec.kind != speckindInsert {
		return errors.IllegalOnConflictBlockError
	}

	onc := new(onconflictblock)
	ns, err := typesutil.GetStruct(field)
	if err != nil {
		return errors.BadOnConflictBlockTypeError
	}

	for i := 0; i < ns.Struct.NumFields(); i++ {
		fld := ns.Struct.Field(i)
		tag := tagutil.New(ns.Struct.Tag(i))

		// In an onconflictblock all fields are expected to be directives
		// with the blank identifier as their name.
		if fld.Name() != "_" {
			continue
		}

		switch dirname := strings.ToLower(typesutil.GetDirectiveName(fld)); dirname {
		case "column":
			if len(onc.column) > 0 || len(onc.index) > 0 || len(onc.constraint) > 0 {
				return errors.ConflictOnConflictBlockTargetProducerError
			}
			list, err := a.collist(tag["sql"], fld)
			if err != nil {
				return err
			}
			onc.column = list.items
		case "index":
			if len(onc.column) > 0 || len(onc.index) > 0 || len(onc.constraint) > 0 {
				return errors.ConflictOnConflictBlockTargetProducerError
			}
			if onc.index = tag.First("sql"); !reident.MatchString(onc.index) {
				return errors.BadIndexIdentifierValueError
			}
		case "constraint":
			if len(onc.column) > 0 || len(onc.index) > 0 || len(onc.constraint) > 0 {
				return errors.ConflictOnConflictBlockTargetProducerError
			}
			if onc.constraint = tag.First("sql"); !reident.MatchString(onc.constraint) {
				return errors.BadConstraintIdentifierValueError
			}
		case "ignore":
			if onc.ignore || onc.update != nil {
				return errors.ConflictOnConflictBlockActionProducerError
			}
			onc.ignore = true
		case "update":
			if onc.ignore || onc.update != nil {
				return errors.ConflictOnConflictBlockActionProducerError
			}
			if onc.update, err = a.collist(tag["sql"], fld); err != nil {
				return err
			}
		default:
			return errors.IllegalOnConflictBlockDirectiveError
		}

	}

	if onc.update != nil && (len(onc.column) == 0 && len(onc.index) == 0 && len(onc.constraint) == 0) {
		return errors.NoOnConflictTargetError
	}

	a.spec.onconflict = onc
	return nil
}

// Parses the given string as a comparison expression and returns the
// individual elements of that expression. The expected format is:
// { column [ comparison-operator [ scalar-operator ] { column | literal } ] }
func (a *analyzer) splitcmpexpr(expr string) (lhs, cop, sop, rhs string) {
	expr = strings.TrimSpace(expr)

	for i := range expr {
		switch expr[i] {
		case '=': // =
			lhs, cop, rhs = expr[:i], expr[i:i+1], expr[i+1:]
		case '!': // !=, !~, !~*
			if len(expr[i:]) > 2 && (expr[i+1] == '~' && expr[i+2] == '*') {
				lhs, cop, rhs = expr[:i], expr[i:i+3], expr[i+3:]
			} else if len(expr[i:]) > 1 && (expr[i+1] == '=' || expr[i+1] == '~') {
				lhs, cop, rhs = expr[:i], expr[i:i+2], expr[i+2:]
			}
		case '<': // <, <=, <>
			if len(expr[i:]) > 1 && (expr[i+1] == '=' || expr[i+1] == '>') {
				lhs, cop, rhs = expr[:i], expr[i:i+2], expr[i+2:]
			} else {
				lhs, cop, rhs = expr[:i], expr[i:i+1], expr[i+1:]
			}
		case '>': // >, >=
			if len(expr[i:]) > 1 && expr[i+1] == '=' {
				lhs, cop, rhs = expr[:i], expr[i:i+2], expr[i+2:]
			} else {
				lhs, cop, rhs = expr[:i], expr[i:i+1], expr[i+1:]
			}
		case '~': // ~, ~*
			if len(expr[i:]) > 1 && expr[i+1] == '*' {
				lhs, cop, rhs = expr[:i], expr[i:i+2], expr[i+2:]
			} else {
				lhs, cop, rhs = expr[:i], expr[i:i+1], expr[i+1:]
			}
		case ' ':

			var (
				j     = i + 1
				x     = strings.ToLower(expr)
				pred1 string // 1st part of predicate (not | is)
				pred2 string // 2nd part of predicate (distinct | true | null | ...)
			)

			if n := len(x[j:]); n > 3 && x[j:j+3] == "not" {
				pred1, pred2 = x[j:j+3], x[j+3:]
			} else if n := len(x[j:]); n > 2 && x[j:j+2] == "is" {
				pred1, pred2 = x[j:j+2], x[j+2:]
			}

			if len(pred2) > 0 {
				for _, adj := range predicateadjectives {
					if pred2[0] != adj[0] {
						continue
					}
					if n := len(adj); len(pred2) >= n && pred2[:n] == adj && (len(pred2) == n || pred2[n] == ' ') {
						lhs = expr[:i]
						cop = pred1 + pred2[:n]
						rhs = expr[j+len(cop):]
						break
					}
				}
			}

			if len(cop) == 0 {
				continue
			}
		default:
			continue
		}

		break // if "continue" wasn't executed, exit the loop
	}

	lhs = strings.TrimSpace(lhs)
	cop = strings.TrimSpace(cop)
	rhs = strings.TrimSpace(rhs)

	if len(rhs) > 0 {
		x := strings.ToLower(rhs)
		switch x[0] {
		case 'a': // ANY or ALL
			n := len("any") // any and all have the same length so we test against both at the same time
			if len(x) >= n && (x[:n] == "any" || x[:n] == "all") && (len(x) == n || x[n] == ' ') {
				sop, rhs = x[:n], rhs[n:]
			}
		case 's': // SOME
			n := len("some")
			if len(x) >= n && x[:n] == "some" && (len(x) == n || x[n] == ' ') {
				sop, rhs = x[:n], rhs[n:]
			}
		}

		sop = strings.TrimSpace(sop)
		rhs = strings.TrimSpace(rhs)
	}

	if len(lhs) == 0 {
		// return expr, "=", "", "" // default
		return expr, "", "", "" // default
	}

	return lhs, cop, sop, rhs
}

func (a *analyzer) limitvar(field *types.Var, tag string) error {
	if a.spec.kind != speckindSelect {
		return errors.IllegalLimitFieldOrDirectiveError
	}
	if a.spec.limit != nil {
		return errors.ConflictLimitProducerError
	}

	limit := new(limitvar)
	if fname := field.Name(); fname != "_" {
		if !a.isint(field.Type()) {
			return errors.BadLimitTypeError
		}
		limit.field = fname
	} else if len(tag) == 0 {
		return errors.NoLimitDirectiveValueError
	}

	if len(tag) > 0 {
		u64, err := strconv.ParseUint(tag, 10, 64)
		if err != nil {
			return errors.BadLimitValueError
		}
		limit.value = u64
	}
	a.spec.limit = limit
	return nil
}

func (a *analyzer) offsetvar(field *types.Var, tag string) error {
	if a.spec.kind != speckindSelect {
		return errors.IllegalOffsetFieldOrDirectiveError
	}
	if a.spec.offset != nil {
		return errors.ConflictOffsetProducerError
	}

	offset := new(offsetvar)
	if fname := field.Name(); fname != "_" {
		if !a.isint(field.Type()) {
			return errors.BadOffsetTypeError
		}
		offset.field = fname
	} else if len(tag) == 0 {
		return errors.NoOffsetDirectiveValueError
	}

	if len(tag) > 0 {
		u64, err := strconv.ParseUint(tag, 10, 64)
		if err != nil {
			return errors.BadOffsetValueError
		}
		offset.value = u64
	}
	a.spec.offset = offset
	return nil
}

func (a *analyzer) orderbydir(tags []string, field *types.Var) (err error) {
	if a.spec.kind != speckindSelect {
		return errors.IllegalOrderByDirectiveError
	} else if len(tags) == 0 {
		return errors.EmptyOrderByListError
	}

	list := new(orderbylist)
	for _, val := range tags {
		val = strings.TrimSpace(val)
		if len(val) == 0 {
			continue
		}

		item := new(orderbyitem)
		if val[0] == '-' {
			item.dir = orderdesc
			val = val[1:]
		}
		if i := strings.Index(val, ":"); i > -1 {
			if val[i+1:] == "nullsfirst" {
				item.nulls = nullsfirst
			} else if val[i+1:] == "nullslast" {
				item.nulls = nullslast
			} else {
				return errors.BadNullsOrderOptionValueError
			}
			val = val[:i]
		}

		if item.col, err = a.colid(val, field); err != nil {
			return err
		}

		list.items = append(list.items, item)
	}

	a.spec.orderby = list
	return nil
}

func (a *analyzer) overridedir(tag string, field *types.Var) error {
	if a.spec.kind != speckindInsert {
		return errors.IllegalOverrideDirectiveError
	}

	val := strings.ToLower(strings.TrimSpace(tag))
	switch val {
	case "system":
		a.spec.override = overridingsystem
	case "user":
		a.spec.override = overridinguser
	default:
		return errors.BadOverrideKindValueError
	}
	return nil
}

func (a *analyzer) resultfield(field *types.Var) error {
	if a.spec.kind != speckindInsert && a.spec.kind != speckindUpdate && a.spec.kind != speckindDelete {
		return errors.IllegalResultFieldError
	}
	if a.spec.returning != nil || a.spec.result != nil || len(a.spec.rowsaffected) > 0 {
		return errors.ConflictResultProducerError
	}

	rel := new(relfield)
	rel.name = field.Name()
	if err := a.relrecordtype(rel, field); err != nil {
		return err
	}

	result := new(resultfield)
	result.name = rel.name
	result.rec = rel.rec
	a.spec.result = result
	return nil
}

func (a *analyzer) rowsaffected(field *types.Var) error {
	if a.spec.kind != speckindInsert && a.spec.kind != speckindUpdate && a.spec.kind != speckindDelete {
		return errors.IllegalRowsAffectedFieldError
	}
	if a.spec.returning != nil || a.spec.result != nil || len(a.spec.rowsaffected) > 0 {
		return errors.ConflictResultProducerError
	}

	if !a.isint(field.Type()) {
		return errors.BadRowsAffectedTypeError
	}
	a.spec.rowsaffected = field.Name()
	return nil
}

func (a *analyzer) textsearch(tag string, field *types.Var) error {
	if a.spec.kind != speckindFilter {
		return errors.IllegalTextSearchDirectiveError
	}

	val := strings.ToLower(strings.TrimSpace(tag))
	cid, err := a.colid(val, field)
	if err != nil {
		return err
	}

	a.spec.textsearch = &cid
	return nil
}

func (a *analyzer) funcname(tagvals []string) funcname {
	for _, v := range tagvals {
		if len(v) > 0 && v[0] == '@' {
			return funcname(strings.ToLower(v[1:]))
		}
	}
	return ""
}

// parses the given string and returns a relid, if the value's format is invalid
// an error will be returned instead. The additional field argument is used only
// for error reporting. The expected format is: "[qualifier.]name[:alias]".
func (a *analyzer) relid(val string, field *types.Var) (id relid, err error) {
	if !rerelid.MatchString(val) {
		return id, errors.BadRelIdError
	}
	if i := strings.LastIndexByte(val, '.'); i > -1 {
		id.qual = val[:i]
		val = val[i+1:]
	}
	if i := strings.LastIndexByte(val, ':'); i > -1 {
		id.alias = val[i+1:]
		val = val[:i]
	}
	id.name = val
	return id, nil
}

func (a *analyzer) iscolid(val string) bool {
	return recolid.MatchString(val) && !rereserved.MatchString(val)
}

// parses the given string and returns a colid, if the value's format is invalid
// an error will be returned instead. The additional field argument is used only
// for error reporting. The expected format is: "[qualifier.]name".
func (a *analyzer) colid(val string, field *types.Var) (id colid, err error) {
	if !a.iscolid(val) {
		return id, errors.BadColIdError
	}
	if i := strings.LastIndexByte(val, '.'); i > -1 {
		id.qual = val[:i]
		val = val[i+1:]
	}
	id.name = val
	return id, nil
}

func (a *analyzer) collist(tag []string, field *types.Var) (*collist, error) {
	if len(tag) == 0 {
		return nil, errors.EmptyColListError
	}

	cl := new(collist)
	if len(tag) == 1 && tag[0] == "*" {
		cl.all = true
		return cl, nil
	}

	cl.items = make([]colid, len(tag))
	for i, val := range tag {
		id, err := a.colid(val, field)
		if err != nil {
			return nil, err
		}
		cl.items[i] = id
	}
	return cl, nil
}

func (a *analyzer) isimported(named *types.Named) bool {
	return named != nil && named.Obj().Pkg().Path() != a.pkg
}

func (a *analyzer) isaccessible(fld *types.Var, named *types.Named) bool {
	return fld.Name() != "_" && (fld.Exported() || !a.isimported(named))
}

// iserrorhandler returns true if the given type implements the ErrorHandler interface.
func (a *analyzer) iserrorhandler(typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	return typesutil.ImplementsErrorHandler(named)
}

// isfilter returns true if the given type is the gosql.Filter type.
func (a *analyzer) isfilter(typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	name := named.Obj().Name()
	if name != "Filter" {
		return false
	}

	path := named.Obj().Pkg().Path()
	return strings.HasSuffix(path, "github.com/frk/gosql")
}

// isint returns true if the given type is one of the basic integer
// types, including the unsigned ones.
func (a *analyzer) isint(typ types.Type) bool {
	basic, ok := typ.(*types.Basic)
	if !ok {
		return false
	}
	kind := basic.Kind()
	return kind == types.Int ||
		kind == types.Int8 ||
		kind == types.Int16 ||
		kind == types.Int32 ||
		kind == types.Int64 ||
		kind == types.Uint ||
		kind == types.Uint8 ||
		kind == types.Uint16 ||
		kind == types.Uint32 ||
		kind == types.Uint64
}

// isbool returns true if the given type is the basic bool type.
func (a *analyzer) isbool(typ types.Type) bool {
	basic, ok := typ.(*types.Basic)
	if !ok {
		return false
	}
	return basic.Kind() == types.Bool
}

type speckind uint

const (
	speckindInsert speckind = iota + 1
	speckindUpdate
	speckindSelect
	speckindDelete
	speckindFilter
)

func (k speckind) String() string {
	switch k {
	case speckindInsert:
		return "Insert"
	case speckindUpdate:
		return "Update"
	case speckindSelect:
		return "Select"
	case speckindDelete:
		return "Delete"
	case speckindFilter:
		return "Filter"
	}
	return "<unknown speckind>"
}

type selectkind uint

const (
	selectfrom      selectkind = iota // the default
	selectcount                       // SELECT COUNT(1) ...
	selectexists                      // SELECT EXISTS( ... )
	selectnotexists                   // SELECT NOT EXISTS( ... )
)

type typespec struct {
	name string   // name of the target struct type
	kind speckind // the kind of the typespec
	// If the typespec is a Select spec this field indicates the
	// specific kind of the select.
	selkind selectkind

	rel        *relfield
	join       *joinblock
	where      *whereblock
	orderby    *orderbylist
	textsearch *colid
	onconflict *onconflictblock

	defaults     *collist
	force        *collist
	returning    *collist
	result       *resultfield
	rowsaffected string

	limit    *limitvar
	offset   *offsetvar
	override overridingkind

	// Indicates that the query to be generated should be executed
	// against all the rows of the relation.
	all bool
	// The name of the ErrorHandler field, if any.
	erh string
	// The name of the Filter type field, if any.
	filter string
}

type relid struct {
	qual  string
	name  string
	alias string
}

type colid struct {
	qual string
	name string
}

func (id colid) isempty() bool {
	return id == colid{}
}

type collist struct {
	all   bool
	items []colid
}

// relfield holds the information on a go struct type and on the
// db relation that's associated with that struct type.
type relfield struct {
	name  string // name of the field that references the relation of the typespec
	relid relid  // the relation identifier extracted from the tag
	isdir bool   // indicates that the gosql.Relation directive was used
	rec   recordtype
}

type resultfield struct {
	name string // name of the field that declares the result of the typespec
	rec  recordtype
}

// recordtype holds information on the type of record a typespec should read from,
// or write to, the associated database relation.
type recordtype struct {
	base      typeinfo // information on the record's base type
	ispointer bool     // indicates whether or not the base type's a pointer type
	isslice   bool     // indicates whether or not the base type's a slice type
	isarray   bool     // indicates whether or not the base type's an array type
	arraylen  int64    // if the base type's an array type, this field will hold the array's length
	// if set, indicates that the recordtype is handled by an iterator
	isiter bool
	// if set the value will hold the method name of the iterator interface
	itermethod string
	// indicates whether or not the type implements the afterscanner interface
	isafterscanner bool
	// fields will hold the info on the recordtype's fields
	fields []*fieldinfo
}

type typeinfo struct {
	name             string   // the name of a named type or empty string for unnamed types
	kind             typekind // the kind of the go type
	pkgpath          string   // the package import path
	pkgname          string   // the package's name
	pkglocal         string   // the local package name (including ".")
	isimported       bool     // indicates whether or not the package is imported
	isscanner        bool     // indicates whether or not the type implements the sql.Scanner interface
	isvaluer         bool     // indicates whether or not the type implements the driver.Valuer interface
	isjsmarshaler    bool     // indicates whether or not the type implements the json.Marshaler interface
	isjsunmarshaler  bool     // indicates whether or not the type implements the json.Unmarshaler interface
	isxmlmarshaler   bool     // indicates whether or not the type implements the xml.Marshaler interface
	isxmlunmarshaler bool     // indicates whether or not the type implements the xml.Unmarshaler interface
	istime           bool     // indicates whether or not the type is time.Time or a type that embeds time.Time
	isbyte           bool     // indicates whether or not the type is the "byte" alias type
	isrune           bool     // indicates whether or not the type is the "rune" alias type
	// if kind is map, key will hold the info on the map's key type
	key *typeinfo
	// if kind is map, elem will hold the info on the map's value type
	// if kind is ptr, elem will hold the info on pointed-to type
	// if kind is slice/array, elem will hold the info on slice/array element type
	elem *typeinfo
	// if kind is array, arraylen will hold the array's length
	arraylen int64
}

// string returns a textual representation of the type that t represents.
// If elideptr is true the "*" will be elided from the output.
func (t *typeinfo) string(elideptr bool) string {
	if t.istime {
		return "time.Time"
	}

	switch t.kind {
	case kindarray:
		return "[" + strconv.FormatInt(t.arraylen, 10) + "]" + t.elem.string(elideptr)
	case kindslice:
		return "[]" + t.elem.string(elideptr)
	case kindmap:
		return "map[" + t.key.string(elideptr) + "]" + t.elem.string(elideptr)
	case kindptr:
		if elideptr {
			return t.elem.string(elideptr)
		} else {
			return "*" + t.elem.string(elideptr)
		}
	case kinduint8:
		if t.isbyte {
			return "byte"
		}
		return "uint8"
	case kindint32:
		if t.isrune {
			return "rune"
		}
		return "int32"
	case kindstruct:
		if len(t.name) > 0 {
			return t.pkgname + "." + t.name
		}
		return "struct{}"
	case kindinterface:
		if len(t.name) > 0 {
			return t.pkgname + "." + t.name
		}
		return "interface{}"
	case kindchan:
		return "chan"
	case kindfunc:
		return "func()"
	default:
		// assume builtin basic
		return typekind2string[t.kind]
	}
	return "<unknown>"
}

// is returns true if t represents a type one of the given kinds or a pointer
// to one of the given kinds.
func (t *typeinfo) is(kk ...typekind) bool {
	for _, k := range kk {
		if t.kind == k || (t.kind == kindptr && t.elem.kind == k) {
			return true
		}
	}
	return false
}

// isslice returns true if t represents a slice type whose elem type is one of
// the given kinds.
func (t *typeinfo) isslice(kk ...typekind) bool {
	if t.kind == kindslice {
		for _, k := range kk {
			if t.elem.kind == k {
				return true
			}
		}
	}
	return false
}

// isslicen returns true if t represents an n dimensional slice type whose
// base elem type is one of the given kinds.
func (t *typeinfo) isslicen(n int, kk ...typekind) bool {
	for ; n > 0; n-- {
		if t.kind != kindslice {
			return false
		}
		t = t.elem
	}
	for _, k := range kk {
		if t.kind == k {
			return true
		}
	}
	return false
}

// isnamed returns true if t represents a named type, or a pointer to a named
// type, whose package path and type name match the given arguments.
func (t *typeinfo) isnamed(pkgpath, name string) bool {
	if t.kind == kindptr {
		t = t.elem
	}
	return t.pkgpath == pkgpath && t.name == name
}

// isnilable returns true if t represents a type that can be nil.
func (t *typeinfo) isnilable() bool {
	return t.is(kindptr, kindslice, kindarray, kindmap, kindinterface)
}

// indicates whether or not the MarshalJSON method can be called on the type.
func (t *typeinfo) canjsonmarshal() bool {
	return t.isjsmarshaler || (t.kind == kindptr && t.elem.isjsmarshaler)
}

// indicates whether or not the UnmarshalJSON method can be called on the type.
func (t *typeinfo) canjsonunmarshal() bool {
	return t.isjsunmarshaler || (t.kind == kindptr && t.elem.isjsunmarshaler)
}

// indicates whether or not the MarshalXML method can be called on the type.
func (t *typeinfo) canxmlmarshal() bool {
	return t.isxmlmarshaler || (t.kind == kindptr && t.elem.isxmlmarshaler)
}

// indicates whether or not the UnmarshalXML method can be called on the type.
func (t *typeinfo) canxmlunmarshal() bool {
	return t.isxmlunmarshaler || (t.kind == kindptr && t.elem.isxmlunmarshaler)
}

type fieldelem struct {
	name         string
	tag          tagutil.Tag
	typename     string // the name of a named type or empty string for unnamed types
	typepkgpath  string // the package import path
	typepkgname  string // the package's name
	typepkglocal string // the local package name (including ".")
	isimported   bool   // indicates whether or not the type is imported
	isembedded   bool   // indicates whether or not the field is embedded
	isexported   bool   // indicates whether or not the field is exported
	ispointer    bool   // indicates whether or not the field type is a pointer type
}

// fieldinfo holds information about a recordtype's field and the corresponding db column.
type fieldinfo struct {
	typ  typeinfo // info about the field's type
	name string   // name of the struct field
	// if the field is nested, path will hold the parent fields' information
	path []*fieldelem
	// indicates whether or not the field is embedded
	isembedded bool
	// indicates whether or not the field is exported
	isexported bool
	// the field's parsed tag
	tag tagutil.Tag
	// the id of the corresponding column
	colid colid
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
	// indicates that the column value should be marshaled/unmarshaled
	// to/from xml before/after being stored/retrieved.
	usexml bool
	// if set to true it indicates that the column value should be wrapped
	// in a COALESCE call when read from the db.
	usecoalesce bool
	coalesceval string
	// for UPDATEs, if set to true, it indicates that the provided field
	// value should be added to the already existing column value.
	binadd bool
	// indicates whether or not an implicit CAST should be allowed.
	cancast bool
}

type joinblock struct {
	// The relid of the top relation in a DELETE-USING / UPDATE-FROM
	// clause, empty in SELECT commands.
	rel   relid
	items []*joinitem
}

type joinitem struct {
	typ   jointype
	rel   relid
	conds []*joincond
}

type joincond struct {
	op   boolop
	col1 colid  // the target column of the join condition
	col2 colid  // the optional 2nd column to be compared to col1
	lit  string // the optional literal value
	cmp  cmpop  // the comparison operator of the join condition
	sop  scalarrop
}

type varinfo struct {
	name string
	typ  typeinfo
}

type whereblock struct {
	name  string
	items []*whereitem
}

type whereitem struct {
	op   boolop
	node interface{} // wherefield, wherecolumn, wherebetween, or whereblock
}

type wherefield struct {
	name  string
	typ   typeinfo //
	colid colid    //
	cmp   cmpop    //
	sop   scalarrop
	// The name of the function to be used to modify the comparison
	// operands' values before comparing them.
	modfunc funcname
}

// wherecolumn is produced from a gosql.Column directive and its tag value.
// wherecolumn represents either a column with a unary comparison predicate,
// a column-to-column comparison, or a column-to-literal comparison.
type wherecolumn struct {
	// The target column of the wherecolumn item.
	colid colid
	// If set, it will hold the id of the column that should be compared
	// to the target column.
	colid2 colid
	// If set, it will hold the literal value that should be compared
	// to the target column.
	lit string
	// If set, it will hold the comparison operator to be used to compare
	// the target column either using a predicate unary operator, or a binary
	// operator comparing against the colid2 column or the lit value.
	cmp cmpop
	sop scalarrop
}

type wherebetween struct {
	name  string
	colid colid
	cmp   cmpop
	x, y  interface{}
}

type onconflictblock struct {
	column     []colid
	index      string
	constraint string
	ignore     bool
	update     *collist
}

// The limitvar is produced from either a Limit directive or from a valid "limit"
// field, it is then, in turn, used to produce a LIMIT clause in a SELECT query.
type limitvar struct {
	// The value provided in the Limit field's / directive's `sql` tag.
	// If the limitvar was produced from a directive the value is used as
	// a constant, but if limitvar was instead produced from a field the
	// value will only be used if the field's actual value is empty during
	// the query's execution, essentially acting as a default fallback.
	value uint64
	// The name of the Limit field, if empty it indicates that the limitvar
	// was produced from the Limit directive.
	field string
}

// The offsetvar is produced from either an Offset directive or from a valid "offset"
// field, it is then, in turn, used to produce an OFFSET clause in a SELECT query.
type offsetvar struct {
	// The value provided in the Offset field's / directive's `sql` tag.
	// If the offsetvar was produced from a directive the value is used as
	// a constant, but if offsetvar was instead produced from a field the
	// value will only be used if the field's actual value is empty during
	// the query's execution, essentially acting as a default fallback.
	value uint64
	// The name of the Offset field, if empty it indicates that the offsetvar
	// was produced from the Offset directive.
	field string
}

type orderbylist struct {
	items []*orderbyitem
}

type orderbyitem struct {
	col   colid
	dir   orderdirection
	nulls nullsposition
}

type boolop uint // boolop operation

const (
	_       boolop = iota // no bool
	booland               // conjunction
	boolor                // disjunction
	boolnot               // negation
)

type cmpop uint // comparison operation

const (
	_ cmpop = iota // no comparison

	// binary comparison operators
	cmpeq  // equals
	cmpne  // not equals
	cmpne2 // not equals
	cmplt  // less than
	cmpgt  // greater than
	cmple  // less than or equal
	cmpge  // greater than or equal

	// binary comparison predicates
	cmpisdistinct  // IS DISTINCT FROM
	cmpnotdistinct // IS NOT DISTINCT FROM

	// pattern matching operators
	cmprexp       // match regular expression
	cmprexpi      // match regular expression (case insensitive)
	cmpnotrexp    // not match regular expression
	cmpnotrexpi   // not match regular expression (case insensitive)
	cmpislike     // LIKE
	cmpnotlike    // NOT LIKE
	cmpisilike    // ILIKE
	cmpnotilike   // NOT ILIKE
	cmpissimilar  // IS SIMILAR TO
	cmpnotsimilar // IS NOT SIMILAR TO

	// array comparison operators
	cmpisin  // IN
	cmpnotin // NOT IN

	// range comparison operators
	cmpisbetween     // BETWEEN x AND y
	cmpnotbetween    // NOT BETWEEN x AND y
	cmpisbetweensym  // BETWEEN SYMMETRIC x AND y
	cmpnotbetweensym // NOT BETWEEN SYMMETRIC x AND y

	// unary comparison predicates
	cmpisnull     // IS NULL
	cmpnotnull    // IS NOT NULL
	cmpistrue     // IS TRUE
	cmpnottrue    // IS NOT TRUE
	cmpisfalse    // IS FALSE
	cmpnotfalse   // IS NOT FALSE
	cmpisunknown  // IS UNKNOWN
	cmpnotunknown // IS NOT UNKNOWN
)

func (op cmpop) is(oo ...cmpop) bool {
	for _, o := range oo {
		if op == o {
			return true
		}
	}
	return false
}

// canusescalar returns true if op can be used together with a scalar array operator.
func (op cmpop) canusescalar() bool {
	return op.isbinary() || op.ispatmatch()
}

// isbinbasic returns true if op is a basic binary comparison operator/predicate.
func (op cmpop) isbinary() bool {
	return op.is(cmpeq, cmpne, cmpne2, cmplt, cmpgt, cmple, cmpge,
		cmpisdistinct, cmpnotdistinct)
}

// ispatmatch returns true if op is a pattern matching comparison operator/predicate.
func (op cmpop) ispatmatch() bool {
	return op.is(cmprexp, cmprexpi, cmpnotrexp, cmpnotrexpi, cmpislike,
		cmpnotlike, cmpisilike, cmpnotilike, cmpissimilar, cmpnotsimilar)
}

// isrange returns true if op is a range-specific comparison predicate.
func (op cmpop) isrange() bool {
	return op.is(cmpisbetween, cmpnotbetween, cmpisbetweensym, cmpnotbetweensym)
}

// isunary returns true if op is a "unary" comparison predicate.
func (op cmpop) isunary() bool {
	return op.is(cmpisnull, cmpnotnull, cmpistrue, cmpnottrue,
		cmpisfalse, cmpnotfalse, cmpisunknown, cmpnotunknown)
}

// isbool returns true if op is one of the "unary" comparison predicates
// that requires a boolean expression.
func (op cmpop) isbool() bool {
	return op.is(cmpistrue, cmpnottrue, cmpisfalse, cmpnotfalse,
		cmpisunknown, cmpnotunknown)
}

// isarr returns true if op is a array comparison predicate.
func (op cmpop) isarr() bool {
	return op.is(cmpisin, cmpnotin)
}

var string2cmpop = map[string]cmpop{
	"=":  cmpeq,
	"<>": cmpne,
	"!=": cmpne2,
	"<":  cmplt,
	">":  cmpgt,
	"<=": cmple,
	">=": cmpge,

	"isdistinct":  cmpisdistinct,
	"notdistinct": cmpnotdistinct,

	"~":          cmprexp,
	"~*":         cmprexpi,
	"!~":         cmpnotrexp,
	"!~*":        cmpnotrexpi,
	"islike":     cmpislike,
	"notlike":    cmpnotlike,
	"isilike":    cmpisilike,
	"notilike":   cmpnotilike,
	"issimilar":  cmpissimilar,
	"notsimilar": cmpnotsimilar,

	"isin":  cmpisin,
	"notin": cmpnotin,

	"isbetween":     cmpisbetween,
	"notbetween":    cmpnotbetween,
	"isbetweensym":  cmpisbetweensym,
	"notbetweensym": cmpnotbetweensym,

	"isnull":     cmpisnull,
	"notnull":    cmpnotnull,
	"istrue":     cmpistrue,
	"nottrue":    cmpnottrue,
	"isfalse":    cmpisfalse,
	"notfalse":   cmpnotfalse,
	"isunknown":  cmpisunknown,
	"notunknown": cmpnotunknown,
}

var cmpop2sql = map[cmpop]string{
	cmpeq:            "=",
	cmpne:            "<>",
	cmpne2:           "!=",
	cmplt:            "<",
	cmpgt:            ">",
	cmple:            "<=",
	cmpge:            ">=",
	cmpisdistinct:    "IS DISTINCT FROM",
	cmpnotdistinct:   "IS NOT DISTINCT FROM",
	cmprexp:          "~",
	cmprexpi:         "~*",
	cmpnotrexp:       "!~",
	cmpnotrexpi:      "!~*",
	cmpislike:        "LIKE",
	cmpnotlike:       "NOT LIKE",
	cmpisilike:       "ILIKE",
	cmpnotilike:      "NOT ILIKE",
	cmpissimilar:     "SIMILAR TO",
	cmpnotsimilar:    "NOT SIMILAR TO",
	cmpisin:          "IN",
	cmpnotin:         "NOT IN",
	cmpisbetween:     "BETWEEN",
	cmpnotbetween:    "NOT BETWEEN",
	cmpisbetweensym:  "BETWEEN SYMMETRIC",
	cmpnotbetweensym: "NOT BETWEEN SYMMETRIC",
	cmpisnull:        "IS NULL",
	cmpnotnull:       "IS NOT NULL",
	cmpistrue:        "IS TRUE",
	cmpnottrue:       "IS NOT TRUE",
	cmpisfalse:       "IS FALSE",
	cmpnotfalse:      "IS NOT FALSE",
	cmpisunknown:     "IS UNKNOWN",
	cmpnotunknown:    "IS NOT UNKNOWN",
}

var predicateadjectives = []string{ // and adverbs
	"between",
	"betweensym",
	"distinct",
	"false",
	"ilike",
	"in",
	"like",
	"null",
	"similar",
	"true",
	"unknown",
}

type scalarrop uint8 // scalar array operator

const (
	_           scalarrop = iota // no operator
	scalarrany                   // ANY
	scalarrsome                  // SOME
	scalarrall                   // ALL
)

var string2scalarrop = map[string]scalarrop{
	"any":  scalarrany,
	"some": scalarrsome,
	"all":  scalarrall,
}

type orderdirection uint8

const (
	orderasc  orderdirection = iota // ASC, default
	orderdesc                       // DESC
)

type nullsposition uint8

const (
	_          nullsposition = iota // none specified, i.e. default
	nullsfirst                      // NULLS FIRST
	nullslast                       // NULLS LAST
)

type overridingkind uint8

const (
	_                overridingkind = iota // no overriding
	overridingsystem                       // OVERRIDING SYSTEM VALUE
	overridinguser                         // OVERRIDING USER VALUE
)

// funcname is the name of a database function that can either be used to modify
// a value, like lower, upper, etc. or a function that can be used as an aggregate.
type funcname string

type jointype uint

const (
	_         jointype = iota // no join
	joinleft                  // LEFT JOIN
	joinright                 // RIGHT JOIN
	joinfull                  // FULL JOIN
	joincross                 // CROSS JOIN
)

var string2jointype = map[string]jointype{
	"leftjoin":  joinleft,
	"rightjoin": joinright,
	"fulljoin":  joinfull,
	"crossjoin": joincross,
}

type typekind uint

const (
	// basic
	kindinvalid typekind = iota
	kindbool
	kindint
	kindint8
	kindint16
	kindint32
	kindint64
	kinduint
	kinduint8
	kinduint16
	kinduint32
	kinduint64
	kinduintptr
	kindfloat32
	kindfloat64
	kindcomplex64
	kindcomplex128
	kindstring
	kindunsafeptr

	// non-basic
	kindarray
	kindinterface
	kindmap
	kindptr
	kindslice
	kindstruct
	kindchan
	kindfunc

	// alisases
	kindbyte = kinduint8
	kindrune = kindint32
)

func (k typekind) String() string {
	if s, ok := typekind2string[k]; ok {
		return s
	}
	return "<invalid>"
}

var basickind2typekind = map[types.BasicKind]typekind{
	types.Invalid:       kindinvalid,
	types.Bool:          kindbool,
	types.Int:           kindint,
	types.Int8:          kindint8,
	types.Int16:         kindint16,
	types.Int32:         kindint32,
	types.Int64:         kindint64,
	types.Uint:          kinduint,
	types.Uint8:         kinduint8,
	types.Uint16:        kinduint16,
	types.Uint32:        kinduint32,
	types.Uint64:        kinduint64,
	types.Uintptr:       kinduintptr,
	types.Float32:       kindfloat32,
	types.Float64:       kindfloat64,
	types.Complex64:     kindcomplex64,
	types.Complex128:    kindcomplex128,
	types.String:        kindstring,
	types.UnsafePointer: kindunsafeptr,
}

var typekind2string = map[typekind]string{
	// builtin basic
	kindbool:       "bool",
	kindint:        "int",
	kindint8:       "int8",
	kindint16:      "int16",
	kindint32:      "int32",
	kindint64:      "int64",
	kinduint:       "uint",
	kinduint8:      "uint8",
	kinduint16:     "uint16",
	kinduint32:     "uint32",
	kinduint64:     "uint64",
	kinduintptr:    "uintptr",
	kindfloat32:    "float32",
	kindfloat64:    "float64",
	kindcomplex64:  "complex64",
	kindcomplex128: "complex128",
	kindstring:     "string",

	// non-basic
	kindarray:     "<array>",
	kindchan:      "<chan>",
	kindfunc:      "<func>",
	kindinterface: "<interface>",
	kindmap:       "<map>",
	kindptr:       "<pointer>",
	kindslice:     "<slice>",
	kindstruct:    "<struct>",
}

const (
	gotypbool         = "bool"
	gotypbools        = "[]bool"
	gotypstring       = "string"
	gotypstrings      = "[]string"
	gotypstringss     = "[][]string"
	gotypstringm      = "map[string]string"
	gotypstringms     = "[]map[string]string"
	gotypbyte         = "byte"
	gotypbytes        = "[]byte"
	gotypbytess       = "[][]byte"
	gotypbytea16      = "[16]byte"
	gotypbytea16s     = "[][16]byte"
	gotyprune         = "rune"
	gotyprunes        = "[]rune"
	gotypruness       = "[][]rune"
	gotypint          = "int"
	gotypints         = "[]int"
	gotypinta2        = "[2]int"
	gotypinta2s       = "[][2]int"
	gotypint8         = "int8"
	gotypint8s        = "[]int8"
	gotypint8ss       = "[][]int8"
	gotypint16        = "int16"
	gotypint16s       = "[]int16"
	gotypint16ss      = "[][]int16"
	gotypint32        = "int32"
	gotypint32s       = "[]int32"
	gotypint32a2      = "[2]int32"
	gotypint32a2s     = "[][2]int32"
	gotypint64        = "int64"
	gotypint64s       = "[]int64"
	gotypint64a2      = "[2]int64"
	gotypint64a2s     = "[][2]int64"
	gotypuint         = "uint"
	gotypuints        = "[]uint"
	gotypuint8        = "uint8"
	gotypuint8s       = "[]uint8"
	gotypuint16       = "uint16"
	gotypuint16s      = "[]uint16"
	gotypuint32       = "uint32"
	gotypuint32s      = "[]uint32"
	gotypuint64       = "uint64"
	gotypuint64s      = "[]uint64"
	gotypfloat32      = "float32"
	gotypfloat32s     = "[]float32"
	gotypfloat64      = "float64"
	gotypfloat64s     = "[]float64"
	gotypfloat64a2    = "[2]float64"
	gotypfloat64a2s   = "[][2]float64"
	gotypfloat64a2ss  = "[][][2]float64"
	gotypfloat64a2a2  = "[2][2]float64"
	gotypfloat64a2a2s = "[][2][2]float64"
	gotypfloat64a3    = "[3]float64"
	gotypfloat64a3s   = "[][3]float64"
	gotypipnet        = "net.IPNet"
	gotypipnets       = "[]net.IPNet"
	gotyptime         = "time.Time"
	gotyptimes        = "[]time.Time"
	gotyptimea2       = "[2]time.Time"
	gotyptimea2s      = "[][2]time.Time"
	gotypbigint       = "big.Int"
	gotypbigints      = "[]big.Int"
	gotypbiginta2     = "[2]big.Int"
	gotypbiginta2s    = "[][2]big.Int"
	gotypbigfloat     = "big.Float"
	gotypbigfloats    = "[]big.Float"
	gotypbigfloata2   = "[2]big.Float"
	gotypbigfloata2s  = "[][2]big.Float"
	gotypnullstringm  = "map[string]sql.NullString"
	gotypnullstringms = "[]map[string]sql.NullString"
)
package main

import (
	"go/types"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/frk/gosql/internal/errors"
	"github.com/frk/gosql/internal/typesutil"
	"github.com/frk/tagutil"
)

var (
	// NOTE(mkopriva): Identifiers MUST begin with a letter (a-z) or an underscore (_).
	// Subsequent characters in an identifier can be letters, underscores, and digits (0-9).

	// Matches a valid identifier.
	rxIdent = regexp.MustCompile(`^[A-Za-z_]\w*$`)

	// Matches a valid db relation identifier.
	// - Valid format: [schema_name.]relation_name[:alias_name]
	rxRelId = regexp.MustCompile(`^(?:[A-Za-z_]\w*\.)?[A-Za-z_]\w*(?:\:[A-Za-z_]\w*)?$`)

	// Matches a valid table column reference.
	// - Valid format: [rel_alias.]column_name
	rxColId = regexp.MustCompile(`^(?:[A-Za-z_]\w*\.)?[A-Za-z_]\w*$`)

	// Matches a few reserved identifiers.
	rxReserved = regexp.MustCompile(`^(?i:true|false|` +
		`current_date|current_time|current_timestamp|` +
		`current_role|current_schema|current_user|` +
		`localtime|localtimestamp|` +
		`session_user)$`)

	// Matches coalesce or coalesce(<value>) where <value> is expected to
	// be a single value literal.
	rxCoalesce = regexp.MustCompile(`(?i)^coalesce$|^coalesce\((.*)\)$`)
)

// TODO(mkopriva): to provide more detailed error messages either pass in the
// details about the file under analysis, or make sure that the caller has that
// information and appends it to the error.
//
// The analyze function runs the analysis on the given named type.
// The result of the analysis is stored in the targetInfo struct.
func analyze(named *types.Named, ti *targetInfo) error {
	structType, ok := named.Underlying().(*types.Struct)
	if !ok {
		panic(named.Obj().Name() + " must be a struct type.") // this shouldn't happen
	}

	name := named.Obj().Name()
	key := strings.ToLower(name)
	if len(key) > 5 {
		key = key[:6]
	}

	if key == "filter" {
		a := new(analyzer)
		a.pkgPath = named.Obj().Pkg().Path()
		a.named = named
		a.target = structType
		a.filter = new(filterStruct)
		a.filter.name = named.Obj().Name()
		if err := a.run(); err != nil {
			return err
		}
		ti.filter = a.filter
		ti.dataField = a.filter.dataField
		return nil
	}

	a := new(analyzer)
	a.pkgPath = named.Obj().Pkg().Path()
	a.named = named
	a.target = structType
	a.query = new(queryStruct)
	a.query.name = named.Obj().Name()

	switch key {
	case "insert":
		a.query.kind = queryKindInsert
	case "update":
		a.query.kind = queryKindUpdate
	case "select":
		a.query.kind = queryKindSelect
	case "delete":
		a.query.kind = queryKindDelete
	default:
		panic(a.query.name + " queryStruct kind has unsupported name prefix.") // this shouldn't happen
	}
	if err := a.run(); err != nil {
		return err
	}
	ti.query = a.query
	ti.dataField = a.query.dataField
	return nil
}

// analyzer holds the state of the analysis
type analyzer struct {
	target  *types.Struct // the types.Struct of the type under analysis
	named   *types.Named  // the types.Named of the type under analysis
	pkgPath string        // the package path of the type under analysis

	// the results
	query  *queryStruct
	filter *filterStruct
}

func (a *analyzer) run() (err error) {
	if a.query != nil {
		return a.queryStruct()
	}
	if a.filter != nil {
		return a.filterStruct()
	}

	panic("nothing to analyze")
	return nil
}

// queryStruct runs the analysis of a queryStruct.
func (a *analyzer) queryStruct() (err error) {
	// Used to track the presence of a field with a `rel` tag. Currently
	// only one "rel field" is allowed, if more than one are found an error
	// will be returned, regarless of whether the tag is empty or not.
	var hasRelTag bool

	for i := 0; i < a.target.NumFields(); i++ {
		fld := a.target.Field(i)
		tag := tagutil.New(a.target.Tag(i))

		// Ensure that there is only one field with the "rel" tag.
		if _, ok := tag["rel"]; ok {
			if hasRelTag {
				return errors.MultipleDataFieldsError
			}
			hasRelTag = true
		}

		if reltag := tag.First("rel"); len(reltag) > 0 {
			rid, err := a.relId(reltag, fld)
			if err != nil {
				return err
			}

			a.query.dataField = new(dataField)
			a.query.dataField.name = fld.Name()
			a.query.dataField.relId = rid

			switch fname := strings.ToLower(a.query.dataField.name); {
			case fname == "count" && isIntegerType(fld.Type()):
				if a.query.kind != queryKindSelect {
					return errors.IllegalCountFieldError

				}
				a.query.kind = queryKindSelectCount
			case fname == "exists" && isBoolType(fld.Type()):
				if a.query.kind != queryKindSelect {
					return errors.IllegalExistsFieldError
				}
				a.query.kind = queryKindSelectExists
			case fname == "notexists" && isBoolType(fld.Type()):
				if a.query.kind != queryKindSelect {
					return errors.IllegalNotExistsFieldError
				}
				a.query.kind = queryKindSelectNotExists
			case fname == "_" && typesutil.IsDirective("Relation", fld.Type()):
				if a.query.kind != queryKindDelete {
					return errors.IllegalRelationDirectiveError
				}
				a.query.dataField.isDirective = true
			default:
				if err := a.dataType(&a.query.dataField.data, fld); err != nil {
					return err
				}
				if (a.query.kind == queryKindInsert || a.query.kind == queryKindUpdate) && a.query.dataField.data.isIter {
					return errors.IllegalIteratorRecordError
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
				if a.query.kind != queryKindUpdate && a.query.kind != queryKindDelete {
					return errors.IllegalAllDirectiveError
				}
				if a.query.all || a.query.whereBlock != nil || len(a.query.filterField) > 0 {
					return errors.ConflictWhereProducerError
				}
				a.query.all = true
			case "default":
				if a.query.kind != queryKindInsert && a.query.kind != queryKindUpdate {
					return errors.IllegalDefaultDirectiveError
				}
				if a.query.defaultList, err = a.colIdList(tag["sql"], fld); err != nil {
					return err
				}
			case "force":
				if a.query.kind != queryKindInsert && a.query.kind != queryKindUpdate {
					return errors.IllegalForceDirectiveError
				}
				if a.query.forceList, err = a.colIdList(tag["sql"], fld); err != nil {
					return err
				}
			case "return":
				if len(a.query.dataField.data.fields) == 0 {
					// TODO test
					return errors.ReturnDirectiveWithNoDataFieldError
				}
				if a.query.kind != queryKindInsert && a.query.kind != queryKindUpdate && a.query.kind != queryKindDelete {
					return errors.IllegalReturnDirectiveError
				}
				if a.query.returnList != nil || a.query.resultField != nil || a.query.rowsAffectedField != nil {
					return errors.ConflictResultProducerError
				}
				if a.query.returnList, err = a.colIdList(tag["sql"], fld); err != nil {
					return err
				}
			case "limit":
				if err := a.limitField(fld, tag.First("sql")); err != nil {
					return err
				}
			case "offset":
				if err := a.offsetField(fld, tag.First("sql")); err != nil {
					return err
				}
			case "orderby":
				if err := a.orderByList(tag["sql"], fld); err != nil {
					return err
				}
			case "override":
				if err := a.overridingKind(tag.First("sql"), fld); err != nil {
					return err
				}
			case "textsearch":
				return errors.IllegalTextSearchDirectiveError
			default:
				return errors.IllegalCommandDirectiveError
			}
			continue
		}

		// fields with specific names
		switch fname := strings.ToLower(fld.Name()); fname {
		case "where":
			if err := a.whereBlock(fld); err != nil {
				return err
			}
		case "join", "from", "using":
			if err := a.joinBlock(fld); err != nil {
				return err
			}
		case "onconflict":
			if err := a.onConflictBlock(fld); err != nil {
				return err
			}
		case "result":
			if err := a.resultField(fld); err != nil {
				return err
			}
		case "limit":
			if err := a.limitField(fld, tag.First("sql")); err != nil {
				return err
			}
		case "offset":
			if err := a.offsetField(fld, tag.First("sql")); err != nil {
				return err
			}
		case "rowsaffected":
			if err := a.rowsAffectedField(fld); err != nil {
				return err
			}
		default:
			// if no match by field name, look for specific field types
			if a.isAccessible(fld, a.named) {
				switch {
				case isFilterType(fld.Type()):
					if !a.query.kind.isSelect() && a.query.kind != queryKindUpdate && a.query.kind != queryKindDelete {
						return errors.IllegalFilterFieldError
					}
					if a.query.all || a.query.whereBlock != nil || len(a.query.filterField) > 0 {
						return errors.ConflictWhereProducerError
					}
					a.query.filterField = fld.Name()
				case isErrorHandler(fld.Type()):
					if a.query.errorHandlerField != nil {
						return errors.ConflictErrorHandlerFieldError
					}
					a.query.errorHandlerField = new(errorHandlerField)
					a.query.errorHandlerField.name = fld.Name()
				case isErrorInfoHandler(fld.Type()):
					if a.query.errorHandlerField != nil {
						return errors.ConflictErrorHandlerFieldError
					}
					a.query.errorHandlerField = new(errorHandlerField)
					a.query.errorHandlerField.name = fld.Name()
					a.query.errorHandlerField.isInfo = true
				}
			}
		}

	}

	if a.query.dataField == nil {
		return errors.NoDataFieldError
	}

	if a.query.kind == queryKindUpdate && a.query.dataField.data.isSlice {
		// If this is an UPDATE with a slice of values, then only matching
		// by primary keys makes sense, therefore a whereBlock, or filter,
		// or the all directive, are disallowed.
		if a.query.whereBlock != nil || len(a.query.filterField) > 0 || a.query.all {
			return errors.IllegalSliceUpdateQueryModifierError
		}
	}

	// TODO if queryKind is Update and the record (single or slice) does not
	// have a primary key AND there's no whereBlock, no filter, no all directive
	// return an error. That case suggests that all records should be updated
	// however the all directive must be provided explicitly, as a way to
	// ensure the programmer does not, by mistake, declare a query that
	// updates all records in a table.

	return nil
}

// filterStruct runs the analysis of a filterStruct.
func (a *analyzer) filterStruct() (err error) {
	// Used to track the presence of a field with a `rel` tag. Currently
	// only one "rel field" is allowed, if more than one are found an error
	// will be returned, regarless of whether the tag is empty or not.
	var hasRelTag bool

	for i := 0; i < a.target.NumFields(); i++ {
		fld := a.target.Field(i)
		tag := tagutil.New(a.target.Tag(i))

		// Ensure that there is only one field with the "rel" tag.
		if _, ok := tag["rel"]; ok {
			if hasRelTag {
				return errors.MultipleDataFieldsError
			}
			hasRelTag = true
		}

		if reltag := tag.First("rel"); len(reltag) > 0 {
			rid, err := a.relId(reltag, fld)
			if err != nil {
				return err
			}

			a.filter.dataField = new(dataField)
			a.filter.dataField.name = fld.Name()
			a.filter.dataField.relId = rid

			if err := a.dataType(&a.filter.dataField.data, fld); err != nil {
				return err
			}
			if a.filter.dataField.data.isIter {
				return errors.IllegalIteratorRecordError
			}
			continue
		}

		// TODO(mkopriva): allow for embedding a struct with "common feature fields",
		// and make sure to also allow imported and local-unexported struct types.

		// fields with gosql directive types
		if dirname := typesutil.GetDirectiveName(fld); fld.Name() == "_" && len(dirname) > 0 {
			if strings.ToLower(dirname) == "textsearch" {
				if err := a.textSearch(tag.First("sql"), fld); err != nil {
					return err
				}
			} else {
				return errors.IllegalCommandDirectiveError
			}
			continue
		}
	}

	if a.filter.dataField == nil {
		return errors.NoDataFieldError
	}

	return nil
}

// dataType
func (a *analyzer) dataType(data *dataType, field *types.Var) error {
	var (
		ftyp  = field.Type()
		named *types.Named
		err   error
		ok    bool
	)
	if named, ok = ftyp.(*types.Named); ok {
		ftyp = named.Underlying()
	}

	// XXX Experimental: Not confident that "go/types.Type.String()" won't
	// produce conflicting values for different types.
	dataTypeKey := ftyp.String()
	dataTypeCache.RLock()
	v := dataTypeCache.m[dataTypeKey]
	dataTypeCache.RUnlock()
	if v != nil {
		*data = *v
		return nil
	}

	// Check whether the relation field's type is an interface or a function,
	// if so, it is then expected to be an iterator, and it is analyzed as such.
	//
	// Failure of the iterator analysis will cause the whole analysis to exit
	// as there's currently no support for non-iterator interfaces nor functions.
	if iface, ok := ftyp.(*types.Interface); ok {
		if named, err = a.iteratorInterface(data, iface, named); err != nil {
			return err
		}
	} else if sig, ok := ftyp.(*types.Signature); ok {
		if named, err = a.iteratorFunction(data, sig); err != nil {
			return err
		}
	} else {
		// If not an iterator, check for slices, arrays, and pointers.
		if slice, ok := ftyp.(*types.Slice); ok { // allows []T / []*T
			ftyp = slice.Elem()
			data.isSlice = true
		} else if array, ok := ftyp.(*types.Array); ok { // allows [N]T / [N]*T
			ftyp = array.Elem()
			data.isArray = true
			data.arrayLen = array.Len()
		}
		if ptr, ok := ftyp.(*types.Pointer); ok { // allows *T
			ftyp = ptr.Elem()
			data.isPointer = true
		}

		// Get the name of the base type, if applicable.
		if data.isSlice || data.isArray || data.isPointer {
			if named, ok = ftyp.(*types.Named); !ok {
				// Fail if the type is a slice, an array, or a pointer
				// while its base type remains unnamed.
				return errors.BadDataFieldTypeError
			}
		}
	}

	if named != nil {
		pkg := named.Obj().Pkg()
		data.typeInfo.name = named.Obj().Name()
		data.typeInfo.pkgPath = pkg.Path()
		data.typeInfo.pkgName = pkg.Name()
		data.typeInfo.pkgLocal = pkg.Name()
		data.typeInfo.isImported = a.isImportedType(named)
		data.isAfterScanner = typesutil.ImplementsAfterScanner(named)
		ftyp = named.Underlying()
	}

	data.typeInfo.kind = a.typeKind(ftyp)
	if data.typeInfo.kind != kindStruct {
		return errors.BadDataFieldTypeError
	}

	styp := ftyp.(*types.Struct)
	if err := a.fieldInfoList(data, styp); err != nil {
		return err
	}

	dataTypeCache.Lock()
	dataTypeCache.m[dataTypeKey] = data
	dataTypeCache.Unlock()
	return nil
}

// fieldInfoList
func (a *analyzer) fieldInfoList(data *dataType, styp *types.Struct) error {
	// The loopstate type holds the state of a loop over a struct's fields.
	type loopstate struct {
		styp *types.Struct // the struct type whose fields are being analyzed
		typ  *typeInfo     // info on the struct type; holds the resulting slice of analyzed fieldInfo
		idx  int           // keeps track of the field index
		pfx  string        // column prefix
		path []*fieldNode
	}

	// LIFO stack of states used for depth first traversal of struct fields.
	stack := []*loopstate{{styp: styp, typ: &data.typeInfo}}

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
				(!fld.Exported() && loop.typ.isImported) {
				continue
			}

			f := new(fieldInfo)
			f.tag = tag
			f.name = fld.Name()
			f.isEmbedded = fld.Embedded()
			f.isExported = fld.Exported()

			// Analyze the field's type.
			ftyp := fld.Type()
			f.typ, ftyp = a.typeInfo(ftyp)

			// If the field's type is a struct and the `sql` tag's
			// value starts with the ">" (descend) marker, then it is
			// considered to be a "parent" field element whose child
			// fields need to be analyzed as well.
			if f.typ.is(kindStruct) && strings.HasPrefix(sqltag, ">") {
				loop2 := new(loopstate)
				loop2.styp = ftyp.(*types.Struct)
				loop2.typ = &f.typ
				loop2.pfx = loop.pfx + strings.TrimPrefix(sqltag, ">")

				// Allocate path of the appropriate size an copy it.
				loop2.path = make([]*fieldNode, len(loop.path))
				_ = copy(loop2.path, loop.path)

				// If the parent node is a pointer to a struct,
				// get the struct type info.
				typ := f.typ
				if typ.kind == kindPtr {
					typ = *typ.elem
				}

				fe := new(fieldNode)
				fe.name = f.name
				fe.tag = f.tag
				fe.isEmbedded = f.isEmbedded
				fe.isExported = f.isExported
				fe.typeName = typ.name
				fe.typePkgPath = typ.pkgPath
				fe.typePkgName = typ.pkgName
				fe.typePkgLocal = typ.pkgLocal
				fe.isImported = typ.isImported
				fe.isPointer = (f.typ.kind == kindPtr)
				loop2.path = append(loop2.path, fe)

				stack = append(stack, loop2)
				continue stackloop
			}

			// If the field is not a struct to be descended,
			// it is considered to be a "leaf" field and as
			// such the analysis of leaf-specific information
			// needs to be carried out.
			f.path = loop.path
			f.isPKey = tag.HasOption("sql", "pk")
			f.nullEmpty = tag.HasOption("sql", "nullempty")
			f.readOnly = tag.HasOption("sql", "ro")
			f.writeOnly = tag.HasOption("sql", "wo")
			f.useJSON = tag.HasOption("sql", "json")
			f.useXML = tag.HasOption("sql", "xml")
			f.useAdd = tag.HasOption("sql", "add")
			f.canCast = tag.HasOption("sql", "cast")
			f.useDefault = tag.HasOption("sql", "default")
			f.useCoalesce, f.coalesceValue = a.coalesceInfo(tag)

			// Resolve the column id.
			colId, err := a.colId(loop.pfx+sqltag, fld)
			if err != nil {
				return err
			}
			f.colId = colId

			// Add the field to the list.
			data.fields = append(data.fields, f)
		}
		stack = stack[:len(stack)-1]
	}
	return nil
}

// The typeInfo method analyzes the given type and returns the result. The analysis
// looks only for information of "named types" and in case of slice, array, map,
// or pointer types it will analyze the element type of those types. The second
// return value is the base element type of the given type.
func (a *analyzer) typeInfo(tt types.Type) (typ typeInfo, base types.Type) {
	base = tt

	if named, ok := base.(*types.Named); ok {
		pkg := named.Obj().Pkg()
		typ.name = named.Obj().Name()
		typ.pkgPath = pkg.Path()
		typ.pkgName = pkg.Name()
		typ.pkgLocal = pkg.Name()
		typ.isImported = a.isImportedType(named)
		typ.isScanner = typesutil.ImplementsScanner(named)
		typ.isValuer = typesutil.ImplementsValuer(named)
		typ.isTime = typesutil.IsTime(named)
		typ.isJSONMarshaler = typesutil.ImplementsJSONMarshaler(named)
		typ.isJSONUnmarshaler = typesutil.ImplementsJSONUnmarshaler(named)
		typ.isXMLMarshaler = typesutil.ImplementsXMLMarshaler(named)
		typ.isXMLUnmarshaler = typesutil.ImplementsXMLUnmarshaler(named)
		base = named.Underlying()
	}

	typ.kind = a.typeKind(base)

	var elem typeInfo // element info
	switch T := base.(type) {
	case *types.Basic:
		typ.isRune = T.Name() == "rune"
		typ.isByte = T.Name() == "byte"
	case *types.Slice:
		elem, base = a.typeInfo(T.Elem())
		typ.elem = &elem
	case *types.Array:
		elem, base = a.typeInfo(T.Elem())
		typ.elem = &elem
		typ.arrayLen = T.Len()
	case *types.Map:
		key, _ := a.typeInfo(T.Key())
		elem, base = a.typeInfo(T.Elem())
		typ.key = &key
		typ.elem = &elem
	case *types.Pointer:
		elem, base = a.typeInfo(T.Elem())
		typ.elem = &elem
	case *types.Interface:
		// If base is an unnamed interface type check at least whether
		// or not it declares, or embeds, one of the relevant methods.
		if typ.name == "" {
			typ.isScanner = typesutil.IsScanner(T)
			typ.isValuer = typesutil.IsValuer(T)
			typ.isJSONMarshaler = typesutil.IsJSONMarshaler(T)
			typ.isJSONUnmarshaler = typesutil.IsJSONUnmarshaler(T)
			typ.isXMLMarshaler = typesutil.IsXMLMarshaler(T)
			typ.isXMLUnmarshaler = typesutil.IsXMLUnmarshaler(T)
		}
	}
	return typ, base
}

// iteratorInterface
func (a *analyzer) iteratorInterface(data *dataType, iface *types.Interface, named *types.Named) (*types.Named, error) {
	if iface.NumExplicitMethods() != 1 {
		return nil, errors.BadIteratorTypeError
	}

	mth := iface.ExplicitMethod(0)
	if !a.isAccessible(mth, named) {
		return nil, errors.BadIteratorTypeError
	}

	sig := mth.Type().(*types.Signature)
	named, err := a.iteratorFunction(data, sig)
	if err != nil {
		return nil, err
	}

	data.iterMethod = mth.Name()
	return named, nil
}

// iteratorFunction
func (a *analyzer) iteratorFunction(data *dataType, sig *types.Signature) (*types.Named, error) {
	// Must take 1 argument and return one value of type error. "func(T) error"
	if sig.Params().Len() != 1 || sig.Results().Len() != 1 || !typesutil.IsError(sig.Results().At(0).Type()) {
		return nil, errors.BadIteratorTypeError
	}

	typ := sig.Params().At(0).Type()
	if ptr, ok := typ.(*types.Pointer); ok { // allows *T
		typ = ptr.Elem()
		data.isPointer = true
	}

	// Make sure that the argument type is a named struct type.
	named, ok := typ.(*types.Named)
	if !ok {
		return nil, errors.BadIteratorTypeError
	} else if _, ok := named.Underlying().(*types.Struct); !ok {
		return nil, errors.BadIteratorTypeError
	}

	data.isIter = true
	return named, nil
}

func (a *analyzer) typeKind(typ types.Type) typeKind {
	switch x := typ.(type) {
	case *types.Basic:
		return basicKindToTypeKind[x.Kind()]
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

// coalesceInfo
func (a *analyzer) coalesceInfo(tag tagutil.Tag) (use bool, val string) {
	if sqltag := tag["sql"]; len(sqltag) > 0 {
		for _, opt := range sqltag[1:] {
			if strings.HasPrefix(opt, "coalesce") {
				use = true
				if match := rxCoalesce.FindStringSubmatch(opt); len(match) > 1 {
					val = match[1]
				}
				break
			}
		}
	}
	return use, val
}

// whereBlock
func (a *analyzer) whereBlock(field *types.Var) (err error) {
	if !a.query.kind.isSelect() && a.query.kind != queryKindUpdate && a.query.kind != queryKindDelete {
		return errors.IllegalWhereBlockError
	}
	if a.query.all || a.query.whereBlock != nil || len(a.query.filterField) > 0 {
		return errors.ConflictWhereProducerError
	}

	ns, err := typesutil.GetStruct(field)
	if err != nil { // fails only if non struct
		return errors.BadWhereBlockTypeError
	}

	// The loopstate type holds the state of a loop over a struct's fields.
	type loopstate struct {
		conds  []*searchCondition
		nested *searchConditionNested
		ns     *typesutil.NamedStruct // the struct type of the whereBlock
		idx    int                    // keeps track of the field index
	}

	// root holds the reference to the root level search conditions
	root := &loopstate{ns: ns}
	// LIFO stack of states used for depth first traversal of struct fields.
	stack := []*loopstate{root}

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
			if fld.Name() != "_" && !a.isAccessible(fld, ns.Named) {
				continue
			}

			item := new(searchCondition)
			loop.conds = append(loop.conds, item)

			// Analyze the bool operation for any but the first
			// item in a whereBlock. Fail if a value was provided
			// but it is not "or" nor "and".
			if len(loop.conds) > 1 {
				item.bool = boolAnd // default to "and"
				if booltag := tag.First("bool"); len(booltag) > 0 {
					v := strings.ToLower(booltag)
					if v == "or" {
						item.bool = boolOr
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

				cond := new(searchConditionNested)
				cond.name = fld.Name()
				item.cond = cond

				loop2 := new(loopstate)
				loop2.ns = ns
				loop2.nested = cond
				stack = append(stack, loop2)
				continue stackloop
			}

			lhs, op, op2, rhs := a.splitPredicateExpr(sqltag)

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
					colId, err := a.colId(lhs, fld)
					if err != nil {
						return err
					}

					cond := new(searchConditionColumn)
					cond.colId = colId
					cond.pred = stringToPredicate[op]
					cond.qua = stringToQuantifier[op2]

					if isColId(rhs) {
						cond.colId2, _ = a.colId(rhs, fld) // ignore error since isColId returned true
					} else {
						cond.literal = rhs // assume literal expression
					}

					if cond.pred.isUnary() {
						// TODO add test
						return errors.IllegalUnaryPredicateError
					} else if cond.qua > 0 && !cond.pred.canQuantify() {
						return errors.BadPredicateComboError
					}

					item.cond = cond
					continue
				}

				// Assume column with unary predicate.
				colId, err := a.colId(lhs, fld)
				if err != nil {
					return err
				}

				// If no operator was provided, default to "istrue"
				if len(op) == 0 {
					op = "istrue"
				}
				pred := stringToPredicate[op]
				if !pred.isUnary() {
					return errors.BadUnaryPredicateError
				}
				if len(op2) > 0 {
					return errors.ExtraQuantifierError
				}

				cond := new(searchConditionColumn)
				cond.colId = colId
				cond.pred = pred
				item.cond = cond
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
					return errors.ExtraQuantifierError
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

						colId, err := a.colId(sqltag, fld)
						if err != nil {
							return err
						}
						if sqltag2 == "x" {
							x = colId
						} else if sqltag2 == "y" {
							y = colId
						}
						continue
					}

					if a.isAccessible(fld, ns.Named) {
						v := new(fieldDatum)
						v.name = fld.Name()
						v.typ, _ = a.typeInfo(fld.Type())

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

				colId, err := a.colId(lhs, fld)
				if err != nil {
					return err
				}

				cond := new(searchConditionBetween)
				cond.name = fld.Name()
				cond.colId = colId
				cond.pred = stringToPredicate[op]
				cond.x = x
				cond.y = y
				item.cond = cond
				continue
			}

			// Analyze field where item.
			colId, err := a.colId(lhs, fld)
			if err != nil {
				return err
			}

			// If no predicate was provided default to "="
			if len(op) == 0 {
				op = "="
			}
			pred := stringToPredicate[op]
			if pred.isUnary() {
				// TODO add test
				return errors.IllegalUnaryPredicateError
			}

			qua := stringToQuantifier[op2]
			if qua > 0 && !pred.canQuantify() {
				return errors.BadPredicateComboError
			}

			cond := new(searchConditionField)
			cond.name = fld.Name()
			cond.typ, _ = a.typeInfo(fld.Type())
			cond.colId = colId
			cond.pred = pred
			cond.qua = qua
			cond.modFunc = a.funcName(tag["sql"][1:])

			if cond.qua > 0 && cond.typ.kind != kindSlice && cond.typ.kind != kindArray {
				return errors.BadQuantifierFieldTypeError
			}

			item.cond = cond
		}

		if loop.nested != nil {
			loop.nested.conds = loop.conds
		}

		stack = stack[:len(stack)-1]
	}

	wb := new(whereBlock)
	wb.name = field.Name()
	wb.conds = root.conds
	a.query.whereBlock = wb
	return nil
}

// joinBlock
func (a *analyzer) joinBlock(field *types.Var) (err error) {
	joinblockname := strings.ToLower(field.Name())
	if joinblockname == "join" && !a.query.kind.isSelect() {
		return errors.IllegalJoinBlockError
	} else if joinblockname == "from" && a.query.kind != queryKindUpdate {
		return errors.IllegalFromBlockError
	} else if joinblockname == "using" && a.query.kind != queryKindDelete {
		return errors.IllegalUsingBlockError
	}

	join := new(joinBlock)
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

		// In a joinBlock all fields are expected to be directives
		// with the blank identifier as their name.
		if fld.Name() != "_" {
			continue
		}

		switch dirname := strings.ToLower(typesutil.GetDirectiveName(fld)); dirname {
		case "relation":
			if joinblockname != "from" && joinblockname != "using" {
				return errors.IllegalJoinBlockRelationDirectiveError
			} else if len(join.relId.name) > 0 {
				return errors.ConflictJoinBlockRelationDirectiveError
			}
			rid, err := a.relId(sqltag, fld)
			if err != nil {
				return err
			}
			join.relId = rid
		case "leftjoin", "rightjoin", "fulljoin", "crossjoin":
			id, err := a.relId(sqltag, fld)
			if err != nil {
				return err
			}

			var conditions []*searchCondition
			for _, val := range tag["sql"][1:] {
				vals := strings.Split(val, ";")
				for i, val := range vals {

					cond := new(searchConditionColumn)
					lhs, op, op2, rhs := a.splitPredicateExpr(val)
					if cond.colId, err = a.colId(lhs, fld); err != nil {
						return err
					}

					// optional right-hand side
					if isColId(rhs) {
						cond.colId2, _ = a.colId(rhs, fld) // ignore error since isColId returned true
					} else {
						cond.literal = rhs
					}

					cond.pred = stringToPredicate[op]
					cond.qua = stringToQuantifier[op2]

					if len(rhs) > 0 {
						if cond.pred.isUnary() {
							// TODO add test
							return errors.IllegalUnaryPredicateError
						} else if cond.qua > 0 && !cond.pred.canQuantify() {
							return errors.BadPredicateComboError
						}
					} else {
						if !cond.pred.isUnary() {
							return errors.BadUnaryPredicateError
						} else if len(op2) > 0 {
							return errors.ExtraQuantifierError
						}
					}

					item := new(searchCondition)
					item.cond = cond
					if len(conditions) > 0 && i == 0 {
						item.bool = boolAnd
					} else if len(conditions) > 0 && i > 0 {
						item.bool = boolOr
					}

					conditions = append(conditions, item)
				}
			}

			item := new(joinItem)
			item.joinType = stringToJoinType[dirname]
			item.relId = id
			item.conds = conditions
			join.items = append(join.items, item)
		default:
			return errors.IllegalJoinBlockDirectiveError
		}

	}

	a.query.joinBlock = join
	return nil
}

// onConflictBlock
func (a *analyzer) onConflictBlock(field *types.Var) (err error) {
	if a.query.kind != queryKindInsert {
		return errors.IllegalOnConflictBlockError
	}

	onc := new(onConflictBlock)
	ns, err := typesutil.GetStruct(field)
	if err != nil {
		return errors.BadOnConflictBlockTypeError
	}

	for i := 0; i < ns.Struct.NumFields(); i++ {
		fld := ns.Struct.Field(i)
		tag := tagutil.New(ns.Struct.Tag(i))

		// In an onConflictBlock all fields are expected to be directives
		// with the blank identifier as their name.
		if fld.Name() != "_" {
			continue
		}

		switch dirname := strings.ToLower(typesutil.GetDirectiveName(fld)); dirname {
		case "column":
			if len(onc.column) > 0 || len(onc.index) > 0 || len(onc.constraint) > 0 {
				return errors.ConflictOnConflictBlockTargetProducerError
			}
			list, err := a.colIdList(tag["sql"], fld)
			if err != nil {
				return err
			}
			onc.column = list.items
		case "index":
			if len(onc.column) > 0 || len(onc.index) > 0 || len(onc.constraint) > 0 {
				return errors.ConflictOnConflictBlockTargetProducerError
			}
			if onc.index = tag.First("sql"); !rxIdent.MatchString(onc.index) {
				return errors.BadIndexIdentifierValueError
			}
		case "constraint":
			if len(onc.column) > 0 || len(onc.index) > 0 || len(onc.constraint) > 0 {
				return errors.ConflictOnConflictBlockTargetProducerError
			}
			if onc.constraint = tag.First("sql"); !rxIdent.MatchString(onc.constraint) {
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
			if onc.update, err = a.colIdList(tag["sql"], fld); err != nil {
				return err
			}
		default:
			return errors.IllegalOnConflictBlockDirectiveError
		}

	}

	if onc.update != nil && (len(onc.column) == 0 && len(onc.index) == 0 && len(onc.constraint) == 0) {
		return errors.NoOnConflictTargetError
	}

	a.query.onConflictBlock = onc
	return nil
}

// Parses the given string as a predicate expression and returns the individual
// elements of that expression. The expected format is:
// { column [ predicate-type [ quantifier ] { column | literal } ] }
func (a *analyzer) splitPredicateExpr(expr string) (lhs, cop, qua, rhs string) {
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
				for _, adj := range predicateAdjectives {
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
				qua, rhs = x[:n], rhs[n:]
			}
		case 's': // SOME
			n := len("some")
			if len(x) >= n && x[:n] == "some" && (len(x) == n || x[n] == ' ') {
				qua, rhs = x[:n], rhs[n:]
			}
		}

		qua = strings.TrimSpace(qua)
		rhs = strings.TrimSpace(rhs)
	}

	if len(lhs) == 0 {
		return expr, "", "", "" // default
	}

	return lhs, cop, qua, rhs
}

// limitField analyzes the given field, which is expected to be either
// the gosql.Limit directive or a plain integer field. The tag argument,
// if not empty, is expected to hold a positive integer.
func (a *analyzer) limitField(field *types.Var, tag string) error {
	if !a.query.kind.isSelect() {
		return errors.IllegalLimitFieldOrDirectiveError
	}
	if a.query.limitField != nil {
		return errors.ConflictLimitProducerError
	}

	f := new(limitField)
	if name := field.Name(); name != "_" {
		if !isIntegerType(field.Type()) {
			return errors.BadLimitTypeError
		}
		f.name = name
	} else if len(tag) == 0 {
		return errors.NoLimitDirectiveValueError
	}

	if len(tag) > 0 {
		u64, err := strconv.ParseUint(tag, 10, 64)
		if err != nil {
			return errors.BadLimitValueError
		}
		f.value = u64
	}
	a.query.limitField = f
	return nil
}

// offsetField analyzes the given field, which is expected to be either
// the gosql.Offset directive or a plain integer field. The tag argument,
// if not empty, is expected to hold a positive integer.
func (a *analyzer) offsetField(field *types.Var, tag string) error {
	if !a.query.kind.isSelect() {
		return errors.IllegalOffsetFieldOrDirectiveError
	}
	if a.query.offsetField != nil {
		return errors.ConflictOffsetProducerError
	}

	f := new(offsetField)
	if name := field.Name(); name != "_" {
		if !isIntegerType(field.Type()) {
			return errors.BadOffsetTypeError
		}
		f.name = name
	} else if len(tag) == 0 {
		return errors.NoOffsetDirectiveValueError
	}

	if len(tag) > 0 {
		u64, err := strconv.ParseUint(tag, 10, 64)
		if err != nil {
			return errors.BadOffsetValueError
		}
		f.value = u64
	}
	a.query.offsetField = f
	return nil
}

// orderByList
func (a *analyzer) orderByList(tags []string, field *types.Var) (err error) {
	if !a.query.kind.isSelect() {
		return errors.IllegalOrderByDirectiveError
	} else if len(tags) == 0 {
		return errors.EmptyOrderByListError
	}

	list := new(orderByList)
	for _, val := range tags {
		val = strings.TrimSpace(val)
		if len(val) == 0 {
			continue
		}

		item := new(orderByItem)
		if val[0] == '-' {
			item.direction = orderDesc
			val = val[1:]
		}
		if i := strings.Index(val, ":"); i > -1 {
			if val[i+1:] == "nullsfirst" {
				item.nulls = nullsFirst
			} else if val[i+1:] == "nullslast" {
				item.nulls = nullsLast
			} else {
				return errors.BadNullsOrderOptionValueError
			}
			val = val[:i]
		}

		if item.colId, err = a.colId(val, field); err != nil {
			return err
		}

		list.items = append(list.items, item)
	}

	a.query.orderByList = list
	return nil
}

func (a *analyzer) overridingKind(tag string, field *types.Var) error {
	if a.query.kind != queryKindInsert {
		return errors.IllegalOverrideDirectiveError
	}

	val := strings.ToLower(strings.TrimSpace(tag))
	switch val {
	case "system":
		a.query.overridingKind = overridingSystem
	case "user":
		a.query.overridingKind = overridingUser
	default:
		return errors.BadOverrideKindValueError
	}
	return nil
}

func (a *analyzer) resultField(field *types.Var) error {
	if a.query.kind != queryKindInsert && a.query.kind != queryKindUpdate && a.query.kind != queryKindDelete {
		return errors.IllegalResultFieldError
	}
	if a.query.returnList != nil || a.query.resultField != nil || a.query.rowsAffectedField != nil {
		return errors.ConflictResultProducerError
	}

	result := new(resultField)
	result.name = field.Name()
	if err := a.dataType(&result.data, field); err != nil {
		return err
	}

	a.query.resultField = result
	return nil
}

func (a *analyzer) rowsAffectedField(field *types.Var) error {
	if a.query.kind != queryKindInsert && a.query.kind != queryKindUpdate && a.query.kind != queryKindDelete {
		return errors.IllegalRowsAffectedFieldError
	}
	if a.query.returnList != nil || a.query.resultField != nil || a.query.rowsAffectedField != nil {
		return errors.ConflictResultProducerError
	}

	ftyp := field.Type()
	if !isIntegerType(ftyp) {
		return errors.BadRowsAffectedTypeError
	}

	a.query.rowsAffectedField = new(rowsAffectedField)
	a.query.rowsAffectedField.name = field.Name()
	a.query.rowsAffectedField.kind = a.typeKind(ftyp)
	return nil
}

func (a *analyzer) textSearch(tag string, field *types.Var) error {
	val := strings.ToLower(strings.TrimSpace(tag))
	cid, err := a.colId(val, field)
	if err != nil {
		return err
	}

	a.filter.textSearchColId = &cid
	return nil
}

func (a *analyzer) funcName(tagvals []string) funcName {
	for _, v := range tagvals {
		if len(v) > 0 && v[0] == '@' {
			return funcName(strings.ToLower(v[1:]))
		}
	}
	return ""
}

// parses the given string and returns a relId, if the value's format is invalid
// an error will be returned instead. The additional field argument is used only
// for error reporting. The expected format is: "[qualifier.]name[:alias]".
func (a *analyzer) relId(val string, field *types.Var) (id relId, err error) {
	if !rxRelId.MatchString(val) {
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

// parses the given string and returns a colId, if the value's format is invalid
// an error will be returned instead. The additional field argument is used only
// for error reporting. The expected format is: "[qualifier.]name".
func (a *analyzer) colId(val string, field *types.Var) (id colId, err error) {
	if !isColId(val) {
		return id, errors.BadColIdError
	}
	if i := strings.LastIndexByte(val, '.'); i > -1 {
		id.qual = val[:i]
		val = val[i+1:]
	}
	id.name = val
	return id, nil
}

// parses the given tag slice as a list of column identifiers, if any of the
// values in the slice is invalid then an error will be returned. The additional
// field argument is used only for error reporting.
func (a *analyzer) colIdList(tag []string, field *types.Var) (*colIdList, error) {
	if len(tag) == 0 {
		return nil, errors.EmptyColListError
	}

	cl := new(colIdList)
	if len(tag) == 1 && tag[0] == "*" {
		cl.all = true
		return cl, nil
	}

	cl.items = make([]colId, len(tag))
	for i, val := range tag {
		id, err := a.colId(val, field)
		if err != nil {
			return nil, err
		}
		cl.items[i] = id
	}
	return cl, nil
}

// isImportedType reports whether or not the given type is imported based on
// on the package in which the target of the analysis is declared.
func (a *analyzer) isImportedType(named *types.Named) bool {
	return named != nil && named.Obj().Pkg().Path() != a.pkgPath
}

// isAccessible reports whether or not the given value is accessible from
// the package in which the target of the analysis is declared.
func (a *analyzer) isAccessible(x exportable, named *types.Named) bool {
	return x.Name() != "_" && (x.Exported() || !a.isImportedType(named))
}

// exportable is implemented by both types.Var and types.Func.
type exportable interface {
	Name() string
	Exported() bool
}

type queryKind uint

const (
	queryKindInsert queryKind = iota + 1
	queryKindUpdate
	queryKindSelect
	queryKindSelectCount
	queryKindSelectExists
	queryKindSelectNotExists
	queryKindDelete
)

// isSelect reports whether or not the query kind is one of the select kinds.
func (k queryKind) isSelect() bool {
	return k == queryKindSelect ||
		k == queryKindSelectCount ||
		k == queryKindSelectExists ||
		k == queryKindSelectNotExists
}

// isNonFromSelect reports whether or not the query kind is one of the non-from-select select kinds.
func (k queryKind) isNonFromSelect() bool {
	return k == queryKindSelectCount ||
		k == queryKindSelectExists ||
		k == queryKindSelectNotExists
}

func (k queryKind) String() string {
	switch k {
	case queryKindInsert:
		return "Insert"
	case queryKindUpdate:
		return "Update"
	case queryKindSelect:
		return "Select"
	case queryKindDelete:
		return "Delete"
	}
	return "<unknown queryKind>"
}

type queryStruct struct {
	name string    // name of the struct type
	kind queryKind // the kind of the queryStruct

	dataField         *dataField
	resultField       *resultField
	joinBlock         *joinBlock
	whereBlock        *whereBlock
	onConflictBlock   *onConflictBlock
	orderByList       *orderByList
	limitField        *limitField
	offsetField       *offsetField
	rowsAffectedField *rowsAffectedField
	defaultList       *colIdList
	forceList         *colIdList
	returnList        *colIdList
	errorHandlerField *errorHandlerField
	overridingKind    overridingKind
	// The name of the Filter type field, if any.
	filterField string
	// Indicates that the query to be generated should be executed
	// against all the rows of the relation.
	all bool
}

type filterStruct struct {
	name            string // name of the struct type
	dataField       *dataField
	textSearchColId *colId
}

// dataField represents the field that holds the information about the target dataType.
type dataField struct {
	name        string // name of the field
	data        dataType
	relId       relId // the relation id as extracted from the field's tag
	isDirective bool  // indicates that the gosql.Relation directive was used
}

type resultField struct {
	name string // name of the field that declares the result of the queryStruct
	data dataType
}

type relId struct {
	qual  string
	name  string
	alias string
}

type colId struct {
	qual string
	name string
}

func (id colId) isEmpty() bool {
	return id == colId{}
}

func (id colId) string() string {
	if len(id.qual) > 0 {
		return id.qual + "." + id.name
	}
	return id.name
}

func (id colId) quoted() string {
	if len(id.qual) > 0 {
		return id.qual + `."` + id.name + `"`
	}
	return `"` + id.name + `"`
}

type colIdList struct {
	all   bool
	items []colId
}

func (cl *colIdList) contains(cid colId) bool {
	for i := 0; i < len(cl.items); i++ {
		if cl.items[i] == cid {
			return true
		}
	}
	return false
}

type rowsAffectedField struct {
	name string // name of the rowsAffectedField field
	kind typeKind
}

type errorHandlerField struct {
	// name of the error handler field
	name string
	// indicates whether or not the field's type implements the ErrorInfoHandler interface.
	isInfo bool
}

// dataType holds information on the type of record a queryStruct should read from,
// or write to, the associated database relation.
type dataType struct {
	typeInfo  typeInfo // information on the record's base type
	isPointer bool     // indicates whether or not the base type's a pointer type
	isSlice   bool     // indicates whether or not the base type's a slice type
	isArray   bool     // indicates whether or not the base type's an array type
	arrayLen  int64    // if the base type's an array type, this field will hold the array's length
	// if set, indicates that the dataType is handled by an iterator
	isIter bool
	// if set the value will hold the method name of the iterator interface
	iterMethod string
	// indicates whether or not the type implements the gosql.AfterScanner interface
	isAfterScanner bool
	// fields will hold the info on the dataType's fields
	fields []*fieldInfo
}

type typeInfo struct {
	name              string   // the name of a named type or empty string for unnamed types
	kind              typeKind // the kind of the go type
	pkgPath           string   // the package import path
	pkgName           string   // the package's name
	pkgLocal          string   // the local package name (including ".")
	isImported        bool     // indicates whether or not the package is imported
	isScanner         bool     // indicates whether or not the type implements the sql.Scanner interface
	isValuer          bool     // indicates whether or not the type implements the driver.Valuer interface
	isJSONMarshaler   bool     // indicates whether or not the type implements the json.Marshaler interface
	isJSONUnmarshaler bool     // indicates whether or not the type implements the json.Unmarshaler interface
	isXMLMarshaler    bool     // indicates whether or not the type implements the xml.Marshaler interface
	isXMLUnmarshaler  bool     // indicates whether or not the type implements the xml.Unmarshaler interface
	isTime            bool     // indicates whether or not the type is time.Time or a type that embeds time.Time
	isByte            bool     // indicates whether or not the type is the "byte" alias type
	isRune            bool     // indicates whether or not the type is the "rune" alias type
	// if kind is map, key will hold the info on the map's key type
	key *typeInfo
	// if kind is map, elem will hold the info on the map's value type
	// if kind is ptr, elem will hold the info on pointed-to type
	// if kind is slice/array, elem will hold the info on slice/array element type
	elem *typeInfo
	// if kind is array, arrayLen will hold the array's length
	arrayLen int64
}

// string returns a textual representation of the type that t represents.
// - If elideptr is true the "*" will be elided from the output.
func (t *typeInfo) string(elideptr bool) string {
	if t.isTime {
		return "time.Time"
	}

	switch t.kind {
	case kindArray:
		return "[" + strconv.FormatInt(t.arrayLen, 10) + "]" + t.elem.string(elideptr)
	case kindSlice:
		return "[]" + t.elem.string(elideptr)
	case kindMap:
		return "map[" + t.key.string(elideptr) + "]" + t.elem.string(elideptr)
	case kindPtr:
		if elideptr {
			return t.elem.string(elideptr)
		} else {
			return "*" + t.elem.string(elideptr)
		}
	case kindUint8:
		if t.isByte {
			return "byte"
		}
		return "uint8"
	case kindInt32:
		if t.isRune {
			return "rune"
		}
		return "int32"
	case kindStruct:
		if len(t.name) > 0 {
			return t.pkgName + "." + t.name
		}
		return "struct{}"
	case kindInterface:
		if len(t.name) > 0 {
			return t.pkgName + "." + t.name
		}
		return "interface{}"
	case kindChan:
		return "chan"
	case kindFunc:
		return "func()"
	default:
		// assume builtin basic
		return typeKindToString[t.kind]
	}
	return "<unknown>"
}

// nameOrLiteral builds and returns the type's name or literal if it's not a named type.
func (t *typeInfo) nameOrLiteral(pkglocal bool) string {
	if len(t.name) > 0 {
		if pkglocal && len(t.pkgLocal) > 0 && t.pkgLocal != "." {
			return t.pkgLocal + "." + t.name
		} else if len(t.pkgName) > 0 {
			return t.pkgName + "." + t.name
		}
		return t.name
	}

	switch t.kind {
	case kindArray:
		return "[" + strconv.FormatInt(t.arrayLen, 10) + "]" + t.elem.nameOrLiteral(pkglocal)
	case kindSlice:
		return "[]" + t.elem.nameOrLiteral(pkglocal)
	case kindMap:
		return "map[" + t.key.nameOrLiteral(pkglocal) + "]" + t.elem.nameOrLiteral(pkglocal)
	case kindPtr:
		return "*" + t.elem.nameOrLiteral(pkglocal)
	case kindUint8:
		if t.isByte {
			return "byte"
		}
		return "uint8"
	case kindInt32:
		if t.isRune {
			return "rune"
		}
		return "int32"
	case kindStruct, kindInterface, kindChan, kindFunc:
		return "<unsupported>"
	default:
		// assume builtin basic
		return typeKindToString[t.kind]
	}
	return "<unknown>"
}

// is reports whether or not t represents a type whose kind matches one of
// the provided typeKinds or a pointer to one of the provided typeKinds.
func (t *typeInfo) is(kk ...typeKind) bool {
	for _, k := range kk {
		if t.kind == k || (t.kind == kindPtr && t.elem.kind == k) {
			return true
		}
	}
	return false
}

// isSlice reports whether or not t represents a slice type whose elem type
// is one of the provided typeKinds.
func (t *typeInfo) isSlice(kk ...typeKind) bool {
	if t.kind == kindSlice {
		for _, k := range kk {
			if t.elem.kind == k {
				return true
			}
		}
	}
	return false
}

// isNilable reports whether or not t represents a type that can be nil.
func (t *typeInfo) isNilable() bool {
	return t.is(kindPtr, kindSlice, kindArray, kindMap, kindInterface)
}

// canJSONMarshal reports whether or not the MarshalJSON method can be invoked
// on an instance of the type represented by t.
func (t *typeInfo) canJSONMarshal() bool {
	return t.isJSONMarshaler || (t.kind == kindPtr && t.elem.isJSONMarshaler)
}

// canJSONUnmarshal reports whether or not the UnmarshalJSON method can
// be invoked on an instance of the type represented by t.
func (t *typeInfo) canJSONUnmarshal() bool {
	return t.isJSONUnmarshaler || (t.kind == kindPtr && t.elem.isJSONUnmarshaler)
}

// canXMLMarshal reports whether or not the MarshalXML method can be invoked
// on an instance of the type represented by t.
func (t *typeInfo) canXMLMarshal() bool {
	return t.isXMLMarshaler || (t.kind == kindPtr && t.elem.isXMLMarshaler)
}

// canXMLUnmarshal reports whether or not the UnmarshalXML method can be
// invoked on an instance of the type represented by t.
func (t *typeInfo) canXMLUnmarshal() bool {
	return t.isXMLUnmarshaler || (t.kind == kindPtr && t.elem.isXMLUnmarshaler)
}

// fieldInfo holds information about a dataType's field.
type fieldInfo struct {
	// Name of the struct field.
	name string
	// Info about the field's type.
	typ typeInfo
	// If the field is nested, path will hold the parent fields' information.
	path []*fieldNode
	// Indicates whether or not the field is embedded.
	isEmbedded bool
	// Indicates whether or not the field is exported.
	isExported bool
	// The field's parsed tag.
	tag tagutil.Tag
	// The identifier of the field's corresponding column.
	colId colId
	// Identifies the field's corresponding column as a primary key.
	isPKey bool
	// Indicates that if the field's value is EMPTY then NULL should
	// be stored in the column during INSERT/UPDATE.
	nullEmpty bool
	// Indicates that field should only be used to read from the database and
	// never to write to it. Can be overriden with the gosql.Force directive.
	readOnly bool
	// Indicates that field should only be used to write to the database and
	// never to read from it. Can be overriden with the gosql.Force directive.
	writeOnly bool
	// Indicates that the DEFAULT marker should be used during INSERT/UPDATE.
	useDefault bool
	// Indicates that the column value should be marshaled/unmarshaled
	// to/from JSON before/after being stored/retrieved.
	useJSON bool
	// Indicates that the column value should be marshaled/unmarshaled
	// to/from XML before/after being stored/retrieved.
	useXML bool
	// If set to true, it indicates that the provided field value should be
	// "added" to the already existing column value.
	// For UPDATEs only.
	useAdd bool
	// Indicates whether or not an implicit CAST should be allowed.
	canCast bool
	// If set to true, it indicates that the column expression should be
	// wrapped in a COALESCE call when read from the db.
	useCoalesce bool
	// If set, it will hold the value literal to be passed as the second
	// argument to the COALESCE call.
	coalesceValue string
}

// fieldNode represents a single node in a nested field's "path". The fieldNode
// is a stripped-down version of fieldInfo that holds only that information that
// is needed by the generator to produce correct Go field selector expressions.
type fieldNode struct {
	// The name of the field.
	name string
	// The tag of the field.
	tag tagutil.Tag
	// The name of the field's type. Empty if the type is unnamed.
	typeName string
	// The package import path for the field's type. Empty if the type is unnamed.
	typePkgPath string
	// The name of the package of the field's type. Empty if the type is unnamed.
	typePkgName string
	// The local name of the imported package of the field's type (including ".").
	// Empty if the type is unnamed.
	typePkgLocal string
	// Indicates whether or not the type is imported.
	isImported bool
	// Indicates whether or not the field is embedded.
	isEmbedded bool
	// Indicates whether or not the field is exported.
	isExported bool
	// Indicates whether or not the field type is a pointer type.
	isPointer bool
}

// fieldDatum holds the bare minimum of information for a field, its name and
// its type, it is used to represent a parameter of a search condition.
type fieldDatum struct {
	// The name of the field.
	name string
	// The type of the field.
	typ typeInfo
}

// joinBlock represents the result of the analysis of a queryStruct's "join block" field.
// The joinBlock is used by the generator to produce a list of table JOIN expressions in
// a SELECT's FROM clause, or an UPDATE's FROM clause, or a DELETE's USING clause.
//
// The joinBlock is declared in 3 different ways:
// - As a "join" field in a select query type
// - As a "from" field in an update query type
// - As a "using" field in a delete query type
type joinBlock struct {
	// The identifier of the top relation in a DELETE-USING / UPDATE-FROM
	// clause, empty in SELECT commands.
	relId relId
	// The list of join items declared in a join block.
	items []*joinItem
}

// joinItem represents the result of parsing the tag of a gosql.JoinXxx directive.
// The joinItem is used by the generator to produce a single JOIN clause.
type joinItem struct {
	// The type of the join.
	joinType joinType
	// The identifier of the relation to be joined.
	relId relId
	// A list of search conditions for the join specification.
	conds []*searchCondition
}

// whereBlock represents the result of the analysis of a queryStruct's "where block"
// field. The whereBlock is used by the generator to produce a WHERE clause.
type whereBlock struct {
	// The name of the "where block" field.
	name string
	// The list of search conditions declared in the whereBlock.
	conds []*searchCondition
}

// searchCondition represents the result of the analysis of a whereBlock's or
// joinBlock's field or directive. The searchCondition is used by the generator
// to produce a search condition in a WHERE clause, or a qualified JOIN ON clause.
type searchCondition struct {
	// If set, the logical connective.
	bool boolean
	// The specific search condition type:
	// - For a whereBlock this can be either searchConditionField, searchConditionColumn,
	//   searchConditionBetween, or searchConditionNested.
	// - For a joinBlock this can only be searchConditionColumn.
	cond interface{}
}

// searchConditionNested represents the result of the analysis of a nested whereBlock.
// The searchConditionNested is used by the generator to produce nested, parenthesized
// search conditions in a WHERE clause.
type searchConditionNested struct {
	// The field's name.
	name string
	// The list of search conditions declared in the nested field.
	conds []*searchCondition
}

// searchConditionField represents the result of the analysis of a whereBlock's field.
// The searchConditionField is used by the generator to produce a column-to-parameter
// comparison in a WHERE clause, passing the field value as the argument to the query.
type searchConditionField struct {
	// The field's name.
	name string
	// The field's type information.
	typ typeInfo
	// The identifier of the associated column.
	colId colId
	// The type of the condition's predicate.
	pred predicate
	// If set, indentifies the quantifier to be used with the predicate.
	qua quantifier
	// If set, the name of the function to be applied to both predicands.
	modFunc funcName
}

// searchConditionColumn represents the result of the analysis of a gosql.Column
// directive's tag, or a gosql.JoinXxx directive's tag. The searchConditionColumn
// is used by the generator to produce a search condition with either a unary column
// comparison, a column-to-column comparison, or a column-to-literal comparison as
// part of a WHERE clause or a qualified JOIN clause.
type searchConditionColumn struct {
	// The column representing the LHS predicand.
	colId colId
	// If set, the column representing the RHS predicand.
	colId2 colId
	// If set, the literal value representing the RHS predicand.
	literal string
	// If set, indentifies the type of the condition's predicate.
	pred predicate
	// If set, indentifies the quantifier to be used with the predicate.
	qua quantifier
}

// searchConditionBetween represents the result of the analysis of a whereBlock's "between" field.
// The searchConditionBetween is used by the generator to produce a BETWEEN predicate in a WHERE clause.
type searchConditionBetween struct {
	// The name of the "between" field.
	name string
	// The type of the BETWEEN predicate.
	pred predicate
	// The primary predicand of the BETWEEN predicate.
	colId colId
	// The lower bound range predicand. Either a colId, or a fieldDatum.
	x interface{}
	// The upper bound range predicand. Either a colId, or a fieldDatum.
	y interface{}
}

// onConflictBlock represents the result of the analysis of a queryStruct's "on conflict" field.
// The onConflictBlock is used by the generator to produce an ON CONFLICT clause in an INSERT query.
type onConflictBlock struct {
	// If set, indicates that the gosql.Column directive was used and the contents
	// of the slice are the column ids parsed from the tag of that directive.
	column []colId
	// If set, indicates that the gosql.Index directive was used. The value
	// is parsed from the directive's tag and represents the index to be used
	// for the on conflict target.
	//
	// NOTE(mkopriva): The index name will be used by the db check to retrive the index
	// expression and the generator will use that to produce the conflict target.
	index string
	// If set, indicates that the gosql.Constraint directive was used. The value
	// is the name of a unique constraint as parsed from the directive's tag.
	//
	// NOTE(mkopriva): The generator will use this value to generate
	// the ON CONFLICT ON CONSTRAINT <constraint_name> clause.
	constraint string
	// If set to true, indicates that the gosql.Ignore directive was used.
	ignore bool
	// If set, indicates that the gosql.Update directive was used, the contents
	// of the colIdList will hold the column ids parsed from the directive's tag.
	update *colIdList
}

// The limitField holds the information extracted from a queryStruct's gosql.Limit
// directive or a valid "limit" field. The information is then used by the generator
// to produce a LIMIT clause in a SELECT query.
type limitField struct {
	// The name of the field, if empty it indicates that the limitField
	// was produced from the gosql.Limit directive.
	name string
	// The value provided in the limit field's / directive's `sql` tag.
	// If the limitField was produced from a directive the value will be
	// used as a constant.
	// If the limitField was produced from a normal field the value will *only*
	// be used if the field's actual value is empty, at runtime during the query's
	// execution, essentially acting as a default fallback.
	value uint64
}

// The offsetField holds the information extracted from a queryStruct's gosql.Offset
// directive or a valid "offset" field. The information is then used by the generator
// to produce an OFFSET clause in a SELECT query.
type offsetField struct {
	// The name of the field, if empty it indicates that the offsetField
	// was produced from the gosql.Offset directive.
	name string
	// The value provided in the offset field's / directive's `sql` tag.
	// If the offsetField was produced from a directive the value will be
	// used as a constant.
	// If the offsetField was produced from a normal field the value will *only*
	// be used if the field's actual value is empty, at runtime during the query's
	// execution, essentially acting as a default fallback.
	value uint64
}

// orderByList represents the result of the analysis of a gosql.OrderBy directive.
// The orderByList is used by the generator to produce the ORDER BY clause.
type orderByList struct {
	// The list of individual orderByItems as parsed from the directive's tag.
	items []*orderByItem
}

// orderByItem represents a single item parsed from the tag of a gosql.OrderBy
// directive. The orderByItem is used by the generator to produce a single
// item in the "sort specification list" of an ORDER BY clause.
type orderByItem struct {
	// The identifier of the column to order by.
	colId colId
	// The direction of the sort order.
	direction orderDirection
	// The position of nulls in the sort order.
	nulls nullsPosition
}

// orderDirection is used to specify the order direction in an ORDER BY clause.
type orderDirection uint8

const (
	orderAsc  orderDirection = iota // ASC, default
	orderDesc                       // DESC
)

// nullsPosition is used to specify the position of NULLs in an ORDER BY clause.
type nullsPosition uint8

const (
	_          nullsPosition = iota // none specified, i.e. default
	nullsFirst                      // NULLS FIRST
	nullsLast                       // NULLS LAST
)

// boolean operation
type boolean uint

const (
	_       boolean = iota // no bool
	boolAnd                // conjunction
	boolOr                 // disjunction
	boolNot                // negation
)

// predicate represents the type of search condition's predicate.
type predicate uint

const (
	_ predicate = iota // no predicate

	// binary comparison predicates
	isEQ        // equals
	notEQ       // not equals
	notEQ2      // not equals
	isLT        // less than
	isGT        // greater than
	isLTE       // less than or equal
	isGTE       // greater than or equal
	isDistinct  // IS DISTINCT FROM
	notDistinct // IS NOT DISTINCT FROM

	// pattern matching predicates
	isMatch    // match regular expression
	isMatchi   // match regular expression (case insensitive)
	notMatch   // not match regular expression
	notMatchi  // not match regular expression (case insensitive)
	isLike     // LIKE
	notLike    // NOT LIKE
	isILike    // ILIKE
	notILike   // NOT ILIKE
	isSimilar  // IS SIMILAR TO
	notSimilar // IS NOT SIMILAR TO

	// array predicates
	isIn  // IN
	notIn // NOT IN

	// range predicates
	isBetween      // BETWEEN x AND y
	notBetween     // NOT BETWEEN x AND y
	isBetweenSym   // BETWEEN SYMMETRIC x AND y
	notBetweenSym  // NOT BETWEEN SYMMETRIC x AND y
	isBetweenAsym  // BETWEEN ASYMMETRIC x AND y
	notBetweenAsym // NOT BETWEEN ASYMMETRIC x AND y

	// null predicates
	isNull  // IS NULL
	notNull // IS NOT NULL

	// truth predicates
	isTrue     // IS TRUE
	notTrue    // IS NOT TRUE
	isFalse    // IS FALSE
	notFalse   // IS NOT FALSE
	isUnknown  // IS UNKNOWN
	notUnknown // IS NOT UNKNOWN
)

// canQuantify reports whether or not the predicate can be used together with a quantifier.
func (p predicate) canQuantify() bool {
	return p.isBinary() || p.isPatternMatch()
}

// isBinary reports whether or not the predicate represents a binary comparison.
func (p predicate) isBinary() bool {
	return p == isEQ || p == notEQ || p == notEQ2 ||
		p == isLT || p == isGT || p == isLTE || p == isGTE ||
		p == isDistinct || p == notDistinct
}

// isUnary reports whether or not the predicate represents a unary comparison.
func (p predicate) isUnary() bool {
	return p == isNull || p == notNull ||
		p == isTrue || p == notTrue ||
		p == isFalse || p == notFalse ||
		p == isUnknown || p == notUnknown
}

// isBoolean reports whether or not the predicate represents a a boolean test.
func (p predicate) isBoolean() bool {
	return p == isTrue || p == notTrue ||
		p == isFalse || p == notFalse ||
		p == isUnknown || p == notUnknown
}

// isPatternMatch reports whether or not the predicate represents a pattern-match comparison.
func (p predicate) isPatternMatch() bool {
	return p == isMatch || p == isMatchi || p == notMatch || p == notMatchi ||
		p == isLike || p == notLike || p == isILike || p == notILike ||
		p == isSimilar || p == notSimilar
}

// isRange reports whether or not the predicate represents a range comparison.
func (p predicate) isRange() bool {
	return p == isBetween || p == notBetween ||
		p == isBetweenSym || p == notBetweenSym ||
		p == isBetweenAsym || p == notBetweenAsym
}

// isQuantified reports whether or not the predicate represents a quantified comparison.
func (p predicate) isQuantified() bool {
	return p == isIn || p == notIn
}

// stringToPredicate is a map of string literals to supported predicates. Used for parsing tags.
var stringToPredicate = map[string]predicate{
	"=":           isEQ,
	"<>":          notEQ,
	"!=":          notEQ2,
	"<":           isLT,
	">":           isGT,
	"<=":          isLTE,
	">=":          isGTE,
	"isdistinct":  isDistinct,
	"notdistinct": notDistinct,

	"~":          isMatch,
	"~*":         isMatchi,
	"!~":         notMatch,
	"!~*":        notMatchi,
	"islike":     isLike,
	"notlike":    notLike,
	"isilike":    isILike,
	"notilike":   notILike,
	"issimilar":  isSimilar,
	"notsimilar": notSimilar,

	"isin":  isIn,
	"notin": notIn,

	"isbetween":      isBetween,
	"notbetween":     notBetween,
	"isbetweensym":   isBetweenSym,
	"notbetweensym":  notBetweenSym,
	"isbetweenasym":  isBetweenAsym,
	"notbetweenasym": notBetweenAsym,

	"isnull":     isNull,
	"notnull":    notNull,
	"istrue":     isTrue,
	"nottrue":    notTrue,
	"isfalse":    isFalse,
	"notfalse":   notFalse,
	"isunknown":  isUnknown,
	"notunknown": notUnknown,
}

// predicateAdjectives is a whitelist of predicate adjectives and adverbs. Used for parsing tags.
var predicateAdjectives = []string{
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

// quantifier represents the type of a comparison predicate quantifier.
type quantifier uint8

const (
	_         quantifier = iota // no operator
	quantAny                    // ANY
	quantSome                   // SOME
	quantAll                    // ALL
)

// stringToQuantifier is a map of string literals to supported quantifiers. Used for parsing of tags.
var stringToQuantifier = map[string]quantifier{
	"any":  quantAny,
	"some": quantSome,
	"all":  quantAll,
}

// overridingKind represents the kind of the optional OVERRIDING clause in an INSERT query.
type overridingKind uint8

const (
	_                overridingKind = iota // no overriding
	overridingSystem                       // OVERRIDING SYSTEM VALUE
	overridingUser                         // OVERRIDING USER VALUE
)

// funcName is the name of a database function that can either be used to modify
// a value, like lower, upper, etc. or a function that can be used as an aggregate.
type funcName string

// joinType represents the type of an SQL JOIN clause.
type joinType uint

const (
	_         joinType = iota // no join
	joinCross                 // CROSS JOIN
	joinInner                 // INNER JOIN
	joinLeft                  // LEFT JOIN
	joinRight                 // RIGHT JOIN
	joinFull                  // FULL JOIN
)

// stringToJoinType is a map of string literals to supported join types. Used for parsing of directives.
var stringToJoinType = map[string]joinType{
	"crossjoin": joinCross,
	"innerjoin": joinInner,
	"leftjoin":  joinLeft,
	"rightjoin": joinRight,
	"fulljoin":  joinFull,
}

// typeKind indicate the specific kind of a Go type.
type typeKind uint

const (
	// basic
	kindInvalid typeKind = iota
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
	kindInterface
	kindMap
	kindPtr
	kindSlice
	kindStruct
	kindChan
	kindFunc

	// alisases
	kindByte = kindUint8
	kindRune = kindInt32
)

func (k typeKind) String() string {
	if s, ok := typeKindToString[k]; ok {
		return s
	}
	return "<invalid>"
}

// basicKindToTypeKind maps basic kinds, as declared in go/types, to typeKind values.
// Used for resolving a type's kind.
var basicKindToTypeKind = map[types.BasicKind]typeKind{
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

var typeKindToString = map[typeKind]string{
	// builtin basic
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

	// non-basic
	kindArray:     "<array>",
	kindChan:      "<chan>",
	kindFunc:      "<func>",
	kindInterface: "<interface>",
	kindMap:       "<map>",
	kindPtr:       "<pointer>",
	kindSlice:     "<slice>",
	kindStruct:    "<struct>",
}

const (
	// A list of common Go types ("identifiers" and "literals")
	// used for resolving type convertability.
	goTypeBool                     = "bool"
	goTypeBoolSlice                = "[]bool"
	goTypeString                   = "string"
	goTypeStringSlice              = "[]string"
	goTypeStringSliceSlice         = "[][]string"
	goTypeStringMap                = "map[string]string"
	goTypeStringMapSlice           = "[]map[string]string"
	goTypeByte                     = "byte"
	goTypeByteSlice                = "[]byte"
	goTypeByteSliceSlice           = "[][]byte"
	goTypeByteArray16              = "[16]byte"
	goTypeByteArray16Slice         = "[][16]byte"
	goTypeRune                     = "rune"
	goTypeRuneSlice                = "[]rune"
	goTypeRuneSliceSlice           = "[][]rune"
	goTypeInt                      = "int"
	goTypeIntSlice                 = "[]int"
	goTypeIntArray2                = "[2]int"
	goTypeIntArray2Slice           = "[][2]int"
	goTypeInt8                     = "int8"
	goTypeInt8Slice                = "[]int8"
	goTypeInt8SliceSlice           = "[][]int8"
	goTypeInt16                    = "int16"
	goTypeInt16Slice               = "[]int16"
	goTypeInt16SliceSlice          = "[][]int16"
	goTypeInt32                    = "int32"
	goTypeInt32Slice               = "[]int32"
	goTypeInt32Array2              = "[2]int32"
	goTypeInt32Array2Slice         = "[][2]int32"
	goTypeInt64                    = "int64"
	goTypeInt64Slice               = "[]int64"
	goTypeInt64Array2              = "[2]int64"
	goTypeInt64Array2Slice         = "[][2]int64"
	goTypeUint                     = "uint"
	goTypeUintSlice                = "[]uint"
	goTypeUint8                    = "uint8"
	goTypeUint8Slice               = "[]uint8"
	goTypeUint16                   = "uint16"
	goTypeUint16Slice              = "[]uint16"
	goTypeUint32                   = "uint32"
	goTypeUint32Slice              = "[]uint32"
	goTypeUint64                   = "uint64"
	goTypeUint64Slice              = "[]uint64"
	goTypeFloat32                  = "float32"
	goTypeFloat32Slice             = "[]float32"
	goTypeFloat64                  = "float64"
	goTypeFloat64Slice             = "[]float64"
	goTypeFloat64Array2            = "[2]float64"
	goTypeFloat64Array2Slice       = "[][2]float64"
	goTypeFloat64Array2SliceSlice  = "[][][2]float64"
	goTypeFloat64Array2Array2      = "[2][2]float64"
	goTypeFloat64Array2Array2Slice = "[][2][2]float64"
	goTypeFloat64Array3            = "[3]float64"
	goTypeFloat64Array3Slice       = "[][3]float64"
	goTypeIPNet                    = "net.IPNet"
	goTypeIPNetSlice               = "[]net.IPNet"
	goTypeTime                     = "time.Time"
	goTypeTimeSlice                = "[]time.Time"
	goTypeTimeArray2               = "[2]time.Time"
	goTypeTimeArray2Slice          = "[][2]time.Time"
	goTypeBigInt                   = "big.Int"
	goTypeBigIntSlice              = "[]big.Int"
	goTypeBigIntArray2             = "[2]big.Int"
	goTypeBigIntArray2Slice        = "[][2]big.Int"
	goTypeBigFloat                 = "big.Float"
	goTypeBigFloatSlice            = "[]big.Float"
	goTypeBigFloatArray2           = "[2]big.Float"
	goTypeBigFloatArray2Slice      = "[][2]big.Float"
	goTypeNullStringMap            = "map[string]sql.NullString"
	goTypeNullStringMapSlice       = "[]map[string]sql.NullString"
)

// isIntegerType reports whether or not the given type is one of the basic (un)signed integer types.
func isIntegerType(typ types.Type) bool {
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

// isBoolType reports whether or not the given type is a boolean.
func isBoolType(typ types.Type) bool {
	basic, ok := typ.(*types.Basic)
	if !ok {
		return false
	}
	return basic.Kind() == types.Bool
}

// isErrorHandler reports whether or not the given type implements the gosql.ErrorHandler interface.
func isErrorHandler(typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	return typesutil.ImplementsErrorHandler(named)
}

// isErrorInfoHandler reports whether or not the given type implements the gosql.ErrorInfoHandler interface.
func isErrorInfoHandler(typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	return typesutil.ImplementsErrorInfoHandler(named)
}

// isFilterType reports whether or not the given type is the gosql.Filter type.
func isFilterType(typ types.Type) bool {
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

// isColId reports whether or not the given value is a valid column identifier.
func isColId(val string) bool {
	return rxColId.MatchString(val) && !rxReserved.MatchString(val)
}

var dataTypeCache = struct {
	sync.RWMutex
	m map[string]*dataType
}{m: make(map[string]*dataType)}
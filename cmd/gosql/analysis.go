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
				// TODO DEFAULTS ALL is INSERT-only, so if a.query.defaults.all==true:
				// - either check if this an update and if so return an error
				// - or set a.query.defaults.all to false and instead
				//   fill the items slice with all columns in the target...
				// - or leave it as is, and have the generator check
				//   whether this is update or not and if it is, instead
				//   of generating DEFAULTS ALL, it would generate the DEAULT
				//   marker for each column in the SET (<column_list>).
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

	// TODO if queryKind is Update, and dataField.record.isSlice == true THEN only
	// matching the items by PKEY makes sense, therefore a whereBlock, or filter,
	// or the all directive, should be disallowed!

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

	// XXX Experimental: Not exactly sure that types.Type.String()
	// will NOT produce conflicting values for different types.
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
			f.useCoalesce, f.coalesceValue = a.coalesceinfo(tag)

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

func (a *analyzer) coalesceinfo(tag tagutil.Tag) (use bool, val string) {
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
		items  []*predicateItem
		nested *predicateNested
		ns     *typesutil.NamedStruct // the struct type of the whereBlock
		idx    int                    // keeps track of the field index
	}

	// root holds the reference to the root level predicate items
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

			item := new(predicateItem)
			loop.items = append(loop.items, item)

			// Analyze the bool operation for any but the first
			// item in a whereBlock. Fail if a value was provided
			// but it is not "or" nor "and".
			if len(loop.items) > 1 {
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

				pred := new(predicateNested)
				pred.name = fld.Name()
				item.pred = pred

				loop2 := new(loopstate)
				loop2.ns = ns
				loop2.nested = pred
				stack = append(stack, loop2)
				continue stackloop
			}

			lhs, op, op2, rhs := a.splitpredicateexpr(sqltag)

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

					pred := new(predicateColumn)
					pred.colId = colId
					pred.kind = stringToPredicateKind[op]
					pred.qua = stringToQuantifier[op2]

					if isColId(rhs) {
						pred.colId2, _ = a.colId(rhs, fld) // ignore error since isColId returned true
					} else {
						pred.literal = rhs // assume literal expression
					}

					if pred.kind.isUnary() {
						// TODO add test
						return errors.IllegalUnaryPredicateError
					} else if pred.qua > 0 && !pred.kind.isQuantifiable() {
						return errors.BadPredicateComboError
					}

					item.pred = pred
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
				predKind := stringToPredicateKind[op]
				if !predKind.isUnary() {
					return errors.BadUnaryPredicateError
				}
				if len(op2) > 0 {
					return errors.ExtraQuantifierError
				}

				pred := new(predicateColumn)
				pred.colId = colId
				pred.kind = predKind
				item.pred = pred
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

				pred := new(predicateBetween)
				pred.name = fld.Name()
				pred.colId = colId
				pred.kind = stringToPredicateKind[op]
				pred.x = x
				pred.y = y
				item.pred = pred
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
			predKind := stringToPredicateKind[op]
			if predKind.isUnary() {
				// TODO add test
				return errors.IllegalUnaryPredicateError
			}

			qua := stringToQuantifier[op2]
			if qua > 0 && !predKind.isQuantifiable() {
				return errors.BadPredicateComboError
			}

			pred := new(predicateField)
			pred.name = fld.Name()
			pred.typ, _ = a.typeInfo(fld.Type())
			pred.colId = colId
			pred.kind = predKind
			pred.qua = qua
			pred.modfunc = a.funcName(tag["sql"][1:])

			if pred.qua > 0 && pred.typ.kind != kindSlice && pred.typ.kind != kindArray {
				return errors.BadQuantifierFieldTypeError
			}

			item.pred = pred
		}

		if loop.nested != nil {
			loop.nested.items = loop.items
		}

		stack = stack[:len(stack)-1]
	}

	wb := new(whereBlock)
	wb.name = field.Name()
	wb.items = root.items
	a.query.whereBlock = wb
	return nil
}

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

			var predicates []*predicateItem
			for _, val := range tag["sql"][1:] {
				vals := strings.Split(val, ";")
				for i, val := range vals {

					pred := new(predicateColumn)
					lhs, op, op2, rhs := a.splitpredicateexpr(val)
					if pred.colId, err = a.colId(lhs, fld); err != nil {
						return err
					}

					// optional right-hand side
					if isColId(rhs) {
						pred.colId2, _ = a.colId(rhs, fld) // ignore error since isColId returned true
					} else {
						pred.literal = rhs
					}

					pred.kind = stringToPredicateKind[op]
					pred.qua = stringToQuantifier[op2]

					if len(rhs) > 0 {
						if pred.kind.isUnary() {
							// TODO add test
							return errors.IllegalUnaryPredicateError
						} else if pred.qua > 0 && !pred.kind.isQuantifiable() {
							return errors.BadPredicateComboError
						}
					} else {
						if !pred.kind.isUnary() {
							return errors.BadUnaryPredicateError
						} else if len(op2) > 0 {
							return errors.ExtraQuantifierError
						}
					}

					item := new(predicateItem)
					item.pred = pred
					if len(predicates) > 0 && i == 0 {
						item.bool = boolAnd
					} else if len(predicates) > 0 && i > 0 {
						item.bool = boolOr
					}

					predicates = append(predicates, item)
				}
			}

			item := new(joinItem)
			item.joinType = stringToJoinType[dirname]
			item.relId = id
			item.predicates = predicates
			join.items = append(join.items, item)
		default:
			return errors.IllegalJoinBlockDirectiveError
		}

	}

	a.query.joinBlock = join
	return nil
}

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
func (a *analyzer) splitpredicateexpr(expr string) (lhs, cop, qua, rhs string) {
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
// If elideptr is true the "*" will be elided from the output.
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

// is returns true if t represents a type one of the given kinds or a pointer
// to one of the given kinds.
func (t *typeInfo) is(kk ...typeKind) bool {
	for _, k := range kk {
		if t.kind == k || (t.kind == kindPtr && t.elem.kind == k) {
			return true
		}
	}
	return false
}

// isSlice returns true if t represents a slice type whose elem type is one of
// the given kinds.
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

// isSliceN returns true if t represents an n dimensional slice type whose
// base elem type is one of the given kinds.
func (t *typeInfo) isSliceN(n int, kk ...typeKind) bool {
	for ; n > 0; n-- {
		if t.kind != kindSlice {
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

// isNamed returns true if t represents a named type, or a pointer to a named
// type, whose package path and type name match the given arguments.
func (t *typeInfo) isNamed(pkgPath, name string) bool {
	if t.kind == kindPtr {
		t = t.elem
	}
	return t.pkgPath == pkgPath && t.name == name
}

// isNilable returns true if t represents a type that can be nil.
func (t *typeInfo) isNilable() bool {
	return t.is(kindPtr, kindSlice, kindArray, kindMap, kindInterface)
}

// indicates whether or not the MarshalJSON method can be called on the type.
func (t *typeInfo) canJSONMarshal() bool {
	return t.isJSONMarshaler || (t.kind == kindPtr && t.elem.isJSONMarshaler)
}

// indicates whether or not the UnmarshalJSON method can be called on the type.
func (t *typeInfo) canJSONUnmarshal() bool {
	return t.isJSONUnmarshaler || (t.kind == kindPtr && t.elem.isJSONUnmarshaler)
}

// indicates whether or not the MarshalXML method can be called on the type.
func (t *typeInfo) canXMLMarshal() bool {
	return t.isXMLMarshaler || (t.kind == kindPtr && t.elem.isXMLMarshaler)
}

// indicates whether or not the UnmarshalXML method can be called on the type.
func (t *typeInfo) canXMLUnmarshal() bool {
	return t.isXMLUnmarshaler || (t.kind == kindPtr && t.elem.isXMLUnmarshaler)
}

// fieldInfo holds information about a dataType's field.
type fieldInfo struct {
	name string   // name of the struct field
	typ  typeInfo // info about the field's type
	// if the field is nested, path will hold the parent fields' information
	path []*fieldNode
	// indicates whether or not the field is embedded
	isEmbedded bool
	// indicates whether or not the field is exported
	isExported bool
	// the field's parsed tag
	tag tagutil.Tag
	// the id of the corresponding column
	colId colId
	// identifies the field's corresponding column as a primary key
	isPKey bool
	// indicates that if the field's value is EMPTY then NULL should
	// be stored in the column during INSERT/UPDATE
	nullEmpty bool
	// indicates that field should only be read from the database and never written
	readOnly bool
	// indicates that field should only be written into the database and never read
	writeOnly bool
	// indicates that the DEFAULT marker should be used during INSERT/UPDATE
	useDefault bool
	// indicates that the column value should be marshaled/unmarshaled
	// to/from json before/after being stored/retrieved.
	useJSON bool
	// indicates that the column value should be marshaled/unmarshaled
	// to/from xml before/after being stored/retrieved.
	useXML bool
	// for UPDATEs, if set to true, it indicates that the provided field
	// value should be added to the already existing column value.
	useAdd bool
	// indicates whether or not an implicit CAST should be allowed.
	canCast bool
	// if set to true it indicates that the column value should be wrapped
	// in a COALESCE call when read from the db.
	useCoalesce   bool
	coalesceValue string
}

type fieldNode struct {
	name         string
	tag          tagutil.Tag
	typeName     string // the name of a named type or empty string for unnamed types
	typePkgPath  string // the package import path
	typePkgName  string // the package's name
	typePkgLocal string // the local package name (including ".")
	isImported   bool   // indicates whether or not the type is imported
	isEmbedded   bool   // indicates whether or not the field is embedded
	isExported   bool   // indicates whether or not the field is exported
	isPointer    bool   // indicates whether or not the field type is a pointer type
}

// fieldDatum holds the bare minimum info of a field, its name and type,
// and it is used to represent a parameter of a predicate.
type fieldDatum struct {
	name string
	typ  typeInfo
}

// joinBlock  ...........
type joinBlock struct {
	// The identifier of the top relation in a DELETE-USING / UPDATE-FROM
	// clause, empty in SELECT commands.
	relId relId
	items []*joinItem
}

// joinItem ....
type joinItem struct {
	joinType   joinType
	relId      relId
	predicates []*predicateItem
}

// whereBlock ....
type whereBlock struct {
	name  string // name of the where block field
	items []*predicateItem
}

// predicateItem.......
type predicateItem struct {
	bool boolean
	pred interface{} // predicateField, predicateColumn, predicateBetween, or predicateNested
}

// predicateNested ....
type predicateNested struct {
	name  string
	items []*predicateItem
}

// predicateField holds the information extracted from a whereBlock's "predicate" field.
type predicateField struct {
	// The name of the field that's one of the two predicands.
	name string
	// Information on the field's type.
	typ typeInfo
	// Identifies the column that's the other of the two predicands.
	colId colId
	// The kind of the predicate.
	kind predicateKind
	// If set, indentifies the quantifier to be used with the predicate.
	qua quantifier
	// If set, the name of the function to be applied to both predicands.
	modfunc funcName
}

// predicateColumn holds the information extracted from a queryStruct's gosql.Column
// directive. The predicateColumn type can represent either a column with a unary comparison
// predicate, a column-to-column comparison, or a column-to-literal comparison.
type predicateColumn struct {
	// The column representing the LHS predicand.
	colId colId
	// If set, the column representing the RHS predicand.
	colId2 colId
	// If set, the literal value representing the RHS predicand.
	literal string
	// If set, indentifies the kind of the predicate.
	kind predicateKind
	// If set, indentifies the quantifier to be used with the predicate.
	qua quantifier
}

// predicateBetween holds the information extracted from a queryStruct's "between" field.
type predicateBetween struct {
	// The name of the field declaring the predicate.
	name string
	// The primary predicand of the predicate.
	colId colId
	// Identifies the kind of the between predicate
	kind predicateKind
	x, y interface{}
}

// onConflictBlock holds the information extracted from a queryStruct's "on conflict" field.
type onConflictBlock struct {
	// If set, indicates that the gosql.Column directive was used and the contents
	// of the slice are the column ids parsed from the tag of that directive.
	column []colId
	// If set, indicates that the gosql.Index directive was used. The value
	// is parsed from the directive's tag and represents the index to be used
	// for the on conflict target.
	//
	// NOTE The index name will be used by the db check to retrive the index
	// expression and the generator will use that to produce the conflict target.
	index string
	// If set, indicates that the gosql.Constraint directive was used. The value
	// is the name of a unique constraint as parsed from the directive's tag.
	//
	// NOTE The generator will use this value to generate
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

// orderByList contains the information extracted from a queryStruct's gosql.OrderBy directive.
type orderByList struct {
	items []*orderByItem
}

// orderByItem holds information about the sort expression to be used in an ORDER BY clause.
type orderByItem struct {
	colId     colId          // the identifier of the column to order by
	direction orderDirection // the direction of the sort order
	nulls     nullsPosition  // the position of nulls in the sort order
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

type predicateKind uint // predicate

const (
	_ predicateKind = iota // no predicate

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

// isQuantifiable returns true if the predicate can be used together with a quantifier.
func (p predicateKind) isQuantifiable() bool {
	return p.isBinary() || p.isPatternPred()
}

// isBinary returns true if the predicate is a binary comparison predicate.
func (p predicateKind) isBinary() bool {
	return p == isEQ || p == notEQ || p == notEQ2 ||
		p == isLT || p == isGT || p == isLTE || p == isGTE ||
		p == isDistinct || p == notDistinct
}

// isUnary returns true if the predicate is a "unary" predicate.
func (p predicateKind) isUnary() bool {
	return p == isNull || p == notNull ||
		p == isTrue || p == notTrue ||
		p == isFalse || p == notFalse ||
		p == isUnknown || p == notUnknown
}

// isPatternPred returns true if the predicate is a pattern matching predicate.
func (p predicateKind) isPatternPred() bool {
	return p == isMatch || p == isMatchi || p == notMatch || p == notMatchi ||
		p == isLike || p == notLike || p == isILike || p == notILike ||
		p == isSimilar || p == notSimilar
}

// isRangePred returns true if the predicate is a range predicate.
func (p predicateKind) isRangePred() bool {
	return p == isBetween || p == notBetween ||
		p == isBetweenSym || p == notBetweenSym ||
		p == isBetweenAsym || p == notBetweenAsym
}

// isTruthPred returns true if the predicate is a "truth" predicate.
func (p predicateKind) isTruthPred() bool {
	return p == isTrue || p == notTrue ||
		p == isFalse || p == notFalse ||
		p == isUnknown || p == notUnknown
}

// isArrayPred returns true if the predicate is an array predicate.
func (p predicateKind) isArrayPred() bool {
	return p == isIn || p == notIn
}

// stringToPredicateKind is a map of string literals to supported predicateKinds.
// Used for parsing tags.
var stringToPredicateKind = map[string]predicateKind{
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

type overridingKind uint8

const (
	_                overridingKind = iota // no overriding
	overridingSystem                       // OVERRIDING SYSTEM VALUE
	overridingUser                         // OVERRIDING USER VALUE
)

// funcName is the name of a database function that can either be used to modify
// a value, like lower, upper, etc. or a function that can be used as an aggregate.
type funcName string

type joinType uint

const (
	_         joinType = iota // no join
	joinLeft                  // LEFT JOIN
	joinRight                 // RIGHT JOIN
	joinFull                  // FULL JOIN
	joinCross                 // CROSS JOIN
)

// stringToJoinType is a map of string literals to supported join types. Used for parsing of directives.
var stringToJoinType = map[string]joinType{
	"leftjoin":  joinLeft,
	"rightjoin": joinRight,
	"fulljoin":  joinFull,
	"crossjoin": joinCross,
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
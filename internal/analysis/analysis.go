package analysis

import (
	"fmt"
	"go/token"
	"go/types"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/frk/gosql/internal/config"
	"github.com/frk/gosql/internal/typesutil"
	"github.com/frk/tagutil"
)

var _ = log.Println
var _ = fmt.Println

var (
	// NOTE(mkopriva): Identifiers MUST begin with a letter (a-z) or an underscore (_).
	// Subsequent characters in an identifier can be letters, underscores, and digits (0-9).

	// Matches a valid identifier.
	rxIdent = regexp.MustCompile(`^[A-Za-z_]\w*$`)

	// Matches a valid db relation identifier.
	// - Valid format: [schema_name.]relation_name[:alias_name]
	rxRelIdent = regexp.MustCompile(`^(?:[A-Za-z_]\w*\.)?[A-Za-z_]\w*(?:\:[A-Za-z_]\w*)?$`)

	// Matches a valid table column reference.
	// - Valid format: [rel_name_or_alias.]column_name
	rxColIdent = regexp.MustCompile(`^(?:[A-Za-z_]\w*\.)?[A-Za-z_]\w*$`)

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

// analysis holds the state of the analyzer.
type analysis struct {
	cfg  config.Config
	fset *token.FileSet
	// The named type under analysis.
	named *types.Named
	// The package path of the type under analysis.
	pkgPath string
	// If the type under analysis a "filter" type this field will hold
	// the result of the analysis, otherwise it will be nil.
	filter *FilterStruct
	// If the type under analysis a "query" type this field will hold
	// the result of the analysis, otherwise it will be nil.
	query *QueryStruct
	// ...
	info *Info
}

// Info holds information related to an analyzed TargetStruct. If the analysis
// returns an error, the collected information will be incomplete.
type Info struct {
	// The FileSet associated with the analyzed TargetStruct.
	FileSet *token.FileSet
	// The package path of the analyzed TargetStruct.
	PkgPath string
	// The type name of the analyzed TargetStruct.
	TypeName string
	// The soruce position of the TargetStruct's type name.
	TypeNamePos token.Pos
	// FieldMap maintains a map of pointers of arbitrary type that represent
	// the result of analyzed fields, to the fields' related go/types specific
	// information. Intended for error reporting by the backend type-checker.
	FieldMap map[FieldPtr]FieldVar
	// RelSpace maintains a set of *unique* relation names or aliases that map
	// onto their respective RelIdent values that were parsed from struct tags.
	RelSpace map[string]RelIdent
	// The analyzed struct.
	Struct TargetStruct
}

// Run analyzes the given named type which is expected to be a struct type whose name
// prefix matches one of the allowed prefixes. It panics if the named type is not actually
// a struct type or if its name does not start with one of the predefined prefixes.
func Run(fset *token.FileSet, named *types.Named, pos token.Pos, cfg config.Config) (*Info, error) {
	structType, ok := named.Underlying().(*types.Struct)
	if !ok {
		panic(named.Obj().Name() + " must be a struct type.") // this shouldn't happen
	}

	a := new(analysis)
	a.cfg = cfg
	a.fset = fset
	a.named = named
	a.pkgPath = named.Obj().Pkg().Path()

	a.info = new(Info)
	a.info.FileSet = fset
	a.info.PkgPath = a.pkgPath
	a.info.TypeName = named.Obj().Name()
	a.info.TypeNamePos = pos
	a.info.FieldMap = make(map[FieldPtr]FieldVar)
	a.info.RelSpace = make(map[string]RelIdent)

	if strings.HasPrefix(strings.ToLower(named.Obj().Name()), "filter") {
		s, err := analyzeFilterStruct(a, structType)
		if err != nil {
			return nil, err
		}
		a.info.Struct = s
	} else {
		s, err := analyzeQueryStruct(a, structType)
		if err != nil {
			return nil, err
		}
		a.info.Struct = s
	}
	return a.info, nil
}

func (a *analysis) error(code errorCode, f *types.Var, blockName, tagString, tagExpr, tagError string) error {
	e := &anError{Code: code, BlockName: blockName, TagString: tagString, TagExpr: tagExpr, TagError: tagError}
	if f != nil {
		p := a.fset.Position(f.Pos())
		e.PkgPath = a.named.Obj().Pkg().Path()
		e.TargetName = a.named.Obj().Name()
		e.FieldType = f.Type().String()
		e.FieldTypeKind = analyzeTypeKind(f.Type()).String()
		e.FieldName = f.Name()
		e.FileName = p.Filename
		e.FileLine = p.Line
	} else {
		p := a.fset.Position(a.named.Obj().Pos())
		e.PkgPath = a.named.Obj().Pkg().Path()
		e.TargetName = a.named.Obj().Name()
		e.FileName = p.Filename
		e.FileLine = p.Line
	}

	if a.query != nil && a.query.Rel != nil {
		e.RelField = a.query.Rel.FieldName
		e.RelType = a.query.Rel.Type
	} else if a.filter != nil && a.filter.Rel != nil {
		e.RelField = a.filter.Rel.FieldName
		e.RelType = a.filter.Rel.Type
	}
	return e
}

// analyzeFilterStruct ...
func analyzeFilterStruct(a *analysis, structType *types.Struct) (*FilterStruct, error) {
	a.filter = new(FilterStruct)
	a.filter.TypeName = a.named.Obj().Name()

	for i := 0; i < structType.NumFields(); i++ {
		fvar := structType.Field(i)
		ftag := structType.Tag(i)
		tag := tagutil.New(ftag)

		// Ensure that there is only one field with the "rel" tag.
		if _, ok := tag["rel"]; ok {
			if a.filter.Rel != nil {
				return nil, a.error(errConflictingRelTag, fvar, "", ftag, "", "")
			}

			rid, ecode := parseRelIdent(tag.First("rel"))
			if ecode > 0 {
				return nil, a.error(ecode, fvar, "", ftag, "", tag.First("rel"))
			} else if ecode, errval := addToRelSpace(a, rid); ecode > 0 {
				// NOTE(mkopriva): Because of the "a.filter.Rel != nil" check above and the
				// fact that FilterXxx types don't accept any other relation-specifying
				// fields, this branch will actually not run, nevertheless it is left here
				// just in case the implementation changes allowing for this error to occur.
				return nil, a.error(ecode, fvar, "", ftag, "", errval)
			}

			a.filter.Rel = new(RelField)
			a.filter.Rel.FieldName = fvar.Name()
			a.filter.Rel.Id = rid
			a.info.FieldMap[a.filter.Rel] = FieldVar{Var: fvar, Tag: ftag}
			if err := analyzeRelType(a, &a.filter.Rel.Type, fvar); err != nil {
				return nil, err
			}

			if a.filter.Rel.Type.IsIter {
				return nil, a.error(errIllegalIteratorField, fvar, "", ftag, "", "")
			}
			continue
		}

		// TODO(mkopriva): allow for embedding a struct with "common feature fields",
		// and make sure to also allow imported and local-unexported struct types.

		if dirname := typesutil.GetDirectiveName(fvar); fvar.Name() == "_" && len(dirname) > 0 {
			// fields with gosql directive types
			if strings.ToLower(dirname) == "textsearch" {
				if err := analyzeTextSearchDirective(a, fvar, ftag); err != nil {
					return nil, err
				}
			} else {
				return nil, a.error(errIllegalQueryField, fvar, "", "", "", "")
			}
		} else {
			// fields with specific names / types
			if typesutil.ImplementsGosqlFilterConstructor(fvar.Type()) {
				if err := analyzeFilterConstructorField(a, fvar, ftag); err != nil {
					return nil, err
				}
			}
		}
	}

	if a.filter.Rel == nil {
		return nil, a.error(errMissingRelField, nil, "", "", "", "") // TODO test
	}
	if a.filter.FilterConstructor == nil {
		return nil, a.error(errMissingFilterConstructor, nil, "", "", "", "") // TODO test
	}
	return a.filter, nil
}

// analyzeQueryStruct runs the analysis of a QueryStruct.
func analyzeQueryStruct(a *analysis, structType *types.Struct) (*QueryStruct, error) {
	a.query = new(QueryStruct)
	a.query.TypeName = a.named.Obj().Name()

	key := tolower(a.query.TypeName)
	if len(key) > 5 {
		key = key[:6]
	}
	switch key {
	case "insert":
		a.query.Kind = QueryKindInsert
	case "update":
		a.query.Kind = QueryKindUpdate
	case "select":
		a.query.Kind = QueryKindSelect
	case "delete":
		a.query.Kind = QueryKindDelete
	default:
		panic(a.query.TypeName + " struct type has unsupported name prefix.") // this shouldn't happen
	}

	// Find and anlyze the "rel" field.
	for i := 0; i < structType.NumFields(); i++ {
		ftag := structType.Tag(i)
		fvar := structType.Field(i)
		tag := tagutil.New(ftag)

		if _, ok := tag["rel"]; ok {
			if err := analyzeQueryStructRelField(a, fvar, ftag, tag.First("rel")); err != nil {
				return nil, err
			}
		}
	}
	if a.query.Rel == nil {
		return nil, a.error(errMissingRelField, nil, "", "", "", "")
	}

	// Analyze the rest of the query struct's fields.
	for i := 0; i < structType.NumFields(); i++ {
		ftag := structType.Tag(i)
		fvar := structType.Field(i)
		tag := tagutil.New(ftag)

		if _, ok := tag["rel"]; ok {
			continue
		}

		if dirname := typesutil.GetDirectiveName(fvar); fvar.Name() == "_" && len(dirname) > 0 {
			// fields with gosql directive types
			if err := analyzeQueryStructDirective(a, fvar, ftag, dirname); err != nil {
				return nil, err
			}
		} else {
			// fields with specific names / types
			if err := analyzeQueryStructField(a, fvar, ftag); err != nil {
				return nil, err
			}
		}
	}

	// TODO(mkopriva): if QueryKind is Select, Update, or Insert, and the analyzed
	// RelType.Fields slice is empty (for Select also check ResultType.Fields), then fail.

	// TODO(mkopriva): allow for embedding a struct with "common feature fields",
	// and make sure to also allow imported and local-unexported struct types.
	//
	// TODO(mkopriva): if QueryKind is Update and the record (single or slice) does not
	// have a primary key AND there's no WhereStruct, no filter, no all directive
	// return an error. That case suggests that all records should be updated
	// however the all directive must be provided explicitly, as a way to
	// ensure the programmer does not, by mistake, declare a query that
	// updates all records in a table.

	return a.query, nil
}

// analyzeQueryStructRelField [ ... ]
func analyzeQueryStructRelField(a *analysis, f *types.Var, ftag, reltag string) error {
	if a.query.Rel != nil {
		return a.error(errConflictingRelTag, f, "", ftag, "", "")
	}
	rid, ecode := parseRelIdent(reltag)
	if ecode > 0 {
		return a.error(ecode, f, "", ftag, "", reltag)
	} else if ecode, errval := addToRelSpace(a, rid); ecode > 0 {
		// NOTE(mkopriva): Because of the "a.query.Rel != nil" check above and the
		// fact that the rel field in query types is intentionally analyzed before
		// any other field, this branch will actually not run, nevertheless it is left here
		// just in case the implementation changes allowing for this error to occur.
		return a.error(ecode, f, "", ftag, "", errval)
	}

	a.query.Rel = new(RelField)
	a.query.Rel.FieldName = f.Name()
	a.query.Rel.Id = rid

	switch fname := strings.ToLower(a.query.Rel.FieldName); {
	default:
		if err := analyzeRelType(a, &a.query.Rel.Type, f); err != nil {
			return err
		}
		if (a.query.Kind == QueryKindInsert || a.query.Kind == QueryKindUpdate) && a.query.Rel.Type.IsIter {
			return a.error(errIllegalQueryField, f, "", ftag, "", "") // TODO test
		}
	case fname == "count" && isIntegerType(f.Type()):
		if a.query.Kind != QueryKindSelect {
			return a.error(errIllegalQueryField, f, "", ftag, "", "")

		}
		a.query.Kind = QueryKindSelectCount
	case fname == "exists" && isBoolType(f.Type()):
		if a.query.Kind != QueryKindSelect {
			return a.error(errIllegalQueryField, f, "", ftag, "", "")
		}
		a.query.Kind = QueryKindSelectExists
	case fname == "notexists" && isBoolType(f.Type()):
		if a.query.Kind != QueryKindSelect {
			return a.error(errIllegalQueryField, f, "", ftag, "", "")
		}
		a.query.Kind = QueryKindSelectNotExists
	case fname == "_" && typesutil.IsDirective("Relation", f.Type()):
		if a.query.Kind != QueryKindDelete {
			return a.error(errIllegalQueryField, f, "", ftag, "", "")
		}
		a.query.Rel.IsDirective = true
	}

	a.info.FieldMap[a.query.Rel] = FieldVar{Var: f, Tag: ftag}
	return nil
}

// analyzeQueryStructDirective [ ... ]
func analyzeQueryStructDirective(a *analysis, f *types.Var, tag string, dirname string) error {
	analyzers := map[string]func(*analysis, *types.Var, string) error{
		"all":      analyzeAllDirective,
		"default":  analyzeDefaultDirective,
		"force":    analyzeForceDirective,
		"optional": analyzeOptionalDirective,
		"return":   analyzeReturnDirective,
		"limit":    analyzeLimitFieldOrDirective,
		"offset":   analyzeOffsetFieldOrDirective,
		"orderby":  analyzeOrderByDirective,
		"override": analyzeOverrideDirective,
	}
	if afunc, ok := analyzers[strings.ToLower(dirname)]; ok {
		return afunc(a, f, tag)
	}

	// illegal directive field
	return a.error(errIllegalQueryField, f, "", "", "", "")
}

// analyzeQueryStructField [ ... ]
func analyzeQueryStructField(a *analysis, f *types.Var, tag string) error {
	analyzers := map[string]func(*analysis, *types.Var, string) error{
		"where":        analyzeWhereStruct,
		"join":         analyzeJoinStruct,
		"from":         analyzeJoinStruct,
		"using":        analyzeJoinStruct,
		"onconflict":   analyzeOnConflictStruct,
		"result":       analyzeResultField,
		"limit":        analyzeLimitFieldOrDirective,
		"offset":       analyzeOffsetFieldOrDirective,
		"rowsaffected": analyzeRowsAffectedField,
	}
	if afunc, ok := analyzers[tolower(f.Name())]; ok {
		return afunc(a, f, tag)
	}

	// if no match by field name, look for specific field types
	if isAccessible(a, f, a.named) {
		switch {
		case isFilterType(f.Type()):
			if err := analyzeFilterField(a, f, tag); err != nil {
				return err
			}
		case isErrorHandler(f.Type()):
			if err := analyzeErrorHandlerField(a, f, tag, false); err != nil {
				return err
			}
		case isErrorInfoHandler(f.Type()):
			if err := analyzeErrorHandlerField(a, f, tag, true); err != nil {
				return err
			}
		case typesutil.IsContext(f.Type()):
			if err := analyzeContextField(a, f, tag); err != nil {
				return err
			}
		}
	}
	return nil
}

// analyzeRelType [ ... ]
func analyzeRelType(a *analysis, rt *RelType, field *types.Var) error {
	rt.FieldMap = make(map[FieldPtr]FieldVar)
	defer func() {
		// NOTE(mkopriva): this step is necessary because of the cache.
		//
		// If there were no cache, each call to analyzeRelType would
		// traverse each field of the rt relType and could therefore
		// store the field info directly into a.info.FieldMap.
		//
		// However, because the cache is in place the fields are not traversed
		// for cached relTypes and the a.info.FieldMap is then not populated.
		for k, v := range rt.FieldMap {
			a.info.FieldMap[k] = v
		}
	}()

	ftyp := field.Type()
	cacheKey := ftyp.String()
	named, ok := ftyp.(*types.Named)
	if ok {
		ftyp = named.Underlying()
		cacheKey = named.String()
	}

	relTypeCache.RLock()
	v := relTypeCache.m[cacheKey]
	relTypeCache.RUnlock()
	if v != nil {
		*rt = *v
		return nil
	}

	// Check whether the relation field's type is an interface or a function,
	// if so, it is then expected to be an iterator, and it is analyzed as such.
	//
	// Failure of the iterator analysis will cause the whole analysis to exit
	// as there's currently no support for non-iterator interfaces nor functions.
	if iface, ok := ftyp.(*types.Interface); ok {
		var isValid bool
		if named, isValid = analyzeIteratorInterface(a, rt, iface, named); !isValid {
			return a.error(errBadIterTypeInterface, field, "", "", "", "")
		}
	} else if sig, ok := ftyp.(*types.Signature); ok {
		var isValid bool
		if named, isValid = analyzeIteratorFunction(a, rt, sig); !isValid {
			return a.error(errBadIterTypeFunc, field, "", "", "", "")
		}
	} else {
		// If not an iterator, check for slices, arrays, and pointers.
		if slice, ok := ftyp.(*types.Slice); ok { // allows []T / []*T
			ftyp = slice.Elem()
			rt.IsSlice = true
		} else if array, ok := ftyp.(*types.Array); ok { // allows [N]T / [N]*T
			ftyp = array.Elem()
			rt.IsArray = true
			rt.ArrayLen = array.Len()
		}
		if ptr, ok := ftyp.(*types.Pointer); ok { // allows *T
			ftyp = ptr.Elem()
			rt.IsPointer = true
		}

		// Get the name of the base type, if applicable.
		if rt.IsSlice || rt.IsArray || rt.IsPointer {
			if named, ok = ftyp.(*types.Named); !ok {
				// Fail if the type is a slice, an array, or a pointer
				// while its base type remains unnamed.
				return a.error(errBadRelType, field, "", "", "", "")
			}
		}
	}

	if named != nil {
		pkg := named.Obj().Pkg()
		rt.Base.Name = named.Obj().Name()
		rt.Base.PkgPath = pkg.Path()
		rt.Base.PkgName = pkg.Name()
		rt.Base.PkgLocal = pkg.Name()
		rt.Base.IsImported = isImportedType(a, named)
		rt.IsAfterScanner = typesutil.ImplementsAfterScanner(named)
		ftyp = named.Underlying()

		relTypeCache.Lock()
		relTypeCache.m[cacheKey] = rt
		relTypeCache.Unlock()
	}

	rt.Base.Kind = analyzeTypeKind(ftyp)
	if rt.Base.Kind != TypeKindStruct {
		return a.error(errBadRelType, field, "", "", "", "")
	}

	styp := ftyp.(*types.Struct)
	return analyzeFieldInfoList(a, rt, styp)
}

// analyzeFieldInfoList
func analyzeFieldInfoList(a *analysis, rt *RelType, styp *types.Struct) error {
	// The loopstate type holds the state of a loop over a struct's fields.
	type loopstate struct {
		styp     *types.Struct // the struct type whose fields are being analyzed
		typ      *TypeInfo     // info on the struct type; holds the resulting slice of analyzed FieldInfo
		idx      int           // keeps track of the field index
		pfx      string        // column prefix
		selector []*FieldSelectorNode
	}

	// LIFO stack of states used for depth first traversal of struct fields.
	stack := []*loopstate{{styp: styp, typ: &rt.Base}}

stackloop:
	for len(stack) > 0 {
		loop := stack[len(stack)-1]
		for loop.idx < loop.styp.NumFields() {
			ftag := loop.styp.Tag(loop.idx)
			fvar := loop.styp.Field(loop.idx)
			tag := tagutil.New(ftag)
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
				fvar.Name() == "_" ||
				// - it's unexported and the field's struct type is imported
				(!fvar.Exported() && loop.typ.IsImported) {
				continue
			}

			f := new(FieldInfo)
			f.Tag = tag
			f.Name = fvar.Name()
			f.IsEmbedded = fvar.Embedded()
			f.IsExported = fvar.Exported()
			rt.FieldMap[f] = FieldVar{Var: fvar, Tag: ftag}

			// Analyze the field's type.
			ftyp := fvar.Type()
			f.Type, ftyp = analyzeTypeInfo(a, ftyp)

			// If the field's type is a struct and the `sql` tag's
			// value starts with the ">" (descend) marker, then it is
			// considered to be a "parent" field element whose child
			// fields need to be analyzed as well.
			if f.Type.Is(TypeKindStruct) && strings.HasPrefix(sqltag, ">") {
				loop2 := new(loopstate)
				loop2.styp = ftyp.(*types.Struct)
				loop2.typ = &f.Type
				loop2.pfx = loop.pfx + strings.TrimPrefix(sqltag, ">")

				// Allocate selector of the appropriate size an copy it.
				loop2.selector = make([]*FieldSelectorNode, len(loop.selector))
				_ = copy(loop2.selector, loop.selector)

				// If the parent node is a pointer to a struct,
				// get the struct type info.
				typ := f.Type
				if typ.Kind == TypeKindPtr {
					typ = *typ.Elem
				}

				node := new(FieldSelectorNode)
				node.Name = f.Name
				node.Tag = f.Tag
				node.IsEmbedded = f.IsEmbedded
				node.IsExported = f.IsExported
				node.TypeName = typ.Name
				node.TypePkgPath = typ.PkgPath
				node.TypePkgName = typ.PkgName
				node.TypePkgLocal = typ.PkgLocal
				node.IsImported = typ.IsImported
				node.IsPointer = (f.Type.Kind == TypeKindPtr)
				node.ReadOnly = tag.HasOption("sql", "ro")
				node.WriteOnly = tag.HasOption("sql", "wo")
				loop2.selector = append(loop2.selector, node)

				stack = append(stack, loop2)
				continue stackloop
			}

			// Resolve the column id.
			cid, ecode, eval := parseColIdent(a, loop.pfx+sqltag)
			if ecode > 0 {
				return a.error(ecode, fvar, "", ftag, "", eval)
			}

			// TODO check the the chan, func, and interface type
			// in association with the write/read?

			// If the field is not a struct to be descended,
			// it is considered to be a "leaf" field and as
			// such the analysis of leaf-specific information
			// needs to be carried out.
			f.ColIdent = cid
			f.Selector = loop.selector
			f.NullEmpty = tag.HasOption("sql", "nullempty")
			f.ReadOnly = tag.HasOption("sql", "ro")
			f.WriteOnly = tag.HasOption("sql", "wo")
			f.UseAdd = tag.HasOption("sql", "add")
			f.UseDefault = tag.HasOption("sql", "default")
			f.UseCoalesce, f.CoalesceValue = parseCoalesceInfo(tag)

			if err := parseFilterColumnKey(a, f); err != nil {
				return err
			}

			// Add the field to the list.
			rt.Fields = append(rt.Fields, f)
			a.info.FieldMap[f] = FieldVar{Var: fvar, Tag: ftag}
		}
		stack = stack[:len(stack)-1]
	}
	return nil
}

// analyzeTypeInfo function analyzes the given type and returns the result. The analysis
// looks only for information of "named types" and in case of slice, array, map, or
// pointer types it will analyze the element type of those types. The second return
// value is the types.Type representation of the base element type of the given type.
func analyzeTypeInfo(a *analysis, tt types.Type) (typ TypeInfo, base types.Type) {
	base = tt

	if named, ok := base.(*types.Named); ok {
		pkg := named.Obj().Pkg()
		typ.Name = named.Obj().Name()
		typ.PkgPath = pkg.Path()
		typ.PkgName = pkg.Name()
		typ.PkgLocal = pkg.Name()
		typ.IsImported = isImportedType(a, named)
		typ.IsScanner = typesutil.ImplementsScanner(named)
		typ.IsValuer = typesutil.ImplementsValuer(named)
		typ.IsJSONMarshaler = typesutil.ImplementsJSONMarshaler(named)
		typ.IsJSONUnmarshaler = typesutil.ImplementsJSONUnmarshaler(named)
		typ.IsXMLMarshaler = typesutil.ImplementsXMLMarshaler(named)
		typ.IsXMLUnmarshaler = typesutil.ImplementsXMLUnmarshaler(named)
		base = named.Underlying()
	}

	typ.Kind = analyzeTypeKind(base)

	var elem TypeInfo // element info
	switch T := base.(type) {
	case *types.Basic:
		typ.IsRune = T.Name() == "rune"
		typ.IsByte = T.Name() == "byte"
	case *types.Slice:
		elem, base = analyzeTypeInfo(a, T.Elem())
		typ.Elem = &elem
	case *types.Array:
		elem, base = analyzeTypeInfo(a, T.Elem())
		typ.Elem = &elem
		typ.ArrayLen = T.Len()
	case *types.Map:
		key, _ := analyzeTypeInfo(a, T.Key())
		elem, base = analyzeTypeInfo(a, T.Elem())
		typ.Key = &key
		typ.Elem = &elem
	case *types.Pointer:
		elem, base = analyzeTypeInfo(a, T.Elem())
		typ.Elem = &elem
	case *types.Interface:
		typ.IsEmptyInterface = typesutil.IsEmptyInterface(T)
		// If base is an unnamed interface type check at least whether
		// or not it declares, or embeds, one of the relevant methods.
		if typ.Name == "" {
			typ.IsScanner = typesutil.IsScanner(T)
			typ.IsValuer = typesutil.IsValuer(T)
		}
	}
	return typ, base
}

// analyzeIteratorInterface [ ... ]
func analyzeIteratorInterface(a *analysis, rt *RelType, iface *types.Interface, named *types.Named) (out *types.Named, isValid bool) {
	if iface.NumExplicitMethods() != 1 {
		return nil, false
	}

	mth := iface.ExplicitMethod(0)
	if !isAccessible(a, mth, named) {
		return nil, false
	}

	sig := mth.Type().(*types.Signature)
	out, isValid = analyzeIteratorFunction(a, rt, sig)
	if !isValid {
		return nil, false
	}

	rt.IterMethod = mth.Name()
	return out, true
}

// analyzeIteratorFunction [ ... ]
func analyzeIteratorFunction(a *analysis, rt *RelType, sig *types.Signature) (out *types.Named, isValid bool) {
	// Must take 1 argument and return one value of type error. "func(T) error"
	if sig.Params().Len() != 1 || sig.Results().Len() != 1 || !typesutil.IsError(sig.Results().At(0).Type()) {
		return nil, false
	}

	typ := sig.Params().At(0).Type()
	if ptr, ok := typ.(*types.Pointer); ok { // allows *T
		typ = ptr.Elem()
		rt.IsPointer = true
	}

	// Make sure that the argument type is a named struct type.
	named, ok := typ.(*types.Named)
	if !ok {
		return nil, false
	} else if _, ok := named.Underlying().(*types.Struct); !ok {
		return nil, false
	}

	rt.IsIter = true
	return named, true
}

////////////////////////////////////////////////////////////////////////////////
// Where Struct Analysis
//

// analyzeWhereStruct
func analyzeWhereStruct(a *analysis, f *types.Var, tag string) (err error) {
	if !a.query.Kind.isSelect() && a.query.Kind != QueryKindUpdate && a.query.Kind != QueryKindDelete {
		return a.error(errIllegalQueryField, f, "", tag, "", "")
	}
	if a.query.Kind == QueryKindUpdate && a.query.Rel.Type.IsSlice {
		return a.error(errIllegalSliceUpdateModifier, f, "", tag, "", "")
	}
	if a.query.All != nil || a.query.Where != nil || a.query.Filter != nil {
		return a.error(errConflictingWhere, f, "", tag, "", "")
	}

	ns, err := typesutil.GetStruct(f)
	if err != nil { // fails only if non struct
		return a.error(errBadFieldTypeStruct, f, "", tag, "", "")
	}

	// The loopstate type holds the state of a loop over a struct's fields.
	type loopstate struct {
		where *WhereStruct
		items []WhereItem
		ns    *typesutil.NamedStruct // the struct type of the WhereStruct
		idx   int                    // keeps track of the field index
	}

	// root holds the reference to the root level search conditions
	root := &loopstate{ns: ns}
	// LIFO stack of states used for depth first traversal of struct fields.
	stack := []*loopstate{root}

stackloop:
	for len(stack) > 0 {
		loop := stack[len(stack)-1]
		for loop.idx < loop.ns.Struct.NumFields() {
			fvar := loop.ns.Struct.Field(loop.idx)
			ftag := loop.ns.Struct.Tag(loop.idx)
			tag := tagutil.New(ftag)
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
			if fvar.Name() != "_" && !isAccessible(a, fvar, ns.Named) {
				continue
			}

			// Analyze the bool operation for any but the first
			// item in a WhereStruct. Fail if a value was provided
			// but it is not "or" nor "and".
			if len(loop.items) > 0 {
				item := new(WhereBoolTag)
				item.Value = BoolAnd // default to "and"
				if val := tolower(tag.First("bool")); len(val) > 0 {
					if val == "or" {
						item.Value = BoolOr
					} else if val != "and" {
						return a.error(errBadBoolTagValue, fvar, f.Name(), ftag, "", val)
					}
				}
				loop.items = append(loop.items, item)
			}

			// Nested wherefields are marked with ">" and should be
			// analyzed before any other fields in the current block.
			if sqltag == ">" {
				ns, err := typesutil.GetStruct(fvar)
				if err != nil {
					return a.error(errBadFieldTypeStruct, fvar, f.Name(), "", "", "")
				}

				loop2 := new(loopstate)
				loop2.ns = ns
				loop2.where = new(WhereStruct)
				loop2.where.FieldName = fvar.Name()
				loop.items = append(loop.items, loop2.where)

				a.info.FieldMap[loop2.where] = FieldVar{Var: fvar, Tag: ftag}
				stack = append(stack, loop2)
				continue stackloop
			}

			lhs, op, op2, rhs := parsePredicateExpr(sqltag)

			// Analyze directive where item.
			if fvar.Name() == "_" {
				if !typesutil.IsDirective("Column", fvar.Type()) {
					continue
				}

				// If the expression in a gosql.Column tag's value
				// contains a right-hand-side, it is expected to be
				// either another column or a value-literal to which
				// the main column should be compared.
				if len(rhs) > 0 {
					cid, ecode, eval := parseColIdent(a, lhs)
					if ecode > 0 {
						return a.error(ecode, fvar, f.Name(), ftag, sqltag, eval)
					}

					item := new(WhereColumnDirective)
					item.LHSColIdent = cid
					item.Predicate = stringToPredicate[op]
					item.Quantifier = stringToQuantifier[op2]

					if cid, ecode, eval := parseColIdent(a, rhs); ecode > 0 {
						if ecode != errBadColIdTagValue {
							return a.error(ecode, fvar, f.Name(), ftag, sqltag, eval)
						}
						// assume literal expression
						item.RHSLiteral = rhs
					} else {
						item.RHSColIdent = cid
					}

					if item.Predicate.IsUnary() {
						return a.error(errIllegalUnaryPredicate, fvar, f.Name(), ftag, sqltag, op)
					} else if item.Quantifier > 0 && !item.Predicate.CanQuantify() {
						return a.error(errIllegalPredicateQuantifier, fvar, f.Name(), ftag, sqltag, op2)
					}

					a.info.FieldMap[item] = FieldVar{Var: fvar, Tag: ftag}
					loop.items = append(loop.items, item)
					continue
				}

				// Assume column with unary predicate.
				cid, ecode, eval := parseColIdent(a, lhs)
				if ecode > 0 {
					return a.error(ecode, fvar, f.Name(), ftag, sqltag, eval)
				}
				// If no operator was provided, default to "istrue"
				if len(op) == 0 {
					op = "istrue"
				}

				item := new(WhereColumnDirective)
				item.LHSColIdent = cid
				item.Predicate = stringToPredicate[op]

				if !item.Predicate.IsUnary() {
					return a.error(errBadDirectiveBooleanExpr, fvar, f.Name(), ftag, "", sqltag)
				} else if len(op2) > 0 {
					return a.error(errIllegalPredicateQuantifier, fvar, f.Name(), ftag, sqltag, op2)
				}

				a.info.FieldMap[item] = FieldVar{Var: fvar, Tag: ftag}
				loop.items = append(loop.items, item)
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
					return a.error(errIllegalPredicateQuantifier, fvar, f.Name(), ftag, sqltag, op2) // TODO test
				}

				ns, err := typesutil.GetStruct(fvar)
				if err != nil {
					return a.error(errBadBetweenPredicate, fvar, f.Name(), ftag, "", "")
				} else if ns.Struct.NumFields() != 2 {
					return a.error(errBadBetweenPredicate, fvar, f.Name(), ftag, "", "")
				}

				var lower, upper RangeBound
				for i := 0; i < 2; i++ {
					f := fvar // for access to the parent "between" struct

					fvar := ns.Struct.Field(i)
					ftag := ns.Struct.Tag(i)
					tag := tagutil.New(ftag)

					if fvar.Name() == "_" && typesutil.IsDirective("Column", fvar.Type()) {
						cid, ecode, eval := parseColIdent(a, tag.First("sql"))
						if ecode > 0 {
							return a.error(ecode, fvar, f.Name(), ftag, tag.First("sql"), eval)
						}

						item := new(BetweenColumnDirective)
						item.ColIdent = cid
						if v := tolower(tag.Second("sql")); v == "x" || v == "lower" {
							lower = item
						} else if v == "y" || v == "upper" {
							upper = item
						}
						a.info.FieldMap[item] = FieldVar{Var: fvar, Tag: ftag}

					} else if isAccessible(a, fvar, ns.Named) {
						item := new(BetweenStructField)
						item.Name = fvar.Name()
						item.Type, _ = analyzeTypeInfo(a, fvar.Type())

						if v := tolower(tag.First("sql")); v == "x" || v == "lower" {
							lower = item
						} else if v == "y" || v == "upper" {
							upper = item
						}
						a.info.FieldMap[item] = FieldVar{Var: fvar, Tag: ftag}
					}
				}

				if lower == nil || upper == nil {
					return a.error(errBadBetweenPredicate, fvar, f.Name(), ftag, "", "")
				}

				cid, ecode, eval := parseColIdent(a, lhs)
				if ecode > 0 {
					return a.error(ecode, fvar, f.Name(), ftag, sqltag, eval)
				}

				item := new(WhereBetweenStruct)
				item.FieldName = fvar.Name()
				item.ColIdent = cid
				item.Predicate = stringToPredicate[op]
				item.LowerBound = lower
				item.UpperBound = upper

				a.info.FieldMap[item] = FieldVar{Var: fvar, Tag: ftag}
				loop.items = append(loop.items, item)
				continue
			}

			// Analyze field where item.
			cid, ecode, eval := parseColIdent(a, lhs)
			if ecode > 0 {
				return a.error(ecode, fvar, f.Name(), ftag, lhs, eval)
			}
			// If no predicate was provided default to "="
			if len(op) == 0 {
				op = "="
			}

			item := new(WhereStructField)
			item.Name = fvar.Name()
			item.Type, _ = analyzeTypeInfo(a, fvar.Type())
			item.ColIdent = cid
			item.Predicate = stringToPredicate[op]
			item.Quantifier = stringToQuantifier[op2]
			item.FuncName = parseFuncName(tag["sql"][1:])

			if item.Predicate.IsUnary() {
				return a.error(errIllegalUnaryPredicate, fvar, f.Name(), ftag, sqltag, op)
			} else if item.Quantifier > 0 && !item.Predicate.CanQuantify() {
				return a.error(errIllegalPredicateQuantifier, fvar, f.Name(), ftag, sqltag, op2)
			} else if item.Quantifier > 0 && !item.Type.IsSequence() {
				return a.error(errIllegalFieldQuantifier, fvar, f.Name(), ftag, sqltag, op2)
			} else if item.Predicate.IsArray() && !item.Type.IsSequence() {
				return a.error(errIllegalListPredicate, fvar, f.Name(), ftag, sqltag, op)
			}

			a.info.FieldMap[item] = FieldVar{Var: fvar, Tag: ftag}
			loop.items = append(loop.items, item)
		}

		if loop.where != nil {
			loop.where.Items = loop.items
		}

		stack = stack[:len(stack)-1]
	}

	a.query.Where = new(WhereStruct)
	a.query.Where.FieldName = f.Name()
	a.query.Where.Items = root.items
	a.info.FieldMap[a.query.Where] = FieldVar{Var: f, Tag: tag}

	// XXX if a.info.TypeName == "DeleteWithUsingJoinBlock1Query" {
	// XXX 	log.Printf("%#v\n", a.query.Where.Items)
	// XXX }
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Join Struct Analysis
//

// analyzeJoinStruct
func analyzeJoinStruct(a *analysis, f *types.Var, tag string) (err error) {
	fname := tolower(f.Name())
	if fname == "join" && !a.query.Kind.isSelect() {
		return a.error(errIllegalQueryField, f, "", tag, "", "")
	} else if fname == "from" && a.query.Kind != QueryKindUpdate {
		return a.error(errIllegalQueryField, f, "", tag, "", "")
	} else if fname == "using" && a.query.Kind != QueryKindDelete {
		return a.error(errIllegalQueryField, f, "", tag, "", "")
	}

	ns, err := typesutil.GetStruct(f)
	if err != nil {
		return a.error(errBadFieldTypeStruct, f, "", tag, "", "")
	}

	join := new(JoinStruct)
	join.FieldName = f.Name()

	for i := 0; i < ns.Struct.NumFields(); i++ {
		ftag := ns.Struct.Tag(i)
		fvar := ns.Struct.Field(i)
		tag := tagutil.New(ftag)
		sqltag := tag.First("sql")

		if sqltag == "-" || sqltag == "" {
			continue
		}

		// In a JoinStruct all fields are expected to be directives
		// with the blank identifier as their name.
		if fvar.Name() != "_" {
			continue
		}

		switch dirName := typesutil.GetDirectiveName(fvar); tolower(dirName) {
		case "relation":
			if err := analyzeJoinStructRelationDirective(a, join, fvar, ftag); err != nil {
				return err
			}
		case "leftjoin", "rightjoin", "fulljoin", "crossjoin", "innerjoin":
			if err := analyzeJoinStructJoinDirective(a, join, dirName, fvar, ftag); err != nil {
				return err
			}
		default:
			return a.error(errIllegalStructDirective, fvar, f.Name(), ftag, "", "")
		}

	}

	a.query.Join = join
	a.info.FieldMap[a.query.Join] = FieldVar{Var: f, Tag: tag}
	return nil
}

// analyzeJoinStructRelationDirective
func analyzeJoinStructRelationDirective(a *analysis, j *JoinStruct, f *types.Var, ftag string) (err error) {
	if kind := tolower(j.FieldName); kind != "from" && kind != "using" {
		return a.error(errIllegalStructDirective, f, j.FieldName, ftag, "", "")
	} else if j.Relation != nil {
		return a.error(errConflictingRelationDirective, f, j.FieldName, ftag, "", "")
	}

	tag := tagutil.New(ftag)
	rid, ecode := parseRelIdent(tag.First("sql"))
	if ecode > 0 {
		return a.error(ecode, f, j.FieldName, ftag, "", tag.First("sql"))
	} else if ecode, errval := addToRelSpace(a, rid); ecode > 0 {
		return a.error(ecode, f, j.FieldName, ftag, "", errval)
	}

	j.Relation = new(RelationDirective)
	j.Relation.RelIdent = rid
	a.info.FieldMap[j.Relation] = FieldVar{Var: f, Tag: ftag}
	return nil
}

// analyzeJoinStructJoinDirective
func analyzeJoinStructJoinDirective(a *analysis, j *JoinStruct, dirName string, f *types.Var, ftag string) (err error) {
	tag := tagutil.New(ftag)
	rid, ecode := parseRelIdent(tag.First("sql"))
	if ecode > 0 {
		return a.error(ecode, f, j.FieldName, ftag, "", tag.First("sql"))
	} else if ecode, errval := addToRelSpace(a, rid); ecode > 0 {
		return a.error(ecode, f, j.FieldName, ftag, "", errval)
	}

	dir := new(JoinDirective)
	dir.RelIdent = rid
	dir.JoinType = stringToJoinType[dirName]

	for _, val := range tag["sql"][1:] {
		vals := strings.Split(val, ";")
		for i, val := range vals {

			// ✅ The left-hand side MUST be a valid column identifier.
			// - If the right-hand side IS present, then:
			//     ✅ The right-hand side MUST be a valid column identifier or a literal.
			//     ✅ The op MUST be present and it MUST be a binary predicate.
			//      - If op2 IS present, then:
			//         ✅ The op MUST be quantifiable.
			//         ✅ The op2 MUST be a valid quantifier.
			//      - If op2 IS NOT present, then:
			// - If the right-hand side IS NOT present, then:
			//     ✅ The op MUST be a valid unary_predicate
			//     ✅ The op2 MUST be empty
			lhs, op, op2, rhs := parsePredicateExpr(val)

			cid, ecode, eval := parseColIdent(a, lhs)
			if ecode > 0 {
				return a.error(ecode, f, j.FieldName, ftag, val, eval)
			}

			// NOTE(mkopriva): At the moment a join condition's left-hand-side
			// column MUST always reference a column of the relation being joined,
			// so to avoid confusion make sure that cid has either no qualifier or,
			// if it has one, it matches the alias of the joined table.
			//
			// TODO(mkopriva): Remove this limitation and properly handle the
			// operands regardless of which side they are positioned in.
			if len(cid.Qualifier) > 0 && (len(rid.Alias) > 0 && rid.Alias != cid.Qualifier) ||
				(len(rid.Alias) == 0 && rid.Name != cid.Qualifier) {
				return a.error(errBadJoinConditionLHS, f, j.FieldName, ftag, val, lhs)
			}

			item := new(JoinConditionTagItem)
			item.LHSColIdent = cid
			item.Predicate = stringToPredicate[op]
			item.Quantifier = stringToQuantifier[op2]

			// binary expression?
			if len(rhs) > 0 {
				if cid, ecode, eval := parseColIdent(a, rhs); ecode > 0 {
					if ecode != errBadColIdTagValue {
						return a.error(ecode, f, j.FieldName, ftag, val, eval)
					}
					// assume literal expression
					item.RHSLiteral = rhs
				} else {
					item.RHSColIdent = cid
				}

				if item.Predicate.IsUnary() {
					return a.error(errIllegalUnaryPredicate, f, j.FieldName, ftag, val, op)
				} else if item.Quantifier > 0 && !item.Predicate.CanQuantify() {
					return a.error(errIllegalPredicateQuantifier, f, j.FieldName, ftag, val, op2)
				}
			} else { // unary expression?

				// If no operator was provided, default to "istrue"
				if len(op) == 0 {
					item.Predicate = stringToPredicate["istrue"]
				}

				// TODO
				if !item.Predicate.IsUnary() {
					return a.error(errBadDirectiveBooleanExpr, f, j.FieldName, ftag, "", val)
				} else if len(op2) > 0 {
					return a.error(errIllegalPredicateQuantifier, f, j.FieldName, ftag, val, op2)
				}
			}

			if len(dir.TagItems) > 0 && i == 0 {
				dir.TagItems = append(dir.TagItems, &JoinBoolTagItem{BoolAnd})
			} else if len(dir.TagItems) > 0 && i > 0 {
				dir.TagItems = append(dir.TagItems, &JoinBoolTagItem{BoolOr})
			}
			dir.TagItems = append(dir.TagItems, item)
		}
	}

	j.Directives = append(j.Directives, dir)
	a.info.FieldMap[dir] = FieldVar{Var: f, Tag: ftag}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// On Conflict Struct Analysis
//

// analyzeOnConflictStruct analyzes the given field as an "onconflict" struct.
// The structTag argument is used for error reporting.
//
// ✅ The kind of the target query MUST be "insert".
// ✅ The type of the given field MUST be a struct type.
// ✅ The struct type MUST contain exactly 1 "conflict_action" directive.
// ✅ The struct type MUST contain exactly 1 "conflict_target" directive, if it
//    contains the gosql.Update "conflict_action" directive.
// ✅ The struct type MAY contain, at most, 1 "conflict_target" directive, if it
//    contains the gosql.Ignore "conflict_action" directive.
func analyzeOnConflictStruct(a *analysis, f *types.Var, tag string) (err error) {
	if a.query.Kind != QueryKindInsert {
		return a.error(errIllegalQueryField, f, "", tag, "", "")
	}
	ns, err := typesutil.GetStruct(f)
	if err != nil {
		return a.error(errBadFieldTypeStruct, f, "", tag, "", "")
	}

	onConflict := new(OnConflictStruct)
	onConflict.FieldName = f.Name()
	for i := 0; i < ns.Struct.NumFields(); i++ {
		fvar := ns.Struct.Field(i)
		ftag := ns.Struct.Tag(i)

		// In an OnConflictStruct all fields are expected to be directives
		// with the blank identifier as their name.
		if fvar.Name() != "_" {
			continue
		}

		switch tolower(typesutil.GetDirectiveName(fvar)) {
		case "column":
			if err = analyzeOnConflictColumnDirective(a, onConflict, fvar, ftag); err != nil {
				return err
			}
		case "index":
			if err = analyzeOnConflictIndexDirective(a, onConflict, fvar, ftag); err != nil {
				return err
			}
		case "constraint":
			if err = analyzeOnConflictConstraintDirective(a, onConflict, fvar, ftag); err != nil {
				return err
			}
		case "ignore":
			if err = analyzeOnConflictIgnoreDirective(a, onConflict, fvar, ftag); err != nil {
				return err
			}
		case "update":
			if err = analyzeOnConflictUpdateDirective(a, onConflict, fvar, ftag); err != nil {
				return err
			}
		default:
			return a.error(errIllegalStructDirective, fvar, f.Name(), ftag, "", "")
		}

	}
	if onConflict.Update != nil && (onConflict.Column == nil && onConflict.Index == nil && onConflict.Constraint == nil) {
		return a.error(errMissingOnConflictTarget, f, "", tag, "", "")
	}

	a.query.OnConflict = onConflict
	a.info.FieldMap[onConflict] = FieldVar{Var: f, Tag: tag}
	return nil
}

// analyzeOnConflictColumnDirective analyzes the given field and its associated
// tag as a "gosql.Column" directive.
//
// ✅ The given OnConflictStruct MUST NOT have any other "conflict_target" fields set.
// ✅ The tag MUST contain a valid identifier.
func analyzeOnConflictColumnDirective(a *analysis, oc *OnConflictStruct, f *types.Var, tag string) (err error) {
	if oc.Column != nil || oc.Index != nil || oc.Constraint != nil {
		return a.error(errConflictingOnConfictTarget, f, oc.FieldName, tag, "", "")
	}

	slice := tagutil.New(tag)["sql"]
	ids, ecode, eval := parseColIdents(a, slice)
	if ecode > 0 {
		return a.error(ecode, f, oc.FieldName, tag, "", eval)
	}

	oc.Column = new(ColumnDirective)
	oc.Column.ColIdents = ids
	a.info.FieldMap[oc.Column] = FieldVar{Var: f, Tag: tag}
	return nil
}

// analyzeOnConflictIndexDirective analyzes the given field and its associated
// tag as a "gosql.Index" directive.
//
// ✅ The given OnConflictStruct MUST NOT have any other "conflict_target" fields set.
// ✅ The tag MUST contain a valid identifier.
func analyzeOnConflictIndexDirective(a *analysis, oc *OnConflictStruct, f *types.Var, tag string) (err error) {
	if oc.Column != nil || oc.Index != nil || oc.Constraint != nil {
		return a.error(errConflictingOnConfictTarget, f, oc.FieldName, tag, "", "")
	}

	name := tagutil.New(tag).First("sql")
	if !rxIdent.MatchString(name) {
		return a.error(errBadIdentTagValue, f, oc.FieldName, tag, "", "")
	}

	oc.Index = new(IndexDirective)
	oc.Index.Name = name
	a.info.FieldMap[oc.Index] = FieldVar{Var: f, Tag: tag}
	return nil
}

// analyzeOnConflictConstraintDirective analyzes the given field and its associated
// tag as a "gosql.Constraint" directive.
//
// ✅ The given OnConflictStruct MUST NOT have any other "conflict_target" fields set.
// ✅ The tag MUST contain a valid identifier.
func analyzeOnConflictConstraintDirective(a *analysis, oc *OnConflictStruct, f *types.Var, tag string) (err error) {
	if oc.Column != nil || oc.Index != nil || oc.Constraint != nil {
		return a.error(errConflictingOnConfictTarget, f, oc.FieldName, tag, "", "")
	}

	name := tagutil.New(tag).First("sql")
	if !rxIdent.MatchString(name) {
		return a.error(errBadIdentTagValue, f, oc.FieldName, tag, "", "")
	}

	oc.Constraint = new(ConstraintDirective)
	oc.Constraint.Name = name
	a.info.FieldMap[oc.Constraint] = FieldVar{Var: f, Tag: tag}
	return nil
}

// analyzeOnConflictIgnoreDirective analyzes the given field as a "gosql.Ignore" directive.
//
// ✅ The given OnConflictStruct MUST NOT have any other "conflict_action" fields set.
func analyzeOnConflictIgnoreDirective(a *analysis, oc *OnConflictStruct, f *types.Var, tag string) (err error) {
	if oc.Ignore != nil || oc.Update != nil {
		return a.error(errConflictingOnConfictAction, f, oc.FieldName, tag, "", "")
	}

	oc.Ignore = new(IgnoreDirective)
	a.info.FieldMap[oc.Ignore] = FieldVar{Var: f, Tag: tag}
	return nil
}

// analyzeOnConflictUpdateDirective analyzes the given field and its associated
// tag as a "gosql.Update" directive.
//
// ✅ The given OnConflictStruct MUST NOT have any other "conflict_action" fields set.
func analyzeOnConflictUpdateDirective(a *analysis, oc *OnConflictStruct, f *types.Var, tag string) (err error) {
	if oc.Ignore != nil || oc.Update != nil {
		return a.error(errConflictingOnConfictAction, f, oc.FieldName, tag, "", "")
	}

	slice := tagutil.New(tag)["sql"]
	list, ecode, eval := parseColIdentList(a, slice)
	if ecode > 0 {
		return a.error(ecode, f, oc.FieldName, tag, "", eval)
	}

	oc.Update = new(UpdateDirective)
	oc.Update.ColIdentList = list
	a.info.FieldMap[oc.Update] = FieldVar{Var: f, Tag: tag}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Plain Field Analysis
//

// analyzeLimitFieldOrDirective analyzes the given field, which is expected to be either
// the gosql.Limit directive or a plain integer field. The tag argument, if not
// empty, is expected to hold a positive integer.
func analyzeLimitFieldOrDirective(a *analysis, f *types.Var, tag string) error {
	if !a.query.Kind.isSelect() {
		return a.error(errIllegalQueryField, f, "", tag, "", "")
	}
	if a.query.Limit != nil {
		return a.error(errConflictingFieldOrDirective, f, "", tag, "", "")
	}

	val := tagutil.New(tag).First("sql")
	limit := new(LimitField)
	if name := f.Name(); name != "_" {
		if !isIntegerType(f.Type()) {
			return a.error(errBadFieldTypeInt, f, "", tag, "", "")
		}
		limit.Name = name
	} else if len(val) == 0 {
		return a.error(errMissingTagValue, f, "", tag, "", "")
	}

	if len(val) > 0 {
		u64, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return a.error(errBadUIntegerTagValue, f, "", tag, "", val)
		}
		limit.Value = u64
	}

	a.query.Limit = limit
	a.info.FieldMap[a.query.Limit] = FieldVar{Var: f, Tag: tag}
	return nil
}

// analyzeOffsetFieldOrDirective analyzes the given field, which is expected to be either
// the gosql.Offset directive or a plain integer field. The tag argument,
// if not empty, is expected to hold a positive integer.
func analyzeOffsetFieldOrDirective(a *analysis, f *types.Var, tag string) error {
	if !a.query.Kind.isSelect() {
		return a.error(errIllegalQueryField, f, "", tag, "", "")
	}
	if a.query.Offset != nil {
		return a.error(errConflictingFieldOrDirective, f, "", tag, "", "")
	}

	val := tagutil.New(tag).First("sql")
	offset := new(OffsetField)
	if name := f.Name(); name != "_" {
		if !isIntegerType(f.Type()) {
			return a.error(errBadFieldTypeInt, f, "", tag, "", "")
		}
		offset.Name = name
	} else if len(val) == 0 {
		return a.error(errMissingTagValue, f, "", tag, "", "")
	}

	if len(val) > 0 {
		u64, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return a.error(errBadUIntegerTagValue, f, "", tag, "", val)
		}
		offset.Value = u64
	}

	a.query.Offset = offset
	a.info.FieldMap[a.query.Offset] = FieldVar{Var: f, Tag: tag}
	return nil
}

func analyzeErrorHandlerField(a *analysis, f *types.Var, tag string, isInfo bool) error {
	if a.query.ErrorHandler != nil {
		return a.error(errConflictingFieldOrDirective, f, "", tag, "", "")
	}

	a.query.ErrorHandler = new(ErrorHandlerField)
	a.query.ErrorHandler.Name = f.Name()
	a.query.ErrorHandler.IsInfo = isInfo
	a.info.FieldMap[a.query.ErrorHandler] = FieldVar{Var: f, Tag: tag}
	return nil
}

func analyzeContextField(a *analysis, f *types.Var, tag string) error {
	if a.query.Context != nil {
		return a.error(errConflictingFieldOrDirective, f, "", tag, "", "")
	}

	a.query.Context = new(ContextField)
	a.query.Context.Name = f.Name()
	a.info.FieldMap[a.query.Context] = FieldVar{Var: f, Tag: tag}
	return nil
}

func analyzeFilterField(a *analysis, f *types.Var, tag string) error {
	if !a.query.Kind.isSelect() && a.query.Kind != QueryKindUpdate && a.query.Kind != QueryKindDelete {
		return a.error(errIllegalQueryField, f, "", tag, "", "")
	}
	if a.query.Kind == QueryKindUpdate && a.query.Rel.Type.IsSlice {
		return a.error(errIllegalSliceUpdateModifier, f, "", tag, "", "")
	}
	if a.query.All != nil || a.query.Where != nil || a.query.Filter != nil {
		return a.error(errConflictingWhere, f, "", tag, "", "")
	}

	a.query.Filter = new(FilterField)
	a.query.Filter.Name = f.Name()
	a.info.FieldMap[a.query.Filter] = FieldVar{Var: f, Tag: tag}
	return nil
}

func analyzeFilterConstructorField(a *analysis, f *types.Var, tag string) error {
	if a.filter.FilterConstructor != nil {
		return a.error(errConflictingFilterConstructor, f, "", tag, "", "")
	}

	a.filter.FilterConstructor = new(FilterConstructorField)
	a.filter.FilterConstructor.Name = f.Name()
	a.info.FieldMap[a.filter.FilterConstructor] = FieldVar{Var: f, Tag: tag}
	return nil
}

func analyzeResultField(a *analysis, f *types.Var, tag string) error {
	if a.query.Kind != QueryKindInsert && a.query.Kind != QueryKindUpdate && a.query.Kind != QueryKindDelete {
		return a.error(errIllegalQueryField, f, "", tag, "", "")
	}
	if a.query.Return != nil || a.query.Result != nil || a.query.RowsAffected != nil {
		return a.error(errConflictingResultTarget, f, "", tag, "", "")
	}

	a.query.Result = new(ResultField)
	a.query.Result.FieldName = f.Name()
	if err := analyzeRelType(a, &a.query.Result.Type, f); err != nil {
		return err
	}

	a.info.FieldMap[a.query.Result] = FieldVar{Var: f, Tag: tag}
	return nil
}

func analyzeRowsAffectedField(a *analysis, f *types.Var, tag string) error {
	if a.query.Kind != QueryKindInsert && a.query.Kind != QueryKindUpdate && a.query.Kind != QueryKindDelete {
		return a.error(errIllegalQueryField, f, "", tag, "", "")
	}
	if a.query.Return != nil || a.query.Result != nil || a.query.RowsAffected != nil {
		return a.error(errConflictingResultTarget, f, "", tag, "", "")
	}

	ftyp := f.Type()
	if !isIntegerType(ftyp) {
		return a.error(errBadFieldTypeInt, f, "", tag, "", "")
	}

	a.query.RowsAffected = new(RowsAffectedField)
	a.query.RowsAffected.Name = f.Name()
	a.query.RowsAffected.TypeKind = analyzeTypeKind(ftyp)
	a.info.FieldMap[a.query.RowsAffected] = FieldVar{Var: f, Tag: tag}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Directive Fields Analysis
//

// analyzeOrderByDirective
func analyzeOrderByDirective(a *analysis, f *types.Var, tag string) (err error) {
	if !a.query.Kind.isSelect() {
		return a.error(errIllegalQueryField, f, "", tag, "", "")
	}

	tags := tagutil.New(tag)["sql"]
	if len(tags) == 0 {
		return a.error(errMissingTagColumnList, f, "", tag, "", "")
	}

	var items []OrderByTagItem
	for _, val := range tags {
		val = strings.TrimSpace(val)
		if len(val) == 0 {
			continue
		}

		var item OrderByTagItem
		if val[0] == '-' {
			item.Direction = OrderDesc
			val = val[1:]
		}
		if i := strings.Index(val, ":"); i > -1 {
			if val[i+1:] == "nullsfirst" {
				item.Nulls = NullsFirst
			} else if val[i+1:] == "nullslast" {
				item.Nulls = NullsLast
			} else {
				return a.error(errBadNullsOrderTagValue, f, "", val, "", val[i+1:])
			}
			val = val[:i]
		}

		cid, ecode, eval := parseColIdent(a, val)
		if ecode > 0 {
			return a.error(ecode, f, "", tag, val, eval)
		}

		item.ColIdent = cid
		items = append(items, item)
	}

	a.query.OrderBy = new(OrderByDirective)
	a.query.OrderBy.Items = items
	a.info.FieldMap[a.query.OrderBy] = FieldVar{Var: f, Tag: tag}
	return nil
}

// analyzeOverrideDirective
func analyzeOverrideDirective(a *analysis, f *types.Var, tag string) (err error) {
	if a.query.Kind != QueryKindInsert {
		return a.error(errIllegalQueryField, f, "", tag, "", "")
	}

	var kind OverridingKind
	switch val := tolower(tagutil.New(tag).First("sql")); val {
	case "system":
		kind = OverridingSystem
	case "user":
		kind = OverridingUser
	default:
		return a.error(errBadOverrideTagValue, f, "", tag, "", val)
	}

	a.query.Override = new(OverrideDirective)
	a.query.Override.Kind = kind
	a.info.FieldMap[a.query.Override] = FieldVar{Var: f, Tag: tag}
	return nil
}

func analyzeReturnDirective(a *analysis, f *types.Var, tag string) error {
	if len(a.query.Rel.Type.Fields) == 0 {
		return a.error(errMissingRelField, f, "", tag, "", "") // TODO test
	}
	if a.query.Kind != QueryKindInsert && a.query.Kind != QueryKindUpdate && a.query.Kind != QueryKindDelete {
		return a.error(errIllegalQueryField, f, "", tag, "", "")
	}
	if a.query.Return != nil || a.query.Result != nil || a.query.RowsAffected != nil {
		return a.error(errConflictingResultTarget, f, "", tag, "", "")
	}

	t := tagutil.New(tag)
	list, ecode, eval := parseColIdentList(a, t["sql"])
	if ecode > 0 {
		return a.error(ecode, f, "", tag, t.Get("sql"), eval)
	}

	// Make sure that the column ids have a matching field.
	for _, id := range list.Items {
		if !a.query.Rel.Type.HasFieldWithColumn(id.Name) {
			return a.error(errColumnFieldUnknown, f, "", tag, t.Get("sql"), id.String())
		}
	}

	a.query.Return = new(ReturnDirective)
	a.query.Return.ColIdentList = list
	a.info.FieldMap[a.query.Return] = FieldVar{Var: f, Tag: tag}
	return nil
}

func analyzeAllDirective(a *analysis, f *types.Var, tag string) error {
	if a.query.Kind != QueryKindUpdate && a.query.Kind != QueryKindDelete {
		return a.error(errIllegalQueryField, f, "", tag, "", "")
	}
	if a.query.Kind == QueryKindUpdate && a.query.Rel.Type.IsSlice {
		return a.error(errIllegalSliceUpdateModifier, f, "", tag, "", "")
	}
	if a.query.All != nil || a.query.Where != nil || a.query.Filter != nil {
		return a.error(errConflictingWhere, f, "", tag, "", "")
	}

	a.query.All = new(AllDirective)
	a.info.FieldMap[a.query.All] = FieldVar{Var: f, Tag: tag}
	return nil
}

func analyzeDefaultDirective(a *analysis, f *types.Var, tag string) error {
	if a.query.Kind != QueryKindInsert && a.query.Kind != QueryKindUpdate {
		return a.error(errIllegalQueryField, f, "", tag, "", "")
	}

	t := tagutil.New(tag)
	list, ecode, eval := parseColIdentList(a, t["sql"])
	if ecode > 0 {
		return a.error(ecode, f, "", tag, t.Get("sql"), eval)
	}

	a.query.Default = new(DefaultDirective)
	a.query.Default.ColIdentList = list
	a.info.FieldMap[a.query.Default] = FieldVar{Var: f, Tag: tag}
	return nil
}

func analyzeForceDirective(a *analysis, f *types.Var, tag string) error {
	if a.query.Kind != QueryKindInsert && a.query.Kind != QueryKindUpdate {
		// TODO test
		return a.error(errIllegalQueryField, f, "", tag, "", "")
	}

	t := tagutil.New(tag)
	list, ecode, eval := parseColIdentList(a, t["sql"])
	if ecode > 0 {
		// TODO test
		return a.error(ecode, f, "", tag, t.Get("sql"), eval)
	}

	// Make sure that the column ids have a matching field.
	for _, id := range list.Items {
		if !a.query.Rel.Type.HasFieldWithColumn(id.Name) {
			// TODO test
			return a.error(errColumnFieldUnknown, f, "", tag, t.Get("sql"), id.String())
		}
	}

	a.query.Force = new(ForceDirective)
	a.query.Force.ColIdentList = list
	a.info.FieldMap[a.query.Force] = FieldVar{Var: f, Tag: tag}
	return nil
}

func analyzeOptionalDirective(a *analysis, f *types.Var, tag string) error {
	if !a.query.Kind.isSelect() {
		return a.error(errIllegalQueryField, f, "", tag, "", "")
	}

	t := tagutil.New(tag)
	list, ecode, eval := parseColIdentList(a, t["sql"])
	if ecode > 0 {
		return a.error(ecode, f, "", tag, t.Get("sql"), eval)
	}

	// Make sure that the column ids have a matching field.
	for _, id := range list.Items {
		if !a.query.Rel.Type.HasFieldWithColumn(id.Name) {
			return a.error(errColumnFieldUnknown, f, "", tag, t.Get("sql"), id.String())
		}
	}

	a.query.Optional = new(OptionalDirective)
	a.query.Optional.ColIdentList = list
	a.info.FieldMap[a.query.Optional] = FieldVar{Var: f, Tag: tag}
	return nil
}

// analyzeTextSearchDirective analyzes the given field and its tag as the gosql.TextSearch
// directive and sets the result to the given analysis' filter.
func analyzeTextSearchDirective(a *analysis, f *types.Var, tag string) error {
	tval := tagutil.New(tag).First("sql")
	tval = strings.ToLower(strings.TrimSpace(tval))

	cid, ecode, eval := parseColIdent(a, tval)
	if ecode > 0 {
		return a.error(ecode, f, "", tag, tval, eval)
	}

	a.filter.TextSearch = new(TextSearchDirective)
	a.filter.TextSearch.ColIdent = cid
	a.info.FieldMap[a.filter.TextSearch] = FieldVar{Var: f, Tag: tag}
	return nil
}

func addToRelSpace(a *analysis, id RelIdent) (ecode errorCode, errval string) {
	if a.info.RelSpace == nil {
		a.info.RelSpace = make(map[string]RelIdent)
	}
	if len(id.Alias) > 0 {
		if _, ok := a.info.RelSpace[id.Alias]; ok {
			return errConflictingRelAlias, id.Alias
		}
		a.info.RelSpace[id.Alias] = id
		return 0, ""
	}
	if _, ok := a.info.RelSpace[id.Name]; ok {
		return errConflictingRelName, id.Name
	}
	a.info.RelSpace[id.Name] = id
	return 0, ""
}

////////////////////////////////////////////////////////////////////////////////
// Misc. Analysis
//

// analyzeTypeKind returns the TypeKind for the given types.Type.
func analyzeTypeKind(typ types.Type) TypeKind {
	switch x := typ.(type) {
	case *types.Basic:
		return typesBasicKindToTypeKind[x.Kind()]
	case *types.Array:
		return TypeKindArray
	case *types.Chan:
		return TypeKindChan
	case *types.Signature:
		return TypeKindFunc
	case *types.Interface:
		return TypeKindInterface
	case *types.Map:
		return TypeKindMap
	case *types.Pointer:
		return TypeKindPtr
	case *types.Slice:
		return TypeKindSlice
	case *types.Struct:
		return TypeKindStruct
	case *types.Named:
		return analyzeTypeKind(x.Underlying())
	}
	return 0 // unsupported / unknown
}

////////////////////////////////////////////////////////////////////////////////
// Parsers
//

// parsePredicateExpr parses the given string as a predicate expression and
// returns the individual elements of that expression. The expected format is:
// { column [ predicate-type [ quantifier ] { column | literal } ] }
func parsePredicateExpr(expr string) (lhs, cop, qua, rhs string) {
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

// parseRelIdent parses the given string as a relation identifier and returns the result.
//
// ✅ The string MUST be in the expected format, which is: "[qualifier.]name[:alias]".
func parseRelIdent(val string) (id RelIdent, ecode errorCode) {
	if !rxRelIdent.MatchString(val) {
		return id, errBadRelIdTagValue
	}
	if i := strings.LastIndexByte(val, '.'); i > -1 {
		id.Qualifier = val[:i]
		val = val[i+1:]
	}
	if i := strings.LastIndexByte(val, ':'); i > -1 {
		id.Alias = val[i+1:]
		val = val[:i]
	}
	id.Name = val
	return id, 0
}

// parseColIdent parses the given string as a column identifier and returns the result.
//
// ✅ The string MUST be in the expected format, which is: "[qualifier.]name".
func parseColIdent(a *analysis, val string) (id ColIdent, ecode errorCode, eval string) {
	if !isColIdent(val) {
		return id, errBadColIdTagValue, val
	}
	if i := strings.LastIndexByte(val, '.'); i > -1 {
		id.Qualifier = val[:i]
		if _, ok := a.info.RelSpace[id.Qualifier]; !ok {
			return id, errUnknownColumnQualifier, id.Qualifier
		}
		val = val[i+1:]
	}
	id.Name = val
	return id, 0, ""
}

// parseColIdents parses the individual strings in the given slice as
// column identifiers and returns the result as []ColIdent.
//
// ✅ The individual strings MUST be in the expected format, which is: "[qualifier.]name".
func parseColIdents(a *analysis, tag []string) (ids []ColIdent, ecode errorCode, eval string) {
	if len(tag) == 0 {
		return nil, errMissingTagColumnList, ""
	}

	ids = make([]ColIdent, len(tag))
	for i, val := range tag {
		id, ecode, eval := parseColIdent(a, val)
		if ecode > 0 {
			return nil, ecode, eval
		}
		ids[i] = id
	}
	return ids, 0, ""
}

// parseColIdentList parses the individual strings in the given slice as
// column identifiers and returns the result as ColIdentList.
//
// ✅ A slice of length=1 holding a "*" string value MAY be use instead of column ids.
// ✅ The individual strings MUST be in the expected format, which is: "[qualifier.]name".
func parseColIdentList(a *analysis, tag []string) (list ColIdentList, ecode errorCode, eval string) {
	if len(tag) == 1 && tag[0] == "*" {
		list.All = true
		return list, 0, ""
	}

	items, ecode, eval := parseColIdents(a, tag)
	if ecode > 0 {
		return list, ecode, eval
	}
	list.Items = items
	return list, 0, ""
}

// parseCoalesceInfo
func parseCoalesceInfo(tag tagutil.Tag) (use bool, val string) {
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

func parseFuncName(tagvals []string) FuncName {
	for _, v := range tagvals {
		if len(v) > 0 && v[0] == '@' {
			return FuncName(strings.ToLower(v[1:]))
		}
	}
	return ""
}

// parseFilterColumnKey
func parseFilterColumnKey(a *analysis, f *FieldInfo) error {
	tag := a.cfg.FilterColumnKeyTag.Value
	sep := a.cfg.FilterColumnKeySeparator.Value
	base := a.cfg.FilterColumnKeyBase.Value

	selector := f.Selector
	fcktag := f.Tag["fck"]
	if len(fcktag) == 0 {
		for i, node := range selector {
			if fcktag = node.Tag["fck"]; len(fcktag) > 0 {
				selector = selector[i:]
				break
			}
		}
	}

	// if present, use the fck tag to override the global config
	if len(fcktag) > 0 {
		for _, opt := range fcktag {
			// TODO(mkopriva): add error reporting for invalid
			// "fck" option keys and/or values.
			var optKey, optVal string
			if i := strings.IndexByte(opt, ':'); i > -1 {
				optKey, optVal = opt[:i], opt[i+1:]
			}

			switch optKey {
			case "tag":
				tag = optVal
			case "sep":
				sep = optVal
			case "base":
				base, _ = strconv.ParseBool(optVal)
			}
		}
	}

	// use field tag
	if len(tag) > 0 {
		if !base {
			f.FilterColumnKey = joinFieldTag(f, selector, tag, sep)
		} else {
			if key := f.Tag.First(tag); key != "-" {
				f.FilterColumnKey = key
			}
		}
	} else {
		// use field name
		if !base {
			f.FilterColumnKey = joinFieldName(f, selector, sep)
		} else {
			f.FilterColumnKey = f.Name
		}
	}
	return nil
}

// joinFieldName
func joinFieldName(f *FieldInfo, sel []*FieldSelectorNode, sep string) (key string) {
	for _, node := range sel {
		key += node.Name + sep
	}
	return key + f.Name
}

// joinFieldTag
func joinFieldTag(f *FieldInfo, sel []*FieldSelectorNode, tag, sep string) (key string) {
	for _, node := range sel {
		k := node.Tag.First(tag)
		if k == "-" || k == "" {
			return ""
		}

		key += k + sep
	}

	k := f.Tag.First(tag)
	if k == "-" || k == "" {
		return ""
	}
	return key + k
}

////////////////////////////////////////////////////////////////////////////////
// Helper Methods & Functions
//

// containts reports whether or not the list contains the given column identifier.
func (cl *ColIdentList) Contains(cid ColIdent) bool {
	if cl.All {
		return true
	}

	for i := 0; i < len(cl.Items); i++ {
		if cl.Items[i].Name == cid.Name {
			return true
		}
	}
	return false
}

// isImportedType reports whether or not the given type is imported based on
// on the package in which the target of the analysis is declared.
func isImportedType(a *analysis, named *types.Named) bool {
	return named != nil && named.Obj().Pkg().Path() != a.pkgPath
}

// isAccessible reports whether or not the given value is accessible from
// the package in which the target of the analysis is declared.
func isAccessible(a *analysis, x exportable, named *types.Named) bool {
	return x.Name() != "_" && (x.Exported() || !isImportedType(a, named))
}

// exportable is implemented by both types.Var and types.Func.
type exportable interface {
	Name() string
	Exported() bool
}

// isIntegerType reports whether or not the given type is one of the basic (un)signed integer types.
func isIntegerType(typ types.Type) bool {
	basic, ok := typ.(*types.Basic)
	if !ok {
		return false
	}
	kind := basic.Kind()
	return types.Int <= kind && kind <= types.Uint64
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

// isColIdent reports whether or not the given value is a valid column identifier.
func isColIdent(val string) bool {
	return rxColIdent.MatchString(val) && !rxReserved.MatchString(val)
}

// tolower normalizes the given string by converting it to lower case and
// also trimming any extra white-space.
func tolower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

////////////////////////////////////////////////////////////////////////////////
// TypeInfo Helper Methods
//

func (t *TypeInfo) GenericLiteral() LiteralType {
	if t.Kind.IsBasic() {
		if t.IsByte {
			return "byte"
		}
		if t.IsRune {
			return "rune"
		}
		return LiteralType(typeKinds[t.Kind])
	}
	return t.literal(false, true)
}

func (t *TypeInfo) Literal() LiteralType {
	return t.literal(false, true)
}

func (t *TypeInfo) literal(pkgLocal, elidePtr bool) LiteralType {
	if len(t.Name) > 0 {
		if pkgLocal && len(t.PkgLocal) > 0 && t.PkgLocal != "." {
			return LiteralType(t.PkgLocal + "." + t.Name)
		} else if len(t.PkgName) > 0 {
			return LiteralType(t.PkgName + "." + t.Name)
		}
		return LiteralType(t.Name)
	}

	switch t.Kind {
	default: // assume builtin basic
		return LiteralType(typeKinds[t.Kind])
	case TypeKindArray:
		return LiteralType("["+strconv.FormatInt(t.ArrayLen, 10)+"]") + t.Elem.literal(pkgLocal, false)
	case TypeKindSlice:
		return "[]" + t.Elem.literal(pkgLocal, false)
	case TypeKindMap:
		return LiteralType("map["+t.Key.literal(pkgLocal, false)+"]") + t.Elem.literal(pkgLocal, false)
	case TypeKindPtr:
		if elidePtr {
			return t.Elem.literal(pkgLocal, false)
		} else {
			return "*" + t.Elem.literal(pkgLocal, false)
		}
		return "*" + t.Elem.literal(pkgLocal, false)
	case TypeKindUint8:
		if t.IsByte {
			return "byte"
		}
		return "uint8"
	case TypeKindInt32:
		if t.IsRune {
			return "rune"
		}
		return "int32"
	case TypeKindInterface:
		if t.IsEmptyInterface {
			return "interface{}"
		}
		return "<unsupported>"
	case TypeKindStruct, TypeKindChan, TypeKindFunc:
		return "<unsupported>"
	}
	return "<unknown>"
}

// Is reports whether or not t represents a type whose kind matches one of
// the provided TypeKinds or a pointer to one of the provided TypeKinds.
func (t *TypeInfo) Is(kk ...TypeKind) bool {
	for _, k := range kk {
		if t.Kind == k || (t.Kind == TypeKindPtr && t.Elem.Kind == k) {
			return true
		}
	}
	return false
}

// IsSliceKind reports whether or not t represents a slice type whose elem type
// is one of the provided TypeKinds.
func (t *TypeInfo) IsSliceKind(kk ...TypeKind) bool {
	if t.Kind == TypeKindSlice {
		for _, k := range kk {
			if t.Elem.Kind == k {
				return true
			}
		}
	}
	return false
}

// IsArray helper reports whether or not the type is of the array kind.
func (t *TypeInfo) IsArray() bool {
	return t.Kind == TypeKindArray
}

// IsSlice helper reports whether or not the type is of the slice kind.
func (t *TypeInfo) IsSlice() bool {
	return t.Kind == TypeKindSlice
}

// IsSequence helper reports whether or not the type is of the slice or array kind.
func (t *TypeInfo) IsSequence() bool {
	return t.IsSlice() || t.IsArray()
}

// isNilable reports whether or not t represents a type that can be nil.
func (t *TypeInfo) IsNilable() bool {
	return t.Is(TypeKindPtr, TypeKindSlice, TypeKindArray, TypeKindMap, TypeKindInterface)
}

// Indicates whether or not the MarshalJSON method can be called on the type.
func (t *TypeInfo) ImplementsJSONMarshaler() bool {
	return t.IsJSONMarshaler || (t.Kind == TypeKindPtr && t.Elem.IsJSONMarshaler)
}

// Indicates whether or not the UnmarshalJSON method can be called on the type.
func (t *TypeInfo) ImplementsJSONUnmarshaler() bool {
	return t.IsJSONUnmarshaler || (t.Kind == TypeKindPtr && t.Elem.IsJSONUnmarshaler)
}

// Indicates whether or not the MarshalXML method can be called on the type.
func (t *TypeInfo) ImplementsXMLMarshaler() bool {
	return t.IsXMLMarshaler || (t.Kind == TypeKindPtr && t.Elem.IsXMLMarshaler)
}

// Indicates whether or not the UnmarshalXML method can be called on the type.
func (t *TypeInfo) ImplementsXMLUnmarshaler() bool {
	return t.IsXMLUnmarshaler || (t.Kind == TypeKindPtr && t.Elem.IsXMLUnmarshaler)
}

// Indicates whether or not an instance of the type's Kind is illegal to be used with encoding/json.
func (t *TypeInfo) IsJSONIllegal() bool {
	return t.Is(TypeKindChan, TypeKindFunc, TypeKindComplex64, TypeKindComplex128)
}

// Indicates whether or not an instance of the type's Kind is illegal to be used with encoding/xml.
func (t *TypeInfo) IsXMLIllegal() bool {
	return t.Is(TypeKindChan, TypeKindFunc, TypeKindMap)
}

////////////////////////////////////////////////////////////////////////////////
// Cache
//
var relTypeCache = struct {
	sync.RWMutex
	m map[string]*RelType
}{m: make(map[string]*RelType)}

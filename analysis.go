package gosql

import (
	"go/types"
	"log"
	"regexp"
	"strings"

	"github.com/frk/gosql/internal/typesutil"
	"github.com/frk/tagutil"
)

// TODO(mkopriva): to provide more detailed error messages either pass in the
// details about the file being analyzed, or make sure that the caller has that
// information and appends it to the error.
func analyze(named *types.Named) (*command, error) {
	a := new(analyzer)
	a.pkg = named.Obj().Pkg().Path()
	a.cmd = &command{name: named.Obj().Name()}

	var ok bool
	if a.cmdtyp, ok = named.Underlying().(*types.Struct); !ok {
		typ := named.Underlying().String()
		return nil, newerr(errBadCmdType, a.cmd.name, typ)
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
		return nil, newerr(errBadCmdName, a.cmd.name)
	}

	if err := a.run(); err != nil {
		return nil, err
	}
	return a.cmd, nil
}

// analyzer holds the state of the analysis
type analyzer struct {
	pkg    string        // the package path of the command under analysis
	cmdtyp *types.Struct // the struct type of the command under analysis
	reltyp *types.Struct // the struct type of the relation under analysis
	cmd    *command      // the result of the analysis
}

func (a *analyzer) run() (err error) {
	for i := 0; i < a.cmdtyp.NumFields(); i++ {
		fld := a.cmdtyp.Field(i)
		tag := tagutil.New(a.cmdtyp.Tag(i))

		if reltag := tag.First("rel"); len(reltag) > 0 {
			rid, err := a.relid(reltag)
			if err != nil {
				return err
			}

			rel := new(relinfo)
			rel.field = fld.Name()
			rel.relid = rid
			if err := a.reldatatype(rel, fld); err != nil {
				return err
			}
			a.cmd.rel = rel
			continue
		}

		switch strings.ToLower(fld.Name()) {
		case "where":
			if err := a.whereblock(fld); err != nil {
				return err
			}
		case "join", "from", "using":
			if err := a.joinblock(fld); err != nil {
				return err
			}
		}

		if dirname := typesutil.GetDirectiveName(fld); len(dirname) > 0 {
			switch strings.ToLower(dirname) {
			case "all":
				a.cmd.all = true
			}
		}

		// errorhandler
		// default columns
		// return columns
		// order by
		// offset
		// override
		// force
	}

	if a.cmd.rel == nil {
		return newerr(errNoRelation, a.cmd.name)
	}
	return nil
}

func (a *analyzer) reldatatype(rel *relinfo, field *types.Var) error {
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
	}

	// If unnamed and not an iterator, check for slices and pointers.
	if named == nil {
		if slice, ok := ftyp.(*types.Slice); ok { // allows []T / []*T
			ftyp = slice.Elem()
			rel.datatype.isslice = true
		} else if array, ok := ftyp.(*types.Array); ok { // allows [N]T / [N]*T
			ftyp = array.Elem()
			rel.datatype.isarray = true
			rel.datatype.arraylen = array.Len()
		}
		if ptr, ok := ftyp.(*types.Pointer); ok { // allows *T
			ftyp = ptr.Elem()
			rel.datatype.ispointer = true
		}
		if named, ok = ftyp.(*types.Named); !ok {
			// Fail if the type is a slice, an array, or a pointer
			// while its base type remains unnamed.
			if rel.datatype.isslice || rel.datatype.isarray || rel.datatype.ispointer {
				return newerr(errBadRelationType, a.cmd.name, rel.field)
			}
		}
	}

	if named != nil {
		pkg := named.Obj().Pkg()
		rel.datatype.name = named.Obj().Name()
		rel.datatype.pkgpath = pkg.Path()
		rel.datatype.pkgname = pkg.Name()
		rel.datatype.pkglocal = pkg.Name()
		rel.datatype.isimported = a.isimported(named)
		rel.datatype.isscanner = typesutil.ImplementsScanner(named)
		rel.datatype.isvaluer = typesutil.ImplementsValuer(named)
		rel.datatype.istime = typesutil.IsTime(named)
		rel.datatype.isafterscanner = typesutil.ImplementsAfterScanner(named)
		ftyp = named.Underlying()
	}

	rel.datatype.kind = a.typekind(ftyp)
	if rel.datatype.kind != kindstruct {
		// Currently only the struct kind is supported as the
		// relation's associated base datatype.
		return newerr(errBadRelationType, a.cmd.name, rel.field)
	}

	styp := ftyp.(*types.Struct)
	return a.relfields(rel, styp)
}

func (a *analyzer) relfields(rel *relinfo, styp *types.Struct) error {
	// The structloop type holds the state of a loop over a struct's fields.
	type structloop struct {
		styp *types.Struct // the struct type whose fields are being analyzed
		typ  *typeinfo     // info on the struct type; holds the resulting slice of analyzed fieldinfo
		idx  int           // keeps track of the field index
		pfx  string        // column prefix
	}

	// LIFO stack of structloops used for depth first traversal of struct fields.
	stack := []*structloop{{styp: styp, typ: &rel.datatype.typeinfo}}

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

			// Add the field to the list.
			loop.typ.fields = append(loop.typ.fields, f)

			// Analyze the field's type.
			ftyp := fld.Type()
			f.typ, ftyp = a.typeinfo(ftyp)

			// If the field's type is a struct and the `sql` tag's
			// value starts with the ">" (descend) marker, then it is
			// considered to be a "branch" field whose child fields
			// need to be analyzed as well.
			if f.typ.kind == kindstruct && strings.HasPrefix(sqltag, ">") && (!f.typ.isslice && !f.typ.isarray) {
				loop2 := new(structloop)
				loop2.styp = ftyp.(*types.Struct)
				loop2.typ = &f.typ
				loop2.pfx = loop.pfx + strings.TrimPrefix(sqltag, ">")
				stack = append(stack, loop2)
				continue stackloop
			}

			// If the field is not a struct to be descended,
			// it is considered to be a "leaf" field and as
			// such the analysis of leaf-specific information
			// needs to be carried out.
			f.auto = tag.HasOption("sql", "auto")
			f.ispkey = tag.HasOption("sql", "pk")
			f.nullempty = tag.HasOption("sql", "nullempty")
			f.readonly = tag.HasOption("sql", "ro")
			f.writeonly = tag.HasOption("sql", "wo")
			f.usejson = tag.HasOption("sql", "json")
			f.binadd = tag.HasOption("sql", "+")
			f.usecoalesce, f.coalesceval = a.coalesceinfo(tag)

			colid, err := a.colid(loop.pfx + sqltag)
			if err != nil {
				return err
			}
			f.colid = colid
		}
		stack = stack[:len(stack)-1]
	}
	return nil
}

// typeinfo analyzes the given type and returns the resulting info.
// The second return value is the base type of the given type.
func (a *analyzer) typeinfo(tt types.Type) (typ typeinfo, base types.Type) {
	base = tt
	if slice, ok := base.(*types.Slice); ok {
		base = slice.Elem()
		typ.isslice = true
	} else if array, ok := base.(*types.Array); ok {
		base = array.Elem()
		typ.isarray = true
		typ.arraylen = array.Len()
	}
	if ptr, ok := base.(*types.Pointer); ok {
		base = ptr.Elem()
		typ.ispointer = true
	}
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
		base = named.Underlying()
	}
	typ.kind = a.typekind(base)
	return typ, base
}

func (a *analyzer) iterator(iface *types.Interface, named *types.Named, rel *relinfo) (*types.Named, error) {
	if iface.NumExplicitMethods() != 1 {
		return nil, newerr(errBadIteratorType, a.cmd.name, rel.field)
	}

	mth := iface.ExplicitMethod(0)

	// Make sure that the method is exported or, if it's not, then at least
	// ensure that the receiver type is local, i.e. not imported, otherwise
	// the method will not be accessible.
	if !mth.Exported() && named != nil && (named.Obj().Pkg().Path() != a.pkg) {
		return nil, newerr(errBadIteratorType, a.cmd.name, rel.field)
	}

	sig := mth.Type().(*types.Signature)
	named, err := a.iteratorfunc(sig, rel)
	if err != nil {
		return nil, err
	}

	rel.datatype.itermethod = mth.Name()
	return named, nil
}

func (a *analyzer) iteratorfunc(sig *types.Signature, rel *relinfo) (*types.Named, error) {
	// Must take 1 argument and return one value of type error. "func(T) error"
	if sig.Params().Len() != 1 || sig.Results().Len() != 1 || !typesutil.IsError(sig.Results().At(0).Type()) {
		return nil, newerr(errBadIteratorType, a.cmd.name, rel.field)
	}

	typ := sig.Params().At(0).Type()
	if ptr, ok := typ.(*types.Pointer); ok { // allows *T
		typ = ptr.Elem()
		rel.datatype.ispointer = true
	}

	// Make sure that the argument type is a named struct type.
	named, ok := typ.(*types.Named)
	if !ok {
		return nil, newerr(errBadIteratorType, a.cmd.name, rel.field)
	} else if _, ok := named.Underlying().(*types.Struct); !ok {
		return nil, newerr(errBadIteratorType, a.cmd.name, rel.field)
	}

	rel.datatype.useiter = true
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

var recoalesce = regexp.MustCompile(`(?i)^coalesce$|^coalesce\((.*)\)$`)

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
	// The structloop type holds the state of a loop over a struct's fields.
	type structloop struct {
		wb  *whereblock
		ns  *typesutil.NamedStruct // the struct type of the whereblock
		idx int                    // keeps track of the field index
	}

	wb := new(whereblock)
	wb.name = field.Name()
	ns, err := typesutil.GetStruct(field)
	if err != nil {
		return err
	}

	// LIFO stack of struct loops used for depth first traversal of struct fields.
	stack := []*structloop{{wb: wb, ns: ns}}

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
						return newerr(errBadBoolTag)
					}
				}
			}

			// Nested whereblocks are marked with ">" and should be
			// analyzed before any other fields in the current block.
			if sqltag == ">" {
				ns, err := typesutil.GetStruct(fld)
				if err != nil {
					return err
				}

				wb := new(whereblock)
				wb.name = fld.Name()
				item.node = wb

				loop2 := new(structloop)
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
					colid, err := a.colid(lhs)
					if err != nil {
						return err
					}

					wn := new(wherecolumn)
					wn.colid = colid
					wn.cmp = string2cmpop[op]
					wn.sop = string2scalarrop[op2]

					if a.iscolid(rhs) {
						colid2, err := a.colid(rhs)
						if err != nil {
							return err
						}
						wn.colid2 = colid2
					} else {
						wn.lit = rhs // assume literal expression
					}

					item.node = wn
					continue
				}

				// Assume column with unary predicate.
				colid, err := a.colid(lhs)
				if err != nil {
					return err
				}

				wn := new(wherecolumn)
				wn.colid = colid
				wn.cmp = string2cmpop[op]
				wn.sop = string2scalarrop[op2]
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
				ns, err := typesutil.GetStruct(fld)
				if err != nil {
					return newerr(errBadBetweenType)
				} else if ns.Struct.NumFields() != 2 {
					return newerr(errBadBetweenType)
				}

				var x, y interface{}
				for i := 0; i < 2; i++ {
					fld := ns.Struct.Field(i)
					tag := tagutil.New(ns.Struct.Tag(i))
					sqltag := tag.First("sql")

					if fld.Name() == "_" && typesutil.IsDirective("Column", fld.Type()) {
						sqltag2 := strings.ToLower(tag.Second("sql"))

						colid, err := a.colid(sqltag)
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

				colid, err := a.colid(lhs)
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
			colid, err := a.colid(lhs)
			if err != nil {
				return err
			}

			wn := new(wherefield)
			wn.name = fld.Name()
			wn.colid = colid
			wn.typ, _ = a.typeinfo(fld.Type())
			wn.cmp = string2cmpop[op]
			wn.sop = string2scalarrop[op2]
			wn.mod = a.modfn(tag["sql"][1:])
			item.node = wn

			// TODO(mkopriva): make sure that, if a scalarrop was
			// provided, the fields type is either slice or array
			// and that the cmpop is an operator that can actually
			// be used with a scalarrop.

		}
		stack = stack[:len(stack)-1]
	}

	a.cmd.where = wb
	return nil
}

func (a *analyzer) joinblock(field *types.Var) (err error) {
	join := new(joinblock)
	ns, err := typesutil.GetStruct(field)
	if err != nil {
		return err
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

		switch dirname := typesutil.GetDirectiveName(fld); dirname {
		case "Relation":
			id, err := a.relid(sqltag)
			if err != nil {
				return err
			}
			join.rel = id
		case "LeftJoin", "RightJoin", "FullJoin", "CrossJoin":
			id, err := a.relid(sqltag)
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
					if cond.col1, err = a.colid(lhs); err != nil {
						return err
					}

					// optional right-hand side
					if a.iscolid(rhs) {
						if cond.col2, err = a.colid(rhs); err != nil {
							return err
						}
					} else {
						cond.lit = rhs
					}

					cond.cmp = string2cmpop[op]
					cond.sop = string2scalarrop[op2]
					conds = append(conds, cond)
				}
			}

			item := new(joinitem)
			item.typ = string2jointype[dirname]
			item.rel = id
			item.conds = conds
			join.items = append(join.items, item)
		}

	}

	a.cmd.join = join
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
		return expr, "=", "", "" // default
	}

	return lhs, cop, sop, rhs
}

func (a *analyzer) modfn(tagvals []string) function {
	for _, v := range tagvals {
		if fn, ok := string2function[strings.ToLower(v)]; ok {
			return fn
		}
	}
	return 0
}

func (a *analyzer) isimported(named *types.Named) bool {
	return named != nil && named.Obj().Pkg().Path() != a.pkg
}

func (a *analyzer) isaccessible(fld *types.Var, named *types.Named) bool {
	return fld.Name() != "_" && (fld.Exported() || !a.isimported(named))
}

var rerelid = regexp.MustCompile(`^(?:[A-Za-z_]\w*\.)?[A-Za-z_]\w*(?:\:[A-Za-z_]\w*)?$`)

// parses the given string and returns a relid, if the value's format is invalid
// an error will be returned instead. The expected format is: "[qualifier.]name[:alias]".
func (a *analyzer) relid(val string) (id relid, err error) {
	if !rerelid.MatchString(val) {
		return id, newerr(errBadObjId)
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

var recolid = regexp.MustCompile(`^(?:[A-Za-z_]\w*\.)?[A-Za-z_]\w*$`)
var rereserved = regexp.MustCompile(`^(?i:true|false|` +
	`current_date|current_time|current_timestamp|` +
	`current_role|current_schema|current_user|` +
	`localtime|localtimestamp|` +
	`session_user)$`)

func (a *analyzer) iscolid(val string) bool {
	return recolid.MatchString(val) && !rereserved.MatchString(val)
}

// parses the given string and returns a colid, if the value's format is invalid
// an error will be returned instead. The expected format is: "[qualifier.]name".
func (a *analyzer) colid(val string) (id colid, err error) {
	if !a.iscolid(val) {
		log.Println("not colid =>", val)
		return id, newerr(errBadColId)
	}
	if i := strings.LastIndexByte(val, '.'); i > -1 {
		id.qual = val[:i]
		val = val[i+1:]
	}
	id.name = val
	return id, nil
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
	name  string  // name of the target struct type
	typ   cmdtype // the type of the command
	rel   *relinfo
	join  *joinblock
	where *whereblock
	// Indicates that the command should be executed against all the rows
	// of the relation.
	all bool
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

// relinfo holds the information on a go struct type and on the
// db relation that's associated with that struct type.
type relinfo struct {
	field    string // name of the field that references the relation in the command
	relid    relid  // the relation identifier
	alias    string
	datatype datatype
	isview   bool // indicates that the relation is a table view
}

// datatype holds information on the type of data a command should read from,
// or write to, the associated database relation.
type datatype struct {
	typeinfo // type info on the datatype
	// if set, indicates that the datatype is handled by an iterator
	useiter bool
	// if set the value will hold the method name of the iterator interface
	itermethod string
	// reports whether or not the type implements the afterscanner interface
	isafterscanner bool
}

type typeinfo struct {
	name       string   // the name of a named type or empty string for unnamed types
	kind       typekind // the kind of the go type
	pkgpath    string   // the package import path
	pkgname    string   // the package's name
	pkglocal   string   // the local package name (including ".")
	isimported bool     // indicates whether or not the package is imported
	isarray    bool     // reports whether or not the type's an array type
	arraylen   int64    // if it's an array type, this field will hold the array's length
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
	// if set to true it indicates that the column value should be wrapped
	// in a COALESCE call when read from the db.
	usecoalesce bool
	coalesceval string
	// for UPDATEs, if set to true, it indicates that the provided field
	// value should be added to the already existing column value.
	binadd bool
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
	node interface{} // wherefield, wherecolumn, whereblock
}

type wherefield struct {
	name  string
	typ   typeinfo //
	colid colid    //
	cmp   cmpop    //
	sop   scalarrop
	mod   function //
}

// wherecolumn is produced from a gosql.Column directive and its tag value.
// wherecolumn represents either a column with a comparison predicate,
// a column-to-column comparison, or a column-to-literal comparison.
type wherecolumn struct {
	// The target column of the wherecolumn item.
	colid colid
	// If set, it will hold the comparison operator to be used to compare
	// the target column either using a predicate unary operator, or a binary
	// operator comparing against the colid2 column or the lit value.
	cmp cmpop
	sop scalarrop
	// If set, it will hold the id of the column that should be compared
	// to the target column.
	colid2 colid
	// If set, it will hold the literal value that should be compared
	// to the target column.
	lit string
}

type wherebetween struct {
	name  string
	colid colid
	cmp   cmpop
	x, y  interface{}
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

	// unary comparison predicates
	cmpisnull     // IS NULL
	cmpnotnull    // IS NOT NULL
	cmpistrue     // IS TRUE
	cmpnottrue    // IS NOT TRUE
	cmpisfalse    // IS FALSE
	cmpnotfalse   // IS NOT FALSE
	cmpisunknown  // IS UNKNOWN
	cmpnotunknown // IS NOT UNKNOWN

	// range comparison operators
	cmpisbetween     // BETWEEN x AND y
	cmpnotbetween    // NOT BETWEEN x AND y
	cmpisbetweensym  // BETWEEN SYMMETRIC x AND y
	cmpnotbetweensym // NOT BETWEEN SYMMETRIC x AND y

	// array comparison operators
	cmpisin  // IN
	cmpnotin // NOT IN

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
)

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

	"isnull":     cmpisnull,
	"notnull":    cmpnotnull,
	"istrue":     cmpistrue,
	"nottrue":    cmpnottrue,
	"isfalse":    cmpisfalse,
	"notfalse":   cmpnotfalse,
	"isunknown":  cmpisunknown,
	"notunknown": cmpnotunknown,

	"isbetween":     cmpisbetween,
	"notbetween":    cmpnotbetween,
	"isbetweensym":  cmpisbetweensym,
	"notbetweensym": cmpnotbetweensym,

	"isin":  cmpisin,
	"notin": cmpnotin,

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

type scalarrop uint // scalar array operator

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

type function uint

const (
	_       function = iota // no function
	fnlower                 // lower
	fnupper                 // upper
)

var string2function = map[string]function{
	"lower": fnlower,
	"upper": fnupper,
}

type jointype uint

const (
	_         jointype = iota // no join
	joinleft                  // LEFT JOIN
	joinright                 // RIGHT JOIN
	joinfull                  // FULL JOIN
	joincross                 // CROSS JOIN
)

var string2jointype = map[string]jointype{
	"LeftJoin":  joinleft,
	"RightJoin": joinright,
	"FullJoin":  joinfull,
	"CrossJoin": joincross,
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
	kindchan
	kindfunc
	kindinterface
	kindmap
	kindptr
	kindslice
	kindstruct
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
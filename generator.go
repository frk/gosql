package gosql

import (
	"bytes"
	"strconv"

	gol "github.com/frk/gosql/internal/golang"
	sql "github.com/frk/gosql/internal/sqlang"
)

const (
	filepreamble = ` DO NOT EDIT. This file was generated by "github.com/frk/gosql".`
	gosqlimport  = `github.com/frk/gosql`
)

var (
	idblank             = gol.Ident{"_"}
	idrecv              = gol.Ident{"q"}
	idconn              = gol.Ident{"c"}
	ididx               = gol.Ident{"i"}
	iderr               = gol.Ident{"err"}
	idnew               = gol.Ident{"new"}
	idnil               = gol.Ident{"nil"}
	idlen               = gol.Ident{"len"}
	idi64               = gol.Ident{"i64"}
	idres               = gol.Ident{"res"}
	idrow               = gol.Ident{"row"}
	idmake              = gol.Ident{"make"}
	idrows              = gol.Ident{"rows"}
	idexec              = gol.Ident{"Exec"}
	iderror             = gol.Ident{"error"}
	idparams            = gol.Ident{"params"}
	idafterscan         = gol.Ident{"AfterScan"}
	idquery             = gol.Ident{"queryString"}
	idiface             = gol.Ident{"interface{}"}
	idifaces            = gol.Ident{"[]interface{}"}
	sxconn              = gol.SelectorExpr{X: gol.Ident{"gosql"}, Sel: gol.Ident{"Conn"}}
	sxerrorinfo         = gol.SelectorExpr{X: gol.Ident{"gosql"}, Sel: gol.Ident{"ErrorInfo"}}
	sxexec              = gol.SelectorExpr{X: gol.Ident{"c"}, Sel: gol.Ident{"Exec"}}
	sxquery             = gol.SelectorExpr{X: gol.Ident{"c"}, Sel: gol.Ident{"Query"}}
	sxqueryrow          = gol.SelectorExpr{X: gol.Ident{"c"}, Sel: gol.Ident{"QueryRow"}}
	sxrowscan           = gol.SelectorExpr{X: gol.Ident{"row"}, Sel: gol.Ident{"Scan"}}
	sxrowsscan          = gol.SelectorExpr{X: gol.Ident{"rows"}, Sel: gol.Ident{"Scan"}}
	sxrowsclose         = gol.SelectorExpr{X: gol.Ident{"rows"}, Sel: gol.Ident{"Close"}}
	callrowserr         = gol.CallExpr{Fun: gol.SelectorExpr{X: gol.Ident{"rows"}, Sel: gol.Ident{"Err"}}}
	callrowsnext        = gol.CallExpr{Fun: gol.SelectorExpr{X: gol.Ident{"rows"}, Sel: gol.Ident{"Next"}}}
	callresrowsaffected = gol.CallExpr{Fun: gol.SelectorExpr{X: gol.Ident{"res"}, Sel: gol.Ident{"RowsAffected"}}}
	callinvaluelist     = gol.CallExpr{Fun: gol.SelectorExpr{X: gol.Ident{"gosql"}, Sel: gol.Ident{"InValueList"}}}
	callmakeparams      = gol.CallExpr{Fun: gol.Ident{"make"}, Args: gol.ArgsList{List: []gol.Expr{idifaces}}}
)

type specinfo struct {
	spec *typespec
	info *pginfo
}

func generate(pkgname string, infos []*specinfo) (*bytes.Buffer, error) {
	g := &generator{infos: infos}
	if err := g.run(pkgname); err != nil {
		return nil, err
	}
	return &g.buf, nil
}

type generator struct {
	infos   []*specinfo
	pkgname string
	buf     bytes.Buffer

	file gol.File

	// spec specific state, needs to be reset on each iteration
	nparam int                // number of parameters
	asvar  bool               // if true, the query string should be declared as a var, not const.
	qargs  []gol.Expr         // query arguments
	insx   []gol.SelectorExpr // slice fields for IN clauses
}

func (g *generator) run(pkgname string) error {
	g.file.PkgName = pkgname
	g.file.Preamble = gol.LineComment{filepreamble}
	g.file.Imports = gol.ImportDecl{{Path: gosqlimport}}

	for _, si := range g.infos {
		g.nparam = 0
		g.asvar = false
		g.qargs = nil
		g.insx = nil

		execdecl := g.execdecl(si)
		g.file.Decls = append(g.file.Decls, execdecl)
	}

	return gol.Write(g.file, &g.buf)
}

func (g *generator) execdecl(si *specinfo) (fn gol.FuncDecl) {
	fn.Name = idexec
	fn.Recv.Name = idrecv
	fn.Recv.Type = gol.StarExpr{X: gol.Ident{si.spec.name}}
	fn.Type.Params = gol.ParamList{{Names: []gol.Ident{idconn}, Type: sxconn}}
	fn.Type.Results = gol.ParamList{{Type: iderror}}

	g.queryargs(si.spec)

	fn.Body.Add(g.querybuild(si))
	fn.Body.Add(gol.NL{})
	fn.Body.Add(g.querydefaults(si))
	fn.Body.Add(g.queryexec(si))
	fn.Body.Add(g.returnstmt(si))
	return fn
}

func (g *generator) querybuild(si *specinfo) (stmt gol.Stmt) {
	sqlnode := g.sqlnode(si)

	token := gol.DECL_CONST
	if g.asvar || len(si.spec.filter) > 0 {
		token = gol.DECL_VAR
	}

	decl := gol.GenDecl{Token: token}
	decl.Specs = []gol.Spec{gol.ValueSpec{
		Names:       []gol.Ident{idquery},
		Values:      []gol.Expr{gol.RawStringNode{sqlnode}},
		LineComment: gol.LineComment{" `"},
	}}

	if len(g.insx) > 0 {
		// prepare the var declarations
		vardecl := gol.GenDecl{Token: gol.DECL_VAR}

		nstatic := gol.ValueSpec{}
		nstatic.Names = []gol.Ident{gol.Ident{"nstatic"}}
		nstatic.Values = []gol.Expr{gol.Int(g.nparam)}
		nstatic.LineComment = gol.LineComment{" number of static parameters"}
		vardecl.Specs = append(vardecl.Specs, nstatic)

		for i, sx := range g.insx {
			num := strconv.Itoa(i + 1)

			lenspec := gol.ValueSpec{}
			lenspec.Names = []gol.Ident{gol.Ident{"len" + num}}
			lenspec.Values = []gol.Expr{gol.CallExpr{Fun: idlen, Args: gol.ArgsList{List: []gol.Expr{sx}}}}
			lenspec.LineComment = gol.LineComment{" length of slice #" + num + " to be unnested"}

			posspec := gol.ValueSpec{}
			posspec.Names = []gol.Ident{gol.Ident{"pos" + num}}
			if i == 0 {
				// the first position is set to the value of nstatic
				posspec.Values = []gol.Expr{gol.Ident{"nstatic"}}
			} else {
				// the rest of the positions are calculated from
				// adding the previous length to the previous position
				prev := strconv.Itoa(i)
				prevlen, prevpos := gol.Ident{"len" + prev}, gol.Ident{"pos" + prev}
				posspec.Values = []gol.Expr{gol.BinaryExpr{X: prevpos, Op: gol.BINARY_ADD, Y: prevlen}}
			}
			posspec.LineComment = gol.LineComment{" starting position of slice #" + num + " parameters"}

			vardecl.Specs = append(vardecl.Specs, lenspec, posspec)
		}

		// next is the query declaration
		list := gol.StmtList{gol.DeclStmt{vardecl}, gol.NL{}, gol.DeclStmt{decl}, gol.NL{}}

		// define the params variable
		asn := gol.AssignStmt{Token: gol.ASSIGN_DEFINE}
		asn.Lhs = []gol.Expr{idparams}
		callmake := callmakeparams
		bin := gol.BinaryExpr{X: gol.Ident{"nstatic"}, Op: gol.BINARY_ADD, Y: gol.Ident{"len1"}}
		for i := 1; i < len(g.insx); i++ {
			y := gol.Ident{"len" + strconv.Itoa(i+1)}
			bin = gol.BinaryExpr{X: bin, Op: gol.BINARY_ADD, Y: y}
		}
		callmake.Args.List = append(callmake.Args.List, bin)
		asn.Rhs = []gol.Expr{callmake}
		list = append(list, asn)

		// directly assign non-slice params
		for i, arg := range g.qargs {
			asn := gol.AssignStmt{Token: gol.ASSIGN}
			asn.Lhs = []gol.Expr{gol.IndexExpr{X: idparams, Index: gol.Int(i)}}
			asn.Rhs = []gol.Expr{arg}
			list = append(list, asn)
		}

		for i, sx := range g.insx {
			lenid := gol.Ident{"len" + strconv.Itoa(i+1)}
			posid := gol.Ident{"pos" + strconv.Itoa(i+1)}

			loop := gol.ForStmt{}
			loop.Init = gol.AssignStmt{Token: gol.ASSIGN_DEFINE, Lhs: []gol.Expr{ididx}, Rhs: []gol.Expr{gol.Int(0)}}
			loop.Cond = gol.BinaryExpr{X: ididx, Op: gol.BINARY_LSS, Y: lenid}
			loop.Post = gol.IncDecStmt{X: ididx, Token: gol.INCDEC_INC}

			asn := gol.AssignStmt{Token: gol.ASSIGN}
			asn.Lhs = []gol.Expr{gol.IndexExpr{X: idparams, Index: gol.BinaryExpr{X: posid, Op: gol.BINARY_ADD, Y: ididx}}}
			asn.Rhs = []gol.Expr{gol.IndexExpr{X: sx, Index: ididx}}

			loop.Body = gol.BlockStmt{List: []gol.Stmt{asn}}

			list = append(list, loop)
		}

		return append(list, gol.NL{})
	} else if len(si.spec.filter) > 0 {
		asn := gol.AssignStmt{Token: gol.ASSIGN_ADD}
		asn.Lhs = []gol.Expr{idquery}
		asn.Rhs = []gol.Expr{gol.CallExpr{Fun: gol.SelectorExpr{
			X:   gol.SelectorExpr{X: idrecv, Sel: gol.Ident{si.spec.filter}},
			Sel: gol.Ident{"ToSQL"},
		}}}

		asn2 := gol.AssignStmt{Token: gol.ASSIGN_DEFINE}
		asn2.Lhs = []gol.Expr{idparams}
		asn2.Rhs = []gol.Expr{gol.CallExpr{Fun: gol.SelectorExpr{
			X:   gol.SelectorExpr{X: idrecv, Sel: gol.Ident{si.spec.filter}},
			Sel: gol.Ident{"Params"},
		}}}
		return gol.StmtList{gol.DeclStmt{decl}, gol.NL{}, asn, asn2, gol.NL{}}
	}
	return gol.DeclStmt{decl}
}

func (g *generator) querydefaults(si *specinfo) (stmt gol.Stmt) {
	var list gol.StmtList
	if l := si.spec.limit; l != nil && l.value > 0 && len(l.field) > 0 {
		sx := gol.SelectorExpr{X: idrecv, Sel: gol.Ident{l.field}}

		asn := gol.AssignStmt{Token: gol.ASSIGN}
		asn.Lhs = []gol.Expr{sx}
		asn.Rhs = []gol.Expr{gol.BasicLit{strconv.FormatUint(l.value, 10)}}

		ifzero := gol.IfStmt{}
		ifzero.Cond = gol.BinaryExpr{X: sx, Op: gol.BINARY_EQL, Y: gol.BasicLit{"0"}}
		ifzero.Body = gol.BlockStmt{List: []gol.Stmt{asn}}
		list = append(list, ifzero)
	}

	if o := si.spec.offset; o != nil && o.value > 0 && len(o.field) > 0 {
		sx := gol.SelectorExpr{X: idrecv, Sel: gol.Ident{o.field}}

		asn := gol.AssignStmt{Token: gol.ASSIGN}
		asn.Lhs = []gol.Expr{sx}
		asn.Rhs = []gol.Expr{gol.BasicLit{strconv.FormatUint(o.value, 10)}}

		ifzero := gol.IfStmt{}
		ifzero.Cond = gol.BinaryExpr{X: sx, Op: gol.BINARY_EQL, Y: gol.BasicLit{"0"}}
		ifzero.Body = gol.BlockStmt{List: []gol.Stmt{asn}}
		list = append(list, ifzero)
	}

	if len(list) == 0 {
		return gol.NoOp{}
	}

	list = append(list, gol.NL{})
	return list
}

func (g *generator) queryexec(si *specinfo) (stmt gol.Stmt) {
	args := gol.ArgsList{List: []gol.Expr{idquery}}
	if len(g.insx) > 0 || len(si.spec.filter) > 0 {
		args.List = append(args.List, idparams)
		args.Ellipsis = true
	} else {
		args.List = append(args.List, g.qargs...)
		if len(args.List) > 3 {
			args.OnePerLine = 2
		}
	}

	// produce c.Exec( ... ) call
	{
		if si.spec.kind != speckindSelect && si.spec.returning == nil && si.spec.result == nil {

			if rafield := si.spec.rowsaffected; rafield != nil {
				// call exec & assign res, err
				asn := gol.AssignStmt{Token: gol.ASSIGN_DEFINE}
				asn.Lhs = []gol.Expr{idres, iderr}
				asn.Rhs = []gol.Expr{gol.CallExpr{Fun: sxexec, Args: args}}

				// check err
				iferr := gol.IfStmt{}
				iferr.Cond = gol.BinaryExpr{X: iderr, Op: gol.BINARY_NEQ, Y: idnil}
				iferr.Body = gol.BlockStmt{List: []gol.Stmt{g.returnerr(si, iderr)}}

				// call RowsAffected & assing i64, err
				asn2 := gol.AssignStmt{Token: gol.ASSIGN_DEFINE}
				asn2.Lhs = []gol.Expr{idi64, iderr}
				asn2.Rhs = []gol.Expr{callresrowsaffected}

				// check err
				iferr2 := gol.IfStmt{}
				iferr2.Cond = gol.BinaryExpr{X: iderr, Op: gol.BINARY_NEQ, Y: idnil}
				iferr2.Body = gol.BlockStmt{List: []gol.Stmt{g.returnerr(si, iderr)}}

				//
				asn3 := gol.AssignStmt{Token: gol.ASSIGN}
				asn3.Lhs = []gol.Expr{gol.SelectorExpr{X: idrecv, Sel: gol.Ident{rafield.name}}}
				if rafield.kind == kindint64 {
					asn3.Rhs = []gol.Expr{idi64}
				} else {
					args := gol.ArgsList{List: []gol.Expr{idi64}}
					asn3.Rhs = []gol.Expr{gol.CallExpr{Fun: gol.Ident{typekind2string[rafield.kind]}, Args: args}}
				}

				return gol.StmtList{asn, iferr, asn2, iferr2, gol.NL{}, asn3}
			} else {
				asn := gol.AssignStmt{Token: gol.ASSIGN_DEFINE}
				asn.Lhs = []gol.Expr{idblank, iderr}
				asn.Rhs = []gol.Expr{gol.CallExpr{Fun: sxexec, Args: args}}
				return asn
			}
		}
	}

	// produce c.QueryRow( ... ) call
	{
		rec := si.spec.rel.rec
		if si.spec.result != nil {
			rec = si.spec.result.rec
		}

		if !rec.isslice && !rec.isiter {
			asn := gol.AssignStmt{Token: gol.ASSIGN_DEFINE}
			asn.Lhs = []gol.Expr{idrow}
			asn.Rhs = []gol.Expr{gol.CallExpr{Fun: sxqueryrow, Args: args}}
			return gol.StmtList{asn, gol.NL{}}
		}
	}

	// produce c.Query( ... ) call with if-err-check, defer-rows-close, and
	// for-rows-next loop to scan the rows
	{
		asn := gol.AssignStmt{Token: gol.ASSIGN_DEFINE}
		asn.Lhs = []gol.Expr{idrows, iderr}
		asn.Rhs = []gol.Expr{gol.CallExpr{Fun: sxquery, Args: args}}

		iferr := gol.IfStmt{}
		iferr.Cond = gol.BinaryExpr{X: iderr, Op: gol.BINARY_NEQ, Y: idnil}
		iferr.Body = gol.BlockStmt{List: []gol.Stmt{g.returnerr(si, iderr)}}

		defclose := gol.DeferStmt{}
		defclose.Call = gol.CallExpr{Fun: sxrowsclose}

		fornext := g.fornext(si)
		return gol.StmtList{asn, iferr, defclose, gol.NL{}, fornext}
	}
	return stmt
}

func (g *generator) returnerr(si *specinfo, errx gol.Expr) gol.ReturnStmt {
	if si.spec.erh == nil {
		return gol.ReturnStmt{errx}
	}
	if si.spec.erh.isinfo {
		lit := gol.CompositeLit{Type: sxerrorinfo, Comma: true, Compact: true}
		lit.Elts = append(lit.Elts, gol.KeyValueExpr{Key: gol.Ident{"Error"}, Value: iderr})
		lit.Elts = append(lit.Elts, gol.KeyValueExpr{Key: gol.Ident{"Query"}, Value: idquery})
		lit.Elts = append(lit.Elts, gol.KeyValueExpr{Key: gol.Ident{"SpecKind"}, Value: gol.String(si.spec.kind.String())})
		lit.Elts = append(lit.Elts, gol.KeyValueExpr{Key: gol.Ident{"SpecName"}, Value: gol.String(si.spec.name)})
		lit.Elts = append(lit.Elts, gol.KeyValueExpr{Key: gol.Ident{"SpecValue"}, Value: idrecv})
		litptr := gol.UnaryExpr{Op: gol.UNARY_AMP, X: lit}

		sx := gol.SelectorExpr{X: idrecv, Sel: gol.Ident{si.spec.erh.name}}
		fun := gol.SelectorExpr{X: sx, Sel: gol.Ident{"HandleErrorInfo"}}
		call := gol.CallExpr{Fun: fun, Args: gol.ArgsList{List: []gol.Expr{litptr}}}
		return gol.ReturnStmt{call}

		// TODO if errx is not iderr, then errx is probably a CallFunc and
		// should first be executed and it's result passed into HandleErrorInfo..
	}
	sx := gol.SelectorExpr{X: idrecv, Sel: gol.Ident{si.spec.erh.name}}
	fun := gol.SelectorExpr{X: sx, Sel: gol.Ident{"HandleError"}}
	call := gol.CallExpr{Fun: fun, Args: gol.ArgsList{List: []gol.Expr{errx}}}
	return gol.ReturnStmt{call}
}

func (g *generator) fornext(si *specinfo) (stmt gol.ForStmt) {
	stmt.Cond = callrowsnext
	// initialize
	{
		rec := si.spec.rel.rec
		if si.spec.result != nil {
			rec = si.spec.result.rec
		}

		if rec.ispointer {
			init := gol.AssignStmt{Token: gol.ASSIGN_DEFINE}
			init.Lhs = []gol.Expr{gol.Ident{"v"}}
			init.Rhs = []gol.Expr{gol.CallExpr{Fun: idnew, Args: gol.ArgsList{List: []gol.Expr{g.rectype(rec)}}}}
			stmt.Body.List = append(stmt.Body.List, init)
		} else {
			vs := gol.ValueSpec{Names: []gol.Ident{{"v"}}, Type: g.rectype(rec)}
			init := gol.GenDecl{Token: gol.DECL_VAR}
			init.Specs = append(init.Specs, vs)
			stmt.Body.List = append(stmt.Body.List, gol.DeclStmt{init})
		}
	}

	// scan & assign error
	{
		var args gol.ArgsList
		if len(si.info.output) > 2 {
			args.OnePerLine = 1
		}

		// The pfieldhandled map is used to keep track of pointer fields
		// that have already been initialized and their types imported.
		var pfieldhandled = make(map[string]bool)

		for _, item := range si.info.output {
			var fx gol.Expr = gol.Ident{"v"}

			var fieldkey string // key for the pfieldhandled map
			for _, pe := range item.field.path {
				fx = gol.SelectorExpr{X: fx, Sel: gol.Ident{pe.name}}

				fieldkey += pe.name
				if pe.ispointer && !pfieldhandled[fieldkey] {
					if pe.isimported {
						g.addimport(pe.typepkgpath, pe.typepkgname, pe.typepkglocal)
					}

					// initialize nested pointer field
					init := gol.AssignStmt{Token: gol.ASSIGN}
					init.Lhs = []gol.Expr{fx}
					init.Rhs = []gol.Expr{gol.CallExpr{Fun: idnew, Args: gol.ArgsList{List: []gol.Expr{g.pathelemtype(pe)}}}}
					stmt.Body.List = append(stmt.Body.List, init)

					pfieldhandled[fieldkey] = true
				}
			}

			fx = gol.SelectorExpr{X: fx, Sel: gol.Ident{item.field.name}}
			args.List = append(args.List, gol.UnaryExpr{Op: gol.UNARY_AMP, X: fx})
		}

		asn := gol.AssignStmt{Token: gol.ASSIGN_DEFINE}
		asn.Lhs = []gol.Expr{iderr}
		asn.Rhs = []gol.Expr{gol.CallExpr{Fun: sxrowsscan, Args: args}}
		stmt.Body.List = append(stmt.Body.List, asn)
	}

	// check error & newline
	{
		iferr := gol.IfStmt{}
		iferr.Cond = gol.BinaryExpr{X: iderr, Op: gol.BINARY_NEQ, Y: idnil}
		iferr.Body = gol.BlockStmt{List: []gol.Stmt{g.returnerr(si, iderr)}}
		stmt.Body.List = append(stmt.Body.List, iferr, gol.NL{})
	}

	// append OR iterate
	{
		rec := si.spec.rel.rec
		fieldname := si.spec.rel.name
		if si.spec.result != nil {
			rec = si.spec.result.rec
			fieldname = si.spec.result.name
		}

		if rec.isafterscanner {
			// call afterscan
			sx := gol.SelectorExpr{X: gol.Ident{"v"}, Sel: idafterscan}
			afterscan := gol.ExprStmt{gol.CallExpr{Fun: sx}}
			stmt.Body.List = append(stmt.Body.List, afterscan)
		}

		if rec.isiter {
			var call gol.CallExpr
			if len(rec.itermethod) > 0 {
				call = gol.CallExpr{Fun: gol.SelectorExpr{
					X:   gol.SelectorExpr{X: idrecv, Sel: gol.Ident{fieldname}},
					Sel: gol.Ident{rec.itermethod}},
					Args: gol.ArgsList{List: []gol.Expr{gol.Ident{"v"}}},
				}
			} else {
				call = gol.CallExpr{Fun: gol.SelectorExpr{
					X:   idrecv,
					Sel: gol.Ident{fieldname}},
					Args: gol.ArgsList{List: []gol.Expr{gol.Ident{"v"}}},
				}
			}

			asn := gol.AssignStmt{Token: gol.ASSIGN_DEFINE}
			asn.Lhs = []gol.Expr{iderr}
			asn.Rhs = []gol.Expr{call}

			iferr := gol.IfStmt{}
			iferr.Init = asn
			iferr.Cond = gol.BinaryExpr{X: iderr, Op: gol.BINARY_NEQ, Y: idnil}
			iferr.Body = gol.BlockStmt{List: []gol.Stmt{g.returnerr(si, iderr)}}
			stmt.Body.List = append(stmt.Body.List, iferr)
		} else {
			appnd := gol.CallExpr{Fun: gol.Ident{"append"}}
			appnd.Args = gol.ArgsList{List: []gol.Expr{
				gol.SelectorExpr{X: idrecv, Sel: gol.Ident{fieldname}},
				gol.Ident{"v"},
			}}

			asn := gol.AssignStmt{Token: gol.ASSIGN}
			asn.Lhs = []gol.Expr{gol.SelectorExpr{X: idrecv, Sel: gol.Ident{fieldname}}}
			asn.Rhs = []gol.Expr{appnd}
			stmt.Body.List = append(stmt.Body.List, asn)
		}
	}
	return stmt
}

func (g *generator) queryargs(spec *typespec) {
	if spec.where != nil && len(spec.where.items) > 0 {
		type loopstate struct {
			items []*predicateitem // the current iteration predicate items
			idx   int              // keeps track of the item index
			sx    gol.SelectorExpr // the selector expression for the current predicate field
		}

		sx := gol.SelectorExpr{X: idrecv, Sel: gol.Ident{spec.where.name}}
		stack := []*loopstate{{items: spec.where.items, sx: sx}}

	stackloop:
		for len(stack) > 0 {
			loop := stack[len(stack)-1]
			for loop.idx < len(loop.items) {
				item := loop.items[loop.idx]
				loop.idx++

				switch node := item.node.(type) {
				case *fieldpredicate:
					if node.pred != isin && node.pred != notin {
						sx := gol.SelectorExpr{X: loop.sx, Sel: gol.Ident{node.field.name}}
						g.qargs = append(g.qargs, sx)
					}
				case *betweenpredicate:
					if x, ok := node.x.(*paramfield); ok {
						sx := gol.SelectorExpr{X: loop.sx, Sel: gol.Ident{node.name}}
						sx = gol.SelectorExpr{X: sx, Sel: gol.Ident{x.name}}
						g.qargs = append(g.qargs, sx)
					}
					if y, ok := node.y.(*paramfield); ok {
						sx := gol.SelectorExpr{X: loop.sx, Sel: gol.Ident{node.name}}
						sx = gol.SelectorExpr{X: sx, Sel: gol.Ident{y.name}}
						g.qargs = append(g.qargs, sx)
					}
				case *nestedpredicate:
					loop2 := new(loopstate)
					loop2.items = node.items
					loop2.sx = gol.SelectorExpr{X: loop.sx, Sel: gol.Ident{node.name}}
					stack = append(stack, loop2)
					continue stackloop
				case *columnpredicate:
					// nothing to do
				}
			}
			stack = stack[:len(stack)-1]
		}
	}

	if spec.limit != nil && len(spec.limit.field) > 0 {
		sx := gol.SelectorExpr{X: idrecv, Sel: gol.Ident{spec.limit.field}}
		g.qargs = append(g.qargs, sx)
	}
	if spec.offset != nil && len(spec.offset.field) > 0 {
		sx := gol.SelectorExpr{X: idrecv, Sel: gol.Ident{spec.offset.field}}
		g.qargs = append(g.qargs, sx)
	}
}

func (g *generator) returnstmt(si *specinfo) (stmt gol.Stmt) {
	if si.spec.rowsaffected != nil {
		return gol.ReturnStmt{idnil}
	}

	if si.spec.kind == speckindSelect || si.spec.returning != nil {
		rel := si.spec.rel

		// does the record type need pre-allocation? and is it imported?
		if rel.rec.base.isimported && (rel.rec.isslice || rel.rec.ispointer) {
			g.addimport(rel.rec.base.pkgpath, rel.rec.base.pkgname, rel.rec.base.pkglocal)
		}

		if rel.rec.isslice || rel.rec.isarray || rel.rec.isiter {
			if si.spec.erh != nil && si.spec.erh.isinfo {
				asn := gol.AssignStmt{Token: gol.ASSIGN_DEFINE}
				asn.Lhs = []gol.Expr{iderr}
				asn.Rhs = []gol.Expr{callrowserr}

				iferr := gol.IfStmt{}
				iferr.Init = asn
				iferr.Cond = gol.BinaryExpr{X: iderr, Op: gol.BINARY_NEQ, Y: idnil}
				iferr.Body = gol.BlockStmt{List: []gol.Stmt{g.returnerr(si, iderr)}}

				return gol.StmtList{iferr, gol.ReturnStmt{idnil}}
			}
			return g.returnerr(si, callrowserr)
		} else {
			var list gol.StmtList // result

			// if the gosql.Return directive was used, make sure that the rel
			// field is properly initialized to avoid "nil pointer" panics.
			if rel.rec.ispointer {
				// initialize rel field
				init := gol.AssignStmt{Token: gol.ASSIGN}
				init.Lhs = []gol.Expr{gol.SelectorExpr{X: idrecv, Sel: gol.Ident{rel.name}}}
				init.Rhs = []gol.Expr{gol.CallExpr{Fun: idnew, Args: gol.ArgsList{List: []gol.Expr{g.rectype(rel.rec)}}}}
				list = append(list, init)
			}

			var args gol.ArgsList
			if len(si.info.output) > 2 {
				args.OnePerLine = 1
			}

			// The pfieldhandled map is used to keep track of pointer fields
			// that have already been initialized and their types imported.
			var pfieldhandled = make(map[string]bool)

			for _, item := range si.info.output {
				fx := gol.SelectorExpr{X: idrecv, Sel: gol.Ident{rel.name}}

				var fieldkey string // key for the pfieldhandled map
				for _, pe := range item.field.path {
					fx = gol.SelectorExpr{X: fx, Sel: gol.Ident{pe.name}}

					fieldkey += pe.name
					if pe.ispointer && !pfieldhandled[fieldkey] {
						if pe.isimported {
							g.addimport(pe.typepkgpath, pe.typepkgname, pe.typepkglocal)
						}

						// initialize nested pointer field
						init := gol.AssignStmt{Token: gol.ASSIGN}
						init.Lhs = []gol.Expr{fx}
						init.Rhs = []gol.Expr{gol.CallExpr{Fun: idnew, Args: gol.ArgsList{List: []gol.Expr{g.pathelemtype(pe)}}}}
						list = append(list, init)

						pfieldhandled[fieldkey] = true
					}
				}

				fx = gol.SelectorExpr{X: fx, Sel: gol.Ident{item.field.name}}
				args.List = append(args.List, gol.UnaryExpr{Op: gol.UNARY_AMP, X: fx})
			}

			if !rel.rec.isafterscanner {
				call := gol.CallExpr{Fun: sxrowscan, Args: args}
				if si.spec.erh != nil && si.spec.erh.isinfo {
					asn := gol.AssignStmt{Token: gol.ASSIGN_DEFINE}
					asn.Lhs = []gol.Expr{iderr}
					asn.Rhs = []gol.Expr{call}

					iferr := gol.IfStmt{}
					iferr.Cond = gol.BinaryExpr{X: iderr, Op: gol.BINARY_NEQ, Y: idnil}
					iferr.Body = gol.BlockStmt{List: []gol.Stmt{g.returnerr(si, iderr)}}

					list = append(list, asn, iferr, gol.ReturnStmt{idnil})
				} else {
					list = append(list, g.returnerr(si, call))
				}
			} else {
				// scan & assing error
				asn := gol.AssignStmt{Token: gol.ASSIGN_DEFINE}
				asn.Lhs = []gol.Expr{iderr}
				asn.Rhs = []gol.Expr{gol.CallExpr{Fun: sxrowscan, Args: args}}

				// check error
				iferr := gol.IfStmt{}
				iferr.Cond = gol.BinaryExpr{X: iderr, Op: gol.BINARY_NEQ, Y: idnil}
				iferr.Body = gol.BlockStmt{List: []gol.Stmt{g.returnerr(si, iderr)}}

				// call afterscan
				sx := gol.SelectorExpr{X: idrecv, Sel: gol.Ident{rel.name}}
				sx = gol.SelectorExpr{X: sx, Sel: idafterscan}
				afterscan := gol.ExprStmt{gol.CallExpr{Fun: sx}}

				// call afterscan
				ret := gol.ReturnStmt{idnil}

				// done
				list = append(list, asn, iferr, gol.NL{}, afterscan, ret)
			}

			return list
		}
	}

	// result field
	if si.spec.result != nil {
		rel := si.spec.result
		// does the record type need pre-allocation? and is it imported?
		if rel.rec.base.isimported && (rel.rec.isslice || rel.rec.ispointer) {
			g.addimport(rel.rec.base.pkgpath, rel.rec.base.pkgname, rel.rec.base.pkglocal)
		}

		if rel.rec.isslice || rel.rec.isarray || rel.rec.isiter {
			if si.spec.erh != nil && si.spec.erh.isinfo {
				asn := gol.AssignStmt{Token: gol.ASSIGN_DEFINE}
				asn.Lhs = []gol.Expr{iderr}
				asn.Rhs = []gol.Expr{callrowserr}

				iferr := gol.IfStmt{}
				iferr.Init = asn
				iferr.Cond = gol.BinaryExpr{X: iderr, Op: gol.BINARY_NEQ, Y: idnil}
				iferr.Body = gol.BlockStmt{List: []gol.Stmt{g.returnerr(si, iderr)}}

				return gol.StmtList{iferr, gol.ReturnStmt{idnil}}
			}
			return g.returnerr(si, callrowserr)
		} else {
			var list gol.StmtList // result

			// if the gosql.Return directive was used, make sure that the rel
			// field is properly initialized to avoid "nil pointer" panics.
			if rel.rec.ispointer {
				// initialize rel field
				init := gol.AssignStmt{Token: gol.ASSIGN}
				init.Lhs = []gol.Expr{gol.SelectorExpr{X: idrecv, Sel: gol.Ident{rel.name}}}
				init.Rhs = []gol.Expr{gol.CallExpr{Fun: idnew, Args: gol.ArgsList{List: []gol.Expr{g.rectype(rel.rec)}}}}
				list = append(list, init)
			}

			var args gol.ArgsList
			if len(si.info.output) > 2 {
				args.OnePerLine = 1
			}

			// The pfieldhandled map is used to keep track of pointer fields
			// that have already been initialized and their types imported.
			var pfieldhandled = make(map[string]bool)

			for _, item := range si.info.output {
				fx := gol.SelectorExpr{X: idrecv, Sel: gol.Ident{rel.name}}

				var fieldkey string // key for the pfieldhandled map
				for _, pe := range item.field.path {
					fx = gol.SelectorExpr{X: fx, Sel: gol.Ident{pe.name}}

					fieldkey += pe.name
					if pe.ispointer && !pfieldhandled[fieldkey] {
						if pe.isimported {
							g.addimport(pe.typepkgpath, pe.typepkgname, pe.typepkglocal)
						}

						// initialize nested pointer field
						init := gol.AssignStmt{Token: gol.ASSIGN}
						init.Lhs = []gol.Expr{fx}
						init.Rhs = []gol.Expr{gol.CallExpr{Fun: idnew, Args: gol.ArgsList{List: []gol.Expr{g.pathelemtype(pe)}}}}
						list = append(list, init)

						pfieldhandled[fieldkey] = true
					}
				}

				fx = gol.SelectorExpr{X: fx, Sel: gol.Ident{item.field.name}}
				args.List = append(args.List, gol.UnaryExpr{Op: gol.UNARY_AMP, X: fx})
			}

			if !rel.rec.isafterscanner {
				call := gol.CallExpr{Fun: sxrowscan, Args: args}
				if si.spec.erh != nil && si.spec.erh.isinfo {
					asn := gol.AssignStmt{Token: gol.ASSIGN_DEFINE}
					asn.Lhs = []gol.Expr{iderr}
					asn.Rhs = []gol.Expr{call}

					iferr := gol.IfStmt{}
					iferr.Cond = gol.BinaryExpr{X: iderr, Op: gol.BINARY_NEQ, Y: idnil}
					iferr.Body = gol.BlockStmt{List: []gol.Stmt{g.returnerr(si, iderr)}}

					list = append(list, asn, iferr, gol.ReturnStmt{idnil})
				} else {
					list = append(list, g.returnerr(si, call))
				}
			} else {
				// scan & assing error
				asn := gol.AssignStmt{Token: gol.ASSIGN_DEFINE}
				asn.Lhs = []gol.Expr{iderr}
				asn.Rhs = []gol.Expr{gol.CallExpr{Fun: sxrowscan, Args: args}}

				// check error
				iferr := gol.IfStmt{}
				iferr.Cond = gol.BinaryExpr{X: iderr, Op: gol.BINARY_NEQ, Y: idnil}
				iferr.Body = gol.BlockStmt{List: []gol.Stmt{g.returnerr(si, iderr)}}

				// call afterscan
				sx := gol.SelectorExpr{X: idrecv, Sel: gol.Ident{rel.name}}
				sx = gol.SelectorExpr{X: sx, Sel: idafterscan}
				afterscan := gol.ExprStmt{gol.CallExpr{Fun: sx}}

				// call afterscan
				ret := gol.ReturnStmt{idnil}

				// done
				list = append(list, asn, iferr, gol.NL{}, afterscan, ret)
			}

			return list
		}
	}

	return g.returnerr(si, iderr)
}

func (g *generator) rectype(rec recordtype) gol.Expr {
	id := gol.Ident{rec.base.name}
	if rec.base.isimported {
		return gol.SelectorExpr{X: gol.Ident{rec.base.pkgname}, Sel: id}
	}
	return id
}

func (g *generator) pathelemtype(pe *pathelem) gol.Expr {
	id := gol.Ident{pe.typename}
	if pe.isimported {
		return gol.SelectorExpr{X: gol.Ident{pe.typepkgname}, Sel: id}
	}
	return id
}

func (g *generator) addimport(path, name, local string) {
	// first check that that package path hasn't yet been added to the imports
	for _, spec := range g.file.Imports {
		if string(spec.Path) == path {
			return
		}
	}

	// if the local name is the same as the package name set it to empty
	if local == name {
		local = ""
	}

	spec := gol.ImportSpec{Path: gol.String(path), Name: gol.Ident{local}}
	g.file.Imports = append(g.file.Imports, spec)
}

func (g *generator) sqlnode(si *specinfo) (node gol.Node) {
	switch si.spec.kind {
	case speckindInsert:
		return g.sqlinsert(si)
	case speckindUpdate:
		return g.sqlupdate(si)
	case speckindSelect:
		return g.sqlselect(si)
	case speckindDelete:
		return g.sqldelete(si)
	}
	return node
}

func (g *generator) sqlinsert(si *specinfo) (insstmt sql.InsertStatement) {
	// TODO
	return insstmt
}

func (g *generator) sqlupdate(si *specinfo) (updstmt sql.UpdateStatement) {
	// TODO
	return updstmt
}

// sqlselect builds and returns an sql.SelectStatement.
func (g *generator) sqlselect(si *specinfo) (selstmt sql.SelectStatement) {
	var columns sql.ValueExprList
	for _, col := range si.info.output {
		columns = append(columns, g.sqlcolexpr(col))
	}

	selstmt.Columns = columns
	selstmt.Table = g.sqlrelid(si.spec.rel.relid)
	selstmt.Join = g.sqljoin(si.spec.join)
	selstmt.Where = g.sqlwhere(si.spec.where)
	selstmt.Order = g.sqlorderby(si.spec)
	selstmt.Limit = g.sqllimit(si.spec)
	selstmt.Offset = g.sqloffset(si.spec)
	return selstmt
}

// sqldelete builds and returns an sql.DeleteStatement.
func (g *generator) sqldelete(si *specinfo) (delstmt sql.DeleteStatement) {
	var returning sql.ReturningClause
	for _, col := range si.info.output {
		returning = append(returning, g.sqlcolexpr(col))
	}

	delstmt.Table = g.sqlrelid(si.spec.rel.relid)
	delstmt.Using = g.sqlusing(si.spec.join)
	delstmt.Where = g.sqlwhere(si.spec.where)
	delstmt.Returning = returning
	return delstmt
}

func (g *generator) sqlwhere(w *whereblock) (where sql.WhereClause) {
	if w != nil {
		sel := gol.SelectorExpr{X: idrecv, Sel: gol.Ident{w.name}}
		where.SearchCondition, _ = g.sqlsearchcond(w.items, sel, false)
	}
	return where
}

func (g *generator) sqlsearchcond(items []*predicateitem, sel gol.SelectorExpr, parenthesized bool) (list sql.BoolValueExprList, count int) {
	for _, item := range items {
		count += 1

		var x sql.BoolValueExpr
		switch node := item.node.(type) {
		// nested: recurse
		case *nestedpredicate:
			var ncount int

			sel := gol.SelectorExpr{X: sel, Sel: gol.Ident{node.name}}
			x, ncount = g.sqlsearchcond(node.items, sel, true)

			count += ncount - 1

		// 3-arg predicate: build & return
		case *betweenpredicate:
			p := sql.BetweenPredicate{}
			p.Predicand = g.sqlcolref(node.colid)
			if x, ok := node.x.(colid); ok {
				p.LowEnd = g.sqlcolref(x)
			} else {
				// assume node.x is *paramfield
				p.LowEnd = g.sqlparam()
			}
			if y, ok := node.y.(colid); ok {
				p.HighEnd = g.sqlcolref(y)
			} else {
				// assume node.x is *paramfield
				p.HighEnd = g.sqlparam()
			}
			x = p

		// 2-arg predicates: prepare first, then build & return
		case *fieldpredicate, *columnpredicate:
			var (
				lhs   sql.ValueExpr
				rhs   sql.ValueExpr
				pred  predicate
				field string
			)

			// prepare
			switch node := node.(type) {
			case *fieldpredicate:
				pred = node.pred
				lhs = g.sqlcolref(node.colid)
				rhs = g.sqlparam()
				if len(node.modfunc) > 0 {
					li := sql.RoutineInvocation{}
					li.Name = string(node.modfunc)
					li.Args = []sql.ValueExpr{lhs}
					lhs = li

					ri := sql.RoutineInvocation{}
					ri.Name = string(node.modfunc)
					ri.Args = []sql.ValueExpr{rhs}
					rhs = ri
				}

				field = node.field.name // needed for isin/notin predicates
			case *columnpredicate:
				pred = node.pred
				lhs = g.sqlcolref(node.colid)
				if !node.colid2.isempty() {
					rhs = g.sqlcolref(node.colid2)
				} else if len(node.lit) > 0 {
					rhs = sql.Literal{node.lit}
				}
			}

			// build & return
			switch pred {
			case iseq, noteq, noteq2, islt, isgt, islte, isgte:
				p := sql.ComparisonPredicate{}
				p.Cmp = predicate2sqlcmpop[pred]
				p.LPredicand = lhs
				p.RPredicand = rhs
				x = p
			case islike, notlike:
				p := sql.LikePredicate{}
				p.Not = (pred == notlike)
				p.Predicand = lhs
				p.Pattern = rhs
				x = p
			case isilike, notilike:
				p := sql.ILikePredicate{}
				p.Not = (pred == notilike)
				p.Predicand = lhs
				p.Pattern = rhs
				x = p
			case issimilar, notsimilar:
				p := sql.SimilarPredicate{}
				p.Not = (pred == notsimilar)
				p.Predicand = lhs
				p.Pattern = rhs
				x = p
			case isdistinct, notdistinct:
				p := sql.DistinctPredicate{}
				p.Not = (pred == notdistinct)
				p.LPredicand = lhs
				p.RPredicand = rhs
				x = p
			case ismatch, ismatchi, notmatch, notmatchi:
				p := sql.RegexPredicate{}
				p.Op = predicate2sqlregexop[pred]
				p.Predicand = lhs
				p.Pattern = rhs
				x = p
			case isin, notin:
				sx := gol.SelectorExpr{X: sel, Sel: gol.Ident{field}}
				g.insx = append(g.insx, sx)
				g.nparam -= 1  // ordinal param won't be used directly
				g.asvar = true // queryString should be var not const

				num := strconv.Itoa(len(g.insx))
				arg1 := gol.Ident{"len" + num}
				arg2 := gol.BinaryExpr{X: gol.Ident{"pos" + num}, Op: gol.BINARY_ADD, Y: gol.Int(1)}

				call := callinvaluelist
				call.Args = gol.ArgsList{List: []gol.Expr{arg1, arg2}}

				p := sql.InPredicate{}
				p.Not = (pred == notin)
				p.Predicand = lhs
				p.ValueList = sql.HostValue{gol.RawStringImplant{call}}
				x = p
			case istrue, nottrue, isfalse, notfalse, isunknown, notunknown:
				p := sql.TruthPredicate{}
				p.Not = (pred == nottrue || pred == notfalse || pred == notunknown)
				p.Truth = predicate2sqltruth[pred]
				p.Predicand = lhs
				x = p
			case isnull, notnull:
				p := sql.NullPredicate{}
				p.Not = (pred == notnull)
				p.Predicand = lhs
				x = p
			default:
				// no predicate, assume lhs is by itself a boolean value expression
				if p, ok := lhs.(sql.BoolValueExpr); ok {
					x = p
				}
			}
		}

		switch item.op {
		default: // initial
			list.Initial = x
		case booland:
			list.Items = append(list.Items, sql.AND{Operand: x})
		case boolor:
			list.Items = append(list.Items, sql.OR{Operand: x})
		}

	}

	if count > 2 {
		list.ListStyle = true
	}

	list.Parenthesized = parenthesized
	return list, count
}

func (g *generator) sqlorderby(spec *typespec) (order sql.OrderClause) {
	if spec.orderby == nil {
		return order
	}

	for _, item := range spec.orderby.items {
		by := sql.OrderBy{}
		by.Column = g.sqlcolref(item.col)
		by.Desc = (item.dir == orderdesc)
		order.List = append(order.List, by)
	}
	return order
}

func (g *generator) sqlusing(jb *joinblock) (using sql.UsingClause) {
	if jb == nil {
		return using
	}

	using.List = []sql.TableExpr{g.sqlrelid(jb.rel)}
	for _, item := range jb.items {
		var join sql.TableJoin
		join.Type = jointype2sqljointype[item.typ]
		join.Rel = g.sqlrelid(item.rel)
		join.Cond = g.sqljoincond(item.conds)
		using.List = append(using.List, join)
	}
	return using
}

func (g *generator) sqljoin(jb *joinblock) (jc sql.JoinClause) {
	if jb == nil {
		return jc
	}

	for _, item := range jb.items {
		var join sql.TableJoin
		join.Type = jointype2sqljointype[item.typ]
		join.Rel = g.sqlrelid(item.rel)
		join.Cond = g.sqljoincond(item.conds)
		jc.List = append(jc.List, join)
	}
	return jc
}

func (g *generator) sqljoincond(items []*predicateitem) (cond sql.JoinCondition) {
	if len(items) > 0 {
		list, _ := g.sqlsearchcond(items, gol.SelectorExpr{}, false)
		list.ListStyle = false

		cond.SearchCondition = list
	}
	return cond
}

// sqllimit generates and returns an sql.LimitClause based on the given spec's "limit" field.
func (g *generator) sqllimit(spec *typespec) (limit sql.LimitClause) {
	if spec.limit != nil {
		if len(spec.limit.field) > 0 {
			limit.Value = g.sqlparam()
		} else if spec.limit.value > 0 {
			limit.Value = sql.LimitUint(spec.limit.value)
		}
		return limit
	}

	// In case the spec doesn't have a "limit" field, but the relation
	// field handles only a single record (i.e. it's not a slice, etc.)
	// then, by default, generate a `LIMIT 1` clause.
	if r := spec.rel.rec; !r.isarray && !r.isslice && !r.isiter {
		limit.Value = sql.LimitInt(1)
		return limit
	}
	return limit
}

// sqloffset generates and returns an sql.OffsetClause based on the given spec's "offset" field.
func (g *generator) sqloffset(spec *typespec) (offset sql.OffsetClause) {
	if spec.offset != nil {
		if len(spec.offset.field) > 0 {
			offset.Value = g.sqlparam()
		} else if spec.offset.value > 0 {
			offset.Value = sql.OffsetUint(spec.offset.value)
		}
		return offset
	}
	return offset
}

func (g *generator) sqlrelid(id relid) sql.Ident {
	return sql.Ident{
		Name:  sql.Name(id.name),
		Qual:  id.qual,
		Alias: id.alias,
	}
}

func (g *generator) sqlcolexpr(fc *fieldcolumn) sql.ValueExpr {
	id := g.sqlcolref(fc.colid)
	// TODO
	//if f.UseCOALESCE || (f.IsNULLable && canCoalesce && !f.Type.IsPointer) {
	//	coalesce := sqlang.Coalesce{}
	//	coalesce.A = col
	//	coalesce.B = _sql_empty_literal(f)

	//	if (f.ColTypeIsEnum || f.ColTypeName == "uuid") && len(f.COALESCEValue) == 0 {
	//		coalesce.A = sqlang.CastExpr{
	//			X:    coalesce.A,
	//			Type: sqlang.Literal("text"),
	//		}
	//	}
	//	return coalesce
	//}
	return id
}

func (g *generator) sqlcolref(id colid) sql.ColumnReference {
	return sql.ColumnReference{
		Qual: id.qual,
		Name: sql.Name(id.name),
	}
}

func (g *generator) sqlparam() sql.OrdinalParameterSpec {
	g.nparam += 1
	return sql.OrdinalParameterSpec{g.nparam}
}

var predicate2sqlcmpop = map[predicate]sql.CMPOP{
	iseq:   sql.EQUAL,
	noteq:  sql.NOT_EQUAL,
	noteq2: sql.NOT_EQUAL2,
	islt:   sql.LESS_THAN,
	isgt:   sql.GREATER_THAN,
	islte:  sql.LESS_THAN_EQUAL,
	isgte:  sql.GREATER_THAN_EQUAL,
}

var predicate2sqlregexop = map[predicate]sql.REGEXOP{
	ismatch:   sql.MATCH,
	ismatchi:  sql.MATCH_CI,
	notmatch:  sql.NOT_MATCH,
	notmatchi: sql.NOT_MATCH_CI,
}

var predicate2sqltruth = map[predicate]sql.TRUTH{
	isunknown:  sql.UNKNOWN,
	notunknown: sql.UNKNOWN,
	istrue:     sql.TRUE,
	nottrue:    sql.TRUE,
	isfalse:    sql.FALSE,
	notfalse:   sql.FALSE,
}

var jointype2sqljointype = map[jointype]sql.JoinType{
	joinleft:  sql.JoinLeft,
	joinright: sql.JoinRight,
	joinfull:  sql.JoinFull,
	joincross: sql.JoinCross,
}

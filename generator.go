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
	idquery             = gol.Ident{"queryString"}
	idparams            = gol.Ident{"params"}
	idi64               = gol.Ident{"i64"}
	idres               = gol.Ident{"res"}
	idrow               = gol.Ident{"row"}
	idrows              = gol.Ident{"rows"}
	idexec              = gol.Ident{"Exec"}
	idiface             = gol.Ident{"interface{}"}
	iderror             = gol.Ident{"error"}
	idafterscan         = gol.Ident{"AfterScan"}
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
	nparam int // number of parameters
}

func (g *generator) run(pkgname string) error {
	g.file.PkgName = pkgname
	g.file.Preamble = gol.LineComment{filepreamble}
	g.file.Imports = gol.ImportDecl{{Path: gosqlimport}}

	for _, si := range g.infos {
		g.nparam = 0

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

	fn.Body.Add(g.querybuild(si))
	fn.Body.Add(gol.NL{})
	fn.Body.Add(g.querydefaults(si))
	fn.Body.Add(g.queryexec(si))
	fn.Body.Add(g.returnstmt(si))
	return fn
}

func (g *generator) querybuild(si *specinfo) (stmt gol.Stmt) {
	decl := gol.GenDecl{Token: gol.DECL_CONST}
	decl.Specs = []gol.Spec{gol.ValueSpec{
		Names:       []gol.Ident{idquery},
		Values:      []gol.Expr{gol.RawStringNode{g.sqlnode(si)}},
		LineComment: gol.LineComment{" `"},
	}}

	if len(si.spec.filter) > 0 {
		decl.Token = gol.DECL_VAR

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
	if len(si.spec.filter) > 0 {
		args.AddExprs(idparams)
		args.Ellipsis = true
	} else {
		args.AddExprs(g.queryargs(si.spec)...)
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
			for _, fe := range item.field.path {
				fx = gol.SelectorExpr{X: fx, Sel: gol.Ident{fe.name}}

				fieldkey += fe.name
				if fe.ispointer && !pfieldhandled[fieldkey] {
					if fe.isimported {
						g.addimport(fe.typepkgpath, fe.typepkgname, fe.typepkglocal)
					}

					// initialize nested pointer field
					init := gol.AssignStmt{Token: gol.ASSIGN}
					init.Lhs = []gol.Expr{fx}
					init.Rhs = []gol.Expr{gol.CallExpr{Fun: idnew, Args: gol.ArgsList{List: []gol.Expr{g.fieldelemtype(fe)}}}}
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
			asn := gol.AssignStmt{Token: gol.ASSIGN_DEFINE}
			asn.Lhs = []gol.Expr{iderr}
			asn.Rhs = []gol.Expr{gol.CallExpr{
				Fun: gol.SelectorExpr{
					X:   gol.SelectorExpr{X: idrecv, Sel: gol.Ident{fieldname}},
					Sel: gol.Ident{rec.itermethod},
				},
				Args: gol.ArgsList{List: []gol.Expr{gol.Ident{"v"}}},
			}}

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

func (g *generator) queryargs(spec *typespec) (args []gol.Expr) {
	if spec.where != nil && len(spec.where.items) > 0 {
		type loopstate struct {
			where *whereblock      // the current iteration whereblock
			idx   int              // keeps track of the item index
			sx    gol.SelectorExpr // the selector expression for the current whereblock
		}

		sx := gol.SelectorExpr{X: idrecv, Sel: gol.Ident{spec.where.name}}
		stack := []*loopstate{{where: spec.where, sx: sx}}

	stackloop:
		for len(stack) > 0 {
			loop := stack[len(stack)-1]
			for loop.idx < len(loop.where.items) {
				item := loop.where.items[loop.idx]
				loop.idx++

				switch node := item.node.(type) {
				case *wherefield:
					sx := gol.SelectorExpr{X: loop.sx, Sel: gol.Ident{node.name}}
					args = append(args, sx)
				case *wherebetween:
					if x, ok := node.x.(*varinfo); ok {
						sx := gol.SelectorExpr{X: loop.sx, Sel: gol.Ident{node.name}}
						sx = gol.SelectorExpr{X: sx, Sel: gol.Ident{x.name}}
						args = append(args, sx)
					}
					if y, ok := node.y.(*varinfo); ok {
						sx := gol.SelectorExpr{X: loop.sx, Sel: gol.Ident{node.name}}
						sx = gol.SelectorExpr{X: sx, Sel: gol.Ident{y.name}}
						args = append(args, sx)
					}
				case *whereblock:
					loop2 := new(loopstate)
					loop2.where = node
					loop2.sx = gol.SelectorExpr{X: loop.sx, Sel: gol.Ident{node.name}}
					stack = append(stack, loop2)
					continue stackloop
				case *wherecolumn:
					// nothing to do
				}
			}
			stack = stack[:len(stack)-1]
		}
	}

	if spec.limit != nil && len(spec.limit.field) > 0 {
		sx := gol.SelectorExpr{X: idrecv, Sel: gol.Ident{spec.limit.field}}
		args = append(args, sx)
	}
	if spec.offset != nil && len(spec.offset.field) > 0 {
		sx := gol.SelectorExpr{X: idrecv, Sel: gol.Ident{spec.offset.field}}
		args = append(args, sx)
	}
	return args
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
				for _, fe := range item.field.path {
					fx = gol.SelectorExpr{X: fx, Sel: gol.Ident{fe.name}}

					fieldkey += fe.name
					if fe.ispointer && !pfieldhandled[fieldkey] {
						if fe.isimported {
							g.addimport(fe.typepkgpath, fe.typepkgname, fe.typepkglocal)
						}

						// initialize nested pointer field
						init := gol.AssignStmt{Token: gol.ASSIGN}
						init.Lhs = []gol.Expr{fx}
						init.Rhs = []gol.Expr{gol.CallExpr{Fun: idnew, Args: gol.ArgsList{List: []gol.Expr{g.fieldelemtype(fe)}}}}
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
				for _, fe := range item.field.path {
					fx = gol.SelectorExpr{X: fx, Sel: gol.Ident{fe.name}}

					fieldkey += fe.name
					if fe.ispointer && !pfieldhandled[fieldkey] {
						if fe.isimported {
							g.addimport(fe.typepkgpath, fe.typepkgname, fe.typepkglocal)
						}

						// initialize nested pointer field
						init := gol.AssignStmt{Token: gol.ASSIGN}
						init.Lhs = []gol.Expr{fx}
						init.Rhs = []gol.Expr{gol.CallExpr{Fun: idnew, Args: gol.ArgsList{List: []gol.Expr{g.fieldelemtype(fe)}}}}
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

func (g *generator) fieldelemtype(fe *fieldelem) gol.Expr {
	id := gol.Ident{fe.typename}
	if fe.isimported {
		return gol.SelectorExpr{X: gol.Ident{fe.typepkgname}, Sel: id}
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
	selstmt.Columns = g.sqlcolumns(si.info.output)
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
	delstmt.Table = g.sqlrelid(si.spec.rel.relid)
	delstmt.Using = g.sqlusing(si.spec.join)
	delstmt.Where = g.sqlwhere(si.spec.where)
	delstmt.Returning = g.sqlreturning(si.info.output)
	return delstmt
}

// sqlcolumns builds and returns a list of sql.ValueExpr values.
func (g *generator) sqlcolumns(columns []*fieldcolumn) (list sql.ValueExprList) {
	for _, c := range columns {
		list = append(list, g.sqlcolexpr(c))
	}
	return list
}

func (g *generator) sqlwhere(w *whereblock) (where sql.WhereClause) {
	if w == nil {
		return where
	}

	where.SearchCondition = g.sqlsearchcond(w, false)
	return where
}

func (g *generator) sqlsearchcond(w *whereblock, nested bool) sql.BoolValueExpr {
	if len(w.items) == 1 {
		return g.sqlboolexpr(w.items[0].node)
	}

	var list sql.BoolValueExprList
	list.Parenthesized = nested
	list.Initial = g.sqlboolexpr(w.items[0].node)

	for _, item := range w.items[1:] {
		x := g.sqlboolexpr(item.node)
		switch item.op {
		case booland:
			list.Items = append(list.Items, sql.AND{Operand: x})
		case boolor:
			list.Items = append(list.Items, sql.OR{Operand: x})
		}

	}
	return list
}

func (g *generator) sqlboolexpr(node interface{}) sql.BoolValueExpr {
	switch n := node.(type) {
	// nested: build & return
	case *whereblock:
		return g.sqlsearchcond(n, true)

	// 3-arg predicate: build & return
	case *wherebetween:
		p := sql.BetweenPredicate{}
		p.Predicand = g.sqlcolref(n.colid)
		if x, ok := n.x.(colid); ok {
			p.LowEnd = g.sqlcolref(x)
		} else {
			// assume n.x is *varinfo
			p.LowEnd = g.sqlparam()
		}
		if y, ok := n.y.(colid); ok {
			p.HighEnd = g.sqlcolref(y)
		} else {
			// assume n.x is *varinfo
			p.HighEnd = g.sqlparam()
		}
		return p

	// 2-arg predicates: prepare first, then build & return
	case *wherefield, *wherecolumn, *joincond:
		var (
			lhs sql.ColumnReference
			rhs sql.ValueExpr
			cmp cmpop
		)

		// prepare
		switch n := node.(type) {
		case *wherefield:
			cmp = n.cmp
			lhs = g.sqlcolref(n.colid)
			rhs = g.sqlparam()
		case *wherecolumn:
			cmp = n.cmp
			lhs = g.sqlcolref(n.colid)
			if !n.colid2.isempty() {
				rhs = g.sqlcolref(n.colid2)
			} else if len(n.lit) > 0 {
				rhs = sql.Literal{n.lit}
			}
		case *joincond:
			cmp = n.cmp
			lhs = g.sqlcolref(n.col1)
			if !n.col2.isempty() {
				rhs = g.sqlcolref(n.col2)
			} else if len(n.lit) > 0 {
				rhs = sql.Literal{n.lit}
			}
		}

		// build & return
		switch cmp {
		case cmpeq, cmpne, cmpne2, cmplt, cmpgt, cmple, cmpge:
			p := sql.ComparisonPredicate{}
			p.Cmp = cmpop2sqlnode[cmp]
			p.LPredicand = lhs
			p.RPredicand = rhs
			return p
		case cmpislike, cmpnotlike:
			p := sql.LikePredicate{}
			p.Not = (cmp == cmpnotlike)
			p.Predicand = lhs
			p.Pattern = rhs
			return p
		case cmpisilike, cmpnotilike:
			p := sql.ILikePredicate{}
			p.Not = (cmp == cmpnotilike)
			p.Predicand = lhs
			p.Pattern = rhs
			return p
		case cmpissimilar, cmpnotsimilar:
			p := sql.SimilarPredicate{}
			p.Not = (cmp == cmpnotsimilar)
			p.Predicand = lhs
			p.Pattern = rhs
			return p
		case cmpisdistinct, cmpnotdistinct:
			p := sql.DistinctPredicate{}
			p.Not = (cmp == cmpnotdistinct)
			p.LPredicand = lhs
			p.RPredicand = rhs
			return p
		case cmprexp, cmprexpi, cmpnotrexp, cmpnotrexpi:
			p := sql.RegexPredicate{}
			p.Op = cmpop2sqlregexop[cmp]
			p.Predicand = lhs
			p.Pattern = rhs
			return p
		case cmpisin, cmpnotin:
			p := sql.InPredicate{}
			p.Not = (cmp == cmpnotin)
			p.Predicand = lhs
			p.ValueList = rhs
			return p
		case cmpistrue, cmpnottrue, cmpisfalse, cmpnotfalse, cmpisunknown, cmpnotunknown:
			p := sql.TruthPredicate{}
			p.Not = (cmp == cmpnottrue || cmp == cmpnotfalse || cmp == cmpnotunknown)
			p.Truth = cmpop2sqltruth[cmp]
			p.Predicand = lhs
			return p
		case cmpisnull, cmpnotnull:
			p := sql.NullPredicate{}
			p.Not = (cmp == cmpnotnull)
			p.Predicand = lhs
			return p
		default:
			return lhs // no comparison
		}
	}

	return nil
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
		join.Type = jointype2sqlnode[item.typ]
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
		join.Type = jointype2sqlnode[item.typ]
		join.Rel = g.sqlrelid(item.rel)
		join.Cond = g.sqljoincond(item.conds)
		jc.List = append(jc.List, join)
	}
	return jc
}

func (g *generator) sqljoincond(cc []*joincond) (cond sql.JoinCondition) {
	if len(cc) == 0 {
		return cond
	}

	if len(cc) == 1 {
		cond.SearchCondition = g.sqlboolexpr(cc[0])
		return cond
	}

	var list sql.BoolValueExprList
	list.Initial = g.sqlboolexpr(cc[0])

	for _, c := range cc[1:] {
		x := g.sqlboolexpr(c)
		switch c.op {
		case booland:
			list.Items = append(list.Items, sql.AND{Operand: x})
		case boolor:
			list.Items = append(list.Items, sql.OR{Operand: x})
		}
	}

	cond.SearchCondition = list
	return cond
}

func (g *generator) sqlreturning(fcs []*fieldcolumn) (returning sql.ReturningClause) {
	if fcs == nil {
		return returning
	}

	for _, fc := range fcs {
		returning = append(returning, g.sqlcolref(fc.colid))
	}
	return returning
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

var cmpop2sqlnode = map[cmpop]sql.CMPOP{
	cmpeq:  sql.EQUAL,
	cmpne:  sql.NOT_EQUAL,
	cmpne2: sql.NOT_EQUAL2,
	cmplt:  sql.LESS_THAN,
	cmpgt:  sql.GREATER_THAN,
	cmple:  sql.LESS_THAN_EQUAL,
	cmpge:  sql.GREATER_THAN_EQUAL,
}

var cmpop2sqlregexop = map[cmpop]sql.REGEXOP{
	cmprexp:     sql.MATCH,
	cmprexpi:    sql.MATCH_CI,
	cmpnotrexp:  sql.NOT_MATCH,
	cmpnotrexpi: sql.NOT_MATCH_CI,
}

var cmpop2sqltruth = map[cmpop]sql.TRUTH{
	cmpisunknown:  sql.UNKNOWN,
	cmpnotunknown: sql.UNKNOWN,
	cmpistrue:     sql.TRUE,
	cmpnottrue:    sql.TRUE,
	cmpisfalse:    sql.FALSE,
	cmpnotfalse:   sql.FALSE,
}

var jointype2sqlnode = map[jointype]sql.JoinType{
	joinleft:  sql.JoinLeft,
	joinright: sql.JoinRight,
	joinfull:  sql.JoinFull,
	joincross: sql.JoinCross,
}

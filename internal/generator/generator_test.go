package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/frk/compare"
	"github.com/frk/gosql/internal/analysis"
	"github.com/frk/gosql/internal/parser"
	"github.com/frk/gosql/internal/postgres"
)

func TestGenerator(t *testing.T) {
	db := &postgres.TestDB{}
	if err := db.Init(); err != nil {
		panic(err)
	}
	defer db.Close()

	type gconf struct {
		fcktag  string
		fcksep  string
		fckbase bool
		quote   bool
	}
	type testcase struct {
		filename string
		gconf    *gconf
	}

	tests := []struct {
		skip      bool
		dirname   string
		testcases []testcase
	}{{
		//skip:    true,
		dirname: "delete",
		testcases: []testcase{
			{filename: "all_directive"},
			{filename: "datatype_1"},
			{filename: "datatype_2"},
			{filename: "filter"},
			{filename: "result_iterator_afterscan"},
			{filename: "result_iterator_errorhandler"},
			{filename: "result_iterator_errorinfohandler"},
			{filename: "result_iterator"},
			{filename: "result_single_afterscan"},
			{filename: "result_single_errorhandler"},
			{filename: "result_single_errorinfohandler"},
			{filename: "result_single"},
			{filename: "result_slice_afterscan"},
			{filename: "result_slice"},
			{filename: "returning_iterator_afterscan"},
			{filename: "returning_iterator"},
			{filename: "returning_single_afterscan"},
			{filename: "returning_single_all"},
			{filename: "returning_single_collist"},
			{filename: "returning_slice_all"},
			{filename: "returning_slice_afterscan"},
			{filename: "returning_slice_collist"},
			{filename: "returning_slice_context"},
			{filename: "returning_slice_errorhandler"},
			{filename: "returning_slice_errorinfohandler"},
			{filename: "rowsaffected"},
			{filename: "rowsaffected_errorhandler"},
			{filename: "rowsaffected_errorinfohandler"},
			{filename: "using_join_block_1"},
			{filename: "using_join_block_2"},
			{filename: "where_block_1"},
			{filename: "where_block_2"},
		},
	}, {
		//skip:    true,
		dirname: "select",
		testcases: []testcase{
			{filename: "afterscan_single"},
			{filename: "afterscan_slice"},
			{filename: "coalesce_table"},
			{filename: "count_basic"},
			{filename: "count_filter"},
			{filename: "count_where"},
			{filename: "exists_filter"},
			{filename: "exists_where"},
			{filename: "iterator_func"},
			{filename: "iterator_func_errorhandler"},
			{filename: "iterator_iface"},
			{filename: "filter_slice"},
			{filename: "filter_iterator"},
			{filename: "joinblock_slice"},
			{filename: "limit_directive"},
			{filename: "limit_field_default"},
			{filename: "limit_field"},
			{filename: "notexists_where"},
			{filename: "notexists_filter"},
			{filename: "offset_directive"},
			{filename: "offset_field_default"},
			{filename: "offset_field"},
			{filename: "orderby_directive"},
			{filename: "record_nested_single"},
			{filename: "record_nested_slice"},
			{filename: "whereblock_array_comparison1"},
			{filename: "whereblock_array_comparison2"},
			{filename: "whereblock_array_comparison3"},
			{filename: "whereblock_between"},
			{filename: "whereblock_isin"},
			{filename: "whereblock_isin2"},
			{filename: "whereblock_isin3"},
			{filename: "whereblock_modifierfunc_single"},
			{filename: "whereblock_nested"},
			{filename: "whereblock_single"},
			{filename: "whereblock_single2"},
			{filename: "whereblock_slice"},
		},
	}, {
		//skip:    true,
		dirname: "insert",
		testcases: []testcase{
			{filename: "basic_single"},
			{filename: "basic_single2"},
			{filename: "basic_slice"},
			{filename: "default_all_returning_single"},
			{filename: "default_all_returning_slice"},
			{filename: "default_all_single"},
			{filename: "default_all_slice"},
			{filename: "default_single"},
			{filename: "default_slice"},
			{filename: "json_single"},
			{filename: "json_slice"},
			{filename: "onconflict_column_ignore_single_1"},
			{filename: "onconflict_column_ignore_single_2"},
			{filename: "onconflict_column_update_single_1"},
			{filename: "onconflict_column_update_returning_slice"},
			{filename: "onconflict_constraint_ignore_single_1"},
			{filename: "onconflict_ignore_single"},
			{filename: "onconflict_ignore_slice"},
			{filename: "onconflict_index_ignore_single_1"},
			{filename: "onconflict_index_ignore_single_2"},
			{filename: "onconflict_index_update_single_1"},
			{filename: "onconflict_index_update_returning_slice"},
			{filename: "result_afterscan_iterator"},
			{filename: "result_afterscan_single"},
			{filename: "result_afterscan_slice"},
			{filename: "result_basic_iterator"},
			{filename: "result_basic_single"},
			{filename: "result_basic_slice"},
			{filename: "result_errorhandler_iterator"},
			{filename: "result_errorhandler_single"},
			{filename: "result_errorinfohandler_iterator"},
			{filename: "result_errorinfohandler_single"},
			{filename: "result_json_single"},
			{filename: "result_json_slice"},
			{filename: "returning_afterscan_single"},
			{filename: "returning_afterscan_slice"},
			{filename: "returning_all_json_single"},
			{filename: "returning_all_json_slice"},
			{filename: "returning_all_single"},
			{filename: "returning_all_slice"},
			{filename: "returning_collist_single"},
			{filename: "returning_collist_slice"},
			{filename: "returning_context_single"},
			{filename: "returning_context_slice"},
			{filename: "returning_errorhandler_slice"},
			{filename: "returning_errorinfohandler_slice"},
			{filename: "rowsaffected_errorhandler_single"},
			{filename: "rowsaffected_errorinfohandler_single"},
			{filename: "rowsaffected_single"},
		},
	}, {
		//skip:    true,
		dirname: "update",
		testcases: []testcase{
			{filename: "all_single"},
			{filename: "filter_single"},
			{filename: "filter_result_slice"},
			{filename: "fromblock_basic_single"},
			{filename: "fromblock_join_single"},
			{filename: "pkey_composite_single"},
			{filename: "pkey_composite_slice"},
			{filename: "pkey_single"},
			{filename: "pkey_slice"},
			{filename: "pkey_returning_all_single"},
			{filename: "whereblock_basic_single_1"},
			{filename: "whereblock_basic_single_2"},
			{filename: "whereblock_result_slice"},
			{filename: "whereblock_returning_all_single"},
		},
	}, {
		//skip:    true,
		dirname: "filter",
		testcases: []testcase{
			{filename: "alias"},
			{filename: "basic"},
			{filename: "basic2", gconf: &gconf{"json", ".", false, true}},
			{filename: "nested"},
			{filename: "textsearch"},
		},
	}, {
		//skip:    true,
		dirname: "pgsql",
		testcases: []testcase{
			{filename: "insert_basic"},
			{filename: "insert_array"},
		},
	}}

	for _, tt := range tests {
		if tt.skip {
			continue
		}

		pkgs, err := parser.Parse("../testdata/generator/"+tt.dirname, false, nil)
		if err != nil {
			t.Fatal(err)
		}
		pkg := pkgs[0]

		for _, tc := range tt.testcases {
			t.Run(tt.dirname+"/"+tc.filename, func(t *testing.T) {
				tinfos := []*postgres.TargetInfo{}
				fileprefix := "../testdata/generator/" + tt.dirname + "/" + tc.filename

				f, err := getFile(pkg, fileprefix+"_in.go")
				if err != nil {
					t.Fatal(err)
				}

				for _, target := range f.Targets {
					// analyze
					ainfo := &analysis.Info{}
					tstruct, err := analysis.Run(pkg.Fset, target.Named, target.Pos, ainfo)
					if err != nil {
						t.Error(err)
						return
					}

					// type check
					targInfo, err := postgres.Check(db.DB, tstruct, ainfo)
					if err != nil {
						t.Error(err)
						return
					}

					tinfos = append(tinfos, targInfo)
				}

				buf := new(bytes.Buffer)
				conf := Config{FilterColumnKeySeparator: ".", QuoteIdentifiers: true} // default
				if tc.gconf != nil {
					conf.FilterColumnKeyTag = tc.gconf.fcktag
					conf.FilterColumnKeySeparator = tc.gconf.fcksep
					conf.FilterColumnKeyBase = tc.gconf.fckbase
					conf.QuoteIdentifiers = tc.gconf.quote
				}
				if err := Write(buf, pkg.Name, tinfos, conf); err != nil {
					t.Error(err)
					return
				}

				got := string(formatBytes(buf))

				out, err := ioutil.ReadFile(fileprefix + "_out.go")
				if err != nil {
					t.Fatal(err)
				}
				want := string(out)

				// compare
				if err := compare.Compare(got, want); err != nil {
					t.Error(err)
				}
			})
		}
	}
}

// helper method...
func getFile(p *parser.Package, filename string) (*parser.File, error) {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	for _, f := range p.Files {
		if f.Path == filename {
			return f, nil
		}
	}
	return nil, fmt.Errorf("file not found: %q", filename)
}

func formatBytes(buf *bytes.Buffer) []byte {
	src, err := format.Source(buf.Bytes())
	if err != nil {
		log.Printf("format error: %s", err)
		return buf.Bytes()
	}
	return src
}

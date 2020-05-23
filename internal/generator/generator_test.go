package generator

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"log"
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

	tests := []struct {
		skip      bool
		dirname   string
		filenames []string
	}{{
		//skip:    true,
		dirname: "delete",
		filenames: []string{
			"all_directive",
			"datatype_1",
			"datatype_2",
			"filter",
			"result_iterator_afterscan",
			"result_iterator_errorhandler",
			"result_iterator_errorinfohandler",
			"result_iterator",
			"result_single_afterscan",
			"result_single_errorhandler",
			"result_single_errorinfohandler",
			"result_single",
			"result_slice_afterscan",
			"result_slice",
			"returning_iterator_afterscan",
			"returning_iterator",
			"returning_single_afterscan",
			"returning_single_all",
			"returning_single_collist",
			"returning_slice_afterscan",
			"returning_slice_errorhandler",
			"returning_slice_errorinfohandler",
			"returning_slice_all",
			"returning_slice_collist",
			"rowsaffected",
			"rowsaffected_errorhandler",
			"rowsaffected_errorinfohandler",
			"using_join_block_1",
			"using_join_block_2",
			"where_block_1",
			"where_block_2",
		},
	}, {
		//skip:    true,
		dirname: "select",
		filenames: []string{
			"afterscan_single",
			"afterscan_slice",
			"count_basic",
			"count_filter",
			"count_where",
			"exists_filter",
			"exists_where",
			"iterator_func",
			"iterator_func_errorhandler",
			"iterator_iface",
			"filter_slice",
			"filter_iterator",
			"joinblock_slice",
			"limit_directive",
			"limit_field_default",
			"limit_field",
			"notexists_where",
			"notexists_filter",
			"offset_directive",
			"offset_field_default",
			"offset_field",
			"orderby_directive",
			"record_nested_single",
			"record_nested_slice",
			"whereblock_array_comparison1",
			"whereblock_array_comparison2",
			"whereblock_array_comparison3",
			"whereblock_between",
			"whereblock_isin",
			"whereblock_isin2",
			"whereblock_isin3",
			"whereblock_modifierfunc_single",
			"whereblock_nested",
			"whereblock_single",
			"whereblock_slice",
		},
	}, {
		//skip:    true,
		dirname: "insert",
		filenames: []string{
			"basic_single",
			"basic_slice",
			"default_all_returning_single",
			"default_all_returning_slice",
			"default_all_single",
			"default_all_slice",
			"default_single",
			"default_slice",
			"json_single",
			"json_slice",
			"onconflict_column_ignore_single_1",
			"onconflict_column_ignore_single_2",
			"onconflict_column_update_single_1",
			"onconflict_column_update_returning_slice",
			"onconflict_constraint_ignore_single_1",
			"onconflict_ignore_single",
			"onconflict_ignore_slice",
			"onconflict_index_ignore_single_1",
			"onconflict_index_ignore_single_2",
			"onconflict_index_update_single_1",
			"onconflict_index_update_returning_slice",
			"result_afterscan_iterator",
			"result_afterscan_single",
			"result_afterscan_slice",
			"result_basic_iterator",
			"result_basic_single",
			"result_basic_slice",
			"result_errorhandler_iterator",
			"result_errorhandler_single",
			"result_errorinfohandler_iterator",
			"result_errorinfohandler_single",
			"result_json_single",
			"result_json_slice",
			"returning_afterscan_single",
			"returning_afterscan_slice",
			"returning_all_json_single",
			"returning_all_json_slice",
			"returning_all_single",
			"returning_all_slice",
			"returning_collist_single",
			"returning_collist_slice",
			"returning_errorhandler_slice",
			"returning_errorinfohandler_slice",
			"rowsaffected_errorhandler_single",
			"rowsaffected_errorinfohandler_single",
			"rowsaffected_single",
		},
	}, {
		//skip:    true,
		dirname: "update",
		filenames: []string{
			"all_single",
			"filter_single",
			"filter_result_slice",
			"fromblock_basic_single",
			"fromblock_join_single",
			"pkey_composite_single",
			"pkey_composite_slice",
			"pkey_single",
			"pkey_slice",
			"pkey_returning_all_single",
			"whereblock_basic_single_1",
			"whereblock_basic_single_2",
			"whereblock_result_slice",
			"whereblock_returning_all_single",
		},
	}, {
		//skip:    true,
		dirname: "filter",
		filenames: []string{
			"alias",
			"basic",
			"nested",
			"textsearch",
		},
	}, {
		//skip:    true,
		dirname: "pgsql",
		filenames: []string{
			"insert_basic",
			"insert_array",
		},
	}}

	for _, tt := range tests {
		if tt.skip {
			continue
		}

		dir, err := parser.ParseDirectory("../testdata/generator/" + tt.dirname)
		if err != nil {
			t.Fatal(err)
		}

		for _, filename := range tt.filenames {
			t.Run(tt.dirname+"/"+filename, func(t *testing.T) {
				tinfos := []*postgres.TargetInfo{}
				fileprefix := "../testdata/generator/" + tt.dirname + "/" + filename

				f := parser.FileWithTargetTypes(dir, fileprefix+"_in.go")
				for _, named := range f.Targets {
					// analyze
					ainfo := &analysis.Info{}
					tstruct, err := analysis.Run(dir.FileSet, named, ainfo)
					if err != nil {
						t.Error(err)
						return
					}

					// type check
					tinfo, err := postgres.Check(db.DB, tstruct, ainfo)
					if err != nil {
						t.Error(err)
						return
					}

					tinfos = append(tinfos, tinfo)
				}

				buf := new(bytes.Buffer)
				conf := Config{FilterKeySeparator: ".", QuotedIdents: true}
				if err := Write(buf, dir.Package.Name, tinfos, conf); err != nil {
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

func formatBytes(buf *bytes.Buffer) []byte {
	src, err := format.Source(buf.Bytes())
	if err != nil {
		log.Printf("format error: %s", err)
		return buf.Bytes()
	}
	return src
}

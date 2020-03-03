package gosql

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"log"
	"testing"

	"github.com/frk/compare"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		dirname   string
		filenames []string
	}{{
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
	}}

	// TODO insert
	// TODO update
	// TODO filter (with textsearch)

	for _, tt := range tests {
		cmd := new(command)
		cmd.pg = testdb.pg

		dir, err := cmd.parsedir("./testdata/generator/" + tt.dirname)
		if err != nil {
			t.Fatal(err)
		}

		for _, filename := range tt.filenames {
			t.Run(filename, func(t *testing.T) {
				fileprefix := "testdata/generator/" + tt.dirname + "/" + filename

				f := cmd.aggtypes(dir, fileprefix+"_in.go")
				buf, err := cmd.run(f)
				if err != nil {
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

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
		filename string
	}{
		// deletes
		{"delete_with_all_directive"},
		//{"delete_with_where_block_1"},

		// selects
		// {"select_with_where_block"},
	}

	cmd := new(command)
	cmd.pg = testdb.pg

	dir, err := cmd.parsedir("./testdata/generator")
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			fileprefix := "testdata/generator/" + tt.filename

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

func formatBytes(buf *bytes.Buffer) []byte {
	src, err := format.Source(buf.Bytes())
	if err != nil {
		log.Printf("format error: %s", err)
		return buf.Bytes()
	}
	return src
}

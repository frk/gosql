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
		{"delete_with_where_block_1"},

		// selects
		// {"select_with_where_block"},
	}

	dir, err := parsedir("./testdata/generator")
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			fileprefix := "testdata/generator/" + tt.filename

			g := &generator{file: fileprefix + "_in.go", dir: dir, pg: testdb.pg}
			if err := g.run(); err != nil {
				t.Error(err)
			}
			got := string(formatBytes(&g.buf))

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

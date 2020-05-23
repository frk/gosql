package postgres

import (
	"os"
	"testing"
)

var testdb *TestDB

func TestMain(m *testing.M) {
	var exitcode int

	func() { // use a func wrapper so we can rely on defer
		testdb = new(TestDB)
		defer testdb.Close()

		if err := testdb.Init(); err != nil {
			panic(err)
		}

		exitcode = m.Run()
	}()
	//
	os.Exit(exitcode)
}

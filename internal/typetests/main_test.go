package typetests

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	var exitcode int

	func() { // use a func wrapper so we can rely on defer
		libpq = new(libpqtest)
		defer libpq.close()

		if err := libpq.init(); err != nil {
			panic(err)
		}

		exitcode = m.Run()
	}()

	os.Exit(exitcode)
}

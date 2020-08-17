package termutil

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	s := new(String)
	s.Color(WHITE, FG_HI, UNDERLINED)
	s.Writeln("hello %d", 123)
	s.Color(RED, FG_HI, BOLD)
	s.Writeln("hello %d", 123)
	fmt.Println(s)
}

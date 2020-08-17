package termutil

import (
	"fmt"
)

type TextMode uint

const (
	PLAIN TextMode = iota
	BOLD
	_
	ITALIC
	UNDERLINED
	_ // blinking (slow)
	_ // blinking (fast)
	_ // reverse
	_ // hide
	CROSS_OUT
)

type FontMode uint

const (
	FG    FontMode = 3
	BG    FontMode = 4
	FG_HI FontMode = 9
	BG_HI FontMode = 10
)

type ColorMode uint

const (
	BLACK ColorMode = iota
	RED
	GREEN
	YELLOW
	BLUE
	PURPLE
	CYAN
	WHITE
)

type String struct {
	value string
	fg    string
	bg    string
}

func (s *String) String() string {
	return s.value
}

func (s *String) Write(format string, a ...interface{}) {
	v := fmt.Sprintf(format, a...)
	if len(s.fg) > 0 {
		v = fmt.Sprintf("%s"+v+"\033[0m", s.fg)
	}
	if len(s.bg) > 0 {
		v = fmt.Sprintf("%s"+v+"\033[0m", s.bg)
	}

	s.value += v
	s.fg, s.bg = "", ""
}

func (s *String) Writeln(format string, a ...interface{}) {
	s.Write(format, a...)
	s.value += "\n"
}

func (s *String) Color(color ColorMode, font FontMode, text TextMode) {
	switch font {
	case FG, FG_HI:
		s.fg = fmt.Sprintf("\033[%d;%d%dm", text, font, color)
	case BG, BG_HI:
		s.bg = fmt.Sprintf("\033[%d;%d%dm", text, font, color)
	}
}

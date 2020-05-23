package writer

import (
	"io"
)

// NewWriter returns a new instance of Writer that wraps around
// the given io.Writer and delegates all of its Write calls to it.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

// Writer is a wrapper around io.Writer that provides a couple methods
// that are intended to help with code generation.
type Writer struct {
	w   io.Writer
	nl  bool   // insert new line on next write
	ind string // indentation
	err error
}

// Write writes the given string s to the underlying io.Writer.
func (w *Writer) Write(s string) {
	if w.err != nil {
		return
	}

	if w.nl {
		w.nl = false
		if _, w.err = w.w.Write([]byte("\n" + w.ind)); w.err != nil {
			return
		}
	}

	_, w.err = w.w.Write([]byte(s))
}

// NewLine tells the Writer to add a new line the next time the Write method is called.
func (w *Writer) NewLine() {
	w.nl = true
}

// NoNewLine tells the Writer to *not* add a new line the next time the Write method is called.
func (w *Writer) NoNewLine() {
	w.nl = false
}

// Indent adds one level to the Writer's indentation setting.
func (w *Writer) Indent() {
	w.ind += "\t"
}

// NoIndent removes the Writer's indentation setting.
func (w *Writer) NoIndent() {
	w.ind = ""
}

// Unindent removes one level from the Writer's indentation setting.
func (w *Writer) Unindent() {
	if l := len(w.ind); l > 0 {
		w.ind = w.ind[:l-1]
	}
}

// Err returns the error encountered by the Writer, if any.
func (w *Writer) Err() error {
	return w.err
}

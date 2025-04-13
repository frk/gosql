package analysis

import (
	"github.com/frk/tagutil"
)

type FieldMode uint

const (
	mode_select FieldMode = 1 << iota
	mode_filter
	mode_insert
	mode_update
	mode_notnull

	mode_read    = mode_select | mode_filter
	mode_write   = mode_insert | mode_update
	mode_default = mode_read | mode_write
)

func (m FieldMode) has(v FieldMode) bool {
	return (m & v) != 0
}

func parseFieldMode(t tagutil.Tag) (m FieldMode) {
	m = mode_default
	tag, ok := t["sql"]
	if !ok || len(tag) <= 1 {
		return m
	}

	for _, opt := range tag[1:] {
		////////////////////////////////////////////////////////////////
		// handle deprecated field options
		switch opt {
		case "ro": // read-only?
			m &= ^mode_write

		case "wo": // write-only?
			m &= ^mode_read

		case "xf": // exclude-filter?
			m &= ^mode_filter

		case "nn": // treat as NOT NULL?
			m |= mode_notnull
		}

		////////////////////////////////////////////////////////////////
		// handle new options
		switch opt[0] {
		case '-':
			for _, c := range opt[1:] {
				switch c {
				case 'r': // "-r" => "can't read"
					m &= ^mode_read
				case 'w': // "-w" => "can't write"
					m &= ^mode_write
				case 's': // "-s" => "can't select"
					m &= ^mode_select
				case 'f': // "-f" => "can't filter"
					m &= ^mode_filter
				case 'i': // "-i" => "can't insert"
					m &= ^mode_insert
				case 'u': // "-u" => "can't update"
					m &= ^mode_update

				case 'N': // "-N" => "NOT NULL" (i.e., not nullable)
					m |= mode_notnull
				}
			}
		case '+':
			for _, c := range opt[1:] {
				switch c {
				case 'r': // "+r" => "can read"
					m |= mode_read
				case 'w': // "+w" => "can write"
					m |= mode_write
				case 's': // "+s" => "can select"
					m |= mode_select
				case 'f': // "+f" => "can filter"
					m |= mode_filter
				case 'i': // "+i" => "can insert"
					m |= mode_insert
				case 'u': // "+u" => "can update"
					m |= mode_update

				case 'N': // "+N" => "NULL" (i.e., nullable)
					m &= ^mode_notnull
				}
			}
		}
	}
	return m
}

////////////////////////////////////////////////////////////////////////////////

func (m FieldMode) IsReadOnly() bool {
	return m.CanSelect() && !m.CanInsert() && !m.CanUpdate()
}

func (m FieldMode) IsWriteOnly() bool {
	return (m.CanInsert() || m.CanUpdate()) && !m.CanSelect()
}

func (m FieldMode) CanInsert() bool {
	return m.has(mode_insert)
}

func (m FieldMode) CanUpdate() bool {
	return m.has(mode_update)
}

func (m FieldMode) CanSelect() bool {
	return m.has(mode_select)
}

func (m FieldMode) CanExcludeFilter() bool {
	return !m.has(mode_filter)
}

func (m FieldMode) TreatAsNotNULL() bool {
	return m.has(mode_notnull)
}

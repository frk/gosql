package config

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// String implements both the flag.Value and the json.Unmarshal interfaces
// enforcing priority of flags over json, meaning that json.Unmarshal will
// not override the value if it was previously set by flag.Var.
type String struct {
	Value string
	IsSet bool
}

// Get implements the flag.Getter interface.
func (s String) Get() interface{} {
	return s.Value
}

// String implements the flag.Value interface.
func (s String) String() string {
	return s.Value
}

// Set implements the flag.Value interface.
func (s *String) Set(value string) error {
	s.Value = value
	s.IsSet = true
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *String) UnmarshalJSON(data []byte) error {
	if !s.IsSet {
		if len(data) == 0 || string(data) == `null` {
			return nil
		}

		var value string
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		s.Value = value
		s.IsSet = true
	}
	return nil
}

// Bool implements both the flag.Value and the json.Unmarshal interfaces
// enforcing priority of flags over json, meaning that json.Unmarshal will
// not override the value if it was previously set by flag.Var.
type Bool struct {
	Value bool
	IsSet bool
}

// IsBoolFlag indicates that the Bool type can be used as a boolean flag.
func (b Bool) IsBoolFlag() bool {
	return true
}

// Get implements the flag.Getter interface.
func (b Bool) Get() interface{} {
	return b.String()
}

// String implements the flag.Value interface.
func (b Bool) String() string {
	return strconv.FormatBool(b.Value)
}

// Set implements the flag.Value interface.
func (b *Bool) Set(value string) error {
	if len(value) > 0 {
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		b.Value = v
		b.IsSet = true
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (b *Bool) UnmarshalJSON(data []byte) error {
	if !b.IsSet {
		if len(data) == 0 || string(data) == `null` {
			return nil
		}

		var value bool
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		b.Value = value
		b.IsSet = true
	}
	return nil
}

// StringSlice implements both the flag.Value and the json.Unmarshal interfaces
// enforcing priority of flags over json, meaning that json.Unmarshal will
// not override the value if it was previously set by flag.Var.
type StringSlice struct {
	Value []string
	IsSet bool
}

// Get implements the flag.Getter interface.
func (ss StringSlice) Get() interface{} {
	return ss.String()
}

// String implements the flag.Value interface.
func (ss StringSlice) String() string {
	return strings.Join(ss.Value, ",")
}

// Set implements the flag.Value interface.
func (ss *StringSlice) Set(value string) error {
	if len(value) > 0 {
		ss.Value = append(ss.Value, value)
		ss.IsSet = true
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ss *StringSlice) UnmarshalJSON(data []byte) error {
	if !ss.IsSet {
		if len(data) == 0 || string(data) == `null` {
			return nil
		}

		var value []string
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		if len(value) > 0 {
			ss.Value = value
			ss.IsSet = true
		}
	}
	return nil
}

// GoType implements both the flag.Value and the json.Unmarshal interfaces
// enforcing priority of flags over json, meaning that json.Unmarshal will
// not override the value if it was previously set by flag.Var.
type GoType struct {
	Name    string
	PkgPath string
	PkgName string
	IsPtr   bool

	IsSet bool
}

// Get implements the flag.Getter interface.
func (t GoType) Get() interface{} {
	return t.String()
}

// String implements the flag.Value interface.
func (t GoType) String() string {
	if t.IsPtr {
		return "*" + t.PkgPath + "." + t.Name
	}
	return t.PkgPath + "." + t.Name
}

// Set implements the flag.Value interface.
func (t *GoType) Set(value string) error {
	if len(value) == 0 {
		return nil
	}

	if i := strings.LastIndex(value, "."); i < 0 {
		return fmt.Errorf("bad method argument type: %q", value)
	} else {
		t.PkgPath, t.Name = value[:i], value[i+1:]
		if len(t.PkgPath) > 0 && t.PkgPath[0] == '*' {
			t.PkgPath = t.PkgPath[1:]
			t.IsPtr = true
		}

		t.IsSet = true
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *GoType) UnmarshalJSON(data []byte) error {
	if !t.IsSet {
		if len(data) == 0 || string(data) == `null` {
			return nil
		}

		var value string
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		return t.Set(value)
	}
	return nil
}

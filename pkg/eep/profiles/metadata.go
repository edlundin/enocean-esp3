package profiles

import "github.com/edlundin/enocean-esp3/pkg/eep"

type EnumValue struct {
	Raw         uint64
	Name        string
	Description string
}

type Field struct {
	Name, Shortcut     string
	BitOff, BitSize    int
	Unit               string
	ScaleMin, ScaleMax float64
	RawMin, RawMax     int
	Enums              []EnumValue
}

func (f Field) Enum(raw uint64) (EnumValue, bool) {
	for _, e := range f.Enums {
		if e.Raw == raw {
			return e, true
		}
	}
	return EnumValue{}, false
}

type Profile struct {
	EEP    eep.EEP
	Title  string
	Fields []Field
}

var Registry = map[string]Profile{}

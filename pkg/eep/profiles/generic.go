package profiles

import (
	"fmt"
	"sort"
	"strings"

	"github.com/edlundin/enocean-esp3/pkg/eep"
)

type Value struct {
	Raw    uint64
	Text   string
	Scaled float64
	Unit   string
}

type Decoded struct {
	Profile Profile
	Values  map[string]Value
	status  byte
}

func Lookup(prof eep.EEP) (Profile, bool) {
	p, ok := Registry[prof.String()]
	return p, ok
}

func Decode(prof eep.EEP, userData []byte, status byte) (Decoded, error) {
	p, ok := Lookup(prof)
	if !ok {
		return Decoded{}, fmt.Errorf("unsupported EEP %s", prof)
	}
	vals := map[string]Value{}
	for i, f := range p.Fields {
		if f.BitSize <= 0 || f.BitSize > 64 || f.BitOff+f.BitSize > len(userData)*8 {
			continue
		}
		raw := getBits(userData, f.BitOff, f.BitSize)
		v := Value{Raw: raw, Unit: f.Unit}
		if ev, ok := f.Enum(raw); ok {
			v.Text = ev.Name
		}
		if f.RawMin != f.RawMax || f.ScaleMin != f.ScaleMax {
			v.Scaled = scale(raw, f.RawMin, f.RawMax, f.ScaleMin, f.ScaleMax)
		}
		vals[fieldKey(f, i)] = v
	}
	return Decoded{Profile: p, Values: vals, status: status}, nil
}

func Encode(prof eep.EEP, values map[string]uint64) ([]byte, byte, error) {
	p, ok := Lookup(prof)
	if !ok {
		return nil, 0, fmt.Errorf("unsupported EEP %s", prof)
	}
	bits := 0
	for _, f := range p.Fields {
		if end := f.BitOff + f.BitSize; end > bits {
			bits = end
		}
	}
	data := make([]byte, (bits+7)/8)
	for i, f := range p.Fields {
		if raw, ok := values[fieldKey(f, i)]; ok {
			setBits(data, f.BitOff, f.BitSize, raw)
		} else if f.Shortcut != "" {
			if raw, ok := values[f.Shortcut]; ok {
				setBits(data, f.BitOff, f.BitSize, raw)
			}
		}
	}
	return data, 0, nil
}

func (d Decoded) EEP() eep.EEP { return d.Profile.EEP }
func (d Decoded) MarshalERP1UserData() ([]byte, byte, error) {
	vals := make(map[string]uint64, len(d.Values))
	for k, v := range d.Values {
		vals[k] = v.Raw
	}
	data, _, err := Encode(d.Profile.EEP, vals)
	return data, d.status, err
}
func (d Decoded) Format() string {
	keys := make([]string, 0, len(d.Values))
	for k := range d.Values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := []string{d.Profile.EEP.String()}
	for _, k := range keys {
		v := d.Values[k]
		if v.Text != "" {
			parts = append(parts, fmt.Sprintf("%s=%s", k, v.Text))
			continue
		}
		if v.Unit != "" {
			parts = append(parts, fmt.Sprintf("%s=%.2f%s", k, v.Scaled, v.Unit))
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=%d", k, v.Raw))
	}
	return strings.Join(parts, " ")
}

func fieldKey(f Field, i int) string {
	if f.Shortcut != "" {
		return f.Shortcut
	}
	if f.Name != "" {
		return f.Name
	}
	return fmt.Sprintf("field%d", i)
}

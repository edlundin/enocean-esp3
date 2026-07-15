package eepgen

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"go/format"
	"html"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf16"
)

type EEP struct {
	Rorgs       []Rorg `xml:"profile>rorg"`
	DirectRorgs []Rorg `xml:"rorg"`
}
type Rorg struct {
	Number string `xml:"number"`
	Title  string `xml:"title"`
	Funcs  []Func `xml:"func"`
}
type Func struct {
	Number string `xml:"number"`
	Title  string `xml:"title"`
	Types  []Type `xml:"type"`
}
type Type struct {
	Number string `xml:"number"`
	Title  string `xml:"title"`
	Cases  []Case `xml:"case"`
}
type Case struct {
	Fields []Field `xml:"datafield"`
}
type Field struct {
	Data        string  `xml:"data"`
	Shortcut    string  `xml:"shortcut"`
	Description string  `xml:"description"`
	BitOff      string  `xml:"bitoffs"`
	BitSize     string  `xml:"bitsize"`
	Unit        string  `xml:"unit"`
	Enums       []Enum  `xml:"enum"`
	Ranges      []Range `xml:"range"`
	Scales      []Scale `xml:"scale"`
}
type Enum struct {
	Items []EnumItem `xml:"item"`
}
type EnumItem struct {
	Value       string  `xml:"value"`
	Min         string  `xml:"min"`
	Max         string  `xml:"max"`
	Unit        string  `xml:"unit"`
	Description string  `xml:"description"`
	Scales      []Scale `xml:"scale"`
}
type Range struct {
	Min string `xml:"min"`
	Max string `xml:"max"`
}
type Scale struct {
	Min string `xml:"min"`
	Max string `xml:"max"`
}

type OutProfile struct {
	Key, Rorg, Func, Type, Title string
	Fields                       []OutField
}
type OutField struct {
	Name, Shortcut, Unit string
	BitOff, BitSize      int
	RawMin, RawMax       int64
	ScaleMin, ScaleMax   float64
	Enums                []OutEnum
}

type OutEnum struct {
	Raw         uint64
	Name        string
	Description string
}

func Generate(xmlPath, outDir string) error {
	profiles, err := Load(xmlPath)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, profiles); err != nil {
		return err
	}
	goSrc, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format generated source: %w", err)
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outDir, "profiles_gen.go"), goSrc, 0o644)
}

func Load(path string) ([]OutProfile, error) {
	root, err := LoadRaw(path)
	if err != nil {
		return nil, err
	}
	rorgs := root.Rorgs
	if len(rorgs) == 0 {
		rorgs = root.DirectRorgs
	}
	var out []OutProfile
	for _, r := range rorgs {
		for _, f := range r.Funcs {
			for _, t := range f.Types {
				p := OutProfile{Rorg: hex2(r.Number), Func: hex2(f.Number), Type: hex2(t.Number), Title: clean(first(t.Title, f.Title, r.Title))}
				p.Key = p.Rorg + "-" + p.Func + "-" + p.Type
				seen := map[string]bool{}
				for _, c := range t.Cases {
					for _, xf := range c.Fields {
						name := clean(xf.Data)
						if name == "" {
							continue
						}
						bo, e1 := strconv.Atoi(strings.TrimSpace(xf.BitOff))
						bs, e2 := strconv.Atoi(strings.TrimSpace(xf.BitSize))
						if e1 != nil || e2 != nil || bs <= 0 {
							continue
						}
						key := name + xf.Shortcut + xf.BitOff + xf.BitSize
						if seen[key] {
							continue
						}
						seen[key] = true
						ranges, scales, unit := xf.Ranges, xf.Scales, xf.Unit
						if item, ok := numericEnumItem(xf); ok {
							if len(ranges) == 0 {
								ranges = []Range{{Min: item.Min, Max: item.Max}}
							}
							if len(scales) == 0 {
								scales = item.Scales
							}
							if strings.TrimSpace(unit) == "" {
								unit = item.Unit
							}
						}
						of := OutField{Name: name, Shortcut: clean(xf.Shortcut), Unit: clean(unit), BitOff: bo, BitSize: bs}
						if len(ranges) > 0 {
							of.RawMin, _ = parseInt(ranges[0].Min)
							of.RawMax, _ = parseRangeMax(ranges[0].Max, xf.Description)
						}
						if len(scales) > 0 {
							of.ScaleMin, _ = parseFloat(scales[0].Min)
							of.ScaleMax, _ = parseFloat(scales[0].Max)
						}
						seenEnums := map[uint64]bool{}
						for _, en := range xf.Enums {
							for _, item := range en.Items {
								desc := clean(item.Description)
								if v, ok := parseEnumValue(item.Value); ok {
									of.Enums = append(of.Enums, OutEnum{Raw: v, Name: enumName(desc, v), Description: desc})
									seenEnums[v] = true
								}
							}
						}
						if v, desc, ok := describedEnum(xf.Description); ok && !seenEnums[v] {
							of.Enums = append(of.Enums, OutEnum{Raw: v, Name: enumName(desc, v), Description: desc})
						}
						p.Fields = append(p.Fields, of)
					}
				}
				if len(p.Fields) > 0 {
					out = append(out, p)
				}
			}
		}
	}
	return out, nil
}

func LoadRaw(path string) (EEP, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return EEP{}, err
	}
	xmlBytes := decodeUTF16(raw)
	dec := xml.NewDecoder(bytes.NewReader(xmlBytes))
	dec.Strict = false
	dec.Entity = xml.HTMLEntity
	var root EEP
	if err := dec.Decode(&root); err != nil && err != io.EOF {
		return EEP{}, err
	}
	return root, nil
}

func decodeUTF16(raw []byte) []byte {
	if len(raw) < 2 || raw[0] != 0xff || raw[1] != 0xfe {
		return raw
	}
	u := make([]uint16, 0, len(raw)/2)
	for i := 2; i+1 < len(raw); i += 2 {
		u = append(u, uint16(raw[i])|uint16(raw[i+1])<<8)
	}
	s := string(utf16.Decode(u))
	s = strings.Replace(s, `encoding="utf-16le"`, `encoding="utf-8"`, 1)
	return []byte(s)
}
func first(v ...string) string {
	for _, s := range v {
		if strings.TrimSpace(s) != "" {
			return s
		}
	}
	return ""
}
func clean(s string) string { return strings.Join(strings.Fields(html.UnescapeString(s)), " ") }
func hex2(s string) string {
	n, _ := strconv.ParseUint(strings.TrimSpace(s), 0, 8)
	return fmt.Sprintf("%02X", n)
}
func parseInt(s string) (int64, bool) {
	n, err := strconv.ParseInt(strings.TrimSpace(s), 0, 64)
	return n, err == nil
}
func parseFloat(s string) (float64, bool) {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f, err == nil
}
func parseRangeMax(s, description string) (int64, bool) {
	if max, ok := parseInt(s); ok {
		return max, true
	}
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return 0, false
	}
	max, maxOK := parseInt(parts[0])
	special, err := strconv.ParseUint(strings.TrimSpace(parts[1]), 0, 64)
	described, _, describedOK := describedEnum(description)
	if maxOK && err == nil && describedOK && special == described {
		return max, true
	}
	return 0, false
}
func parseEnumValue(s string) (uint64, bool) {
	s = strings.TrimSpace(s)
	if i := strings.Index(s, " ("); i >= 0 && strings.HasSuffix(s, ")") {
		s = s[:i]
	}
	v, err := strconv.ParseUint(s, 0, 64)
	return v, err == nil
}
func numericEnumItem(f Field) (EnumItem, bool) {
	var found EnumItem
	ok := false
	for _, enum := range f.Enums {
		for _, item := range enum.Items {
			if _, minOK := parseInt(item.Min); !minOK {
				continue
			}
			if _, maxOK := parseInt(item.Max); !maxOK || len(item.Scales) == 0 && strings.TrimSpace(item.Unit) == "" {
				continue
			}
			if ok {
				return EnumItem{}, false
			}
			found, ok = item, true
		}
	}
	return found, ok
}
func describedEnum(s string) (uint64, string, bool) {
	const marker = "value "
	i := strings.Index(strings.ToLower(s), marker)
	if i < 0 {
		return 0, "", false
	}
	rest := s[i+len(marker):]
	eq := strings.Index(rest, "=")
	if eq < 0 {
		return 0, "", false
	}
	v, err := strconv.ParseUint(strings.TrimSpace(rest[:eq]), 0, 64)
	if err != nil {
		return 0, "", false
	}
	desc := strings.TrimSpace(rest[eq+1:])
	if i := strings.IndexAny(desc, ",;."); i >= 0 {
		desc = desc[:i]
	}
	desc = clean(desc)
	return v, desc, desc != ""
}
func enumName(desc string, raw uint64) string {
	name := desc
	for _, cut := range []string{":", " or ", " (", " - ", ","} {
		if i := strings.Index(name, cut); i >= 0 {
			name = name[:i]
		}
	}
	var b strings.Builder
	upperNext := true
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if upperNext {
				r = unicode.ToUpper(r)
			}
			b.WriteRune(r)
			upperNext = false
		} else {
			upperNext = true
		}
	}
	out := b.String()
	if out == "" || unicode.IsDigit([]rune(out)[0]) || len(out) > 32 {
		return fmt.Sprintf("Value%d", raw)
	}
	return out
}

var tmpl = template.Must(template.New("profiles").Parse(`// Code generated by eepgen; DO NOT EDIT.
package profiles

import (
	"github.com/edlundin/enocean-esp3/pkg/eep"
	"github.com/edlundin/enocean-esp3/pkg/enums"
)

func init() {
{{- range . }}
	Registry["{{ .Key }}"] = Profile{EEP: eep.EEP{Rorg: enums.Rorg(0x{{ .Rorg }}), Func: 0x{{ .Func }}, Type: 0x{{ .Type }}}, Title: {{ printf "%q" .Title }}, Fields: []Field{
		{{- range .Fields }}
		{Name: {{ printf "%q" .Name }}, Shortcut: {{ printf "%q" .Shortcut }}, BitOff: {{ .BitOff }}, BitSize: {{ .BitSize }}, Unit: {{ printf "%q" .Unit }}, ScaleMin: {{ printf "%g" .ScaleMin }}, ScaleMax: {{ printf "%g" .ScaleMax }}, RawMin: {{ .RawMin }}, RawMax: {{ .RawMax }}{{ if .Enums }}, Enums: []EnumValue{ {{- range .Enums }}{Raw: {{ .Raw }}, Name: {{ printf "%q" .Name }}, Description: {{ printf "%q" .Description }}}, {{- end }} }{{ end }}},
		{{- end }}
	}}
{{- end }}
}
`))

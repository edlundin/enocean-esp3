package eepgen

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGenerateTinyXML verifies GenerateTinyXML behavior.
func TestGenerateTinyXML(t *testing.T) {
	dir := t.TempDir()
	xml := filepath.Join(dir, "eep.xml")
	if err := os.WriteFile(xml, []byte(`<eep><rorg><number>0xF6</number><title>RPS</title><func><number>0x01</number><title>Switches</title><type><number>0x01</number><title>Rocker</title><case><datafield><data>Push button</data><shortcut>PB</shortcut><bitoffs>3</bitoffs><bitsize>1</bitsize><range><min>0</min><max>1</max></range><scale><min>0</min><max>1</max></scale><unit>bool</unit><enum><item><value>0</value><description>released</description></item><item><value>1</value><description>pressed</description></item></enum></datafield></case></type></func></rorg></eep>`), 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "out")
	if err := Generate(xml, out); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(out, "profiles_gen.go")); err != nil {
		t.Fatal(err)
	}
	profiles, err := Load(xml)
	if err != nil {
		t.Fatal(err)
	}
	if got := profiles[0].Fields[0]; got.Unit != "bool" || got.Enums[1].Description != "pressed" || got.Enums[1].Name != "Pressed" {
		t.Fatalf("field = %#v", got)
	}
}

// TestParseCompositeNumbers verifies ParseCompositeNumbers behavior.
func TestParseCompositeNumbers(t *testing.T) {
	if got, ok := parseRangeMax("100, 255", "Control variable; value 255 = auto"); !ok || got != 100 {
		t.Fatalf("parseRangeMax() = %d, %t", got, ok)
	}
	if got, ok := parseRangeMax("100, 254", "Control variable; value 255 = auto"); ok || got != 0 {
		t.Fatalf("mismatched parseRangeMax() = %d, %t", got, ok)
	}
	for _, value := range []string{"100 or +40", "4095 (409,5)", "20 .. 80"} {
		if _, ok := parseFloat(value); ok {
			t.Fatalf("ambiguous value %q parsed as one number", value)
		}
	}
}

// TestEnumNameGeneralizesLongDescriptions verifies EnumNameGeneralizesLongDescriptions behavior.
func TestEnumNameGeneralizesLongDescriptions(t *testing.T) {
	if got := enumName(`Button AI: "Switch light on" or "Dim light down" or "Move blind closed"`, 0); got != "ButtonAI" {
		t.Fatalf("enumName = %q", got)
	}
	if got := enumName("Pressed & Hold", 1); got != "PressedHold" {
		t.Fatalf("enumName = %q", got)
	}
}

func TestDescribedEnum(t *testing.T) {
	if value, description, ok := describedEnum("Control; value 255 = auto, use default"); !ok || value != 255 || description != "auto" {
		t.Fatalf("describedEnum() = %d, %q, %v", value, description, ok)
	}
	if _, _, ok := describedEnum("Control; value nope = auto"); ok {
		t.Fatal("invalid enum value accepted")
	}
}

// TestLoadRealEEP268 verifies LoadRealEEP268 behavior.
func TestLoadRealEEP268(t *testing.T) {
	path := filepath.Join("..", "..", "eep268.xml")
	if _, err := os.Stat(path); err != nil {
		t.Skip("eep268.xml not present")
	}
	root, err := LoadRaw(path)
	if err != nil {
		t.Fatal(err)
	}
	rorgs := root.Rorgs
	if len(rorgs) != 4 {
		t.Fatalf("profile RORGs = %d", len(rorgs))
	}
	funcs, types, cases, fields := 0, 0, 0, 0
	for _, r := range rorgs {
		funcs += len(r.Funcs)
		for _, f := range r.Funcs {
			types += len(f.Types)
			for _, typ := range f.Types {
				cases += len(typ.Cases)
				for _, c := range typ.Cases {
					fields += len(c.Fields)
				}
			}
		}
	}
	if funcs != 45 || types != 270 || cases != 355 || fields != 2435 {
		t.Fatalf("counts funcs/types/cases/fields = %d/%d/%d/%d", funcs, types, cases, fields)
	}

	profiles, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	byKey := map[string]OutProfile{}
	for _, p := range profiles {
		byKey[p.Key] = p
	}
	assertField(t, byKey["D5-00-01"], "CO", 7, 1)
	assertField(t, byKey["D5-00-01"], "LRN", 4, 1)
	assertField(t, byKey["F6-01-01"], "PB", 3, 1)
	assertField(t, byKey["A5-02-01"], "TMP", 16, 8)
	if got := findField(byKey["A5-02-01"], "TMP"); got.RawMin != 255 || got.RawMax != 0 || got.ScaleMin != -40 || got.ScaleMax != 0 || got.Unit != "°C" {
		t.Fatalf("A5-02-01 TMP range/scale = %#v", got)
	}
	if got := findField(byKey["F6-01-01"], "PB").Enums[1]; got.Description != "Pressed & Hold" || got.Name != "PressedHold" {
		t.Fatalf("F6-01-01 PB enum 1 = %#v", got)
	}
	for _, shortcut := range []string{"POS", "ANG"} {
		if got := findField(byKey["D2-05-00"], shortcut); got.RawMin != 0 || got.RawMax != 100 || got.ScaleMin != 0 || got.ScaleMax != 100 || got.Unit != "%" || len(got.Enums) != 1 || got.Enums[0].Raw != 127 {
			t.Fatalf("D2-05-00 %s range/scale/enum = %#v", shortcut, got)
		}
	}
	if got := findField(byKey["A5-20-10"], "CVAR"); got.RawMin != 0 || got.RawMax != 100 || got.ScaleMin != 0 || got.ScaleMax != 100 || got.Unit != "%" || len(got.Enums) != 1 || got.Enums[0].Raw != 255 || got.Enums[0].Name != "Auto" {
		t.Fatalf("A5-20-10 CVAR range/scale/enum = %#v", got)
	}
	if got := findField(byKey["D2-05-00"], "VERT"); len(got.Enums) != 1 || got.Enums[0].Raw != 32767 || got.Enums[0].Name != "NoChange" {
		t.Fatalf("D2-05-00 VERT enum = %#v", got)
	}
}

// assertField checks a generated field's bit position and size.
func assertField(t *testing.T, p OutProfile, shortcut string, off, size int) {
	t.Helper()
	f := findField(p, shortcut)
	if f.Shortcut == "" {
		t.Fatalf("field %s missing in %s", shortcut, p.Key)
	}
	if f.BitOff != off || f.BitSize != size {
		t.Fatalf("%s %s offset/size = %d/%d", p.Key, shortcut, f.BitOff, f.BitSize)
	}
}

// findField finds a generated field by shortcut.
func findField(p OutProfile, shortcut string) OutField {
	for _, f := range p.Fields {
		if f.Shortcut == shortcut {
			return f
		}
	}
	return OutField{}
}

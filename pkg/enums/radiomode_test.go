package enums

import "testing"

func TestRadioMode(t *testing.T) {
	cases := []struct{ b byte; v RadioMode; s string }{
		{0x00, RadioModeERP1, "ERP1"},
		{0x01, RadioModeERP2, "ERP2"},
	}
	for _, c := range cases {
		v, err := ParseRadioModeFromByte(c.b)
		if err != nil || v != c.v || v.String() != c.s || !v.Valid() { t.Fatalf("%#x => %v %v", c.b, v, err) }
	}
	if _, err := ParseRadioModeFromByte(0xff); err == nil { t.Fatal("expected error") }
	if RadioMode(0xff).String() != "UNKNOWN" || RadioMode(0xff).Valid() { t.Fatal("invalid radio mode accepted") }
}

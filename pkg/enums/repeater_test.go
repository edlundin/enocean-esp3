package enums

import "testing"

// TestRepeaterMode verifies RepeaterMode behavior.
func TestRepeaterMode(t *testing.T) {
	cases := []struct{ b byte; v RepeaterMode; s string }{
		{0x00, RepeaterModeOFF, "OFF"}, {0x01, RepeaterModeON, "ON"}, {0x02, RepeaterModeSELECTIVE, "SELECTIVE"},
	}
	for _, c := range cases {
		v, err := ParseRepeaterModeFromByte(c.b)
		if err != nil || v != c.v || v.String() != c.s || !v.Valid() { t.Fatalf("%#x => %v %v", c.b, v, err) }
	}
	if _, err := ParseRepeaterModeFromByte(0xff); err == nil { t.Fatal("expected error") }
	if RepeaterMode(0xff).String() != "UNKNOWN" || RepeaterMode(0xff).Valid() { t.Fatal("invalid repeater mode accepted") }
}

// TestRepeaterLevel verifies RepeaterLevel behavior.
func TestRepeaterLevel(t *testing.T) {
	cases := []struct{ b byte; v RepeaterLevel; s string }{
		{0x00, RepeaterLevelNO_REPETITION, "NO_REPEATING"}, {0x01, RepeaterLevel1_REPETITION, "1_REPEAT"}, {0x02, RepeaterLevel2_REPETITION, "2_REPEAT"},
	}
	for _, c := range cases {
		v, err := ParseRepeaterLevelFromByte(c.b)
		if err != nil || v != c.v || v.String() != c.s || !v.Valid() { t.Fatalf("%#x => %v %v", c.b, v, err) }
	}
	if _, err := ParseRepeaterLevelFromByte(0xff); err == nil { t.Fatal("expected error") }
	if RepeaterLevel(0xff).String() != "UNKNOWN" || RepeaterLevel(0xff).Valid() { t.Fatal("invalid repeater level accepted") }
}

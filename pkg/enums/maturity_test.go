package enums

import "testing"

func TestMaturity(t *testing.T) {
	cases := []struct{ b byte; v Maturity; s string }{
		{0x00, MaturityFORWARDED_IMMEDIATELY, "FORWARDED_IMMEDIATELY"},
		{0x01, MaturityFORWARDED_ON_TIMEOUT, "FORWARDED_ON_TIMEOUT"},
		{0x02, MaturityFORWARD_SUBTELEGRAMS, "FORWARD_SUBTELEGRAMS"},
	}
	for _, c := range cases {
		v, err := ParseMaturityFromByte(c.b)
		if err != nil || v != c.v || v.String() != c.s || !v.Valid() { t.Fatalf("%#x => %v %v", c.b, v, err) }
	}
	if _, err := ParseMaturityFromByte(0xff); err == nil { t.Fatal("expected error") }
	if Maturity(0xff).String() != "UNKNOWN" || Maturity(0xff).Valid() { t.Fatal("invalid maturity accepted") }
}

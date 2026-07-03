package profiles

import "testing"

func TestBitHelpers(t *testing.T) {
	b := []byte{0, 0}
	setBits(b, 3, 5, 0b10101)
	if got := getBits(b, 3, 5); got != 0b10101 { t.Fatalf("bits = %b", got) }
	if got := scale(128, 255, 0, -40, 0); got < -21 || got > -19 { t.Fatalf("scale = %f", got) }
	if got := unscale(-20, 255, 0, -40, 0); got < 127 || got > 128 { t.Fatalf("unscale = %d", got) }
}

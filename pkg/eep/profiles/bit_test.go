package profiles

import "testing"

// TestBitHelpers verifies BitHelpers behavior.
func TestBitHelpers(t *testing.T) {
	b := []byte{0, 0}
	setBits(b, 3, 5, 0b10101)
	if got := getBits(b, 3, 5); got != 0b10101 {
		t.Fatalf("bits = %b", got)
	}
}

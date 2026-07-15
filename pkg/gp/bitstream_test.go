package gp

import "testing"

// TestBitstreamEdges verifies BitstreamEdges behavior.
func TestBitstreamEdges(t *testing.T) {
	buf := []byte{0, 0}
	if err := writeUnsigned(buf, 3, 5, 0b10101); err != nil { t.Fatal(err) }
	if got, err := readUnsigned(buf, 3, 5); err != nil || got != 0b10101 { t.Fatalf("read unsigned = %b, %v", got, err) }
	if err := writeSigned(buf, 0, 4, -1); err != nil { t.Fatal(err) }
	if got, err := readSigned(buf, 0, 4); err != nil || got != -1 { t.Fatalf("read signed = %d, %v", got, err) }
	if got, err := readSigned(buf, 0, 0); err != nil || got != 0 { t.Fatalf("zero signed = %d, %v", got, err) }
	if got := bytesForBits(0); got != nil { t.Fatalf("zero bytes = %v", got) }
	if got := len(bytesForBits(9)); got != 2 { t.Fatalf("bytes = %d", got) }
	if err := writeUnsigned(buf, -1, 1, 0); err != ErrBitstreamOutOfRange { t.Fatalf("range err = %v", err) }
}

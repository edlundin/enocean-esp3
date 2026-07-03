package eep

import (
	"errors"
	"math"
	"testing"
)

func TestReadWriteBits(t *testing.T) {
	if got, err := ReadBits([]byte{0x80}, 7, 1); err != nil || got != 1 {
		t.Fatalf("ReadBits MSB = %d, %v", got, err)
	}
	payload := []byte{0x00, 0x00, 0x80, 0x10}
	if got, err := ReadBits(payload, 16, 8); err != nil || got != 0x80 {
		t.Fatalf("ReadBits temp = %d, %v", got, err)
	}
	if got, err := ReadBits(payload, 28, 1); err != nil || got != 1 {
		t.Fatalf("ReadBits LRN = %d, %v", got, err)
	}

	b := []byte{0, 0}
	if err := WriteBits(b, 5, 6, 0x2a); err != nil {
		t.Fatal(err)
	}
	if got, err := ReadBits(b, 5, 6); err != nil || got != 0x2a {
		t.Fatalf("cross-byte round trip = %d, %v", got, err)
	}
}

func TestReadWriteBitsBounds(t *testing.T) {
	if _, err := ReadBits([]byte{0}, 4, 5); !errors.Is(err, ErrBitfieldOutOfRange) {
		t.Fatalf("ReadBits error = %v", err)
	}
	if err := WriteBits([]byte{0}, 0, 65, 0); !errors.Is(err, ErrBitfieldOutOfRange) {
		t.Fatalf("WriteBits error = %v", err)
	}
}

func TestScaleRaw(t *testing.T) {
	if got := ScaleRaw(255, 255, 0, -40, 0); got != -40 {
		t.Fatalf("raw 255 = %v", got)
	}
	if got := ScaleRaw(0, 255, 0, -40, 0); got != 0 {
		t.Fatalf("raw 0 = %v", got)
	}
	if got := ScaleRaw(128, 255, 0, -40, 0); math.Abs(got-(-20.0784314)) > 0.001 {
		t.Fatalf("raw 128 = %v", got)
	}
	if got := UnscaleRaw(-20.0784314, 255, 0, -40, 0); got != 128 {
		t.Fatalf("unscale = %v", got)
	}
}

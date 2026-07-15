package profiles

import eepcore "github.com/edlundin/enocean-esp3/pkg/eep"

// getBits returns Bits.
func getBits(b []byte, off, size int) uint64 {
	v, _ := eepcore.ReadBits(b, off, size)
	return v
}

// setBits updates Bits.
func setBits(b []byte, off, size int, v uint64) {
	_ = eepcore.WriteBits(b, off, size, v)
}

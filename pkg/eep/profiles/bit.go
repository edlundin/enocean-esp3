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

// scale maps a raw field value into its engineering range.
func scale(raw uint64, rmin, rmax int, smin, smax float64) float64 {
	return eepcore.ScaleRaw(raw, rmin, rmax, smin, smax)
}

// unscale maps an engineering value back to its raw field value.
func unscale(v float64, rmin, rmax int, smin, smax float64) uint64 {
	return eepcore.UnscaleRaw(v, rmin, rmax, smin, smax)
}

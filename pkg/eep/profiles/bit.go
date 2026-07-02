package profiles

import eepcore "github.com/edlundin/enocean-esp3/pkg/eep"

func getBits(b []byte, off, size int) uint64 {
	v, _ := eepcore.ReadBits(b, off, size)
	return v
}

func setBits(b []byte, off, size int, v uint64) {
	_ = eepcore.WriteBits(b, off, size, v)
}

func scale(raw uint64, rmin, rmax int, smin, smax float64) float64 {
	return eepcore.ScaleRaw(raw, rmin, rmax, smin, smax)
}

func unscale(v float64, rmin, rmax int, smin, smax float64) uint64 {
	return eepcore.UnscaleRaw(v, rmin, rmax, smin, smax)
}

package eep

import (
	"errors"
	"math"
)

var ErrBitfieldOutOfRange = errors.New("bitfield out of range")

// ReadBits reads an eep268.xml bit field. Offsets are zero-based, with bit 0
// at the least significant bit of data[0], matching the checked-in XML facts.
func ReadBits(data []byte, bitOffset, bitSize int) (uint64, error) {
	if bitOffset < 0 || bitSize < 0 || bitSize > 64 || bitOffset+bitSize > len(data)*8 {
		return 0, ErrBitfieldOutOfRange
	}
	var v uint64
	for i := 0; i < bitSize; i++ {
		if data[(bitOffset+i)/8]&(1<<uint((bitOffset+i)%8)) != 0 {
			v |= 1 << uint(i)
		}
	}
	return v, nil
}

// WriteBits writes an eep268.xml bit field. Bits outside the field are kept.
func WriteBits(data []byte, bitOffset, bitSize int, value uint64) error {
	if bitOffset < 0 || bitSize < 0 || bitSize > 64 || bitOffset+bitSize > len(data)*8 {
		return ErrBitfieldOutOfRange
	}
	for i := 0; i < bitSize; i++ {
		mask := byte(1 << uint((bitOffset+i)%8))
		if value&(1<<uint(i)) != 0 {
			data[(bitOffset+i)/8] |= mask
		} else {
			data[(bitOffset+i)/8] &^= mask
		}
	}
	return nil
}

func ScaleRaw(raw uint64, rawMin, rawMax int, scaleMin, scaleMax float64) float64 {
	return scaleMin + (float64(raw)-float64(rawMin))*(scaleMax-scaleMin)/float64(rawMax-rawMin)
}

func UnscaleRaw(value float64, rawMin, rawMax int, scaleMin, scaleMax float64) uint64 {
	return uint64(math.Round(float64(rawMin) + (value-scaleMin)*float64(rawMax-rawMin)/(scaleMax-scaleMin)))
}

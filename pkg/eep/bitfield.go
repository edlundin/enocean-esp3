package eep

import (
	"errors"
	"math"
)

var ErrBitfieldOutOfRange = errors.New("bitfield out of range")

// ReadBits reads an eep268.xml bit field. Offset 0 is the first transmitted
// bit: the most significant bit of data[0].
func ReadBits(data []byte, bitOffset, bitSize int) (uint64, error) {
	if bitOffset < 0 || bitSize < 0 || bitSize > 64 || bitOffset+bitSize > len(data)*8 {
		return 0, ErrBitfieldOutOfRange
	}
	var v uint64
	for i := 0; i < bitSize; i++ {
		v <<= 1
		if data[(bitOffset+i)/8]&(1<<uint(7-(bitOffset+i)%8)) != 0 {
			v |= 1
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
		mask := byte(1 << uint(7-(bitOffset+i)%8))
		if value&(1<<uint(bitSize-1-i)) != 0 {
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

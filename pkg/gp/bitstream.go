package gp

import "errors"

var ErrBitstreamOutOfRange = errors.New("gp bitstream out of range")

// readUnsigned reads Unsigned.
func readUnsigned(data []byte, bitOffset, bitSize int) (uint64, error) {
	if bitOffset < 0 || bitSize < 0 || bitSize > 64 || bitOffset+bitSize > len(data)*8 {
		return 0, ErrBitstreamOutOfRange
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

// readSigned reads Signed.
func readSigned(data []byte, bitOffset, bitSize int) (int64, error) {
	v, err := readUnsigned(data, bitOffset, bitSize)
	if err != nil || bitSize == 0 {
		return int64(v), err
	}
	if bitSize == 64 {
		return int64(v), nil
	}
	if v&(1<<uint(bitSize-1)) == 0 {
		return int64(v), nil
	}
	return int64(v) - int64(uint64(1)<<uint(bitSize)), nil
}

// writeUnsigned writes Unsigned.
func writeUnsigned(data []byte, bitOffset, bitSize int, value uint64) error {
	if bitOffset < 0 || bitSize < 0 || bitSize > 64 || bitOffset+bitSize > len(data)*8 {
		return ErrBitstreamOutOfRange
	}
	for i := 0; i < bitSize; i++ {
		mask := byte(1 << uint(7-((bitOffset+i)%8)))
		if value&(1<<uint(bitSize-1-i)) != 0 {
			data[(bitOffset+i)/8] |= mask
		} else {
			data[(bitOffset+i)/8] &^= mask
		}
	}
	return nil
}

// writeSigned writes Signed.
func writeSigned(data []byte, bitOffset, bitSize int, value int64) error {
	return writeUnsigned(data, bitOffset, bitSize, uint64(value))
}

// bytesForBits returns the bytes required for a bit count.
func bytesForBits(bits int) []byte {
	if bits <= 0 {
		return nil
	}
	return make([]byte, (bits+7)/8)
}

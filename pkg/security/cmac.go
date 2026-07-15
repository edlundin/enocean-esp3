package security

import "crypto/aes"

// cmac computes an AES-CMAC.
func cmac(key [16]byte, msg []byte) ([]byte, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	l := make([]byte, aes.BlockSize)
	block.Encrypt(l, l)
	k1 := dbl(l)
	k2 := dbl(k1)
	n := (len(msg) + aes.BlockSize - 1) / aes.BlockSize
	if n == 0 {
		n = 1
	}
	lastComplete := len(msg) != 0 && len(msg)%aes.BlockSize == 0
	last := make([]byte, aes.BlockSize)
	if lastComplete {
		copy(last, msg[(n-1)*aes.BlockSize:])
		xor(last, k1)
	} else {
		copy(last, msg[(n-1)*aes.BlockSize:])
		last[len(msg)%aes.BlockSize] = 0x80
		xor(last, k2)
	}
	x := make([]byte, aes.BlockSize)
	for i := 0; i < n-1; i++ {
		xor(x, msg[i*aes.BlockSize:(i+1)*aes.BlockSize])
		block.Encrypt(x, x)
	}
	xor(x, last)
	block.Encrypt(x, x)
	return x, nil
}

// dbl doubles a CMAC block in the finite field.
func dbl(in []byte) []byte {
	out := make([]byte, aes.BlockSize)
	carry := byte(0)
	for i := aes.BlockSize - 1; i >= 0; i-- {
		out[i] = in[i]<<1 | carry
		carry = in[i] >> 7
	}
	if carry != 0 {
		out[aes.BlockSize-1] ^= 0x87
	}
	return out
}

// xor returns the byte-wise XOR of two blocks.
func xor(dst, src []byte) {
	for i := range dst {
		dst[i] ^= src[i]
	}
}

package security

import "crypto/aes"

func cmac(key [16]byte, msg []byte) ([]byte, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	zero := make([]byte, 16)
	l := make([]byte, 16)
	block.Encrypt(l, zero)
	k1 := dbl(l)
	k2 := dbl(k1)
	n := (len(msg) + 15) / 16
	if n == 0 {
		n = 1
	}
	lastComplete := len(msg) != 0 && len(msg)%16 == 0
	last := make([]byte, 16)
	if lastComplete {
		copy(last, msg[(n-1)*16:])
		xor(last, k1)
	} else {
		copy(last, msg[(n-1)*16:])
		last[len(msg)%16] = 0x80
		xor(last, k2)
	}
	x := make([]byte, 16)
	for i := 0; i < n-1; i++ {
		xor(x, msg[i*16:(i+1)*16])
		block.Encrypt(x, x)
	}
	xor(x, last)
	block.Encrypt(x, x)
	return x, nil
}

func dbl(in []byte) []byte {
	out := make([]byte, 16)
	carry := byte(0)
	for i := 15; i >= 0; i-- {
		out[i] = in[i]<<1 | carry
		carry = in[i] >> 7
	}
	if carry != 0 {
		out[15] ^= 0x87
	}
	return out
}

func xor(dst, src []byte) {
	for i := range dst {
		dst[i] ^= src[i]
	}
}

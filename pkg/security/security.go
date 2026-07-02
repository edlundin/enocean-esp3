package security

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"

	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/erp1"
)

const (
	RLCExplicit32CMAC32VAES byte = 0xf3
)

var vaesIV = []byte{0x34, 0x10, 0xde, 0x8f, 0x1a, 0xba, 0x3e, 0xff, 0x9f, 0x5a, 0x11, 0x71, 0x72, 0xea, 0xca, 0xbd}

type SLF byte

func (s SLF) RLCLength() int {
	switch (byte(s) >> 5) & 7 {
	case 4, 5, 6, 7:
		return 4
	case 2, 3:
		return 2
	default:
		return 0
	}
}

func (s SLF) CMACLength() int {
	switch (byte(s) >> 3) & 3 {
	case 1:
		return 3
	case 2, 3:
		return 4
	default:
		return 0
	}
}

func (s SLF) Encrypted() bool { return byte(s)&7 != 0 }

type Secure struct {
	Rorg enums.Rorg
	Data []byte
	RLC  []byte
	CMAC []byte
}

func EncodeSEC_R(key [16]byte, slf SLF, rlc []byte, rorg enums.Rorg, data []byte) (erp1.Packet, error) {
	if len(rlc) != slf.RLCLength() {
		return erp1.Packet{}, fmt.Errorf("RLC length %d, want %d", len(rlc), slf.RLCLength())
	}
	plain := append([]byte{byte(rorg)}, data...)
	enc, err := vaes(key, plain)
	if err != nil {
		return erp1.Packet{}, err
	}
	mac, err := cmac(key, append(append([]byte{byte(enums.RorgSEC_R)}, enc...), rlc...))
	if err != nil {
		return erp1.Packet{}, err
	}
	payload := append(enc, rlc...)
	payload = append(payload, mac[:slf.CMACLength()]...)
	return erp1.Packet{Rorg: enums.RorgSEC_R, UserData: payload, SecurityLevel: byte(slf)}, nil
}

func DecodeSEC_R(key [16]byte, slf SLF, p erp1.Packet) (Secure, error) {
	if p.Rorg != enums.RorgSEC_R {
		return Secure{}, errors.New("not SEC_R")
	}
	rl, ml := slf.RLCLength(), slf.CMACLength()
	if len(p.UserData) < 1+rl+ml {
		return Secure{}, errors.New("SEC_R payload too short")
	}
	encEnd := len(p.UserData) - rl - ml
	enc, rlc, got := p.UserData[:encEnd], p.UserData[encEnd:encEnd+rl], p.UserData[encEnd+rl:]
	mac, err := cmac(key, append(append([]byte{byte(enums.RorgSEC_R)}, enc...), rlc...))
	if err != nil {
		return Secure{}, err
	}
	if !equal(got, mac[:ml]) {
		return Secure{}, errors.New("invalid CMAC")
	}
	plain, err := vaes(key, enc)
	if err != nil {
		return Secure{}, err
	}
	return Secure{Rorg: enums.Rorg(plain[0]), Data: plain[1:], RLC: append([]byte(nil), rlc...), CMAC: append([]byte(nil), got...)}, nil
}

func vaes(key [16]byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	out := append([]byte(nil), data...)
	stream := cipher.NewOFB(block, vaesIV)
	stream.XORKeyStream(out, out)
	return out, nil
}

func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var v byte
	for i := range a {
		v |= a[i] ^ b[i]
	}
	return v == 0
}

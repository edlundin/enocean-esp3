package security

import (
	"crypto/aes"
	"errors"
	"fmt"

	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/erp1"
)

const (
	RLCImplicit24CMAC24VAES     byte = 0x8b
	RLCExplicit24CMAC24VAES     byte = 0xab
	RLCExplicit24Of32CMAC24VAES byte = 0xcb
	RLCExplicit32CMAC32VAES     byte = 0xf3
)

var (
	ErrNeedsChaining    = errors.New("SEC_R payload requires SEC_CDM")
	ErrRLCStateRequired = errors.New("full RLC state required")
	vaesIV              = [aes.BlockSize]byte{0x34, 0x10, 0xde, 0x8f, 0x1a, 0xba, 0x3e, 0xff, 0x9f, 0x5a, 0x11, 0x71, 0x72, 0xea, 0xca, 0xbd}
)

type SLF byte

// Validate validates the security level format.
func (s SLF) Validate() error {
	rlcType, cmacType, encryptionType := byte(s)>>5, (byte(s)>>3)&3, byte(s)&7
	if rlcType < 4 || cmacType < 1 || cmacType > 2 || encryptionType != 3 {
		return fmt.Errorf("unsupported SLF 0x%02x", byte(s))
	}
	return nil
}

// RLCLength is the full rolling-code size used for encryption and CMAC.
func (s SLF) RLCLength() int {
	switch byte(s) >> 5 {
	case 4, 5:
		return 3
	case 6, 7:
		return 4
	default:
		return 0
	}
}

// TransmittedRLCLength returns the transmitted rolling-code length.
func (s SLF) TransmittedRLCLength() int {
	switch byte(s) >> 5 {
	case 4:
		return 0
	case 5, 6:
		return 3
	case 7:
		return 4
	default:
		return 0
	}
}

// CMACLength returns the configured CMAC length.
func (s SLF) CMACLength() int {
	switch (byte(s) >> 3) & 3 {
	case 1:
		return 3
	case 2:
		return 4
	default:
		return 0
	}
}

// Encrypted reports whether encryption is enabled.
func (s SLF) Encrypted() bool { return byte(s)&7 == 3 }

type Secure struct {
	Rorg enums.Rorg
	Data []byte
	RLC  []byte
	CMAC []byte
}

// EncodeSEC_R encodes SEC_R.
func EncodeSEC_R(key [16]byte, slf SLF, rlc []byte, rorg enums.Rorg, data []byte) (erp1.Packet, error) {
	payload, err := encodeSecurePayload(key, slf, rlc, rorg, data)
	if err != nil {
		return erp1.Packet{}, err
	}
	if len(payload) > 15 {
		return erp1.Packet{}, ErrNeedsChaining
	}
	return erp1.Packet{Rorg: enums.RorgSEC_R, UserData: payload, SecurityLevel: byte(slf)}, nil
}

// EncodeSEC_CDM secures data and splits the resulting payload into SEC_CDM packets.
func EncodeSEC_CDM(key [16]byte, slf SLF, rlc []byte, seq byte, rorg enums.Rorg, data []byte) ([]erp1.Packet, error) {
	payload, err := encodeSecurePayload(key, slf, rlc, rorg, data)
	if err != nil {
		return nil, err
	}
	packets, err := SplitSEC_CDM(seq, payload)
	if err != nil {
		return nil, err
	}
	for i := range packets {
		packets[i].SecurityLevel = byte(slf)
	}
	return packets, nil
}

// encodeSecurePayload encodes SecurePayload.
func encodeSecurePayload(key [16]byte, slf SLF, rlc []byte, rorg enums.Rorg, data []byte) ([]byte, error) {
	if err := slf.Validate(); err != nil {
		return nil, err
	}
	if len(rlc) != slf.RLCLength() {
		return nil, fmt.Errorf("RLC length %d, want %d", len(rlc), slf.RLCLength())
	}
	plain := append([]byte{byte(rorg)}, data...)
	enc, err := vaes(key, rlc, plain)
	if err != nil {
		return nil, err
	}
	mac, err := cmac(key, append(append([]byte{byte(enums.RorgSEC_R)}, enc...), rlc...))
	if err != nil {
		return nil, err
	}
	txLen := slf.TransmittedRLCLength()
	payload := append([]byte(nil), enc...)
	payload = append(payload, rlc[len(rlc)-txLen:]...)
	return append(payload, mac[:slf.CMACLength()]...), nil
}

// DecodeSEC_R decodes SEC_R.
func DecodeSEC_R(key [16]byte, slf SLF, p erp1.Packet) (Secure, error) {
	return DecodeSEC_RWithRLC(key, slf, nil, p)
}

// DecodeSEC_RWithRLC accepts the receiver's full rolling-code state for SLFs
// that transmit only part (or none) of the RLC.
func DecodeSEC_RWithRLC(key [16]byte, slf SLF, rlc []byte, p erp1.Packet) (Secure, error) {
	if p.Rorg != enums.RorgSEC_R {
		return Secure{}, errors.New("not SEC_R")
	}
	if err := slf.Validate(); err != nil {
		return Secure{}, err
	}
	txLen, macLen := slf.TransmittedRLCLength(), slf.CMACLength()
	encEnd := len(p.UserData) - txLen - macLen
	if encEnd < 1 {
		return Secure{}, errors.New("SEC_R payload too short")
	}
	enc := p.UserData[:encEnd]
	transmittedRLC := p.UserData[encEnd : encEnd+txLen]
	got := p.UserData[encEnd+txLen:]
	if rlc == nil {
		if txLen != slf.RLCLength() {
			return Secure{}, ErrRLCStateRequired
		}
		rlc = append([]byte(nil), transmittedRLC...)
	} else {
		if len(rlc) != slf.RLCLength() {
			return Secure{}, fmt.Errorf("RLC length %d, want %d", len(rlc), slf.RLCLength())
		}
		if !equal(transmittedRLC, rlc[len(rlc)-txLen:]) {
			return Secure{}, errors.New("transmitted RLC does not match receiver state")
		}
	}
	mac, err := cmac(key, append(append([]byte{byte(enums.RorgSEC_R)}, enc...), rlc...))
	if err != nil {
		return Secure{}, err
	}
	if !equal(got, mac[:macLen]) {
		return Secure{}, errors.New("invalid CMAC")
	}
	plain, err := vaes(key, rlc, enc)
	if err != nil {
		return Secure{}, err
	}
	return Secure{Rorg: enums.Rorg(plain[0]), Data: plain[1:], RLC: append([]byte(nil), rlc...), CMAC: append([]byte(nil), got...)}, nil
}

// vaes applies the VAES transformation.
func vaes(key [16]byte, rlc, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	base := vaesIV
	for i, b := range rlc {
		base[i] ^= b
	}
	out := make([]byte, len(data))
	var previous, input, stream [aes.BlockSize]byte
	for offset := 0; offset < len(data); offset += aes.BlockSize {
		for i := range input {
			input[i] = base[i] ^ previous[i]
		}
		block.Encrypt(stream[:], input[:])
		n := min(aes.BlockSize, len(data)-offset)
		for i := 0; i < n; i++ {
			out[offset+i] = data[offset+i] ^ stream[i]
		}
		previous = stream
	}
	return out, nil
}

// equal reports whether two byte slices are equal.
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

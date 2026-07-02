package security

import (
	"bytes"
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
)

func TestCMACRFC4493(t *testing.T) {
	key := [16]byte{0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6, 0xab, 0xf7, 0x15, 0x88, 0x09, 0xcf, 0x4f, 0x3c}
	mac, err := cmac(key, []byte{})
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0xbb, 0x1d, 0x69, 0x29, 0xe9, 0x59, 0x37, 0x28, 0x7f, 0xa3, 0x7d, 0x12, 0x9b, 0x75, 0x67, 0x46}
	if !bytes.Equal(mac, want) {
		t.Fatalf("% x", mac)
	}
}

func TestSECRoundTrip(t *testing.T) {
	key := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	raw := []byte{1, 2, 3, 4}
	p, err := EncodeSEC_R(key, SLF(RLCExplicit32CMAC32VAES), []byte{0, 0, 0, 1}, enums.RorgSYS_EX, raw)
	if err != nil {
		t.Fatal(err)
	}
	got, err := DecodeSEC_R(key, SLF(RLCExplicit32CMAC32VAES), p)
	if err != nil {
		t.Fatal(err)
	}
	if got.Rorg != enums.RorgSYS_EX || !bytes.Equal(got.Data, raw) || !bytes.Equal(got.RLC, []byte{0, 0, 0, 1}) {
		t.Fatalf("%#v", got)
	}
	p.UserData[0] ^= 1
	if _, err := DecodeSEC_R(key, SLF(RLCExplicit32CMAC32VAES), p); err == nil {
		t.Fatal("expected CMAC error")
	}
}

func TestSLF(t *testing.T) {
	s := SLF(RLCExplicit32CMAC32VAES)
	if s.RLCLength() != 4 || s.CMACLength() != 4 || !s.Encrypted() {
		t.Fatalf("bad slf")
	}
}

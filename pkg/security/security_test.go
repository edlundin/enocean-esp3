package security

import (
	"bytes"
	"errors"
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
)

// TestCMACRFC4493 verifies CMACRFC4493 behavior.
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

// TestSECAppendixA41 verifies SECAppendixA41 behavior.
func TestSECAppendixA41(t *testing.T) {
	key := [16]byte{0x45, 0x6e, 0x4f, 0x63, 0x65, 0x61, 0x6e, 0x20, 0x47, 0x6d, 0x62, 0x48, 0x2e, 0x31, 0x33, 0x00}
	rlc := []byte{0xc0, 0xff, 0xee}
	want := []byte{0x3e, 0xea, 0xc4, 0xa2, 0xdf, 0xc0, 0xff, 0xee, 0xea, 0xf2, 0x0e}
	p, err := EncodeSEC_R(key, SLF(RLCExplicit24CMAC24VAES), rlc, enums.Rorg4BS, []byte{0x08, 0x27, 0xff, 0x80})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(p.UserData, want) {
		t.Fatalf("SEC_R payload = %x, want %x", p.UserData, want)
	}
	got, err := DecodeSEC_R(key, SLF(RLCExplicit24CMAC24VAES), p)
	if err != nil || got.Rorg != enums.Rorg4BS || !bytes.Equal(got.Data, []byte{0x08, 0x27, 0xff, 0x80}) {
		t.Fatalf("decoded = %#v, %v", got, err)
	}
}

// TestSECRoundTrip verifies SECRoundTrip behavior.
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

// TestSECAppendixA43 verifies SECAppendixA43 behavior.
func TestSECAppendixA43(t *testing.T) {
	key := [16]byte{0xe5, 0x08, 0x80, 0xcf, 0x67, 0x79, 0x0d, 0x5d, 0x66, 0xaa, 0x7f, 0x3b, 0x7a, 0xd7, 0x7a, 0x3f}
	data := make([]byte, 30)
	for i := range data {
		data[i] = byte(i)
	}
	want := []byte{0xbb, 0x17, 0xc1, 0x7a, 0x05, 0xca, 0xf5, 0x57, 0x5d, 0xe2, 0x08, 0x30, 0x2f, 0xb5, 0x72, 0xa0, 0xfd, 0x3a, 0x44, 0x34, 0xa4, 0x10, 0x96, 0xf1, 0x02, 0xe6, 0x0d, 0xc2, 0x0d, 0x77, 0x7a, 0x01, 0x02, 0x03, 0x04, 0x3b, 0x4c, 0x38, 0x0f}
	packets, err := EncodeSEC_CDM(key, SLF(RLCExplicit32CMAC32VAES), []byte{1, 2, 3, 4}, 1, enums.RorgMSC, data)
	if err != nil {
		t.Fatal(err)
	}
	parts := make([]ChainPart, len(packets))
	for i, packet := range packets {
		parts[i], err = ParseSEC_CDM(packet)
		if err != nil {
			t.Fatal(err)
		}
	}
	got, done, err := MergeSEC_CDM(parts)
	if err != nil || !done || !bytes.Equal(got, want) {
		t.Fatalf("secure chained payload = %x, want %x, done=%v err=%v", got, want, done, err)
	}
}

// TestSLF verifies SLF behavior.
func TestSLF(t *testing.T) {
	tests := []struct {
		slf                   byte
		rlc, transmitted, mac int
	}{
		{RLCImplicit24CMAC24VAES, 3, 0, 3},
		{RLCExplicit24CMAC24VAES, 3, 3, 3},
		{RLCExplicit24Of32CMAC24VAES, 4, 3, 3},
		{RLCExplicit32CMAC32VAES, 4, 4, 4},
	}
	for _, tc := range tests {
		s := SLF(tc.slf)
		if err := s.Validate(); err != nil || s.RLCLength() != tc.rlc || s.TransmittedRLCLength() != tc.transmitted || s.CMACLength() != tc.mac || !s.Encrypted() {
			t.Fatalf("SLF 0x%02x: rlc=%d tx=%d mac=%d encrypted=%v err=%v", tc.slf, s.RLCLength(), s.TransmittedRLCLength(), s.CMACLength(), s.Encrypted(), err)
		}
	}
	for _, invalid := range []SLF{0x00, 0x73, 0xf2, 0xfb} {
		if err := invalid.Validate(); err == nil {
			t.Fatalf("SLF 0x%02x accepted", byte(invalid))
		}
	}
}

// TestSECNeedsChainingAndPartialRLCState verifies SECNeedsChainingAndPartialRLCState behavior.
func TestSECNeedsChainingAndPartialRLCState(t *testing.T) {
	key := [16]byte{1}
	if _, err := EncodeSEC_R(key, SLF(RLCExplicit32CMAC32VAES), []byte{0, 0, 0, 1}, enums.RorgSYS_EX, make([]byte, 6)); err != nil {
		t.Fatalf("six-byte SEC_R payload rejected: %v", err)
	}
	if _, err := EncodeSEC_R(key, SLF(RLCExplicit32CMAC32VAES), []byte{0, 0, 0, 1}, enums.RorgSYS_EX, make([]byte, 7)); !errors.Is(err, ErrNeedsChaining) {
		t.Fatalf("oversized SEC_R error = %v", err)
	}

	rlc := []byte{1, 2, 3, 4}
	p, err := EncodeSEC_R(key, SLF(RLCExplicit24Of32CMAC24VAES), rlc, enums.RorgRPS, []byte{1})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DecodeSEC_R(key, SLF(RLCExplicit24Of32CMAC24VAES), p); !errors.Is(err, ErrRLCStateRequired) {
		t.Fatalf("partial RLC decode error = %v", err)
	}
	if got, err := DecodeSEC_RWithRLC(key, SLF(RLCExplicit24Of32CMAC24VAES), rlc, p); err != nil || got.Rorg != enums.RorgRPS {
		t.Fatalf("stateful decode = %#v, %v", got, err)
	}
}

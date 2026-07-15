package security

import (
	"bytes"
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/erp1"
)

// TestSECCDMChain verifies SECCDMChain behavior.
func TestSECCDMChain(t *testing.T) {
	data := bytes.Repeat([]byte{0xab}, 40)
	packets, err := SplitSEC_CDM(2, data)
	if err != nil {
		t.Fatal(err)
	}
	if len(packets) != 4 {
		t.Fatalf("parts=%d", len(packets))
	}
	var parts []ChainPart
	for _, p := range packets {
		part, err := ParseSEC_CDM(p)
		if err != nil {
			t.Fatal(err)
		}
		parts = append(parts, part)
	}
	got, done, err := MergeSEC_CDM(parts)
	if err != nil || !done || !bytes.Equal(got, data) {
		t.Fatalf("done=%v err=%v len=%d", done, err, len(got))
	}
}

// TestSECCDMMaxLength verifies SECCDMMaxLength behavior.
func TestSECCDMMaxLength(t *testing.T) {
	packets, err := SplitSEC_CDM(1, make([]byte, MaxChainData))
	if err != nil || len(packets) != MaxChainParts {
		t.Fatalf("parts=%d err=%v", len(packets), err)
	}
	if _, err := SplitSEC_CDM(1, make([]byte, MaxChainData+1)); err == nil {
		t.Fatal("oversized chain accepted")
	}
}

// TestSECCDMAppendixA43 verifies SECCDMAppendixA43 behavior.
func TestSECCDMAppendixA43(t *testing.T) {
	secure := []byte{0xbb, 0x17, 0xc1, 0x7a, 0x05, 0xca, 0xf5, 0x57, 0x5d, 0xe2, 0x08, 0x30, 0x2f, 0xb5, 0x72, 0xa0, 0xfd, 0x3a, 0x44, 0x34, 0xa4, 0x10, 0x96, 0xf1, 0x02, 0xe6, 0x0d, 0xc2, 0x0d, 0x77, 0x7a, 0x01, 0x02, 0x03, 0x04, 0x3b, 0x4c, 0x38, 0x0f}
	want := [][]byte{
		{0x40, 0x00, 0x27, 0xbb, 0x17, 0xc1, 0x7a, 0x05, 0xca, 0xf5, 0x57, 0x5d, 0xe2},
		{0x41, 0x08, 0x30, 0x2f, 0xb5, 0x72, 0xa0, 0xfd, 0x3a, 0x44, 0x34, 0xa4, 0x10, 0x96},
		{0x42, 0xf1, 0x02, 0xe6, 0x0d, 0xc2, 0x0d, 0x77, 0x7a, 0x01, 0x02, 0x03, 0x04, 0x3b},
		{0x43, 0x4c, 0x38, 0x0f},
	}
	packets, err := SplitSEC_CDM(1, secure)
	if err != nil || len(packets) != len(want) {
		t.Fatalf("parts=%d err=%v", len(packets), err)
	}
	for i := range want {
		if !bytes.Equal(packets[i].UserData, want[i]) {
			t.Fatalf("part %d = %x, want %x", i, packets[i].UserData, want[i])
		}
	}
}

// TestSECCDMNeedsMoreDuplicate verifies SECCDMNeedsMoreDuplicate behavior.
func TestSECCDMNeedsMoreDuplicate(t *testing.T) {
	packets, _ := SplitSEC_CDM(1, bytes.Repeat([]byte{1}, 20))
	p0, _ := ParseSEC_CDM(packets[0])
	if _, done, err := MergeSEC_CDM([]ChainPart{p0}); err != nil || done {
		t.Fatalf("done=%v err=%v", done, err)
	}
	if _, _, err := MergeSEC_CDM([]ChainPart{p0, p0}); err == nil {
		t.Fatal("expected duplicate")
	}
	p2 := p0
	p2.Index = 2
	p2.Data = bytes.Repeat([]byte{2}, 20)
	if _, done, err := MergeSEC_CDM([]ChainPart{p0, p2}); err != nil || done {
		t.Fatalf("gapped chain done=%v err=%v", done, err)
	}
}

func TestSECCDMMalformedInputs(t *testing.T) {
	for _, tc := range []struct {
		name string
		data []byte
	}{
		{"short packet", []byte{0x40}},
		{"long packet", make([]byte, 15)},
		{"short first part", []byte{0x40, 0}},
		{"oversized declared length", []byte{0x40, 0xff, 0xff}},
		{"oversized first part", append([]byte{0x40}, make([]byte, 13)...)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ParseSEC_CDM(erp1.Packet{Rorg: enums.RorgSEC_CDM, UserData: tc.data}); err == nil {
				t.Fatal("malformed packet accepted")
			}
		})
	}

	parts := []ChainPart{{Seq: 1, Index: 0, Length: 1, Data: []byte{1, 2}}}
	if _, _, err := MergeSEC_CDM(parts); err == nil {
		t.Fatal("data exceeding declared length accepted")
	}
}

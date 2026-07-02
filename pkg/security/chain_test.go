package security

import (
	"bytes"
	"testing"
)

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

func TestSECCDMNeedsMoreDuplicate(t *testing.T) {
	packets, _ := SplitSEC_CDM(1, bytes.Repeat([]byte{1}, 20))
	p0, _ := ParseSEC_CDM(packets[0])
	if _, done, err := MergeSEC_CDM([]ChainPart{p0}); err != nil || done {
		t.Fatalf("done=%v err=%v", done, err)
	}
	if _, _, err := MergeSEC_CDM([]ChainPart{p0, p0}); err == nil {
		t.Fatal("expected duplicate")
	}
}

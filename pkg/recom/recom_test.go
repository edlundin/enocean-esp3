package recom

import (
	"bytes"
	"reflect"
	"testing"
)

func TestProductID(t *testing.T) {
	p := ProductID{Manufacturer: 0x0123, Product: 0x456789ab}
	if got := p.MarshalBinary(); !bytes.Equal(got, []byte{1, 0x23, 0x45, 0x67, 0x89, 0xab}) {
		t.Fatalf("% x", got)
	}
	back, err := ParseProductID(p.MarshalBinary())
	if err != nil || back != p {
		t.Fatalf("%#v %v", back, err)
	}
}

func TestParamRecords(t *testing.T) {
	recs := []ParamRecord{{Index: 0x1234, Value: []byte{1, 2, 3}}, {Index: 0xabcd}}
	b, err := MarshalParamRecords(recs)
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x12, 0x34, 3, 1, 2, 3, 0xab, 0xcd, 0}; !bytes.Equal(b, want) {
		t.Fatalf("% x", b)
	}
	back, err := ParseParamRecords(b)
	if err != nil || !reflect.DeepEqual(back, recs) {
		t.Fatalf("%#v %v", back, err)
	}
}

func TestParamRecordLimits(t *testing.T) {
	if _, err := MarshalParamRecords([]ParamRecord{{Value: make([]byte, 65)}}); err == nil {
		t.Fatal("expected value length error")
	}
	if _, err := ParseParamRecords([]byte{0, 1, 65}); err == nil {
		t.Fatal("expected parse length error")
	}
}

func TestLinkEntry(t *testing.T) {
	e := LinkEntry{EEP: [3]byte{0xa5, 2, 1}, DeviceID: 0x01020304, Data: [2]byte{0xaa, 0xbb}}
	if got := e.MarshalBinary(); !bytes.Equal(got, []byte{0xa5, 2, 1, 1, 2, 3, 4, 0xaa, 0xbb}) {
		t.Fatalf("% x", got)
	}
	back, err := ParseLinkEntry(e.MarshalBinary())
	if err != nil || back != e {
		t.Fatalf("%#v %v", back, err)
	}
}

package srm

import (
	"bytes"
	"testing"
)

func TestSYSExHeader(t *testing.T) {
	b, err := (Message{Function: FuncPing, Payload: []byte{1, 2}}).MarshalSYSEx()
	if err != nil || !bytes.Equal(b, []byte{0, 6, 1, 2}) {
		t.Fatalf("% x %v", b, err)
	}
	back, err := ParseSYSEx(b)
	if err != nil || back.Function != FuncPing || back.ManufacturerID != nil || !bytes.Equal(back.Payload, []byte{1, 2}) {
		t.Fatalf("%#v %v", back, err)
	}
}

func TestManufacturerSYSExHeader(t *testing.T) {
	mid := uint16(0x123)
	b, err := (Message{ManufacturerID: &mid, Function: 0x456, Payload: []byte{9}}).MarshalSYSEx()
	if err != nil || !bytes.Equal(b, []byte{0x92, 0x34, 0x56, 9}) {
		t.Fatalf("% x %v", b, err)
	}
	back, err := ParseSYSEx(b)
	if err != nil || back.ManufacturerID == nil || *back.ManufacturerID != mid || back.Function != 0x456 {
		t.Fatalf("%#v %v", back, err)
	}
}

func TestQueryStatusAnswer(t *testing.T) {
	a := QueryStatusAnswer{LastFunction: FuncMemoryRead, Return: ReturnTooMuchData}
	if got := a.Payload(); !bytes.Equal(got, []byte{0x02, 0x04, 0x0e}) {
		t.Fatalf("% x", got)
	}
	back, err := ParseQueryStatusAnswer(a.Payload())
	if err != nil || back != a {
		t.Fatalf("%#v %v", back, err)
	}
}

func TestRPCPayloads(t *testing.T) {
	if got := RemoteLearnPayload(true); !bytes.Equal(got, []byte{1}) {
		t.Fatalf("% x", got)
	}
	if got := MemoryReadPayload(0x01020304, 5); !bytes.Equal(got, []byte{1, 2, 3, 4, 5}) {
		t.Fatalf("% x", got)
	}
	got, err := MemoryWritePayload(0x01020304, []byte{0xaa, 0xbb})
	if err != nil || !bytes.Equal(got, []byte{1, 2, 3, 4, 2, 0xaa, 0xbb}) {
		t.Fatalf("% x %v", got, err)
	}
}

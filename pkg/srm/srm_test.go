package srm

import (
	"bytes"
	"testing"
)

// TestSYSExHeader verifies SYSExHeader behavior.
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

// TestManufacturerSYSExHeader verifies ManufacturerSYSExHeader behavior.
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

// TestQueryStatusAnswer verifies QueryStatusAnswer behavior.
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

// TestRejectReservedHeaderBits verifies RejectReservedHeaderBits behavior.
func TestRejectReservedHeaderBits(t *testing.T) {
	if _, err := ParseSYSEx([]byte{0x10, 0x06}); err == nil {
		t.Fatal("reserved Alliance header bits accepted")
	}
	for _, b := range [][]byte{{0x00, 0x00}, {0x80, 0x00, 0x00}} {
		if _, err := ParseSYSEx(b); err == nil {
			t.Fatalf("reserved function zero accepted in % x", b)
		}
	}
	if _, err := ParseQueryStatusAnswer([]byte{0xf2, 0x04, 0x00}); err == nil {
		t.Fatal("reserved query status bits accepted")
	}
	mid := uint16(0x800)
	if _, err := (Message{ManufacturerID: &mid, Function: FuncPing}).MarshalSYSEx(); err == nil {
		t.Fatal("oversized manufacturer ID accepted")
	}
}

// TestRPCPayloads verifies RPCPayloads behavior.
func TestRPCPayloads(t *testing.T) {
	if got := RemoteLearnPayload(true); !bytes.Equal(got, []byte{1}) {
		t.Fatalf("% x", got)
	}
	if got := RemoteLearnPayload(false); !bytes.Equal(got, []byte{3}) {
		t.Fatalf("% x", got)
	}
	for _, tc := range []struct {
		payload []byte
		valid   bool
	}{{[]byte{1}, true}, {[]byte{6}, true}, {[]byte{0}, false}, {[]byte{7}, false}, {nil, false}} {
		_, err := ParseRemoteLearnPayload(tc.payload)
		if (err == nil) != tc.valid {
			t.Fatalf("payload %x valid=%v err=%v", tc.payload, tc.valid, err)
		}
	}
	if got := MemoryReadPayload(0x01020304, 5); !bytes.Equal(got, []byte{1, 2, 3, 4, 5}) {
		t.Fatalf("% x", got)
	}
	got, err := MemoryWritePayload(0x01020304, []byte{0xaa, 0xbb})
	if err != nil || !bytes.Equal(got, []byte{1, 2, 3, 4, 2, 0xaa, 0xbb}) {
		t.Fatalf("% x %v", got, err)
	}
}

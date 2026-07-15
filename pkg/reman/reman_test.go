package reman

import (
	"bytes"
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/deviceid"
)

// TestPacketRoundTrip verifies PacketRoundTrip behavior.
func TestPacketRoundTrip(t *testing.T) {
	msg := Message{Seq: 1, ManufacturerID: ManufacturerID, Function: FuncQueryID, Payload: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, SourceID: 0x01020304, DestinationID: deviceid.BroadcastId()}
	packets, err := msg.Packets()
	if err != nil {
		t.Fatal(err)
	}
	if len(packets) != 2 || packets[0].UserData[0] != 0x40 || packets[1].UserData[0] != 0x41 {
		t.Fatalf("%#v", packets)
	}
	var parts []Part
	for _, p := range packets {
		part, err := ParsePacket(p)
		if err != nil {
			t.Fatal(err)
		}
		parts = append(parts, part)
	}
	back, done, err := Merge(parts)
	if err != nil || !done {
		t.Fatalf("done=%v err=%v", done, err)
	}
	if back.Function != msg.Function || back.ManufacturerID != msg.ManufacturerID || !bytes.Equal(back.Payload, msg.Payload) {
		t.Fatalf("%#v", back)
	}
}

// TestHeaderPacking verifies HeaderPacking behavior.
func TestHeaderPacking(t *testing.T) {
	b := make([]byte, 4)
	putHeader(b, 508, 0x7ff, 0x804)
	l, m, f := getHeader(b)
	if l != 508 || m != 0x7ff || f != 0x804 {
		t.Fatalf("%d %x %x", l, m, f)
	}
}

// TestMergeNeedsMoreAndDuplicate verifies MergeNeedsMoreAndDuplicate behavior.
func TestMergeNeedsMoreAndDuplicate(t *testing.T) {
	msg := Message{Seq: 2, ManufacturerID: 1, Function: 2, Payload: []byte{1, 2, 3, 4, 5}, SourceID: 1}
	packets, _ := msg.Packets()
	p0, _ := ParsePacket(packets[0])
	if _, done, err := Merge([]Part{p0}); err != nil || done {
		t.Fatalf("done=%v err=%v", done, err)
	}
	if _, _, err := Merge([]Part{p0, p0}); err == nil {
		t.Fatal("expected duplicate error")
	}
	p2, _ := ParsePacket(packets[1])
	p2.Index = 2
	if _, done, err := Merge([]Part{p0, p2}); err != nil || done {
		t.Fatalf("gapped chain done=%v err=%v", done, err)
	}
}

// TestPacketsValidateIdentifiers verifies PacketsValidateIdentifiers behavior.
func TestPacketsValidateIdentifiers(t *testing.T) {
	valid := Message{Seq: 1, ManufacturerID: 0x7ff, Function: 0xfff}
	if _, err := valid.Packets(); err != nil {
		t.Fatalf("valid identifiers rejected: %v", err)
	}
	for _, invalid := range []Message{
		{Seq: 1, ManufacturerID: 0x800, Function: 1},
		{Seq: 1, ManufacturerID: 1, Function: 0},
		{Seq: 1, ManufacturerID: 1, Function: 0x1000},
	} {
		if _, err := invalid.Packets(); err == nil {
			t.Fatalf("invalid message accepted: %#v", invalid)
		}
	}
}

// TestCodePayload verifies CodePayload behavior.
func TestCodePayload(t *testing.T) {
	b, err := CodePayload(0x12345678)
	if err != nil || !bytes.Equal(b, []byte{0x12, 0x34, 0x56, 0x78}) {
		t.Fatalf("% x %v", b, err)
	}
	if _, err := CodePayload(0); err == nil {
		t.Fatal("reserved code accepted")
	}
	if s, err := ParseStatusAnswer([]byte{byte(ReturnSessionClosed)}); err != nil || s.Return != ReturnSessionClosed {
		t.Fatalf("%#v %v", s, err)
	}
}

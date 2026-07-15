package smartack

import (
	"reflect"
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/deviceid"
	"github.com/edlundin/enocean-esp3/pkg/eep"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/erp1"
)

// TestRoundTripMessages verifies RoundTripMessages behavior.
func TestRoundTripMessages(t *testing.T) {
	prof, _ := eep.FromTriplet(enums.Rorg4BS, 0x02, 0x01)
	msgs := []Message{
		LearnRequest{RequestCode: RequestDefaultSensor, ManufacturerID: 0x123, EEP: prof, RSSI: 0x44, RepeaterID: 0x01020304},
		LearnReply{ResponseTime: 0x1234, AckCode: 0x20, SensorID: 0x11223344},
		LearnAcknowledge{ResponseTime: 0x4567, AckCode: 0x01, MailboxIndex: 0x7f},
		LearnReclaim{},
		DataReclaim{MailboxIndex: 0x42},
		Signal{Index: SignalReset},
	}
	for _, msg := range msgs {
		pkt := msg.ERP1(deviceid.DeviceID(0xaabbccdd))
		got, err := Parse(pkt)
		if err != nil {
			t.Fatalf("%T: %v", msg, err)
		}
		if !reflect.DeepEqual(got, msg) {
			t.Fatalf("%T got %#v want %#v", msg, got, msg)
		}
	}
}

// TestExactPayloads verifies ExactPayloads behavior.
func TestExactPayloads(t *testing.T) {
	prof, _ := eep.FromTriplet(enums.Rorg4BS, 2, 1)
	lr := LearnRequest{RequestCode: RequestDefaultSensor, ManufacturerID: 0x123, EEP: prof, RSSI: 0x44, RepeaterID: 0x01020304}.ERP1(0)
	wantLR := []byte{0xf9, 0x23, 0xa5, 0x02, 0x01, 0x44, 0x01, 0x02, 0x03, 0x04}
	if !reflect.DeepEqual(lr.UserData, wantLR) || lr.Rorg != enums.RorgSM_LRN_REQ || lr.SubTelNum != 3 {
		t.Fatalf("% x", lr.UserData)
	}

	la := LearnAcknowledge{ResponseTime: 0x1234, AckCode: 0x10, MailboxIndex: 0x05}.ERP1(0)
	if want := []byte{0x02, 0x12, 0x34, 0x10, 0x05}; !reflect.DeepEqual(la.UserData, want) || la.SubTelNum != 1 {
		t.Fatalf("% x", la.UserData)
	}

	dr := DataReclaim{MailboxIndex: 0x7f}.ERP1(0)
	if dr.UserData[0] != 0xff {
		t.Fatalf("% x", dr.UserData)
	}
}

// TestRejectBadPackets verifies RejectBadPackets behavior.
func TestRejectBadPackets(t *testing.T) {
	bad := []erp1.Packet{
		{Rorg: enums.Rorg4BS},
		{Rorg: enums.RorgSM_LRN_REQ, UserData: []byte{1}},
		{Rorg: enums.RorgSM_LRN_ANS, UserData: []byte{0x09}},
		{Rorg: enums.RorgSIGNAL, UserData: []byte{0x04}},
	}
	for _, p := range bad {
		if _, err := Parse(p); err == nil {
			t.Fatalf("expected error for %#v", p)
		}
	}
}

// TestAckCodeClass verifies AckCodeClass behavior.
func TestAckCodeClass(t *testing.T) {
	if AckCode(0x12).Class() != "failed learn-in" {
		t.Fatal(AckCode(0x12).Class())
	}
	if AckCode(0xff).Class() != "application-specific" {
		t.Fatal(AckCode(0xff).Class())
	}
}

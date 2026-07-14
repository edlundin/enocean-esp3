package profiles

import (
	"math"
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/eep"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/erp1"
)

func TestBits(t *testing.T) {
	if got := getBits([]byte{0x01}, 7, 1); got != 1 {
		t.Fatal(got)
	}
	if got := getBits([]byte{0, 0, 0x80, 0x08}, 16, 8); got != 0x80 {
		t.Fatal(got)
	}
	if got := getBits([]byte{0, 0, 0x80, 0x08}, 28, 1); got != 1 {
		t.Fatal(got)
	}
	b := []byte{0, 0}
	setBits(b, 5, 6, 0x2a)
	if got := getBits(b, 5, 6); got != 0x2a {
		t.Fatal(got)
	}
}

func TestProfiles(t *testing.T) {
	d, err := ParseUserData(mustEEP(enums.Rorg1BS, 0, 1), []byte{0x09}, 0)
	if err != nil || !d.(D50001).ContactClosed || d.(D50001).LearnButton {
		t.Fatalf("%#v %v", d, err)
	}

	f, err := ParseUserData(mustEEP(enums.RorgRPS, 1, 1), []byte{0x10}, 0)
	if err != nil || !f.(F60101).Pressed {
		t.Fatalf("%#v %v", f, err)
	}

	a, err := ParseUserData(mustEEP(enums.Rorg4BS, 2, 1), []byte{0, 0, 128, 0x08}, 0)
	if err != nil {
		t.Fatal(err)
	}
	at := a.(A50201)
	if at.TemperatureRaw != 128 || math.Abs(at.TemperatureC-(-20.0784314)) > 0.001 || at.LearnButton {
		t.Fatalf("%#v", at)
	}
	out, _, err := at.MarshalERP1UserData()
	if err != nil || out[2] != 128 || out[3]&0x08 == 0 {
		t.Fatalf("% x %v", out, err)
	}
}

func TestParsePacketRorgMismatch(t *testing.T) {
	prof, _ := eep.FromTriplet(enums.Rorg4BS, 2, 1)
	_, err := ParsePacket(erp1.Packet{Rorg: enums.Rorg1BS}, prof)
	if err == nil {
		t.Fatal("expected mismatch error")
	}
}

package profiles

import (
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/eep"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/erp1"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
)

func TestGenericDecodeEncodeGeneratedRegistry(t *testing.T) {
	prof := mustEEP(enums.RorgRPS, 0x02, 0x01) // generated-only path
	got, err := ParseUserData(prof, []byte{0x11}, 0)
	if err != nil {
		t.Fatal(err)
	}
	d, ok := got.(Decoded)
	if !ok {
		t.Fatalf("got %T", got)
	}
	if d.Values["EB"].Raw != 1 || d.Values["SA"].Raw != 1 {
		t.Fatalf("%#v", d.Values)
	}
	data, _, err := Encode(prof, map[string]uint64{"EB": 1, "SA": 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 1 || data[0] != 0x11 {
		t.Fatalf("% x", data)
	}
}

// eep268.xml defines D2-00-01 Message Type A DB_1=01, DB_0=81
// for MI=1, KP=Presence, and CV=Configuration data valid. The fixed ESP3
// frame wraps those bytes in a valid ERP1 telegram.
func TestDecodeD20001SpecVector(t *testing.T) {
	profile, err := eep.FromString("D2-00-01")
	if err != nil {
		t.Fatal(err)
	}
	telegram, err := esp3.NewEsp3TelegramFromHexString("55000807013dd20181010203040001ffffffff4000de")
	if err != nil {
		t.Fatal(err)
	}
	packet, err := erp1.NewPacketFromEsp3(telegram)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParsePacket(packet, profile)
	if err != nil {
		t.Fatal(err)
	}
	got := parsed.(Decoded)
	for field, want := range map[string]uint64{"MI": 1, "KP": 1, "CV": 1} {
		if actual := got.Values[field].Raw; actual != want {
			t.Errorf("%s = %d, want %d", field, actual, want)
		}
	}
}

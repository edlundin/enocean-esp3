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
	got, err := ParseUserData(prof, []byte{0x11}, 0xa5)
	if err != nil {
		t.Fatal(err)
	}
	d, ok := got.(Decoded)
	if !ok {
		t.Fatalf("got %T", got)
	}
	if d.EEP() != prof || d.Values["EB"].Raw != 1 || d.Values["SA"].Raw != 1 {
		t.Fatalf("%#v", d)
	}
	data, status, err := d.MarshalERP1UserData()
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 1 || data[0] != 0x11 || status != 0xa5 {
		t.Fatalf("data=% x status=%02x", data, status)
	}
}

func TestGenericFormattingAndErrors(t *testing.T) {
	prof := mustEEP(enums.RorgRPS, 0x02, 0x01)
	d := Decoded{
		Profile: Profile{EEP: prof},
		Values: map[string]Value{
			"zraw":  {Raw: 7},
			"atext": {Text: "on"},
			"munit": {Scaled: 12.345, Unit: "°C"},
		},
	}
	if got, want := d.Format(), prof.String()+" atext=on munit=12.35°C zraw=7"; got != want {
		t.Fatalf("Format() = %q, want %q", got, want)
	}

	unknown := eep.EEP{Rorg: enums.Rorg(0xff), Func: 0xff, Type: 0xff}
	if _, err := Decode(unknown, nil, 0); err == nil {
		t.Fatal("Decode accepted an unsupported EEP")
	}
	if _, _, err := Encode(unknown, nil); err == nil {
		t.Fatal("Encode accepted an unsupported EEP")
	}
}

func TestFieldKeyFallbacks(t *testing.T) {
	for _, tc := range []struct {
		name  string
		field Field
		index int
		want  string
	}{
		{name: "shortcut", field: Field{Name: "name", Shortcut: "shortcut"}, want: "shortcut"},
		{name: "name", field: Field{Name: "name"}, want: "name"},
		{name: "index", index: 3, want: "field3"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := fieldKey(tc.field, tc.index); got != tc.want {
				t.Fatalf("fieldKey() = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestDecodeD20001SpecVector verifies the D2-00-01 Message Type A vector.
// eep268.xml defines DB_1=01 and DB_0=81 for MI=1, KP=Presence, and
// CV=Configuration data valid. The fixed ESP3 frame wraps those bytes in a
// valid ERP1 telegram.
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

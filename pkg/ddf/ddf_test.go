package ddf

import (
	"strings"
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
)

// TestParseDDF verifies ParseDDF behavior.
func TestParseDDF(t *testing.T) {
	f, err := Parse(strings.NewReader(`<Enocean_Devices schemaVersion="2.0"><Device Product_ID="0x001122334455"><RX><EEP Rorg="0xA5" Func="0x02" Type="0x01"/></RX><TX><EURID><EEP Rorg="0xF6" Func="0x01" Type="0x01"/></EURID></TX></Device></Enocean_Devices>`))
	if err != nil {
		t.Fatal(err)
	}
	if f.Version != "2.0" || len(f.Devices) != 1 {
		t.Fatalf("%#v", f)
	}
	prof, err := f.Devices[0].RX.EEP[0].EEP()
	if err != nil || prof.Rorg != enums.Rorg4BS || prof.Func != 2 || prof.Type != 1 {
		t.Fatalf("%#v %v", prof, err)
	}
}

// TestRejectMissingRequiredDDFContent verifies RejectMissingRequiredDDFContent behavior.
func TestRejectMissingRequiredDDFContent(t *testing.T) {
	for _, input := range []string{
		`<Enocean_Devices><Device Product_ID="0x001122334455"/></Enocean_Devices>`,
		`<Enocean_Devices schemaVersion="   "><Device Product_ID="0x001122334455"/></Enocean_Devices>`,
		`<Enocean_Devices schemaVersion="2.0"/>`,
	} {
		if _, err := Parse(strings.NewReader(input)); err == nil {
			t.Fatalf("invalid DDF accepted: %s", input)
		}
	}
}

// TestRejectBadProductID verifies RejectBadProductID behavior.
func TestRejectBadProductID(t *testing.T) {
	_, err := Parse(strings.NewReader(`<Enocean_Devices schemaVersion="2.0"><Device Product_ID="0x1234"/></Enocean_Devices>`))
	if err == nil {
		t.Fatal("expected error")
	}
}

// TestEEPRefOptionalFuncType verifies EEPRefOptionalFuncType behavior.
func TestEEPRefOptionalFuncType(t *testing.T) {
	prof, err := (EEPRef{Rorg: "0xD1"}).EEP()
	if err != nil || prof.Rorg != enums.RorgMSC || prof.Func != 0 || prof.Type != 0 {
		t.Fatalf("%#v %v", prof, err)
	}
}

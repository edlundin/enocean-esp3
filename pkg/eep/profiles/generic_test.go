package profiles

import (
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
)

func TestGenericDecodeEncodeGeneratedRegistry(t *testing.T) {
	prof := mustEEP(enums.RorgRPS, 0x02, 0x01) // generated-only path
	got, err := ParseUserData(prof, []byte{0x88}, 0)
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
	if len(data) != 1 || data[0] != 0x88 {
		t.Fatalf("% x", data)
	}
}

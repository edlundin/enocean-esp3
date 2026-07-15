package security

import (
	"encoding/hex"
	"testing"
)

// TestCMACVectors verifies CMACVectors behavior.
func TestCMACVectors(t *testing.T) {
	key := [16]byte{0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6, 0xab, 0xf7, 0x15, 0x88, 0x09, 0xcf, 0x4f, 0x3c}
	cases := []struct{ msg, mac string }{
		{"", "bb1d6929e95937287fa37d129b756746"},
		{"6bc1bee22e409f96e93d7e117393172a", "070a16b46b4d4144f79bdd9dd04a287c"},
		{"6bc1bee22e409f96e93d7e117393172aae2d8a57", "7d85449ea6ea19c823a7bf78837dfade"},
	}
	for _, c := range cases {
		msg, _ := hex.DecodeString(c.msg)
		got, err := cmac(key, msg)
		if err != nil { t.Fatal(err) }
		if hex.EncodeToString(got) != c.mac { t.Fatalf("cmac(%s) = %x", c.msg, got) }
	}
}

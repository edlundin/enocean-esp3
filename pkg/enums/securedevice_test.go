package enums

import "testing"

// TestSecureDeviceDirection verifies SecureDeviceDirection behavior.
func TestSecureDeviceDirection(t *testing.T) {
	cases := []struct{ b byte; v SecureDeviceDirection; s string }{
		{0x00, SecureDeviceDirectionINBOUND_TABLE, "INBOUND_TABLE"},
		{0x01, SecureDeviceDirectionOUTBOUND_TABLE, "OUTBOUND_TABLE"},
		{0x02, SecureDeviceDirectionOUTBOUND_BROADCAST_TABLE, "OUTBOUND_BROADCAST_TABLE"},
		{0x03, SecureDeviceDirectionALL, "ALL_OR_REMAN_TABLE"},
		{0xff, SecureDeviceDirectionNONE, "NONE"},
	}
	for _, c := range cases {
		v, err := ParseSecureDeviceDirectionFromByte(c.b)
		if err != nil || v != c.v || v.String() != c.s || !v.Valid() { t.Fatalf("%#x => %v %v", c.b, v, err) }
	}
	if _, err := ParseSecureDeviceDirectionFromByte(0xfe); err == nil { t.Fatal("expected error") }
	if SecureDeviceDirection(0xfe).String() != "UNKNOWN" || SecureDeviceDirection(0xfe).Valid() { t.Fatal("invalid direction accepted") }
}

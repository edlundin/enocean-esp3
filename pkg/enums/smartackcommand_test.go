package enums

import "testing"

func TestSmartAckCommand(t *testing.T) {
	cases := []struct{ b byte; v SmartAckCommand; s string }{
		{0x01, SmartAckCommandWR_LEARN_MODE, "WR_LEARN_MODE"},
		{0x02, SmartAckCommandRD_LEARN_MODE, "RD_LEARN_MODE"},
		{0x03, SmartAckCommandWR_LEARN_CONFIRM, "WR_LEARN_CONFIRM"},
		{0x04, SmartAckCommandWR_CLIENT_LEARN_RQ, "WR_CLIENT_LEARN_RQ"},
		{0x05, SmartAckCommandWR_RESET, "WR_RESET"},
		{0x06, SmartAckCommandWR_RD_LEARNED_CLIENTS, "WR_RD_LEARNED_CLIENTS"},
		{0x07, SmartAckCommandWR_RECLAIMS, "WR_RECLAIMS"},
		{0x08, SmartAckCommandWR_WR_POSTMASTER, "WR_WR_POSTMASTER"},
	}
	for _, c := range cases {
		v, err := ParseSmartAckCommandFromByte(c.b)
		if err != nil || v != c.v || v.String() != c.s || !v.Valid() { t.Fatalf("%#x => %v %v", c.b, v, err) }
	}
	if _, err := ParseSmartAckCommandFromByte(0xff); err == nil { t.Fatal("expected error") }
	if SmartAckCommand(0xff).String() != "UNKNOWN" || SmartAckCommand(0xff).Valid() { t.Fatal("invalid smart ack command accepted") }
}

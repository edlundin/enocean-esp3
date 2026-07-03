package enums

import "testing"

func TestFilterCriterion(t *testing.T) {
	cases := []struct{ b byte; v FilterCriterion; s string }{
		{0x01, FilterCriterionSENDER_ID, "SENDER_ID"},
		{0x02, FilterCriterionRORG, "RORG"},
		{0x03, FilterCriterionRSSI, "RSSI"},
		{0x04, FilterCriterionDESTINATION_ID, "DESTINATION_ID"},
	}
	for _, c := range cases {
		v, err := ParseFilterFromByte(c.b)
		if err != nil || v != c.v || v.String() != c.s || !v.Valid() { t.Fatalf("%#x => %v %v", c.b, v, err) }
	}
	if _, err := ParseFilterFromByte(0xff); err == nil { t.Fatal("expected error") }
	if FilterCriterion(0xff).String() != "UNKNOWN" || FilterCriterion(0xff).Valid() { t.Fatal("invalid filter accepted") }
}

func TestFilterActionMask(t *testing.T) {
	cases := []struct{ b byte; v FilterActionMask; s string }{
		{0x00, FilterActionNO_FORWARD, "NO_FORWARD"},
		{0x40, FilterActionNO_REPEAT, "NO_REPEAT"},
		{0x80, FilterActionFORWARD, "FORWARD"},
		{0xc0, FilterActionREPEAT, "REPEAT"},
	}
	for _, c := range cases {
		v, err := ParseFilterActionMaskFromByte(c.b)
		if err != nil || v != c.v || v.String() != c.s || !v.Valid() { t.Fatalf("%#x => %v %v", c.b, v, err) }
	}
	if _, err := ParseFilterActionMaskFromByte(0xff); err == nil { t.Fatal("expected error") }
	if FilterActionMask(0xff).String() != "UNKNOWN" || FilterActionMask(0xff).Valid() { t.Fatal("invalid action accepted") }
}

func TestFilerOperator(t *testing.T) {
	cases := []struct{ b byte; v FilerOperator; s string }{
		{0x00, FilerOperatorOR_ALL_FILTERS, "OR_ALL_FILTERS"},
		{0x01, FilerOperatorAND_ALL_FILTERS, "AND_ALL_FILTERS"},
		{0x08, FilerOperatorOR_FOR_RECEIVE_AND_FOR_REPEAT, "OR_FOR_RECEIVE_AND_FOR_REPEAT"},
		{0x09, FilerOperatorAND_FOR_RECEIVE_OR_FOR_REPEAT, "AND_FOR_RECEIVE_OR_FOR_REPEAT"},
	}
	for _, c := range cases {
		v, err := ParseFilerOperatorFromByte(c.b)
		if err != nil || v != c.v || v.String() != c.s || !v.Valid() { t.Fatalf("%#x => %v %v", c.b, v, err) }
	}
	if _, err := ParseFilerOperatorFromByte(0xff); err == nil { t.Fatal("expected error") }
	if FilerOperator(0xff).String() != "UNKNOWN" || FilerOperator(0xff).Valid() { t.Fatal("invalid operator accepted") }
}

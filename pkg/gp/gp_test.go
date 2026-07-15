package gp

import (
	"encoding/hex"
	"errors"
	"fmt"
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
)

// TestRorgs verifies Rorgs behavior.
func TestRorgs(t *testing.T) {
	for _, r := range []enums.Rorg{enums.RorgGP_TI, enums.RorgGP_TR, enums.RorgGP_CD, enums.RorgGP_SD} {
		if !IsRorg(r) || !r.Valid() {
			t.Fatalf("%s should be a valid GP RORG", r)
		}
	}
}

// TestTeachInHeaders verifies TeachInHeaders behavior.
func TestTeachInHeaders(t *testing.T) {
	b, err := EncodeRequestHeader(RequestHeader{ManufacturerID: 0x7ff, Bidirectional: true, Purpose: PurposeTeachIn})
	if err != nil || hex.EncodeToString(b) != "fff0" {
		t.Fatalf("request header = %x, %v", b, err)
	}
	h, err := DecodeRequestHeader(b)
	if err != nil || h.ManufacturerID != 0x7ff || !h.Bidirectional || h.Purpose != PurposeTeachIn {
		t.Fatalf("decoded request header = %#v, %v", h, err)
	}

	b, err = EncodeRequestHeader(RequestHeader{ManufacturerID: 0x7ff, Bidirectional: true, Purpose: PurposeToggle})
	if err != nil || hex.EncodeToString(b) != "fff8" {
		t.Fatalf("toggle request header = %x, %v", b, err)
	}

	b, err = EncodeResponseHeader(ResponseHeader{ManufacturerID: 0x7ff, Result: ResultRejectedChannels})
	if err != nil || hex.EncodeToString(b) != "fff8" {
		t.Fatalf("response header = %x, %v", b, err)
	}
	r, err := DecodeResponseHeader(b)
	if err != nil || r.ManufacturerID != 0x7ff || r.Result != ResultRejectedChannels {
		t.Fatalf("decoded response header = %#v, %v", r, err)
	}
}

// TestChannelDefinitions verifies ChannelDefinitions behavior.
func TestChannelDefinitions(t *testing.T) {
	data := Channel{Type: ChannelData, SignalType: 0x06, ValueType: ValueCurrent, ResolutionCode: 0x5, EngineeringMin: 0, ScalingMin: 1, EngineeringMax: 5, ScalingMax: 1}
	b, bits, err := EncodeChannelDefinition(data)
	if err != nil || bits != 40 || hex.EncodeToString(b) != "4195001051" {
		t.Fatalf("data channel = %x/%d, %v", b, bits, err)
	}
	got, used, err := DecodeChannelDefinition(b, 0)
	if err != nil || used != 40 || got != data {
		t.Fatalf("decoded data channel = %#v/%d, %v", got, used, err)
	}

	signed := Channel{Type: ChannelData, SignalType: 0x06, ValueType: ValueCurrent, ResolutionCode: 0x5, EngineeringMin: 0x80, ScalingMin: 1, EngineeringMax: 0xff, ScalingMax: 1}
	b, bits, err = EncodeChannelDefinition(signed)
	if err != nil || bits != 40 || hex.EncodeToString(b) != "4195801ff1" {
		t.Fatalf("signed data channel = %x/%d, %v", b, bits, err)
	}
	got, used, err = DecodeChannelDefinition(b, 0)
	if err != nil || used != 40 || got != signed {
		t.Fatalf("decoded signed channel = %#v/%d, %v", got, used, err)
	}
	if min, max := got.EngineeringRange(); min != -128 || max != -1 {
		t.Fatalf("signed engineering range = %d..%d", min, max)
	}

	flag := Channel{Type: ChannelFlag, SignalType: 0x09, ValueType: ValueSetPointAbsolute}
	b, bits, err = EncodeChannelDefinition(flag)
	if err != nil || bits != 12 || bitsHex(b, bits) != "826" {
		t.Fatalf("flag channel = %x/%d bits %s, %v", b, bits, bitsHex(b, bits), err)
	}
	got, used, err = DecodeChannelDefinition(b, 0)
	if err != nil || used != 12 || got != flag {
		t.Fatalf("decoded flag channel = %#v/%d, %v", got, used, err)
	}
}

// TestCompleteAndSelectiveData verifies CompleteAndSelectiveData behavior.
func TestCompleteAndSelectiveData(t *testing.T) {
	channels := []Channel{
		{Type: ChannelData, ResolutionCode: 0x5},        // 6 bit
		{Type: ChannelData, ResolutionCode: 0x6},        // 8 bit
		{Type: ChannelEnumeration, ResolutionCode: 0x4}, // 5 bit
	}
	complete, err := EncodeCompleteData(channels, []uint64{0x2a, 0xbc, 0x15})
	if err != nil || hex.EncodeToString(complete) != "aaf2a0" {
		t.Fatalf("complete = %x, %v", complete, err)
	}
	values, err := DecodeCompleteData(channels, complete)
	if err != nil || len(values) != 3 || values[0] != 0x2a || values[1] != 0xbc || values[2] != 0x15 {
		t.Fatalf("decoded complete = %#v, %v", values, err)
	}

	selectiveChannels := []Channel{
		{Type: ChannelData, ResolutionCode: 0x5},        // 6 bit
		{Type: ChannelData, ResolutionCode: 0x1},        // 2 bit
		{Type: ChannelEnumeration, ResolutionCode: 0x4}, // 5 bit
	}
	selective, err := EncodeSelectiveData(selectiveChannels, []SelectedValue{{Index: 1, Value: 0x2}})
	if err != nil || hex.EncodeToString(selective) != "1060" {
		t.Fatalf("selective = %x, %v", selective, err)
	}
	selected, err := DecodeSelectiveData(selectiveChannels, selective)
	if err != nil || len(selected) != 1 || selected[0] != (SelectedValue{Index: 1, Value: 0x2}) {
		t.Fatalf("decoded selective = %#v, %v", selected, err)
	}
}

// TestBitstreamSignedAndBounds verifies BitstreamSignedAndBounds behavior.
func TestBitstreamSignedAndBounds(t *testing.T) {
	b := []byte{0}
	if err := writeSigned(b, 0, 4, -2); err != nil {
		t.Fatal(err)
	}
	if b[0] != 0xe0 {
		t.Fatalf("signed bits = %08b", b[0])
	}
	if got, err := readSigned(b, 0, 4); err != nil || got != -2 {
		t.Fatalf("read signed = %d, %v", got, err)
	}
	if _, err := readUnsigned([]byte{0}, 7, 2); !errors.Is(err, ErrBitstreamOutOfRange) {
		t.Fatalf("bounds error = %v", err)
	}
}

// bitsHex formats a bit field as hexadecimal.
func bitsHex(b []byte, bits int) string {
	v, _ := readUnsigned(b, 0, bits)
	return fmt.Sprintf("%0*x", (bits+3)/4, v)
}

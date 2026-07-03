package gp

import "testing"

func TestDataErrors(t *testing.T) {
	channels := []Channel{{Type: ChannelTeachInInformation}, {Type: ChannelFlag}, {Type: ChannelData, ResolutionCode: 6}}
	if _, err := operationalChannels([]Channel{{Type: ChannelData}}); err == nil { t.Fatal("bad resolution accepted") }
	if _, err := EncodeCompleteData(channels, []uint64{1}); err == nil { t.Fatal("bad complete count accepted") }
	if _, err := DecodeCompleteData(channels, make([]byte, MaxMessageLength+1)); err == nil { t.Fatal("oversized complete data accepted") }
	if _, err := EncodeSelectiveData(channels, make([]SelectedValue, 16)); err == nil { t.Fatal("too many selective values accepted") }
	if _, err := EncodeSelectiveData(channels, []SelectedValue{{Index: -1}}); err == nil { t.Fatal("negative index accepted") }
	if _, err := DecodeSelectiveData(channels, []byte{0x10}); err == nil { t.Fatal("truncated selective data accepted") }
}

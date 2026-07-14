package srm

import (
	"bytes"
	"testing"
)

func TestSRMPayloadHelpers(t *testing.T) {
	if PingPayload() != nil {
		t.Fatal("ping payload should be nil")
	}
	if got := PingResponsePayload(0x42); !bytes.Equal(got, []byte{0x42}) {
		t.Fatalf("ping response = %x", got)
	}
	if got := RemoteLearnPayload(true); !bytes.Equal(got, []byte{1}) {
		t.Fatalf("remote learn true = %x", got)
	}
	if got := RemoteLearnPayload(false); !bytes.Equal(got, []byte{3}) {
		t.Fatalf("remote learn false = %x", got)
	}
}

package profiles

import (
	"bytes"
	"testing"
)

func TestD201SetOutput(t *testing.T) {
	for _, test := range []struct {
		value uint8
		want  []byte
	}{
		{100, []byte{0x01, 0x01, 0x64}},
		{0, []byte{0x01, 0x01, 0x00}},
	} {
		got, err := D201SetOutput(1, test.value)
		if err != nil || !bytes.Equal(got, test.want) {
			t.Fatalf("D201SetOutput(1, %d) = %x, %v; want %x, nil", test.value, got, err, test.want)
		}
	}
	if _, err := D201SetOutput(0, 100); err == nil {
		t.Fatal("zero channel succeeded")
	}
	if _, err := D201SetOutput(1, 101); err == nil {
		t.Fatal("out-of-range output succeeded")
	}
}

func TestParseD201Status(t *testing.T) {
	status, err := ParseD201Status([]byte{0x04, 0x01, 0x64})
	if err != nil || status.Channel != 1 || status.Output != 100 {
		t.Fatalf("ParseD201Status = %+v, %v", status, err)
	}
	for _, data := range [][]byte{{0x04, 0x01}, {0x01, 0x01, 0x64}, {0x04, 0x00, 0x64}, {0x04, 0x01, 0x7f}} {
		if _, err := ParseD201Status(data); err == nil {
			t.Fatalf("ParseD201Status(%x) succeeded", data)
		}
	}
}

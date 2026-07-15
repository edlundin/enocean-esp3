package eep

import (
	"strings"
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
)

// TestFromTriplet verifies FromTriplet behavior.
func TestFromTriplet(t *testing.T) {
	tests := []struct {
		name string
		rorg enums.Rorg
		fn   byte
		typ  byte
		want EEP
	}{
		{"minimum_values", 0x00, 0x00, 0x00, EEP{Rorg: 0x00, Func: 0x00, Type: 0x00}},
		{"maximum_values", 0xff, 0xb0, 0x7f, EEP{Rorg: 0xff, Func: 0xb0, Type: 0x7f}},
		{"common_values", 0xf6, 0x01, 0x01, EEP{Rorg: 0xf6, Func: 0x01, Type: 0x01}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromTriplet(tt.rorg, tt.fn, tt.typ)
			if err != nil {
				t.Fatalf("FromTriplet() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("FromTriplet() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFromString verifies FromString behavior.
func TestFromString(t *testing.T) {
	valid := map[string]EEP{
		"00-00-00": {Rorg: 0x00, Func: 0x00, Type: 0x00},
		"FF-B0-7F": {Rorg: 0xff, Func: 0xb0, Type: 0x7f},
		"D2-A0-01": {Rorg: enums.RorgVLD, Func: 0xa0, Type: 0x01},
		"f6-01-01": {Rorg: enums.RorgRPS, Func: 0x01, Type: 0x01},
	}
	for input, want := range valid {
		t.Run(input, func(t *testing.T) {
			got, err := FromString(input)
			if err != nil {
				t.Fatalf("FromString(%q) error = %v", input, err)
			}
			if got != want {
				t.Fatalf("FromString(%q) = %v, want %v", input, got, want)
			}
		})
	}

	invalid := []struct {
		input string
		want  string
	}{
		{"", "invalid format"},
		{"00-01", "invalid format"},
		{"00-01-02-03", "invalid format"},
		{"100-01-01", "invalid RORG"},
		{"GG-01-01", "invalid RORG"},
		{"00-B1-01", "invalid FUNC"},
		{"00-100-01", "invalid FUNC"},
		{"00-GG-01", "invalid FUNC"},
		{"00-01-80", "invalid TYPE"},
		{"00-01-100", "invalid TYPE"},
		{"00-01-GG", "invalid TYPE"},
	}
	for _, tt := range invalid {
		t.Run(tt.input, func(t *testing.T) {
			_, err := FromString(tt.input)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("FromString(%q) error = %v, want %q", tt.input, err, tt.want)
			}
		})
	}
}

// TestEEPString verifies EEPString behavior.
func TestEEPString(t *testing.T) {
	e := EEP{Rorg: 0xff, Func: 0xb0, Type: 0x7f}
	if got := e.String(); got != "FF-B0-7F" {
		t.Fatalf("String() = %q", got)
	}
	parsed, err := FromString(e.String())
	if err != nil {
		t.Fatal(err)
	}
	if parsed != e {
		t.Fatalf("round trip = %v, want %v", parsed, e)
	}
}

// TestEEPConstants verifies EEPConstants behavior.
func TestEEPConstants(t *testing.T) {
	if minFunc != 0x00 || maxFunc != 0xb0 || minType != 0x00 || maxType != 0x7f {
		t.Fatalf("unexpected bounds: func %02x..%02x type %02x..%02x", minFunc, maxFunc, minType, maxType)
	}
}

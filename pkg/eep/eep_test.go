package eep

import (
	"strings"
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
)

func TestFromTriplet(t *testing.T) {
	t.Run("valid_eep_values", func(t *testing.T) {
		tests := []struct {
			name     string
			rorg     enums.Rorg
			funcVal  byte
			typeVal  byte
			expected EEP
		}{
			{"minimum_values", 0x00, 0x00, 0x00, EEP{Rorg: 0x00, Func: 0x00, Type: 0x00}},
			{"maximum_values", 0xff, 0x60, 0x7f, EEP{Rorg: 0xff, Func: 0x60, Type: 0x7f}},
			{"middle_values", 0x80, 0x30, 0x40, EEP{Rorg: 0x80, Func: 0x30, Type: 0x40}},
			{"common_values", 0xf6, 0x01, 0x01, EEP{Rorg: 0xf6, Func: 0x01, Type: 0x01}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := FromTriplet(tt.rorg, tt.funcVal, tt.typeVal)
				if err != nil {
					t.Errorf("FromTriplet() error = %v, want no error", err)
					return
				}
				if result != tt.expected {
					t.Errorf("FromTriplet() = %v, want %v", result, tt.expected)
				}
			})
		}
	})

	t.Run("invalid_func_out_of_bounds", func(t *testing.T) {
		tests := []struct {
			name    string
			rorg    enums.Rorg
			funcVal byte
			typeVal byte
		}{
			{"func_too_low", 0x00, 0x61, 0x00},
			{"func_too_high", 0x00, 0x61, 0x00},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := FromTriplet(tt.rorg, tt.funcVal, tt.typeVal)
				if err == nil {
					t.Errorf("FromTriplet() expected error for invalid FUNC %v, got no error", tt.funcVal)
					return
				}
				if !strings.Contains(err.Error(), "invalid FUNC: out of bounds") {
					t.Errorf("FromTriplet() error = %v, want 'invalid FUNC: out of bounds'", err)
				}
			})
		}
	})

	t.Run("invalid_type_out_of_bounds", func(t *testing.T) {
		tests := []struct {
			name    string
			rorg    enums.Rorg
			funcVal byte
			typeVal byte
		}{
			{"type_too_low", 0x00, 0x00, 0x80},
			{"type_too_high", 0x00, 0x00, 0x7f + 1},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := FromTriplet(tt.rorg, tt.funcVal, tt.typeVal)
				if err == nil {
					t.Errorf("FromTriplet() expected error for invalid TYPE %v, got no error", tt.typeVal)
					return
				}
				if !strings.Contains(err.Error(), "invalid TYPE: out of bounds") {
					t.Errorf("FromTriplet() error = %v, want 'invalid TYPE: out of bounds'", err)
				}
			})
		}
	})
}

func TestFromString(t *testing.T) {
	t.Run("valid_eep_strings", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected EEP
		}{
			{"minimum_values", "00-00-00", EEP{Rorg: 0x00, Func: 0x00, Type: 0x00}},
			{"maximum_values", "FF-60-7F", EEP{Rorg: 0xff, Func: 0x60, Type: 0x7f}},
			{"middle_values", "80-30-40", EEP{Rorg: 0x80, Func: 0x30, Type: 0x40}},
			{"common_values", "F6-01-01", EEP{Rorg: 0xf6, Func: 0x01, Type: 0x01}},
			{"lowercase", "f6-01-01", EEP{Rorg: 0xf6, Func: 0x01, Type: 0x01}},
			{"mixed_case", "F6-01-01", EEP{Rorg: 0xf6, Func: 0x01, Type: 0x01}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := FromString(tt.input)
				if err != nil {
					t.Errorf("FromString(%q) error = %v, want no error", tt.input, err)
					return
				}
				if result != tt.expected {
					t.Errorf("FromString(%q) = %v, want %v", tt.input, result, tt.expected)
				}
			})
		}
	})

	t.Run("invalid_format", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
		}{
			{"empty_string", ""},
			{"single_field", "00"},
			{"two_fields", "00-01"},
			{"four_fields", "00-01-02-03"},
			{"wrong_separator", "00:01:02"},
			{"no_separator", "000102"},
			{"extra_dashes", "00--01--02"},
			{"leading_dash", "-00-01-02"},
			{"trailing_dash", "00-01-02-"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := FromString(tt.input)
				if err == nil {
					t.Errorf("FromString(%q) expected error for invalid format, got no error", tt.input)
					return
				}
				if !strings.Contains(err.Error(), "invalid format (RR-FF-TT)") {
					t.Errorf("FromString(%q) error = %v, want 'invalid format (RR-FF-TT)'", tt.input, err)
				}
			})
		}
	})

	t.Run("invalid_rorg", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
		}{
			{"non_hex_rorg", "GG-01-01"},
			{"rorg_too_large", "100-01-01"},   // This will be parsed as 0x100 which is > 0xFF
			{"rorg_very_large", "1000-01-01"}, // This will be parsed as 0x1000 which is > 0xFF

			{"empty_rorg", "-01-01"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := FromString(tt.input)
				if err == nil {
					t.Errorf("FromString(%q) expected error for invalid RORG, got no error", tt.input)
					return
				}
				if !strings.Contains(err.Error(), "invalid RORG") {
					t.Errorf("FromString(%q) error = %v, want 'invalid RORG'", tt.input, err)
				}
			})
		}
	})

	t.Run("invalid_func", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
		}{
			{"non_hex_func", "00-GG-01"},
			{"func_too_large", "00-FF-01"},    // This will be parsed as 0xFF which is > 0x60
			{"func_very_large", "00-1000-01"}, // This will be parsed as 0x1000 which is > 0xFF
			{"func_too_high", "00-61-01"},
			{"empty_func", "00--01"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := FromString(tt.input)
				if err == nil {
					t.Errorf("FromString(%q) expected error for invalid FUNC, got no error", tt.input)
					return
				}
				if !strings.Contains(err.Error(), "invalid FUNC") {
					t.Errorf("FromString(%q) error = %v, want 'invalid FUNC'", tt.input, err)
				}
			})
		}
	})

	t.Run("invalid_type", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
		}{
			{"non_hex_type", "00-01-GG"},
			{"type_too_large", "00-01-FF"},    // This will be parsed as 0xFF which is > 0x7F
			{"type_very_large", "00-01-1000"}, // This will be parsed as 0x1000 which is > 0xFF
			{"type_too_high", "00-01-80"},
			{"empty_type", "00-01-"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := FromString(tt.input)
				if err == nil {
					t.Errorf("FromString(%q) expected error for invalid TYPE, got no error", tt.input)
					return
				}
				if !strings.Contains(err.Error(), "invalid TYPE") {
					t.Errorf("FromString(%q) error = %v, want 'invalid TYPE'", tt.input, err)
				}
			})
		}
	})

	t.Run("boundary_values", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected EEP
		}{
			{"rorg_min", "00-01-01", EEP{Rorg: 0x00, Func: 0x01, Type: 0x01}},
			{"rorg_max", "FF-01-01", EEP{Rorg: 0xff, Func: 0x01, Type: 0x01}},
			{"func_min", "01-00-01", EEP{Rorg: 0x01, Func: 0x00, Type: 0x01}},
			{"func_max", "01-60-01", EEP{Rorg: 0x01, Func: 0x60, Type: 0x01}},
			{"type_min", "01-01-00", EEP{Rorg: 0x01, Func: 0x01, Type: 0x00}},
			{"type_max", "01-01-7F", EEP{Rorg: 0x01, Func: 0x01, Type: 0x7f}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := FromString(tt.input)
				if err != nil {
					t.Errorf("FromString(%q) error = %v, want no error", tt.input, err)
					return
				}
				if result != tt.expected {
					t.Errorf("FromString(%q) = %v, want %v", tt.input, result, tt.expected)
				}
			})
		}
	})
}

func TestEEP_String(t *testing.T) {
	t.Run("formatting", func(t *testing.T) {
		tests := []struct {
			name     string
			eep      EEP
			expected string
		}{
			{"minimum_values", EEP{Rorg: 0x00, Func: 0x00, Type: 0x00}, "00-00-00"},
			{"maximum_values", EEP{Rorg: 0xff, Func: 0x60, Type: 0x7f}, "FF-60-7F"},
			{"middle_values", EEP{Rorg: 0x80, Func: 0x30, Type: 0x40}, "80-30-40"},
			{"common_values", EEP{Rorg: 0xf6, Func: 0x01, Type: 0x01}, "F6-01-01"},
			{"single_digits", EEP{Rorg: 0x0a, Func: 0x0b, Type: 0x0c}, "0A-0B-0C"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := tt.eep.String()
				if result != tt.expected {
					t.Errorf("EEP.String() = %v, want %v", result, tt.expected)
				}
			})
		}
	})

	t.Run("round_trip", func(t *testing.T) {
		tests := []struct {
			name string
			eep  EEP
		}{
			{"round_trip_1", EEP{Rorg: 0xf6, Func: 0x01, Type: 0x01}},
			{"round_trip_2", EEP{Rorg: 0xd5, Func: 0x00, Type: 0x01}},
			{"round_trip_3", EEP{Rorg: 0xa5, Func: 0x02, Type: 0x05}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Convert to string
				str := tt.eep.String()

				// Parse back from string
				parsed, err := FromString(str)
				if err != nil {
					t.Errorf("Failed to parse EEP string %q: %v", str, err)
					return
				}

				// Verify they're equal
				if parsed != tt.eep {
					t.Errorf("Round trip failed: original %v, parsed %v", tt.eep, parsed)
				}
			})
		}
	})
}

func TestEEP_Constants(t *testing.T) {
	t.Run("constant_values", func(t *testing.T) {
		if minFunc != 0x00 {
			t.Errorf("minFunc = %v, want 0x00", minFunc)
		}
		if maxFunc != 0x60 {
			t.Errorf("maxFunc = %v, want 0x60", maxFunc)
		}
		if minType != 0x00 {
			t.Errorf("minType = %v, want 0x00", minType)
		}
		if maxType != 0x7f {
			t.Errorf("maxType = %v, want 0x7f", maxType)
		}
	})
}

func TestEEP_Integration(t *testing.T) {
	t.Run("create_and_format", func(t *testing.T) {
		// Create EEP from triplet
		eep, err := FromTriplet(0xf6, 0x01, 0x01)
		if err != nil {
			t.Fatalf("Failed to create EEP from triplet: %v", err)
		}

		// Verify fields
		if eep.Rorg != 0xf6 {
			t.Errorf("EEP.Rorg = %v, want 0xf6", eep.Rorg)
		}
		if eep.Func != 0x01 {
			t.Errorf("EEP.Func = %v, want 0x01", eep.Func)
		}
		if eep.Type != 0x01 {
			t.Errorf("EEP.Type = %v, want 0x01", eep.Type)
		}

		// Format to string
		str := eep.String()
		expected := "F6-01-01"
		if str != expected {
			t.Errorf("EEP.String() = %v, want %v", str, expected)
		}

		// Parse back from string
		parsed, err := FromString(str)
		if err != nil {
			t.Errorf("Failed to parse EEP string %q: %v", str, err)
			return
		}

		// Verify they're equal
		if parsed != eep {
			t.Errorf("Integration test failed: original %v, parsed %v", eep, parsed)
		}
	})
}

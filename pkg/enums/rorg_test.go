package enums

import (
	"testing"
)

// TestParseRorgFromByte verifies ParseRorgFromByte behavior.
func TestParseRorgFromByte(t *testing.T) {
	t.Run("parses all valid rorg types correctly", func(t *testing.T) {
		testCases := []struct {
			input    byte
			expected Rorg
		}{
			{0xf6, RorgRPS},
			{0xd5, Rorg1BS},
			{0xa5, Rorg4BS},
			{0xd2, RorgVLD},
			{0xd1, RorgMSC},
			{0xa6, RorgADT},
			{0xb0, RorgGP_TI},
			{0xb1, RorgGP_TR},
			{0xb2, RorgGP_CD},
			{0xb3, RorgGP_SD},
			{0xc6, RorgSM_LRN_REQ},
			{0xc7, RorgSM_LRN_ANS},
			{0xa7, RorgSM_REC},
			{0xc5, RorgSYS_EX},
			{0x30, RorgSEC},
			{0x31, RorgSEC_R},
			{0x32, RorgSEC_D},
			{0x33, RorgSEC_CDM},
			{0x34, RorgSEC_MAN},
			{0x35, RorgSEC_TI},
			{0xd0, RorgSIGNAL},
			{0xd4, RorgUTE},
		}

		for _, tc := range testCases {
			t.Run(tc.expected.String(), func(t *testing.T) {
				result, err := ParseRorgFromByte(tc.input)
				if err != nil {
					t.Errorf("expected no error for input 0x%02x, got: %s", tc.input, err)
				}
				if result != tc.expected {
					t.Errorf("expected %v, got %v", tc.expected, result)
				}
			})
		}
	})

	t.Run("returns error for invalid rorg type", func(t *testing.T) {
		invalidInputs := []byte{0x00, 0x01, 0x02, 0x10, 0x20, 0x40, 0x80, 0xff}

		for _, input := range invalidInputs {
			t.Run(t.Name(), func(t *testing.T) {
				result, err := ParseRorgFromByte(input)
				if err == nil {
					t.Errorf("expected error for input 0x%02x, got nil", input)
				}
				if err.Error() != "invalid rorg" {
					t.Errorf("expected error 'invalid rorg', got '%s'", err.Error())
				}
				if result != 0 {
					t.Errorf("expected result 0, got %v", result)
				}
			})
		}
	})
}

// TestRorgString verifies RorgString behavior.
func TestRorgString(t *testing.T) {
	t.Run("returns correct string for all valid rorg types", func(t *testing.T) {
		testCases := []struct {
			input    Rorg
			expected string
		}{
			{RorgRPS, "RPS"},
			{Rorg1BS, "1BS"},
			{Rorg4BS, "4BS"},
			{RorgVLD, "VLD"},
			{RorgMSC, "MSC"},
			{RorgADT, "ADT"},
			{RorgGP_TI, "GP_TI"},
			{RorgGP_TR, "GP_TR"},
			{RorgGP_CD, "GP_CD"},
			{RorgGP_SD, "GP_SD"},
			{RorgSM_LRN_REQ, "SM_LRN_REQ"},
			{RorgSM_LRN_ANS, "SM_LRN_ANS"},
			{RorgSM_REC, "SM_REC"},
			{RorgSYS_EX, "SYS_EX"},
			{RorgSEC, "SEC"},
			{RorgSEC_R, "SEC_R"},
			{RorgSEC_D, "SEC_D"},
			{RorgSEC_CDM, "SEC_CDM"},
			{RorgSEC_MAN, "SEC_MAN"},
			{RorgSEC_TI, "SEC_TI"},
			{RorgSIGNAL, "SIGNAL"},
			{RorgUTE, "UTE"},
		}

		for _, tc := range testCases {
			t.Run(tc.expected, func(t *testing.T) {
				result := tc.input.String()
				if result != tc.expected {
					t.Errorf("expected '%s', got '%s'", tc.expected, result)
				}
			})
		}
	})

	t.Run("returns UNKNOWN for invalid rorg types", func(t *testing.T) {
		invalidTypes := []Rorg{0x00, 0x01, 0x02, 0x10, 0x20, 0x40, 0x80, 0xff}

		for _, input := range invalidTypes {
			t.Run(t.Name(), func(t *testing.T) {
				result := input.String()
				if result != "UNKNOWN" {
					t.Errorf("expected 'UNKNOWN' for input %v, got '%s'", input, result)
				}
			})
		}
	})
}

// TestRorgValid verifies RorgValid behavior.
func TestRorgValid(t *testing.T) {
	t.Run("Rorg_Valid", func(t *testing.T) {
		// Test valid rorg types
		validRorgs := []Rorg{
			RorgRPS,
			Rorg1BS,
			Rorg4BS,
			RorgVLD,
			RorgMSC,
			RorgGP_TI,
			RorgGP_TR,
			RorgGP_CD,
			RorgGP_SD,
		}
		for _, rorg := range validRorgs {
			if !rorg.Valid() {
				t.Errorf("Rorg %v should be valid", rorg)
			}
		}

		// Test invalid rorg types
		invalidRorgs := []Rorg{0x00, 0x01, 0x02, 0xFF}
		for _, rorg := range invalidRorgs {
			if rorg.Valid() {
				t.Errorf("Rorg %v should not be valid", rorg)
			}
		}
	})
}

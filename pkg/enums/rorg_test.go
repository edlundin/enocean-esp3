package enums

import (
	"testing"
)

func TestParseRorgFromByte(t *testing.T) {
	t.Run("parses all valid rorg types correctly", func(t *testing.T) {
		testCases := []struct {
			input    uint8
			expected Rorg
		}{
			{0xf6, RorgRPS},
			{0xd5, Rorg1BS},
			{0xa5, Rorg4BS},
			{0xd2, RorgVLD},
			{0xd1, RorgMSC},
			{0xa6, RorgADT},
			{0xc6, RorgSM_LRN_REQ},
			{0xc7, RorgSM_LRN_ANS},
			{0xa7, RorgSM_REC},
			{0xc5, RorgSYS_EX},
			{0x30, RorgSEC},
			{0x31, RorgSEC_ENCAPS},
			{0x34, RorgSEC_MAN},
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
		invalidInputs := []uint8{0x00, 0x01, 0x02, 0x10, 0x20, 0x40, 0x80, 0xff}

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
			{RorgSM_LRN_REQ, "SM_LRN_REQ"},
			{RorgSM_LRN_ANS, "SM_LRN_ANS"},
			{RorgSM_REC, "SM_REC"},
			{RorgSYS_EX, "SYS_EX"},
			{RorgSEC, "SEC"},
			{RorgSEC_ENCAPS, "SEC_ENCAPS"},
			{RorgSEC_MAN, "SEC_MAN"},
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

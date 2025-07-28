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
			{0xf6, RORG_RPS},
			{0xd5, RORG_1BS},
			{0xa5, RORG_4BS},
			{0xd2, RORG_VLD},
			{0xd1, RORG_MSC},
			{0xa6, RORG_ADT},
			{0xc6, RORG_SM_LRN_REQ},
			{0xc7, RORG_SM_LRN_ANS},
			{0xa7, RORG_SM_REC},
			{0xc5, RORG_SYS_EX},
			{0x30, RORG_SEC},
			{0x31, RORG_SEC_ENCAPS},
			{0x34, RORG_SEC_MAN},
			{0xd0, RORG_SIGNAL},
			{0xd4, RORG_UTE},
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

func TestRorg_String(t *testing.T) {
	t.Run("returns correct string for all valid rorg types", func(t *testing.T) {
		testCases := []struct {
			input    Rorg
			expected string
		}{
			{RORG_RPS, "RPS"},
			{RORG_1BS, "1BS"},
			{RORG_4BS, "4BS"},
			{RORG_VLD, "VLD"},
			{RORG_MSC, "MSC"},
			{RORG_ADT, "ADT"},
			{RORG_SM_LRN_REQ, "SM_LRN_REQ"},
			{RORG_SM_LRN_ANS, "SM_LRN_ANS"},
			{RORG_SM_REC, "SM_REC"},
			{RORG_SYS_EX, "SYS_EX"},
			{RORG_SEC, "SEC"},
			{RORG_SEC_ENCAPS, "SEC_ENCAPS"},
			{RORG_SEC_MAN, "SEC_MAN"},
			{RORG_SIGNAL, "SIGNAL"},
			{RORG_UTE, "UTE"},
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

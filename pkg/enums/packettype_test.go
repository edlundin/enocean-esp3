package enums

import (
	"testing"
)

// TestParsePacketTypeFromByte verifies ParsePacketTypeFromByte behavior.
func TestParsePacketTypeFromByte(t *testing.T) {
	t.Run("parses all valid packet types correctly", func(t *testing.T) {
		testCases := []struct {
			input    byte
			expected PacketType
		}{
			{0x01, PacketTypeRADIO_ERP1},
			{0x02, PacketTypeRESPONSE},
			{0x03, PacketTypeRADIO_SUB_TEL},
			{0x04, PacketTypeEVENT},
			{0x05, PacketTypeCOMMON_COMMAND},
			{0x06, PacketTypeSMART_ACK_COMMAND},
			{0x07, PacketTypeREMOTE_MAN_COMMAND},
			{0x09, PacketTypeRADIO_MESSAGE},
			{0x0a, PacketTypeRADIO_ERP2},
			{0x0b, PacketTypeCONFIG_COMMAND},
			{0x0c, PacketTypeCOMMAND_ACCEPTED},
			{0x10, PacketTypeRADIO_802_15_4},
			{0x11, PacketTypeCOMMAND_2_4},
		}

		for _, tc := range testCases {
			t.Run(tc.expected.String(), func(t *testing.T) {
				result, err := ParsePacketTypeFromByte(tc.input)
				if err != nil {
					t.Errorf("expected no error for input 0x%02x, got: %s", tc.input, err)
				}
				if result != tc.expected {
					t.Errorf("expected %v, got %v", tc.expected, result)
				}
			})
		}
	})

	t.Run("returns error for invalid packet type", func(t *testing.T) {
		invalidInputs := []byte{0x00, 0x08, 0x0d, 0x0e, 0x0f, 0x12, 0x13, 0xff}

		for _, input := range invalidInputs {
			t.Run(t.Name(), func(t *testing.T) {
				result, err := ParsePacketTypeFromByte(input)
				if err == nil {
					t.Errorf("expected error for input 0x%02x, got nil", input)
				}
				if err.Error() != "invalid packet type" {
					t.Errorf("expected error 'invalid packet type', got '%s'", err.Error())
				}
				if result != 0 {
					t.Errorf("expected result 0, got %v", result)
				}
			})
		}
	})
}

// TestPacketType_String verifies PacketType_String behavior.
func TestPacketType_String(t *testing.T) {
	t.Run("returns correct string for all valid packet types", func(t *testing.T) {
		testCases := []struct {
			input    PacketType
			expected string
		}{
			{PacketTypeRADIO_ERP1, "RADIO_ERP1"},
			{PacketTypeRESPONSE, "RESPONSE"},
			{PacketTypeRADIO_SUB_TEL, "RADIO_SUB_TEL"},
			{PacketTypeEVENT, "EVENT"},
			{PacketTypeCOMMON_COMMAND, "COMMON_COMMAND"},
			{PacketTypeSMART_ACK_COMMAND, "SMART_ACK_COMMAND"},
			{PacketTypeREMOTE_MAN_COMMAND, "REMOTE_MAN_COMMAND"},
			{PacketTypeRADIO_MESSAGE, "RADIO_MESSAGE"},
			{PacketTypeRADIO_ERP2, "RADIO_ERP2"},
			{PacketTypeCONFIG_COMMAND, "CONFIG_COMMAND"},
			{PacketTypeCOMMAND_ACCEPTED, "COMMAND_ACCEPTED"},
			{PacketTypeRADIO_802_15_4, "RADIO_802_15_4"},
			{PacketTypeCOMMAND_2_4, "COMMAND_2_4"},
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

	t.Run("returns UNKNOWN for invalid packet types", func(t *testing.T) {
		invalidTypes := []PacketType{0x00, 0x08, 0x0d, 0x0e, 0x0f, 0x12, 0x13, 0xff}

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

// TestPacketTypeValid verifies PacketTypeValid behavior.
func TestPacketTypeValid(t *testing.T) {
	t.Run("PacketType_Valid", func(t *testing.T) {
		// Test valid packet types
		validTypes := []PacketType{
			PacketTypeRADIO_ERP1,
			PacketTypeRESPONSE,
			PacketTypeRADIO_SUB_TEL,
			PacketTypeEVENT,
			PacketTypeCOMMON_COMMAND,
		}
		for _, pt := range validTypes {
			if !pt.Valid() {
				t.Errorf("PacketType %v should be valid", pt)
			}
		}

		// Test invalid packet types
		invalidTypes := []PacketType{0x08, 0x0D, 0x0E, 0x0F, 0x12}
		for _, pt := range invalidTypes {
			if pt.Valid() {
				t.Errorf("PacketType %v should not be valid", pt)
			}
		}
	})
}

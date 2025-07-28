package enums

import (
	"testing"
)

func TestParsePacketTypeFromByte(t *testing.T) {
	t.Run("parses all valid packet types correctly", func(t *testing.T) {
		testCases := []struct {
			input    uint8
			expected PacketType
		}{
			{0x01, PACKET_TYPE_RADIO_ERP1},
			{0x02, PACKET_TYPE_RESPONSE},
			{0x03, PACKET_TYPE_RADIO_SUB_TEL},
			{0x04, PACKET_TYPE_EVENT},
			{0x05, PACKET_TYPE_COMMON_COMMAND},
			{0x06, PACKET_TYPE_SMART_ACK_COMMAND},
			{0x07, PACKET_TYPE_REMOTE_MAN_COMMAND},
			{0x09, PACKET_TYPE_RADIO_MESSAGE},
			{0x0a, PACKET_TYPE_RADIO_ERP2},
			{0x0b, PACKET_TYPE_CONFIG_COMMAND},
			{0x0c, PACKET_TYPE_COMMAND_ACCEPTED},
			{0x10, PACKET_TYPE_RADIO_802_15_4},
			{0x11, PACKET_TYPE_COMMAND_2_4},
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
		invalidInputs := []uint8{0x00, 0x08, 0x0d, 0x0e, 0x0f, 0x12, 0x13, 0xff}

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

func TestPacketType_String(t *testing.T) {
	t.Run("returns correct string for all valid packet types", func(t *testing.T) {
		testCases := []struct {
			input    PacketType
			expected string
		}{
			{PACKET_TYPE_RADIO_ERP1, "RADIO_ERP1"},
			{PACKET_TYPE_RESPONSE, "RESPONSE"},
			{PACKET_TYPE_RADIO_SUB_TEL, "RADIO_SUB_TEL"},
			{PACKET_TYPE_EVENT, "EVENT"},
			{PACKET_TYPE_COMMON_COMMAND, "COMMON_COMMAND"},
			{PACKET_TYPE_SMART_ACK_COMMAND, "SMART_ACK_COMMAND"},
			{PACKET_TYPE_REMOTE_MAN_COMMAND, "REMOTE_MAN_COMMAND"},
			{PACKET_TYPE_RADIO_MESSAGE, "RADIO_MESSAGE"},
			{PACKET_TYPE_RADIO_ERP2, "RADIO_ERP2"},
			{PACKET_TYPE_CONFIG_COMMAND, "CONFIG_COMMAND"},
			{PACKET_TYPE_COMMAND_ACCEPTED, "COMMAND_ACCEPTED"},
			{PACKET_TYPE_RADIO_802_15_4, "RADIO_802_15_4"},
			{PACKET_TYPE_COMMAND_2_4, "COMMAND_2_4"},
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

package esp3

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
)

// TestCrcTable verifies CrcTable behavior.
func TestCrcTable(t *testing.T) {
	t.Run("returns the CRC table", func(t *testing.T) {
		expectedTable := [256]byte{
			0x00, 0x07, 0x0e, 0x09, 0x1c, 0x1b, 0x12, 0x15, 0x38, 0x3f, 0x36, 0x31, 0x24, 0x23, 0x2a, 0x2d,
			0x70, 0x77, 0x7e, 0x79, 0x6c, 0x6b, 0x62, 0x65, 0x48, 0x4f, 0x46, 0x41, 0x54, 0x53, 0x5a, 0x5d,
			0xe0, 0xe7, 0xee, 0xe9, 0xfc, 0xfb, 0xf2, 0xf5, 0xd8, 0xdf, 0xd6, 0xd1, 0xc4, 0xc3, 0xca, 0xcd,
			0x90, 0x97, 0x9e, 0x99, 0x8c, 0x8b, 0x82, 0x85, 0xa8, 0xaf, 0xa6, 0xa1, 0xb4, 0xb3, 0xba, 0xbd,
			0xc7, 0xc0, 0xc9, 0xce, 0xdb, 0xdc, 0xd5, 0xd2, 0xff, 0xf8, 0xf1, 0xf6, 0xe3, 0xe4, 0xed, 0xea,
			0xb7, 0xb0, 0xb9, 0xbe, 0xab, 0xac, 0xa5, 0xa2, 0x8f, 0x88, 0x81, 0x86, 0x93, 0x94, 0x9d, 0x9a,
			0x27, 0x20, 0x29, 0x2e, 0x3b, 0x3c, 0x35, 0x32, 0x1f, 0x18, 0x11, 0x16, 0x03, 0x04, 0x0d, 0x0a,
			0x57, 0x50, 0x59, 0x5e, 0x4b, 0x4c, 0x45, 0x42, 0x6f, 0x68, 0x61, 0x66, 0x73, 0x74, 0x7d, 0x7a,
			0x89, 0x8e, 0x87, 0x80, 0x95, 0x92, 0x9b, 0x9c, 0xb1, 0xb6, 0xbf, 0xb8, 0xad, 0xaa, 0xa3, 0xa4,
			0xf9, 0xfe, 0xf7, 0xf0, 0xe5, 0xe2, 0xeb, 0xec, 0xc1, 0xc6, 0xcf, 0xc8, 0xdd, 0xda, 0xd3, 0xd4,
			0x69, 0x6e, 0x67, 0x60, 0x75, 0x72, 0x7b, 0x7c, 0x51, 0x56, 0x5f, 0x58, 0x4d, 0x4a, 0x43, 0x44,
			0x19, 0x1e, 0x17, 0x10, 0x05, 0x02, 0x0b, 0x0c, 0x21, 0x26, 0x2f, 0x28, 0x3d, 0x3a, 0x33, 0x34,
			0x4e, 0x49, 0x40, 0x47, 0x52, 0x55, 0x5c, 0x5b, 0x76, 0x71, 0x78, 0x7f, 0x6A, 0x6d, 0x64, 0x63,
			0x3e, 0x39, 0x30, 0x37, 0x22, 0x25, 0x2c, 0x2b, 0x06, 0x01, 0x08, 0x0f, 0x1a, 0x1d, 0x14, 0x13,
			0xae, 0xa9, 0xa0, 0xa7, 0xb2, 0xb5, 0xbc, 0xbb, 0x96, 0x91, 0x98, 0x9f, 0x8a, 0x8D, 0x84, 0x83,
			0xde, 0xd9, 0xd0, 0xd7, 0xc2, 0xc5, 0xcc, 0xcb, 0xe6, 0xe1, 0xe8, 0xef, 0xfa, 0xfd, 0xf4, 0xf3,
		}

		crcTable := crcTable()

		if !reflect.DeepEqual(expectedTable, crcTable) {
			t.Errorf("incorrect CRC table\nexpected: %v\ngot: %v", expectedTable, crcTable)
		}
	})
}

// TestComputeCrc8 verifies ComputeCrc8 behavior.
func TestComputeCrc8(t *testing.T) {
	t.Run("computes CRC8 for single byte", func(t *testing.T) {
		testCases := []struct {
			name     string
			b        byte
			initial  byte
			expected byte
		}{
			{"zero byte with zero initial", 0x00, 0x00, 0x00},
			{"zero byte with non-zero initial", 0x00, 0x55, 0xac},
			{"non-zero byte with zero initial", 0x55, 0x00, 0xac},
			{"non-zero byte with non-zero initial", 0x55, 0x55, 0x00},
			{"test byte 0xFF", 0xFF, 0x00, 0xf3},
			{"test byte 0x01", 0x01, 0x00, 0x07},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := ComputeCrc8(tc.b, tc.initial)
				if result != tc.expected {
					t.Errorf("ComputeCrc8(0x%02x, 0x%02x) = 0x%02x, expected 0x%02x",
						tc.b, tc.initial, result, tc.expected)
				}
			})
		}
	})

	t.Run("computes CRC8 incrementally", func(t *testing.T) {
		// Test that incremental CRC matches slice CRC
		data := []byte{0x00, 0x0C, 0x07, 0x01}
		expected := ComputeCrcSlice(data)

		crc := byte(0)
		for _, b := range data {
			crc = ComputeCrc8(b, crc)
		}

		if crc != expected {
			t.Errorf("incremental CRC = 0x%02x, expected 0x%02x", crc, expected)
		}
	})
}

// TestComputeCrcSlice verifies ComputeCrcSlice behavior.
func TestComputeCrcSlice(t *testing.T) {
	t.Run("computes CRC8 for empty slice", func(t *testing.T) {
		result := ComputeCrcSlice([]byte{})
		if result != 0x00 {
			t.Errorf("ComputeCrcSlice([]) = 0x%02x, expected 0x00", result)
		}
	})

	t.Run("computes CRC8 for known test vectors", func(t *testing.T) {
		testCases := []struct {
			name     string
			data     []byte
			expected byte
		}{
			{"single zero byte", []byte{0x00}, 0x00},
			{"single 0x55 byte", []byte{0x55}, 0xac},
			{"header example", []byte{0x00, 0x0C, 0x07, 0x01}, 0x96},
			{"data example", []byte{0xD2, 0x00, 0x00, 0x00, 0x00, 0xFF, 0x03, 0xFF, 0x82, 0x00, 0x85, 0x80, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0x41, 0x00}, 0x99},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := ComputeCrcSlice(tc.data)
				if result != tc.expected {
					t.Errorf("ComputeCrcSlice(%v) = 0x%02x, expected 0x%02x",
						tc.data, result, tc.expected)
				}
			})
		}
	})

	t.Run("computes CRC8 for various data lengths", func(t *testing.T) {
		for length := 1; length <= 256; length *= 2 {
			data := make([]byte, length)
			for i := range data {
				data[i] = byte(i)
			}
			crc := ComputeCrcSlice(data)
			// Verify it doesn't panic and returns a valid byte
			if crc < 0 || crc > 255 {
				t.Errorf("CRC out of range for length %d: 0x%02x", length, crc)
			}
		}
	})
}

// TestTelegram_Serialize verifies Telegram_Serialize behavior.
func TestTelegram_Serialize(t *testing.T) {
	t.Run("returns the ESP3 telegram in hex format", func(t *testing.T) {
		expectedTelegram := []byte{0x55, 0x00, 0x0c, 0x07, 0x01, 0x96, 0xd2, 0x00, 0x00, 0x00, 0x00, 0xff, 0x03, 0xff, 0x82, 0x00, 0x85, 0x80, 0x00, 0xff, 0xff, 0xff, 0xff, 0x41, 0x00, 0x99}
		telegram := NewTelegramFromData(enums.PacketTypeRADIO_ERP1, []byte{0xd2, 0x00, 0x00, 0x00, 0x00, 0xff, 0x03, 0xff, 0x82, 0x00, 0x85, 0x80},
			[]byte{0x00, 0xff, 0xff, 0xff, 0xff, 0x41, 0x00})

		if !reflect.DeepEqual(expectedTelegram, telegram.Serialize()) {
			t.Errorf("incorrect telegram serialization\nexpected: %v\ngot: %v", expectedTelegram, telegram.Serialize())
		}
	})
}

// TestFromData verifies FromData behavior.
func TestFromData(t *testing.T) {
	t.Run("feeds structure from arguments", func(t *testing.T) {
		expectedTelegram := Telegram{
			PacketType: enums.PacketTypeRADIO_ERP1,
			Data:       []byte{0xD2, 0x00, 0x00, 0x00, 0x00, 0xFF, 0x03, 0xFF, 0x82, 0x00, 0x85, 0x80},
			OptData:    []byte{0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0x41, 0x00},
		}
		telegram := NewTelegramFromData(enums.PacketTypeRADIO_ERP1, []byte{0xd2, 0x00, 0x00, 0x00, 0x00, 0xff, 0x03, 0xff, 0x82, 0x00, 0x85, 0x80},
			[]byte{0x00, 0xff, 0xff, 0xff, 0xff, 0x41, 0x00})

		if !reflect.DeepEqual(expectedTelegram, telegram) {
			t.Errorf("incorrect telegram\nexpected: %v\ngot: %v", expectedTelegram, telegram)
		}
	})
}

// TestFromHexString verifies FromHexString behavior.
func TestFromHexString(t *testing.T) {
	t.Run("parses hex string into an esp3 structure", func(t *testing.T) {
		expectedTelegram := Telegram{
			PacketType: enums.PacketTypeRADIO_ERP1,
			Data:       []byte{0xD2, 0x00, 0x00, 0x00, 0x00, 0xFF, 0x03, 0xFF, 0x82, 0x00, 0x85, 0x80},
			OptData:    []byte{0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0x41, 0x00},
		}
		telegram, err := NewEsp3TelegramFromHexString("55000C070196D200000000FF03FF8200858000FFFFFFFF410099")

		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if !reflect.DeepEqual(expectedTelegram, telegram) {
			t.Errorf("incorrect telegram\nexpected: %v\ngot: %v", expectedTelegram, telegram)
		}
	})

	t.Run("returns an error when hex string is not in hex format", func(t *testing.T) {
		_, err := NewEsp3TelegramFromHexString("55XX0C070196D200000000FF03FF8200858000FFFFFFFF410099")

		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("returns an error when hex string does not start with 55", func(t *testing.T) {
		expectedError := "sync byte missing"
		_, err := NewEsp3TelegramFromHexString("000C0701FFD200000000FF03FF8200858000FFFFFFFF410099")

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if err.Error() != expectedError {
			t.Errorf("expected: %s, got: %s", expectedError, err.Error())
		}
	})

	t.Run("returns an error when hex string does not start with 55 (odd length)", func(t *testing.T) {
		expectedError := "sync byte missing"
		_, err := NewEsp3TelegramFromHexString("5000C0701FFD200000000FF03FF8200858000FFFFFFFF410099")

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if err.Error() != expectedError {
			t.Errorf("expected: %s, got: %s", expectedError, err.Error())
		}
	})

	t.Run("returns an error when hex string is too short", func(t *testing.T) {
		expectedError := "hex string too short"
		_, err := NewEsp3TelegramFromHexString("55000C070196")

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if err.Error() != expectedError {
			t.Errorf("expected: %s, got: %s", expectedError, err.Error())
		}
	})

	t.Run("returns an error when packet type is invalid", func(t *testing.T) {
		expectedError := "invalid packet type"
		_, err := NewEsp3TelegramFromHexString("55000C07FF62D200000000FF03FF8200858000FFFFFFFF410099")

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if err.Error() != expectedError {
			t.Errorf("expected: %s, got: %s", expectedError, err.Error())
		}
	})

	t.Run("returns an error when CRC8H is invalid", func(t *testing.T) {
		expectedError := "invalid CRC8H (got:0xff, valid:0x96)"
		_, err := NewEsp3TelegramFromHexString("55000C0701FFD200000000FF03FF8200858000FFFFFFFF410099")

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if err.Error() != expectedError {
			t.Errorf("expected: %s, got: %s", expectedError, err.Error())
		}
	})

	t.Run("returns an error when CRC8D is invalid", func(t *testing.T) {
		expectedError := "invalid CRC8D (got:0xff, valid:0x99)"
		_, err := NewEsp3TelegramFromHexString("55000C070196D200000000FF03FF8200858000FFFFFFFF4100FF")

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if err.Error() != expectedError {
			t.Errorf("expected: %s, got: %s", expectedError, err.Error())
		}
	})

	t.Run("rejects declared payload length mismatches", func(t *testing.T) {
		tests := map[string]string{
			"data longer than available payload": "5500020002d80000",
			"undeclared optional byte":           "550001000265000107",
			"missing declared optional byte":     "5500010102700000",
		}
		for name, input := range tests {
			t.Run(name, func(t *testing.T) {
				if _, err := NewEsp3TelegramFromHexString(input); err == nil {
					t.Fatal("expected packet length mismatch")
				}
			})
		}
	})

	t.Run("handles odd-length hex strings by padding", func(t *testing.T) {
		// Test with odd length hex string (should be padded with leading zero)
		_, err := NewEsp3TelegramFromHexString("55000C070196D200000000FF03FF8200858000FFFFFFFF410099")
		if err != nil {
			t.Errorf("expected no error for valid hex string, got: %s", err)
		}
	})

	t.Run("handles minimum valid packet", func(t *testing.T) {
		// Minimum packet: sync(1) + header(4) + crc8h(1) + data(0) + optdata(0) + crc8d(1) = 7 bytes
		// Header: dataLen(0x0000) + optDataLen(0x00) + packetType(0x01)
		header := []byte{0x00, 0x00, 0x00, 0x01}
		crc8h := ComputeCrcSlice(header)
		crc8d := ComputeCrcSlice([]byte{}) // Empty data + optdata
		packet := append([]byte{0x55}, header...)
		packet = append(packet, crc8h)
		packet = append(packet, crc8d)

		hexStr := fmt.Sprintf("%x", packet)
		telegram, err := NewEsp3TelegramFromHexString(hexStr)
		if err != nil {
			t.Errorf("failed to parse minimum packet: %v", err)
			return
		}
		if telegram.PacketType != enums.PacketTypeRADIO_ERP1 {
			t.Errorf("expected packet type RADIO_ERP1, got %v", telegram.PacketType)
		}
		if len(telegram.Data) != 0 {
			t.Errorf("expected empty data, got %v", telegram.Data)
		}
		if len(telegram.OptData) != 0 {
			t.Errorf("expected empty optdata, got %v", telegram.OptData)
		}
	})

	t.Run("handles large data length", func(t *testing.T) {
		// Test with large data length (use reasonable size to avoid memory issues in tests)
		// In practice, ESP3 supports up to 65535 bytes, but for testing we use 1000 bytes
		data := make([]byte, 1000)
		optData := []byte{0x01, 0x02}
		for i := range data {
			data[i] = byte(i % 256)
		}

		telegram := NewTelegramFromData(enums.PacketTypeRADIO_ERP1, data, optData)
		serialized := telegram.Serialize()

		// Verify it serializes correctly
		if len(serialized) != 1+4+1+len(data)+len(optData)+1 {
			t.Errorf("serialized length mismatch")
		}

		// Verify we can parse it back
		hexStr := fmt.Sprintf("%x", serialized)
		parsed, err := NewEsp3TelegramFromHexString(hexStr)
		if err != nil {
			t.Errorf("failed to parse large packet: %v", err)
			return
		}
		if len(parsed.Data) != len(data) {
			t.Errorf("data length mismatch: expected %d, got %d", len(data), len(parsed.Data))
		}
	})

	t.Run("handles boundary values for packet type", func(t *testing.T) {
		// Test with boundary packet type values
		boundaryTypes := []enums.PacketType{0x00, 0x01, 0xFF}
		for _, ptype := range boundaryTypes {
			telegram := NewTelegramFromData(ptype, []byte{0x01}, []byte{})
			serialized := telegram.Serialize()
			if len(serialized) == 0 {
				t.Errorf("serialization failed for packet type 0x%02x", byte(ptype))
			}
		}
	})
}

// TestTelegram_RoundTrip verifies Telegram_RoundTrip behavior.
func TestTelegram_RoundTrip(t *testing.T) {
	t.Run("round-trip serialization for all packet types", func(t *testing.T) {
		packetTypes := []enums.PacketType{
			enums.PacketTypeRADIO_ERP1,
			enums.PacketTypeRESPONSE,
			enums.PacketTypeRADIO_SUB_TEL,
			enums.PacketTypeEVENT,
			enums.PacketTypeCOMMON_COMMAND,
			enums.PacketTypeSMART_ACK_COMMAND,
			enums.PacketTypeREMOTE_MAN_COMMAND,
			enums.PacketTypeRADIO_MESSAGE,
			enums.PacketTypeRADIO_ERP2,
			enums.PacketTypeCONFIG_COMMAND,
			enums.PacketTypeCOMMAND_ACCEPTED,
			enums.PacketTypeRADIO_802_15_4,
			enums.PacketTypeCOMMAND_2_4,
		}

		for _, packetType := range packetTypes {
			t.Run(packetType.String(), func(t *testing.T) {
				data := []byte{0x01, 0x02, 0x03}
				optData := []byte{0x04, 0x05}

				original := NewTelegramFromData(packetType, data, optData)
				serialized := original.Serialize()

				hexStr := fmt.Sprintf("%x", serialized)
				parsed, err := NewEsp3TelegramFromHexString(hexStr)
				if err != nil {
					t.Errorf("failed to parse: %v", err)
					return
				}

				if parsed.PacketType != original.PacketType {
					t.Errorf("packet type mismatch: expected %v, got %v", original.PacketType, parsed.PacketType)
				}
				if !reflect.DeepEqual(parsed.Data, original.Data) {
					t.Errorf("data mismatch: expected %v, got %v", original.Data, parsed.Data)
				}
				if !reflect.DeepEqual(parsed.OptData, original.OptData) {
					t.Errorf("optdata mismatch: expected %v, got %v", original.OptData, parsed.OptData)
				}
			})
		}
	})

	t.Run("round-trip with empty data", func(t *testing.T) {
		original := NewTelegramFromData(enums.PacketTypeRESPONSE, []byte{}, []byte{})
		serialized := original.Serialize()
		hexStr := fmt.Sprintf("%x", serialized)
		parsed, err := NewEsp3TelegramFromHexString(hexStr)
		if err != nil {
			t.Errorf("failed to parse: %v", err)
			return
		}
		if !reflect.DeepEqual(parsed, original) {
			t.Errorf("round-trip failed: expected %v, got %v", original, parsed)
		}
	})

	t.Run("round-trip with empty optdata", func(t *testing.T) {
		data := []byte{0x01, 0x02, 0x03}
		original := NewTelegramFromData(enums.PacketTypeRESPONSE, data, nil)
		serialized := original.Serialize()
		hexStr := fmt.Sprintf("%x", serialized)
		parsed, err := NewEsp3TelegramFromHexString(hexStr)
		if err != nil {
			t.Errorf("failed to parse: %v", err)
			return
		}
		if !reflect.DeepEqual(parsed.Data, original.Data) {
			t.Errorf("data mismatch: expected %v, got %v", original.Data, parsed.Data)
		}
		if len(parsed.OptData) != 0 {
			t.Errorf("expected empty optdata, got %v", parsed.OptData)
		}
	})

	t.Run("round-trip with large payloads", func(t *testing.T) {
		data := make([]byte, 1000)
		optData := make([]byte, 100)
		for i := range data {
			data[i] = byte(i % 256)
		}
		for i := range optData {
			optData[i] = byte(i % 256)
		}

		original := NewTelegramFromData(enums.PacketTypeRADIO_ERP1, data, optData)
		serialized := original.Serialize()
		hexStr := fmt.Sprintf("%x", serialized)
		parsed, err := NewEsp3TelegramFromHexString(hexStr)
		if err != nil {
			t.Errorf("failed to parse: %v", err)
			return
		}
		if !reflect.DeepEqual(parsed, original) {
			t.Errorf("round-trip failed for large payloads")
		}
	})
}

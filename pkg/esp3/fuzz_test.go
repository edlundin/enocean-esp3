//go:build go1.18
// +build go1.18

package esp3

import (
	"encoding/hex"
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
)

// FuzzNewEsp3TelegramFromHexString fuzzes the packet parser with random hex strings
func FuzzNewEsp3TelegramFromHexString(f *testing.F) {
	// Add seed corpus with known-good packets
	seedCorpus := []string{
		"55000C070196D200000000FF03FF8200858000FFFFFFFF410099",
		"550001020100",
		"550001040100",
		"550001050100",
	}

	for _, seed := range seedCorpus {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, hexStr string) {
		// Ensure no panics occur
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic occurred: %v", r)
			}
		}()

		// Try to parse the hex string
		telegram, err := NewEsp3TelegramFromHexString(hexStr)

		// If parsing succeeds, verify the telegram is valid
		if err == nil {
			// Verify packet type is within valid range (byte is always 0-255)
			_ = telegram.PacketType

			// Verify data lengths are reasonable
			if len(telegram.Data) > 65535 {
				t.Errorf("data length too large: %d", len(telegram.Data))
			}
			if len(telegram.OptData) > 255 {
				t.Errorf("optdata length too large: %d", len(telegram.OptData))
			}

			// Try to serialize back - should not panic
			serialized := telegram.Serialize()
			if len(serialized) == 0 {
				t.Errorf("serialization returned empty data")
			}

			// Verify serialized data starts with sync byte
			if len(serialized) > 0 && serialized[0] != 0x55 {
				t.Errorf("serialized data missing sync byte")
			}
		}
		// Errors are expected for malformed input, so we don't fail on them
	})
}

// FuzzComputeCrc8 fuzzes the CRC8 computation
func FuzzComputeCrc8(f *testing.F) {
	// Add seed corpus
	f.Add(byte(0x00), byte(0x00))
	f.Add(byte(0x55), byte(0x00))
	f.Add(byte(0xFF), byte(0xFF))
	f.Add(byte(0x01), byte(0x07))

	f.Fuzz(func(t *testing.T, b byte, crc byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic occurred: %v", r)
			}
		}()

		result := ComputeCrc8(b, crc)
		// Result is always a valid byte (0-255)
		_ = result
	})
}

// FuzzComputeCrcSlice fuzzes the CRC8 slice computation
func FuzzComputeCrcSlice(f *testing.F) {
	// Add seed corpus
	f.Add([]byte{0x00, 0x01, 0x02})
	f.Add([]byte{0x55})
	f.Add([]byte{})
	f.Add([]byte{0xFF, 0xFF, 0xFF, 0xFF})

	f.Fuzz(func(t *testing.T, data []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic occurred: %v", r)
			}
		}()

		result := ComputeCrcSlice(data)
		// Result is always a valid byte (0-255)
		_ = result
	})
}

// FuzzTelegramSerialize fuzzes telegram serialization
func FuzzTelegramSerialize(f *testing.F) {
	// Add seed corpus with various packet types and data sizes
	seedTelegrams := []struct {
		packetType byte
		data       []byte
		optData    []byte
	}{
		{0x01, []byte{0x01}, []byte{}},
		{0x02, []byte{0x00}, []byte{}},
		{0x04, []byte{0x01}, []byte{}},
		{0x01, []byte{0xd2, 0x00, 0x00, 0x00, 0x00, 0xff, 0x03, 0xff, 0x82, 0x00, 0x85, 0x80}, []byte{0x00, 0xff, 0xff, 0xff, 0xff, 0x41, 0x00}},
	}

	for _, seed := range seedTelegrams {
		f.Add(seed.packetType, seed.data, seed.optData)
	}

	f.Fuzz(func(t *testing.T, packetType byte, data []byte, optData []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic occurred: %v", r)
			}
		}()

		// Limit data sizes to prevent excessive memory usage
		if len(data) > 10000 {
			return
		}
		if len(optData) > 255 {
			return
		}

		// Create telegram (packetType might be invalid, but that's OK for fuzzing)
		telegram := Telegram{
			PacketType: enums.PacketType(packetType),
			Data:       data,
			OptData:    optData,
		}

		// Serialize should not panic
		serialized := telegram.Serialize()

		// Verify basic structure
		if len(serialized) < 7 {
			t.Errorf("serialized packet too short: %d bytes", len(serialized))
		}

		// Verify sync byte
		if serialized[0] != 0x55 {
			t.Errorf("missing sync byte, got 0x%02x", serialized[0])
		}

		// Try to parse back if it looks valid
		if len(serialized) >= 7 {
			hexStr := hex.EncodeToString(serialized)
			_, err := NewEsp3TelegramFromHexString(hexStr)
			// Errors are OK for fuzzing - we're just checking for panics
			_ = err
		}
	})
}

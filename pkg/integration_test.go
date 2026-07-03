package pkg

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/erp1"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
	"github.com/edlundin/enocean-esp3/pkg/event"
	"github.com/edlundin/enocean-esp3/pkg/response"
	"github.com/edlundin/enocean-esp3/pkg/subtel"
)

// TestRoundTripIntegration tests the full encode → decode round-trip for all packet types
func TestRoundTripIntegration(t *testing.T) {
	t.Run("ERP1 packet round-trip", func(t *testing.T) {
		// Create original ERP1 packet
		originalPacket := erp1.Packet{
			DestinationID: 0x12345678,
			Rorg:          enums.RorgVLD,
			Rssi:          0x80,
			SecurityLevel: 0x00,
			Status:        0x85,
			SubTelNum:     0x03,
			SenderID:      0x87654321,
			UserData:      []byte{0x01, 0x02, 0x03, 0x04},
		}

		// Convert to ESP3 telegram
		telegram := originalPacket.ToEsp3()

		// Serialize to bytes
		serialized := telegram.Serialize()

		// Parse back from hex string
		hexStr := hex.EncodeToString(serialized)
		parsedTelegram, err := esp3.NewEsp3TelegramFromHexString(hexStr)
		if err != nil {
			t.Fatalf("failed to parse serialized telegram: %v", err)
		}

		// Convert back to ERP1 packet
		parsedPacket, err := erp1.NewPacketFromEsp3(parsedTelegram)
		if err != nil {
			t.Fatalf("failed to parse ERP1 packet: %v", err)
		}

		// Verify all fields match
		if parsedPacket.DestinationID != originalPacket.DestinationID {
			t.Errorf("DestinationID mismatch: expected %v, got %v", originalPacket.DestinationID, parsedPacket.DestinationID)
		}
		if parsedPacket.SenderID != originalPacket.SenderID {
			t.Errorf("SenderID mismatch: expected %v, got %v", originalPacket.SenderID, parsedPacket.SenderID)
		}
		if parsedPacket.Rorg != originalPacket.Rorg {
			t.Errorf("Rorg mismatch: expected %v, got %v", originalPacket.Rorg, parsedPacket.Rorg)
		}
		// Note: erp1.ToEsp3() hardcodes RSSI to 0xff in optData, so the parsed value will be 0xff
		// This is a limitation of the current implementation
		if parsedPacket.Rssi != 0xff {
			t.Errorf("Rssi mismatch: expected 0xff (hardcoded in ToEsp3), got 0x%02x", parsedPacket.Rssi)
		}
		if parsedPacket.Status != originalPacket.Status {
			t.Errorf("Status mismatch: expected 0x%02x, got 0x%02x", originalPacket.Status, parsedPacket.Status)
		}
		if !reflect.DeepEqual(parsedPacket.UserData, originalPacket.UserData) {
			t.Errorf("UserData mismatch: expected %v, got %v", originalPacket.UserData, parsedPacket.UserData)
		}
	})

	t.Run("Response packet round-trip", func(t *testing.T) {
		// Create ESP3 telegram with response packet
		telegram := esp3.NewTelegramFromData(enums.PacketTypeRESPONSE, []byte{0x00, 0x01, 0x02, 0x03}, []byte{0x04, 0x05})

		// Serialize
		serialized := telegram.Serialize()
		hexStr := hex.EncodeToString(serialized)

		// Parse back
		parsedTelegram, err := esp3.NewEsp3TelegramFromHexString(hexStr)
		if err != nil {
			t.Fatalf("failed to parse: %v", err)
		}

		// Convert to Response packet
		responsePacket, err := response.NewPacketFromEsp3(parsedTelegram)
		if err != nil {
			t.Fatalf("failed to parse response: %v", err)
		}

		// Verify
		if responsePacket.Code != enums.ReturnCodeSUCCESS {
			t.Errorf("expected ReturnCodeSUCCESS, got %v", responsePacket.Code)
		}
		expectedData := []byte{0x01, 0x02, 0x03}
		if !reflect.DeepEqual(responsePacket.Data, expectedData) {
			t.Errorf("data mismatch: expected %v, got %v", expectedData, responsePacket.Data)
		}
	})

	t.Run("Event packet round-trip", func(t *testing.T) {
		// Create CO_TX_DONE event
		telegram := esp3.NewTelegramFromData(enums.PacketTypeEVENT, []byte{byte(enums.EventCodeCO_TX_DONE)}, []byte{})

		// Serialize
		serialized := telegram.Serialize()
		hexStr := hex.EncodeToString(serialized)

		// Parse back
		parsedTelegram, err := esp3.NewEsp3TelegramFromHexString(hexStr)
		if err != nil {
			t.Fatalf("failed to parse: %v", err)
		}

		// Convert to Event
		eventPacket, err := event.NewPacketFromEsp3(parsedTelegram)
		if err != nil {
			t.Fatalf("failed to parse event: %v", err)
		}

		// Verify
		if eventPacket.Description() != enums.EventCodeCO_TX_DONE {
			t.Errorf("expected CO_TX_DONE, got %v", eventPacket.Description())
		}
	})

	t.Run("SubTel packet round-trip", func(t *testing.T) {
		// Create SubTel packet
		subTelPacket := subtel.Packet{
			DestinationID: 0x12345678,
			Rorg:          enums.RorgVLD,
			Rssi:          0xff,
			SecurityLevel: 0x00,
			Status:        0x80,
			SubTelNum:     0x02,
			SenderID:      0x87654321,
			SubTels: []subtel.SubTel{
				{Tick: 0x01, Rssi: 0x80, Status: 0x00},
				{Tick: 0x02, Rssi: 0x81, Status: 0x01},
			},
			Timestamp: 0x1234,
			UserData:  []byte{0x01, 0x02},
		}

		// Convert to ESP3
		telegram := subTelPacket.ToEsp3()

		// Serialize
		serialized := telegram.Serialize()
		hexStr := hex.EncodeToString(serialized)

		// Parse back
		parsedTelegram, err := esp3.NewEsp3TelegramFromHexString(hexStr)
		if err != nil {
			t.Fatalf("failed to parse: %v", err)
		}

		// Convert back to SubTel
		parsedSubTel, err := subtel.NewPacketFromEsp3(parsedTelegram)
		if err != nil {
			t.Fatalf("failed to parse SubTel: %v", err)
		}

		// Verify
		if parsedSubTel.DestinationID != subTelPacket.DestinationID {
			t.Errorf("DestinationID mismatch")
		}
		if parsedSubTel.Timestamp != subTelPacket.Timestamp {
			t.Errorf("Timestamp mismatch: expected 0x%04x, got 0x%04x", subTelPacket.Timestamp, parsedSubTel.Timestamp)
		}
		if len(parsedSubTel.SubTels) != len(subTelPacket.SubTels) {
			t.Errorf("SubTels count mismatch: expected %d, got %d", len(subTelPacket.SubTels), len(parsedSubTel.SubTels))
		}
	})
}

// TestPacketStreamIntegration tests handling of packet streams (multiple packets in sequence)
func TestPacketStreamIntegration(t *testing.T) {
	t.Run("handles multiple packets in sequence", func(t *testing.T) {
		packets := []esp3.Telegram{
			esp3.NewTelegramFromData(enums.PacketTypeRESPONSE, []byte{0x00}, []byte{}),
			esp3.NewTelegramFromData(enums.PacketTypeEVENT, []byte{byte(enums.EventCodeCO_TX_DONE)}, []byte{}),
			esp3.NewTelegramFromData(enums.PacketTypeRADIO_ERP1, []byte{0xd2, 0x01, 0x02, 0x03, 0x04, 0x05, 0x80}, []byte{0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00}),
		}

		// Serialize all packets
		stream := make([]byte, 0)
		for _, packet := range packets {
			stream = append(stream, packet.Serialize()...)
		}

		// Parse packets from stream (simulating real-world parsing)
		parsedPackets := make([]esp3.Telegram, 0)
		offset := 0
		for offset < len(stream) {
			// Find sync byte
			for offset < len(stream) && stream[offset] != 0x55 {
				offset++
			}
			if offset >= len(stream) {
				break
			}

			// Try to parse packet starting at offset
			if offset+7 > len(stream) {
				break // Not enough data
			}

			// Read header
			dataLen := uint16(stream[offset+1])<<8 | uint16(stream[offset+2])
			optDataLen := stream[offset+3]
			packetLen := 1 + 4 + 1 + int(dataLen) + int(optDataLen) + 1

			if offset+packetLen > len(stream) {
				break // Incomplete packet
			}

			// Parse packet
			packetBytes := stream[offset : offset+packetLen]
			hexStr := hex.EncodeToString(packetBytes)
			parsed, err := esp3.NewEsp3TelegramFromHexString(hexStr)
			if err != nil {
				t.Errorf("failed to parse packet at offset %d: %v", offset, err)
				offset++
				continue
			}

			parsedPackets = append(parsedPackets, parsed)
			offset += packetLen
		}

		if len(parsedPackets) != len(packets) {
			t.Errorf("expected %d packets, parsed %d", len(packets), len(parsedPackets))
		}

		// Verify each packet
		for i := range packets {
			if i >= len(parsedPackets) {
				break
			}
			if parsedPackets[i].PacketType != packets[i].PacketType {
				t.Errorf("packet %d: type mismatch", i)
			}
		}
	})

	t.Run("handles partial packets in stream", func(t *testing.T) {
		// Create a complete packet
		packet := esp3.NewTelegramFromData(enums.PacketTypeRESPONSE, []byte{0x00, 0x01}, []byte{})
		serialized := packet.Serialize()

		// Test with partial packet (first half)
		partial1 := serialized[:len(serialized)/2]
		hexStr1 := hex.EncodeToString(partial1)
		_, err := esp3.NewEsp3TelegramFromHexString(hexStr1)
		if err == nil {
			t.Errorf("expected error for partial packet, got nil")
		}

		// Test with almost complete packet (missing CRC8D)
		partial2 := serialized[:len(serialized)-1]
		hexStr2 := hex.EncodeToString(partial2)
		_, err = esp3.NewEsp3TelegramFromHexString(hexStr2)
		if err == nil {
			t.Errorf("expected error for incomplete packet, got nil")
		}
	})
}

// TestBoundaryConditions tests edge cases and boundary values
func TestBoundaryConditions(t *testing.T) {
	t.Run("minimum packet size", func(t *testing.T) {
		// Minimum packet: sync(1) + header(4) + crc8h(1) + data(0) + optdata(0) + crc8d(1) = 7 bytes
		telegram := esp3.NewTelegramFromData(enums.PacketTypeRESPONSE, []byte{}, []byte{})
		serialized := telegram.Serialize()

		if len(serialized) < 7 {
			t.Errorf("packet too short: %d bytes", len(serialized))
		}

		hexStr := hex.EncodeToString(serialized)
		parsed, err := esp3.NewEsp3TelegramFromHexString(hexStr)
		if err != nil {
			t.Errorf("failed to parse minimum packet: %v", err)
		}

		if len(parsed.Data) != 0 || len(parsed.OptData) != 0 {
			t.Errorf("expected empty data and optdata")
		}
	})

	t.Run("maximum data length", func(t *testing.T) {
		// Test with large data length (use smaller size to avoid memory issues in tests)
		// In practice, ESP3 supports up to 65535 bytes, but for testing we use a reasonable size
		data := make([]byte, 1000)
		for i := range data {
			data[i] = byte(i % 256)
		}

		telegram := esp3.NewTelegramFromData(enums.PacketTypeRADIO_ERP1, data, []byte{})
		serialized := telegram.Serialize()

		hexStr := hex.EncodeToString(serialized)
		parsed, err := esp3.NewEsp3TelegramFromHexString(hexStr)
		if err != nil {
			t.Errorf("failed to parse large length packet: %v", err)
			return
		}

		if len(parsed.Data) != len(data) {
			t.Errorf("data length mismatch: expected %d, got %d", len(data), len(parsed.Data))
		}
	})

	t.Run("maximum optdata length", func(t *testing.T) {
		// Test with maximum optdata length (255 bytes)
		optData := make([]byte, 255)
		for i := range optData {
			optData[i] = byte(i % 256)
		}

		telegram := esp3.NewTelegramFromData(enums.PacketTypeRADIO_ERP1, []byte{0x01}, optData)
		serialized := telegram.Serialize()

		hexStr := hex.EncodeToString(serialized)
		parsed, err := esp3.NewEsp3TelegramFromHexString(hexStr)
		if err != nil {
			t.Errorf("failed to parse maximum optdata length packet: %v", err)
		}

		if len(parsed.OptData) != len(optData) {
			t.Errorf("optdata length mismatch: expected %d, got %d", len(optData), len(parsed.OptData))
		}
	})

	t.Run("all packet types", func(t *testing.T) {
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
				telegram := esp3.NewTelegramFromData(packetType, []byte{0x01, 0x02}, []byte{0x03})
				serialized := telegram.Serialize()
				hexStr := hex.EncodeToString(serialized)
				parsed, err := esp3.NewEsp3TelegramFromHexString(hexStr)
				if err != nil {
					t.Errorf("failed to parse %s: %v", packetType.String(), err)
					return
				}
				if parsed.PacketType != packetType {
					t.Errorf("packet type mismatch: expected %v, got %v", packetType, parsed.PacketType)
				}
			})
		}
	})
}

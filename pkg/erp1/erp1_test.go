package erp1

import (
	"slices"
	"testing"

	device_id "github.com/edlundin/enocean-esp3/pkg/device-id"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
)

func TestNewErp1PacketFromEsp3(t *testing.T) {
	t.Run("successfully creates Erp1Packet from valid ESP3 telegram", func(t *testing.T) {
		telegram := esp3.Esp3Telegram{
			PacketType: enums.PACKET_TYPE_RADIO_ERP1,
			Data:       []byte{0xd2, 0x00, 0x00, 0x00, 0x00, 0xff, 0x03, 0xff, 0x82, 0x00, 0x85, 0x80},
			OptData:    []byte{0x03, 0x12, 0x34, 0x56, 0x78, 0xff, 0x00},
		}

		packet, err := NewErp1PacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		expectedDestID, _ := device_id.FromByteArray([]byte{0x12, 0x34, 0x56, 0x78})
		expectedSenderID, _ := device_id.FromByteArray([]byte{0xff, 0x82, 0x00, 0x85})

		if packet.DestinationID != expectedDestID {
			t.Errorf("expected DestinationID %v, got %v", expectedDestID, packet.DestinationID)
		}

		if packet.Rorg != enums.Rorg(0xd2) {
			t.Errorf("expected Rorg %v, got %v", enums.Rorg(0xd2), packet.Rorg)
		}

		if packet.Rssi != 0xff {
			t.Errorf("expected Rssi %v, got %v", byte(0xff), packet.Rssi)
		}

		if packet.SecurityLevel != 0x00 {
			t.Errorf("expected SecurityLevel %v, got %v", byte(0x00), packet.SecurityLevel)
		}

		if packet.Status != 0x80 {
			t.Errorf("expected Status %v, got %v", byte(0x80), packet.Status)
		}

		if packet.SubTelNum != 0x03 {
			t.Errorf("expected SubTelNum %v, got %v", byte(0x03), packet.SubTelNum)
		}

		expectedSenderID, _ = device_id.FromByteArray([]byte{0xff, 0x82, 0x00, 0x85})
		if packet.SenderID != expectedSenderID {
			t.Errorf("expected SenderID %v, got %v", expectedSenderID, packet.SenderID)
		}

		expectedUserData := []byte{0x00, 0x00, 0x00, 0x00, 0xff, 0x03}
		if !slices.Equal(packet.UserData, expectedUserData) {
			t.Errorf("expected UserData %v, got %v", expectedUserData, packet.UserData)
		}
	})

	t.Run("returns error for invalid packet type", func(t *testing.T) {
		telegram := esp3.Esp3Telegram{
			PacketType: enums.PACKET_TYPE_RESPONSE,
			Data:       []byte{0xd2, 0x00, 0x00, 0x00, 0x00, 0xff, 0x03, 0xff, 0x82, 0x00, 0x85, 0x80},
			OptData:    []byte{0x03, 0x12, 0x34, 0x56, 0x78, 0xff, 0x00},
		}

		_, err := NewErp1PacketFromEsp3(telegram)
		if err == nil {
			t.Errorf("expected error, got nil")
		}

		expectedError := "invalid packet type"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("returns error for data too short", func(t *testing.T) {
		telegram := esp3.Esp3Telegram{
			PacketType: enums.PACKET_TYPE_RADIO_ERP1,
			Data:       []byte{0xd2, 0x12, 0x34, 0x56},
			OptData:    []byte{0x03, 0x12, 0x34, 0x56, 0x78, 0xff, 0x00},
		}

		_, err := NewErp1PacketFromEsp3(telegram)
		if err == nil {
			t.Errorf("expected error, got nil")
		}

		expectedError := "data too short"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("returns error for optData too short for destination ID", func(t *testing.T) {
		telegram := esp3.Esp3Telegram{
			PacketType: enums.PACKET_TYPE_RADIO_ERP1,
			Data:       []byte{0xd2, 0x12, 0x34, 0x56, 0x78, 0x85},
			OptData:    []byte{0x03, 0x12, 0x34},
		}

		_, err := NewErp1PacketFromEsp3(telegram)
		if err == nil {
			t.Errorf("expected error, got nil")
		}

		expectedError := "optData too short for destination ID"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("handles minimum data length correctly", func(t *testing.T) {
		telegram := esp3.Esp3Telegram{
			PacketType: enums.PACKET_TYPE_RADIO_ERP1,
			Data:       []byte{0xd2, 0x12, 0x34, 0x56, 0x78, 0x85},
			OptData:    []byte{0x03, 0x12, 0x34, 0x56, 0x78, 0xff, 0x00},
		}

		packet, err := NewErp1PacketFromEsp3(telegram)

		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if len(packet.UserData) != 0 {
			t.Errorf("expected empty UserData, got %v", packet.UserData)
		}
	})

	t.Run("handles edge case with exact minimum lengths", func(t *testing.T) {
		telegram := esp3.Esp3Telegram{
			PacketType: enums.PACKET_TYPE_RADIO_ERP1,
			Data:       []byte{0xd2, 0x12, 0x34, 0x56, 0x78, 0x85},
			OptData:    []byte{0x03, 0x12, 0x34, 0x56, 0x78, 0xff, 0x00},
		}

		packet, err := NewErp1PacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		expectedDestID, _ := device_id.FromByteArray([]byte{0x12, 0x34, 0x56, 0x78})
		expectedSenderID, _ := device_id.FromByteArray([]byte{0x12, 0x34, 0x56, 0x78})

		if packet.DestinationID != expectedDestID {
			t.Errorf("expected DestinationID %v, got %v", expectedDestID, packet.DestinationID)
		}

		if packet.SenderID != expectedSenderID {
			t.Errorf("expected SenderID %v, got %v", expectedSenderID, packet.SenderID)
		}

		if packet.Rorg != enums.Rorg(0xd2) {
			t.Errorf("expected Rorg %v, got %v", enums.Rorg(0xd2), packet.Rorg)
		}

		if packet.Status != 0x85 {
			t.Errorf("expected Status %v, got %v", byte(0x85), packet.Status)
		}

		if packet.SubTelNum != 0x03 {
			t.Errorf("expected SubTelNum %v, got %v", byte(0x03), packet.SubTelNum)
		}

		if packet.Rssi != 0xff {
			t.Errorf("expected Rssi %v, got %v", byte(0xff), packet.Rssi)
		}

		if packet.SecurityLevel != 0x00 {
			t.Errorf("expected SecurityLevel %v, got %v", byte(0x00), packet.SecurityLevel)
		}
	})

	t.Run("handles case with maximum data lengths", func(t *testing.T) {
		telegram := esp3.Esp3Telegram{
			PacketType: enums.PACKET_TYPE_RADIO_ERP1,
			Data:       []byte{0xd2, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x12, 0x34, 0x56, 0x78, 0x85},
			OptData:    []byte{0x03, 0x12, 0x34, 0x56, 0x78, 0xff, 0x00},
		}

		packet, err := NewErp1PacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		expectedUserData := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}
		if !slices.Equal(packet.UserData, expectedUserData) {
			t.Errorf("expected UserData %v, got %v", expectedUserData, packet.UserData)
		}

		expectedSenderID, _ := device_id.FromByteArray([]byte{0x12, 0x34, 0x56, 0x78})
		if packet.SenderID != expectedSenderID {
			t.Errorf("expected SenderID %v, got %v", expectedSenderID, packet.SenderID)
		}

		if packet.Status != 0x85 {
			t.Errorf("expected Status %v, got %v", byte(0x85), packet.Status)
		}
	})

	t.Run("handles case with different rorg values", func(t *testing.T) {
		telegram := esp3.Esp3Telegram{
			PacketType: enums.PACKET_TYPE_RADIO_ERP1,
			Data:       []byte{0xa5, 0x12, 0x34, 0x56, 0x78, 0x85},
			OptData:    []byte{0x03, 0x12, 0x34, 0x56, 0x78, 0xff, 0x00},
		}

		packet, err := NewErp1PacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if packet.Rorg != enums.Rorg(0xa5) {
			t.Errorf("expected Rorg %v, got %v", enums.Rorg(0xa5), packet.Rorg)
		}
	})

	t.Run("handles case with different status values", func(t *testing.T) {
		telegram := esp3.Esp3Telegram{
			PacketType: enums.PACKET_TYPE_RADIO_ERP1,
			Data:       []byte{0xd2, 0x12, 0x34, 0x56, 0x78, 0xff},
			OptData:    []byte{0x03, 0x12, 0x34, 0x56, 0x78, 0xff, 0x00},
		}

		packet, err := NewErp1PacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if packet.Status != 0xff {
			t.Errorf("expected Status %v, got %v", byte(0xff), packet.Status)
		}
	})

	t.Run("handles case with different subTelNum values", func(t *testing.T) {
		telegram := esp3.Esp3Telegram{
			PacketType: enums.PACKET_TYPE_RADIO_ERP1,
			Data:       []byte{0xd2, 0x12, 0x34, 0x56, 0x78, 0x85},
			OptData:    []byte{0x05, 0x12, 0x34, 0x56, 0x78, 0xff, 0x00},
		}

		packet, err := NewErp1PacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if packet.SubTelNum != 0x05 {
			t.Errorf("expected SubTelNum %v, got %v", byte(0x05), packet.SubTelNum)
		}
	})
}

func TestErp1Packet_ToEsp3(t *testing.T) {
	t.Run("converts Erp1Packet to ESP3 telegram correctly", func(t *testing.T) {
		destID, _ := device_id.FromByteArray([]byte{0x12, 0x34, 0x56, 0x78})
		senderID, _ := device_id.FromByteArray([]byte{0xff, 0x82, 0x00, 0x85})

		packet := Erp1Packet{
			DestinationID: destID,
			Rorg:          enums.Rorg(0xd2),
			Rssi:          0x80,
			SecurityLevel: 0x00,
			Status:        0x85,
			SubTelNum:     0x03,
			SenderID:      senderID,
			UserData:      []byte{0x00, 0x00, 0x00, 0x00, 0xff, 0x03, 0xff, 0x82, 0x00},
		}

		telegram := packet.ToEsp3()

		if telegram.PacketType != enums.PACKET_TYPE_RADIO_ERP1 {
			t.Errorf("expected PacketType %v, got %v", enums.PACKET_TYPE_RADIO_ERP1, telegram.PacketType)
		}

		expectedData := []byte{0xd2, 0x00, 0x00, 0x00, 0x00, 0xff, 0x03, 0xff, 0x82, 0x00, 0xff, 0x82, 0x00, 0x85, 0x85}
		if !slices.Equal(telegram.Data, expectedData) {
			t.Errorf("expected Data %v, got %v", expectedData, telegram.Data)
		}

		expectedOptData := []byte{0x03, 0x12, 0x34, 0x56, 0x78, 0xff, 0x03}
		if !slices.Equal(telegram.OptData, expectedOptData) {
			t.Errorf("expected OptData %v, got %v", expectedOptData, telegram.OptData)
		}
	})

	t.Run("handles empty UserData correctly", func(t *testing.T) {
		destinationID, _ := device_id.FromByteArray([]byte{0x12, 0x34, 0x56, 0x78})
		senderID, _ := device_id.FromByteArray([]byte{0xff, 0x82, 0x00, 0x85})

		packet := Erp1Packet{
			DestinationID: destinationID,
			Rorg:          enums.Rorg(0xd2),
			Rssi:          0x80,
			SecurityLevel: 0x00,
			Status:        0x85,
			SubTelNum:     0x03,
			SenderID:      senderID,
			UserData:      []byte{},
		}

		telegram := packet.ToEsp3()

		expectedData := []byte{0xd2, 0xff, 0x82, 0x00, 0x85, 0x85}
		if !slices.Equal(telegram.Data, expectedData) {
			t.Errorf("expected Data %v, got %v", expectedData, telegram.Data)
		}
	})

	t.Run("handles broadcast destination ID correctly", func(t *testing.T) {
		packet := Erp1Packet{
			DestinationID: device_id.BroadcastId(),
			Rorg:          enums.Rorg(0xd2),
			Rssi:          0x80,
			SecurityLevel: 0x00,
			Status:        0x85,
			SubTelNum:     0x03,
			SenderID:      device_id.DeviceID(0x12345678),
			UserData:      []byte{0x01, 0x02, 0x03},
		}

		telegram := packet.ToEsp3()

		expectedOptData := []byte{0x03, 0xff, 0xff, 0xff, 0xff, 0xff, 0x03}
		if !slices.Equal(telegram.OptData, expectedOptData) {
			t.Errorf("expected OptData %v, got %v", expectedOptData, telegram.OptData)
		}
	})
}

func TestErp1Packet_Serialize(t *testing.T) {
	t.Run("serializes Erp1Packet to byte array", func(t *testing.T) {
		destinationID, _ := device_id.FromByteArray([]byte{0x12, 0x34, 0x56, 0x78})
		senderID, _ := device_id.FromByteArray([]byte{0xff, 0x82, 0x00, 0x85})

		packet := Erp1Packet{
			DestinationID: destinationID,
			Rorg:          enums.Rorg(0xd2),
			Rssi:          0x80,
			SecurityLevel: 0x00,
			Status:        0x85,
			SubTelNum:     0x03,
			SenderID:      senderID,
			UserData:      []byte{0x01, 0x02, 0x03},
		}

		serialized := packet.Serialize()
		if len(serialized) == 0 {
			t.Errorf("expected non-empty serialized data")
		}

		if serialized[0] != 0x55 {
			t.Errorf("expected sync byte 0x55, got 0x%02x", serialized[0])
		}

		expectedDataLen := uint16(len(packet.UserData) + 6)
		actualDataLen := uint16(serialized[1])<<8 | uint16(serialized[2])
		if actualDataLen != expectedDataLen {
			t.Errorf("expected DataLen %d, got %d", expectedDataLen, actualDataLen)
		}
	})

	t.Run("serializes packet with empty UserData", func(t *testing.T) {
		packet := Erp1Packet{
			DestinationID: device_id.DeviceID(0x12345678),
			Rorg:          enums.Rorg(0xd2),
			Rssi:          0x80,
			SecurityLevel: 0x00,
			Status:        0x85,
			SubTelNum:     0x03,
			SenderID:      device_id.DeviceID(0x87654321),
			UserData:      []byte{},
		}

		serialized := packet.Serialize()
		if len(serialized) == 0 {
			t.Errorf("expected non-empty serialized data")
		}

		const expectedDataLen = 6

		actualDataLen := uint16(serialized[1])<<8 | uint16(serialized[2])
		if actualDataLen != expectedDataLen {
			t.Errorf("expected DataLen %d, got %d", expectedDataLen, actualDataLen)
		}
	})
}

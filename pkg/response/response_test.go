package response

import (
	"reflect"
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
)

func TestNewPacketFromEsp3(t *testing.T) {
	t.Run("successfully creates Response packet from valid ESP3 telegram", func(t *testing.T) {
		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeRESPONSE,
			Data:       []byte{0x00, 0x01, 0x02, 0x03},
			OptData:    []byte{0x04, 0x05},
		}

		packet, err := NewPacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if packet.Code != enums.ReturnCodeSUCCESS {
			t.Errorf("expected ReturnCodeSUCCESS, got %v", packet.Code)
		}

		expectedData := []byte{0x01, 0x02, 0x03}
		if !reflect.DeepEqual(packet.Data, expectedData) {
			t.Errorf("expected Data %v, got %v", expectedData, packet.Data)
		}

		expectedOptData := []byte{0x04, 0x05}
		if !reflect.DeepEqual(packet.OptData, expectedOptData) {
			t.Errorf("expected OptData %v, got %v", expectedOptData, packet.OptData)
		}
	})

	t.Run("returns error for invalid packet type", func(t *testing.T) {
		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeRADIO_ERP1,
			Data:       []byte{0x00, 0x01, 0x02},
			OptData:    []byte{},
		}

		_, err := NewPacketFromEsp3(telegram)
		if err == nil {
			t.Errorf("expected error, got nil")
		}

		expectedError := "invalid packet type"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("handles all return codes", func(t *testing.T) {
		returnCodes := []struct {
			code    enums.ReturnCode
			byteVal byte
			name    string
		}{
			{enums.ReturnCodeSUCCESS, 0x00, "SUCCESS"},
			{enums.ReturnCodeERROR, 0x01, "ERROR"},
			{enums.ReturnCodeNOT_SUPPORTED, 0x02, "NOT_SUPPORTED"},
			{enums.ReturnCodeWRONG_ARGUMENT, 0x03, "WRONG_ARGUMENT"},
			{enums.ReturnCodeOPERATION_DENIED, 0x04, "OPERATION_DENIED"},
			{enums.ReturnCodeLOCK_SET, 0x05, "LOCK_SET"},
			{enums.ReturnCodeBUFFER_TO_SMALL, 0x06, "BUFFER_TO_SMALL"},
			{enums.ReturnCodeNO_FREE_BUFFER, 0x07, "NO_FREE_BUFFER"},
			{enums.ReturnCodeBASEID_OUT_OF_RANGE, 0x90, "BASEID_OUT_OF_RANGE"},
			{enums.ReturnCodeBASEID_MAX_REACHED, 0x91, "BASEID_MAX_REACHED"},
		}

		for _, tc := range returnCodes {
			t.Run(tc.name, func(t *testing.T) {
				telegram := esp3.Telegram{
					PacketType: enums.PacketTypeRESPONSE,
					Data:       []byte{tc.byteVal},
					OptData:    []byte{},
				}

				packet, err := NewPacketFromEsp3(telegram)
				if err != nil {
					t.Errorf("failed to parse return code %s: %v", tc.name, err)
					return
				}

				if packet.Code != tc.code {
					t.Errorf("expected return code %v, got %v", tc.code, packet.Code)
				}
			})
		}
	})

	t.Run("handles empty data", func(t *testing.T) {
		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeRESPONSE,
			Data:       []byte{0x00},
			OptData:    []byte{},
		}

		packet, err := NewPacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if len(packet.Data) != 0 {
			t.Errorf("expected empty data, got %v", packet.Data)
		}
	})

	t.Run("handles empty optdata", func(t *testing.T) {
		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeRESPONSE,
			Data:       []byte{0x00, 0x01, 0x02},
			OptData:    []byte{},
		}

		packet, err := NewPacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if len(packet.OptData) != 0 {
			t.Errorf("expected empty optdata, got %v", packet.OptData)
		}
	})

	t.Run("handles large data payloads", func(t *testing.T) {
		data := make([]byte, 100)
		data[0] = 0x00 // Return code
		for i := 1; i < len(data); i++ {
			data[i] = byte(i)
		}

		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeRESPONSE,
			Data:       data,
			OptData:    []byte{0x01, 0x02, 0x03},
		}

		packet, err := NewPacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if len(packet.Data) != len(data)-1 {
			t.Errorf("expected data length %d, got %d", len(data)-1, len(packet.Data))
		}
	})

	t.Run("handles invalid return code", func(t *testing.T) {
		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeRESPONSE,
			Data:       []byte{0xFF}, // Invalid return code
			OptData:    []byte{},
		}

		_, err := NewPacketFromEsp3(telegram)
		if err == nil {
			t.Errorf("expected error for invalid return code, got nil")
		}
	})

	t.Run("extracts data after return code correctly", func(t *testing.T) {
		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeRESPONSE,
			Data:       []byte{0x00, 0xAA, 0xBB, 0xCC, 0xDD},
			OptData:    []byte{0xEE, 0xFF},
		}

		packet, err := NewPacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		expectedData := []byte{0xAA, 0xBB, 0xCC, 0xDD}
		if !reflect.DeepEqual(packet.Data, expectedData) {
			t.Errorf("expected Data %v, got %v", expectedData, packet.Data)
		}

		expectedOptData := []byte{0xEE, 0xFF}
		if !reflect.DeepEqual(packet.OptData, expectedOptData) {
			t.Errorf("expected OptData %v, got %v", expectedOptData, packet.OptData)
		}
	})
}





























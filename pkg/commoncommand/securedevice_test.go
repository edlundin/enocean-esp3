package commoncommand

import (
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/deviceid"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

// TestNewWrSecureDeviceAdd verifies NewWrSecureDeviceAdd behavior.
func TestNewWrSecureDeviceAdd(t *testing.T) {
	t.Run("creates secure device add command with valid teach info", func(t *testing.T) {
		securityKey := [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
		rollingCode := [3]byte{0x01, 0x02, 0x03}
		cmd, err := NewWrSecureDeviceAdd(
			0x01, // security level format
			deviceid.DeviceID(0x12345678),
			securityKey,
			rollingCode,
			enums.SecureDeviceDirectionOUTBOUND_TABLE,
			0x01, // ptm module
			0x0F, // teach info (max valid value)
		)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_SECUREDEVICE_ADD {
			t.Errorf("expected CommandCode WR_SECUREDEVICE_ADD, got 0x%02x", cmd.CommandCode)
		}

		if cmd.SecurityLevelFormat != 0x01 {
			t.Errorf("expected SecurityLevelFormat 1, got %d", cmd.SecurityLevelFormat)
		}

		if cmd.DeviceID != deviceid.DeviceID(0x12345678) {
			t.Errorf("expected DeviceID 0x12345678, got 0x%08x", cmd.DeviceID)
		}
	})

	t.Run("returns error for teach info out of range", func(t *testing.T) {
		securityKey := [16]byte{}
		rollingCode := [3]byte{}
		_, err := NewWrSecureDeviceAdd(
			0x01,
			deviceid.DeviceID(0x12345678),
			securityKey,
			rollingCode,
			enums.SecureDeviceDirectionOUTBOUND_TABLE,
			0x01,
			0x10, // invalid: > 0x0F
		)
		if err == nil {
			t.Fatal("expected error for teach info > 0x0F, got nil")
		}

		if err.Error() != "teach info out of range: only half a byte is allowed, use NewWrSecureDeviceV2Add 1-byte teach-in info" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})
}

// TestWrSecureDeviceAdd_Serialize verifies WrSecureDeviceAdd_Serialize behavior.
func TestWrSecureDeviceAdd_Serialize(t *testing.T) {
	t.Run("serializes secure device add command", func(t *testing.T) {
		securityKey := [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
		rollingCode := [3]byte{0x01, 0x02, 0x03}
		cmd, _ := NewWrSecureDeviceAdd(
			0x01,
			deviceid.DeviceID(0x12345678),
			securityKey,
			rollingCode,
			enums.SecureDeviceDirectionOUTBOUND_TABLE,
			0x01,
			0x05,
		)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + SecurityLevelFormat(1) + DeviceID(4) + SecurityKey(16) + RollingCode(3) = 25 bytes
		if len(telegram.Data) != 25 {
			t.Errorf("expected Data length 25, got %d", len(telegram.Data))
		}

		// OptData: Direction(1) + PTMModule(1) + TeachInInfo(1) = 3 bytes
		if len(telegram.OptData) != 3 {
			t.Errorf("expected OptData length 3, got %d", len(telegram.OptData))
		}
	})
}

// TestNewWrSecureDeviceDel verifies NewWrSecureDeviceDel behavior.
func TestNewWrSecureDeviceDel(t *testing.T) {
	t.Run("creates secure device delete command", func(t *testing.T) {
		cmd, err := NewWrSecureDeviceDel(
			deviceid.DeviceID(0x12345678),
			enums.SecureDeviceDirectionOUTBOUND_TABLE,
		)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_SECUREDEVICE_DEL {
			t.Errorf("expected CommandCode WR_SECUREDEVICE_DEL, got 0x%02x", cmd.CommandCode)
		}

		if cmd.DeviceID != deviceid.DeviceID(0x12345678) {
			t.Errorf("expected DeviceID 0x12345678, got 0x%08x", cmd.DeviceID)
		}

		if cmd.Direction != enums.SecureDeviceDirectionOUTBOUND_TABLE {
			t.Errorf("expected Direction OUTBOUND_TABLE, got %v", cmd.Direction)
		}
	})
}

// TestWrSecureDeviceDel_Serialize verifies WrSecureDeviceDel_Serialize behavior.
func TestWrSecureDeviceDel_Serialize(t *testing.T) {
	t.Run("serializes secure device delete command", func(t *testing.T) {
		cmd, _ := NewWrSecureDeviceDel(
			deviceid.DeviceID(0x12345678),
			enums.SecureDeviceDirectionOUTBOUND_TABLE,
		)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + DeviceID(4) = 5 bytes
		if len(telegram.Data) != 5 {
			t.Errorf("expected Data length 5, got %d", len(telegram.Data))
		}

		// OptData: Direction(1) = 1 byte
		if len(telegram.OptData) != 1 {
			t.Errorf("expected OptData length 1, got %d", len(telegram.OptData))
		}
	})
}

// TestNewRdSecureDeviceByIndex verifies NewRdSecureDeviceByIndex behavior.
func TestNewRdSecureDeviceByIndex(t *testing.T) {
	t.Run("creates read secure device by index command", func(t *testing.T) {
		cmd, err := NewRdSecureDeviceByIndex(0x05, enums.SecureDeviceDirectionOUTBOUND_TABLE)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandRD_SECUREDEVICE_BY_INDEX {
			t.Errorf("expected CommandCode RD_SECUREDEVICE_BY_INDEX, got 0x%02x", cmd.CommandCode)
		}

		if cmd.Index != 0x05 {
			t.Errorf("expected Index 5, got %d", cmd.Index)
		}

		if cmd.Direction != enums.SecureDeviceDirectionOUTBOUND_TABLE {
			t.Errorf("expected Direction OUTBOUND_TABLE, got %v", cmd.Direction)
		}
	})

	t.Run("returns error for index out of range", func(t *testing.T) {
		_, err := NewRdSecureDeviceByIndex(0xFF, enums.SecureDeviceDirectionOUTBOUND_TABLE)
		if err == nil {
			t.Fatal("expected error for index 0xFF, got nil")
		}

		if err.Error() != "index must be between 0 and 254" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})
}

// TestRdSecureDeviceByIndex_Serialize verifies RdSecureDeviceByIndex_Serialize behavior.
func TestRdSecureDeviceByIndex_Serialize(t *testing.T) {
	t.Run("serializes read secure device by index command", func(t *testing.T) {
		cmd, _ := NewRdSecureDeviceByIndex(0x05, enums.SecureDeviceDirectionOUTBOUND_TABLE)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + Index(1) = 2 bytes
		if len(telegram.Data) != 2 {
			t.Errorf("expected Data length 2, got %d", len(telegram.Data))
		}

		// OptData: Direction(1) = 1 byte
		if len(telegram.OptData) != 1 {
			t.Errorf("expected OptData length 1, got %d", len(telegram.OptData))
		}
	})
}

// TestParseRdSecureDeviceByIndexResponseOK verifies ParseRdSecureDeviceByIndexResponseOK behavior.
func TestParseRdSecureDeviceByIndexResponseOK(t *testing.T) {
	t.Run("parses secure device by index response", func(t *testing.T) {
		// Response from Data: SecurityLevelFormat(1) + DeviceID(4) + PrivateKey(16) = 21 bytes
		// Response from OptData: RollingCode(3) + PSK(16) + TeachInInfo(1) = 20 bytes
		resp := response.Packet{
			Code: enums.ReturnCodeSUCCESS,
			Data: []byte{
				0x01,                   // SecurityLevelFormat
				0x12, 0x34, 0x56, 0x78, // DeviceID
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, // PrivateKey
			},
			OptData: []byte{
				0x01, 0x02, 0x03, // RollingCode
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20, // PSK
				0x05, // TeachInInfo
			},
		}

		result, err := ParseRdSecureDeviceByIndexResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.SecurityLevelFormat != 0x01 {
			t.Errorf("expected SecurityLevelFormat 1, got %d", result.SecurityLevelFormat)
		}

		if result.DeviceID != deviceid.DeviceID(0x12345678) {
			t.Errorf("expected DeviceID 0x12345678, got 0x%08x", result.DeviceID)
		}

		expectedPrivateKey := [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
		if result.PrivateKey != expectedPrivateKey {
			t.Errorf("expected PrivateKey %v, got %v", expectedPrivateKey, result.PrivateKey)
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x01, 0x12, 0x34, 0x56, 0x78},
			OptData: nil,
		}

		_, err := ParseRdSecureDeviceByIndexResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for non-success return code, got nil")
		}

		if err.Error() != "invalid return code" {
			t.Errorf("expected error 'invalid return code', got '%s'", err.Error())
		}
	})

	t.Run("returns error for insufficient data", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{},
			OptData: []byte{},
		}

		_, err := ParseRdSecureDeviceByIndexResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}

		if err.Error() != "failed to deserialize response" {
			t.Errorf("expected error 'failed to deserialize response', got '%s'", err.Error())
		}
	})
}

// TestNewRdNumSecureDevices verifies NewRdNumSecureDevices behavior.
func TestNewRdNumSecureDevices(t *testing.T) {
	t.Run("creates read number of secure devices command", func(t *testing.T) {
		cmd, err := NewRdNumSecureDevices(enums.SecureDeviceDirectionOUTBOUND_TABLE)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandRD_NUMSECUREDEVICES {
			t.Errorf("expected CommandCode RD_NUMSECUREDEVICES, got 0x%02x", cmd.CommandCode)
		}

		if cmd.Direction != enums.SecureDeviceDirectionOUTBOUND_TABLE {
			t.Errorf("expected Direction OUTBOUND_TABLE, got %v", cmd.Direction)
		}
	})
}

// TestRdNumSecureDevices_Serialize verifies RdNumSecureDevices_Serialize behavior.
func TestRdNumSecureDevices_Serialize(t *testing.T) {
	t.Run("serializes read number of secure devices command", func(t *testing.T) {
		cmd, _ := NewRdNumSecureDevices(enums.SecureDeviceDirectionOUTBOUND_TABLE)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) = 1 byte
		if len(telegram.Data) != 1 {
			t.Errorf("expected Data length 1, got %d", len(telegram.Data))
		}

		// OptData: Direction(1) = 1 byte
		if len(telegram.OptData) != 1 {
			t.Errorf("expected OptData length 1, got %d", len(telegram.OptData))
		}
	})
}

// TestParseRdNumSecureDevicesResponseOK verifies ParseRdNumSecureDevicesResponseOK behavior.
func TestParseRdNumSecureDevicesResponseOK(t *testing.T) {
	t.Run("parses number of secure devices response", func(t *testing.T) {
		// Response: NumSecureDevices(1) = 1 byte
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x10},
			OptData: nil,
		}

		result, err := ParseRdNumSecureDevicesResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.NumSecureDevices != 0x10 {
			t.Errorf("expected NumSecureDevices 16, got %d", result.NumSecureDevices)
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x10},
			OptData: nil,
		}

		_, err := ParseRdNumSecureDevicesResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for non-success return code, got nil")
		}

		if err.Error() != "invalid return code" {
			t.Errorf("expected error 'invalid return code', got '%s'", err.Error())
		}
	})

	t.Run("returns error for insufficient data", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{},
			OptData: nil,
		}

		_, err := ParseRdNumSecureDevicesResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}

		if err.Error() != "failed to deserialize response" {
			t.Errorf("expected error 'failed to deserialize response', got '%s'", err.Error())
		}
	})
}

// TestNewRdSecureDeviceByID verifies NewRdSecureDeviceByID behavior.
func TestNewRdSecureDeviceByID(t *testing.T) {
	t.Run("creates read secure device by ID command", func(t *testing.T) {
		cmd, err := NewRdSecureDeviceByID(
			deviceid.DeviceID(0x12345678),
			enums.SecureDeviceDirectionOUTBOUND_TABLE,
		)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandRD_SECUREDEVICE_BY_ID {
			t.Errorf("expected CommandCode RD_SECUREDEVICE_BY_ID, got 0x%02x", cmd.CommandCode)
		}

		if cmd.DeviceID != deviceid.DeviceID(0x12345678) {
			t.Errorf("expected DeviceID 0x12345678, got 0x%08x", cmd.DeviceID)
		}

		if cmd.Direction != enums.SecureDeviceDirectionOUTBOUND_TABLE {
			t.Errorf("expected Direction OUTBOUND_TABLE, got %v", cmd.Direction)
		}
	})
}

// TestRdSecureDeviceByID_Serialize verifies RdSecureDeviceByID_Serialize behavior.
func TestRdSecureDeviceByID_Serialize(t *testing.T) {
	t.Run("serializes read secure device by ID command", func(t *testing.T) {
		cmd, _ := NewRdSecureDeviceByID(
			deviceid.DeviceID(0x12345678),
			enums.SecureDeviceDirectionOUTBOUND_TABLE,
		)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + DeviceID(4) = 5 bytes
		if len(telegram.Data) != 5 {
			t.Errorf("expected Data length 5, got %d", len(telegram.Data))
		}

		// OptData: Direction(1) = 1 byte
		if len(telegram.OptData) != 1 {
			t.Errorf("expected OptData length 1, got %d", len(telegram.OptData))
		}
	})
}

// TestParseRdSecureDeviceByIDResponseOK verifies ParseRdSecureDeviceByIDResponseOK behavior.
func TestParseRdSecureDeviceByIDResponseOK(t *testing.T) {
	t.Run("parses secure device by ID response", func(t *testing.T) {
		// Response: SecurityLevelFormat(1) + Index(1) = 2 bytes
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x01, 0x05},
			OptData: nil,
		}

		result, err := ParseRdSecureDeviceByIDResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.SecurityLevelFormat != 0x01 {
			t.Errorf("expected SecurityLevelFormat 1, got %d", result.SecurityLevelFormat)
		}

		if result.Index != 0x05 {
			t.Errorf("expected Index 5, got %d", result.Index)
		}
	})

	t.Run("returns error for index out of range in response", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x01, 0xFF},
			OptData: nil,
		}

		_, err := ParseRdSecureDeviceByIDResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for index 0xFF, got nil")
		}

		if err.Error() != "index out of range" {
			t.Errorf("expected error 'index out of range', got '%s'", err.Error())
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x01, 0x05},
			OptData: nil,
		}

		_, err := ParseRdSecureDeviceByIDResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for non-success return code, got nil")
		}

		if err.Error() != "invalid return code" {
			t.Errorf("expected error 'invalid return code', got '%s'", err.Error())
		}
	})

	t.Run("returns error for insufficient data", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{},
			OptData: nil,
		}

		_, err := ParseRdSecureDeviceByIDResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}

		if err.Error() != "failed to deserialize response" {
			t.Errorf("expected error 'failed to deserialize response', got '%s'", err.Error())
		}
	})
}

// TestNewWrSecureDeviceAddPSK verifies NewWrSecureDeviceAddPSK behavior.
func TestNewWrSecureDeviceAddPSK(t *testing.T) {
	t.Run("creates secure device add PSK command", func(t *testing.T) {
		psk := [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
		cmd, err := NewWrSecureDeviceAddPSK(deviceid.DeviceID(0x12345678), psk)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_SECUREDEVICE_ADD_PSK {
			t.Errorf("expected CommandCode WR_SECUREDEVICE_ADD_PSK, got 0x%02x", cmd.CommandCode)
		}

		if cmd.DeviceID != deviceid.DeviceID(0x12345678) {
			t.Errorf("expected DeviceID 0x12345678, got 0x%08x", cmd.DeviceID)
		}

		if cmd.PSK != psk {
			t.Errorf("expected PSK %v, got %v", psk, cmd.PSK)
		}
	})
}

// TestWrSecureDeviceAddPSK_Serialize verifies WrSecureDeviceAddPSK_Serialize behavior.
func TestWrSecureDeviceAddPSK_Serialize(t *testing.T) {
	t.Run("serializes secure device add PSK command", func(t *testing.T) {
		psk := [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
		cmd, _ := NewWrSecureDeviceAddPSK(deviceid.DeviceID(0x12345678), psk)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + DeviceID(4) + PSK(16) = 21 bytes
		if len(telegram.Data) != 21 {
			t.Errorf("expected Data length 21, got %d", len(telegram.Data))
		}
	})
}

// TestNewWrSecureDeviceSendTeachIn verifies NewWrSecureDeviceSendTeachIn behavior.
func TestNewWrSecureDeviceSendTeachIn(t *testing.T) {
	t.Run("creates secure device send teach-in command", func(t *testing.T) {
		cmd, err := NewWrSecureDeviceSendTeachIn(deviceid.DeviceID(0x12345678), 0x05)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_SECUREDEVICE_SENDTEACHIN {
			t.Errorf("expected CommandCode WR_SECUREDEVICE_SENDTEACHIN, got 0x%02x", cmd.CommandCode)
		}

		if cmd.DeviceID != deviceid.DeviceID(0x12345678) {
			t.Errorf("expected DeviceID 0x12345678, got 0x%08x", cmd.DeviceID)
		}

		if cmd.TeachInInfo != 0x05 {
			t.Errorf("expected TeachInInfo 5, got %d", cmd.TeachInInfo)
		}
	})
}

// TestWrSecureDeviceSendTeachIn_Serialize verifies WrSecureDeviceSendTeachIn_Serialize behavior.
func TestWrSecureDeviceSendTeachIn_Serialize(t *testing.T) {
	t.Run("serializes secure device send teach-in command", func(t *testing.T) {
		cmd, _ := NewWrSecureDeviceSendTeachIn(deviceid.DeviceID(0x12345678), 0x05)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + DeviceID(4) = 5 bytes
		if len(telegram.Data) != 5 {
			t.Errorf("expected Data length 5, got %d", len(telegram.Data))
		}

		// OptData: TeachInInfo(1) = 1 byte
		if len(telegram.OptData) != 1 {
			t.Errorf("expected OptData length 1, got %d", len(telegram.OptData))
		}
	})
}

// TestNewWrTemporaryRLCWindow verifies NewWrTemporaryRLCWindow behavior.
func TestNewWrTemporaryRLCWindow(t *testing.T) {
	t.Run("creates write temporary RLC window command", func(t *testing.T) {
		cmd, err := NewWrTemporaryRLCWindow(true, 0x00010000)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_TEMPORARY_RLC_WINDOW {
			t.Errorf("expected CommandCode WR_TEMPORARY_RLC_WINDOW, got 0x%02x", cmd.CommandCode)
		}

		if !cmd.Enable {
			t.Errorf("expected Enable = true, got false")
		}

		if cmd.RLCWindow != 0x00010000 {
			t.Errorf("expected RLCWindow 0x00010000, got 0x%08x", cmd.RLCWindow)
		}
	})

	t.Run("creates write temporary RLC window command with disable", func(t *testing.T) {
		cmd, err := NewWrTemporaryRLCWindow(false, 0x00010000)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.Enable {
			t.Errorf("expected Enable = false, got true")
		}
	})
}

// TestWrTemporaryRLCWindow_Serialize verifies WrTemporaryRLCWindow_Serialize behavior.
func TestWrTemporaryRLCWindow_Serialize(t *testing.T) {
	t.Run("serializes write temporary RLC window command", func(t *testing.T) {
		cmd, _ := NewWrTemporaryRLCWindow(true, 0x00010000)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + Enable(1) + RLCWindow(4) = 6 bytes
		if len(telegram.Data) != 6 {
			t.Errorf("expected Data length 6, got %d", len(telegram.Data))
		}
	})
}

// TestNewRdSecureDevicePSK verifies NewRdSecureDevicePSK behavior.
func TestNewRdSecureDevicePSK(t *testing.T) {
	t.Run("creates read secure device PSK command", func(t *testing.T) {
		cmd, err := NewRdSecureDevicePSK(deviceid.DeviceID(0x12345678), enums.SecureDeviceDirectionOUTBOUND_TABLE)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandRD_SECUREDEVICE_PSK {
			t.Errorf("expected CommandCode RD_SECUREDEVICE_PSK, got 0x%02x", cmd.CommandCode)
		}

		if cmd.DeviceID != deviceid.DeviceID(0x12345678) {
			t.Errorf("expected DeviceID 0x12345678, got 0x%08x", cmd.DeviceID)
		}
	})

	t.Run("creates read secure device PSK command for current device", func(t *testing.T) {
		cmd, err := NewRdSecureDevicePSK(deviceid.DeviceID(0x00000000), enums.SecureDeviceDirectionOUTBOUND_TABLE)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.DeviceID != 0x00000000 {
			t.Errorf("expected DeviceID 0x00000000, got 0x%08x", cmd.DeviceID)
		}
	})
}

// TestRdSecureDevicePSK_Serialize verifies RdSecureDevicePSK_Serialize behavior.
func TestRdSecureDevicePSK_Serialize(t *testing.T) {
	t.Run("serializes read secure device PSK command", func(t *testing.T) {
		cmd, _ := NewRdSecureDevicePSK(deviceid.DeviceID(0x12345678), enums.SecureDeviceDirectionOUTBOUND_TABLE)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + DeviceID(4) = 5 bytes
		if len(telegram.Data) != 5 {
			t.Errorf("expected Data length 5, got %d", len(telegram.Data))
		}
	})
}

// TestParseRdSecureDevicePSKResponseOK verifies ParseRdSecureDevicePSKResponseOK behavior.
func TestParseRdSecureDevicePSKResponseOK(t *testing.T) {
	t.Run("parses secure device PSK response", func(t *testing.T) {
		// Response: PSK(16) = 16 bytes
		resp := response.Packet{
			Code: enums.ReturnCodeSUCCESS,
			Data: []byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
			},
			OptData: nil,
		}

		result, err := ParseRdSecureDevicePSKResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		expectedPSK := [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
		if result.PSK != expectedPSK {
			t.Errorf("expected PSK %v, got %v", expectedPSK, result.PSK)
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code: enums.ReturnCodeERROR,
			Data: []byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
			},
			OptData: nil,
		}

		_, err := ParseRdSecureDevicePSKResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for non-success return code, got nil")
		}

		if err.Error() != "invalid return code" {
			t.Errorf("expected error 'invalid return code', got '%s'", err.Error())
		}
	})

	t.Run("returns error for insufficient data", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{},
			OptData: nil,
		}

		_, err := ParseRdSecureDevicePSKResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}

		if err.Error() != "failed to deserialize response" {
			t.Errorf("expected error 'failed to deserialize response', got '%s'", err.Error())
		}
	})
}

// TestNewWrRLCSavePeriod verifies NewWrRLCSavePeriod behavior.
func TestNewWrRLCSavePeriod(t *testing.T) {
	t.Run("creates write RLC save period command", func(t *testing.T) {
		cmd, err := NewWrRLCSavePeriod(0x10)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_RLC_SAVE_PERIOD {
			t.Errorf("expected CommandCode WR_RLC_SAVE_PERIOD, got 0x%02x", cmd.CommandCode)
		}

		if cmd.SavePeriod != 0x10 {
			t.Errorf("expected SavePeriod 16, got %d", cmd.SavePeriod)
		}
	})
}

// TestWrRLCSavePeriod_Serialize verifies WrRLCSavePeriod_Serialize behavior.
func TestWrRLCSavePeriod_Serialize(t *testing.T) {
	t.Run("serializes write RLC save period command", func(t *testing.T) {
		cmd, _ := NewWrRLCSavePeriod(0x10)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + SavePeriod(1) = 2 bytes
		if len(telegram.Data) != 2 {
			t.Errorf("expected Data length 2, got %d", len(telegram.Data))
		}
	})
}

// TestNewWrRLCLegacyMode verifies NewWrRLCLegacyMode behavior.
func TestNewWrRLCLegacyMode(t *testing.T) {
	t.Run("creates write RLC legacy mode command", func(t *testing.T) {
		cmd, err := NewWrRLCLegacyMode(enums.RLCModeSTANDARD)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_RLC_LEGACY_MODE {
			t.Errorf("expected CommandCode WR_RLC_LEGACY_MODE, got 0x%02x", cmd.CommandCode)
		}

		if cmd.RLCMode != enums.RLCModeSTANDARD {
			t.Errorf("expected RLCMode STANDARD, got %v", cmd.RLCMode)
		}
	})

	t.Run("returns error for invalid RLC mode", func(t *testing.T) {
		_, err := NewWrRLCLegacyMode(enums.RLCMode(0xFF))
		if err == nil {
			t.Fatal("expected error for invalid RLC mode, got nil")
		}

		if err.Error() != "invalid RLC mode" {
			t.Errorf("expected error 'invalid RLC mode', got '%s'", err.Error())
		}
	})
}

// TestWrRLCLegacyMode_Serialize verifies WrRLCLegacyMode_Serialize behavior.
func TestWrRLCLegacyMode_Serialize(t *testing.T) {
	t.Run("serializes write RLC legacy mode command", func(t *testing.T) {
		cmd, _ := NewWrRLCLegacyMode(enums.RLCModeSTANDARD)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + RLCMode(1) = 2 bytes
		if len(telegram.Data) != 2 {
			t.Errorf("expected Data length 2, got %d", len(telegram.Data))
		}
	})
}

// TestNewWrSecureDeviceV2Add verifies NewWrSecureDeviceV2Add behavior.
func TestNewWrSecureDeviceV2Add(t *testing.T) {
	t.Run("creates secure device V2 add command", func(t *testing.T) {
		privateKey := [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
		cmd, err := NewWrSecureDeviceV2Add(
			0x01, // security level format
			deviceid.DeviceID(0x12345678),
			privateKey,
			0x00010000, // rolling code
			enums.SecureDeviceDirectionOUTBOUND_TABLE,
		)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_SECUREDEVICEV2_ADD {
			t.Errorf("expected CommandCode WR_SECUREDEVICEV2_ADD, got 0x%02x", cmd.CommandCode)
		}

		if cmd.SecurityLevelFormat != 0x01 {
			t.Errorf("expected SecurityLevelFormat 1, got %d", cmd.SecurityLevelFormat)
		}

		if cmd.DeviceID != deviceid.DeviceID(0x12345678) {
			t.Errorf("expected DeviceID 0x12345678, got 0x%08x", cmd.DeviceID)
		}
	})
}

// TestWrSecureDeviceV2Add_Serialize verifies WrSecureDeviceV2Add_Serialize behavior.
func TestWrSecureDeviceV2Add_Serialize(t *testing.T) {
	t.Run("serializes secure device V2 add command", func(t *testing.T) {
		privateKey := [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
		cmd, _ := NewWrSecureDeviceV2Add(
			0x01,
			deviceid.DeviceID(0x12345678),
			privateKey,
			0x00010000,
			enums.SecureDeviceDirectionOUTBOUND_TABLE,
		)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + SecurityLevelFormat(1) + DeviceID(4) + PrivateKey(16) + RollingCode(4) = 26 bytes
		if len(telegram.Data) != 26 {
			t.Errorf("expected Data length 26, got %d", len(telegram.Data))
		}

		// OptData: Direction(1) = 1 byte
		if len(telegram.OptData) != 1 {
			t.Errorf("expected OptData length 1, got %d", len(telegram.OptData))
		}
	})

	t.Run("returns error for invalid direction", func(t *testing.T) {
		privateKey := [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
		cmd, _ := NewWrSecureDeviceV2Add(
			0x01,
			deviceid.DeviceID(0x12345678),
			privateKey,
			0x00010000,
			enums.SecureDeviceDirection(0xFF), // Invalid direction
		)
		_, err := cmd.Serialize()
		if err == nil {
			t.Fatal("expected error for invalid direction, got nil")
		}
		if err.Error() != "direction must be INBOUND_TABLE, OUTBOUND_TABLE or OUTBOUND_BROADCAST_TABLE" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})
}

// TestNewRdSecureDeviceV2ByIndex verifies NewRdSecureDeviceV2ByIndex behavior.
func TestNewRdSecureDeviceV2ByIndex(t *testing.T) {
	t.Run("creates read secure device V2 by index command", func(t *testing.T) {
		cmd, err := NewRdSecureDeviceV2ByIndex(0x05, enums.SecureDeviceDirectionOUTBOUND_TABLE)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandRD_SECUREDEVICEV2_BY_INDEX {
			t.Errorf("expected CommandCode RD_SECUREDEVICEV2_BY_INDEX, got 0x%02x", cmd.CommandCode)
		}

		if cmd.Index != 0x05 {
			t.Errorf("expected Index 5, got %d", cmd.Index)
		}

		if cmd.Direction != enums.SecureDeviceDirectionOUTBOUND_TABLE {
			t.Errorf("expected Direction OUTBOUND_TABLE, got %v", cmd.Direction)
		}
	})

	t.Run("returns error for index out of range", func(t *testing.T) {
		_, err := NewRdSecureDeviceV2ByIndex(0xFF, enums.SecureDeviceDirectionOUTBOUND_TABLE)
		if err == nil {
			t.Fatal("expected error for index 0xFF, got nil")
		}

		if err.Error() != "index must be between 0 and 254" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})
}

// TestRdSecureDeviceV2ByIndex_Serialize verifies RdSecureDeviceV2ByIndex_Serialize behavior.
func TestRdSecureDeviceV2ByIndex_Serialize(t *testing.T) {
	t.Run("serializes read secure device V2 by index command", func(t *testing.T) {
		cmd, _ := NewRdSecureDeviceV2ByIndex(0x05, enums.SecureDeviceDirectionOUTBOUND_TABLE)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + Index(1) = 2 bytes
		if len(telegram.Data) != 2 {
			t.Errorf("expected Data length 2, got %d", len(telegram.Data))
		}

		// OptData: Direction(1) = 1 byte
		if len(telegram.OptData) != 1 {
			t.Errorf("expected OptData length 1, got %d", len(telegram.OptData))
		}
	})

	t.Run("returns error for invalid index", func(t *testing.T) {
		cmd := RdSecureDeviceV2ByIndex{
			CommandCode: enums.CommonCommandRD_SECUREDEVICEV2_BY_INDEX,
			Index:       0xFF, // Invalid: > 0xFE
			Direction:   enums.SecureDeviceDirectionOUTBOUND_TABLE,
		}
		_, err := cmd.Serialize()
		if err == nil {
			t.Fatal("expected error for invalid index, got nil")
		}
		if err.Error() != "index must be between 0 and 254" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})

	t.Run("returns error for invalid direction", func(t *testing.T) {
		cmd := RdSecureDeviceV2ByIndex{
			CommandCode: enums.CommonCommandRD_SECUREDEVICEV2_BY_INDEX,
			Index:       0x05,
			Direction:   enums.SecureDeviceDirection(0xFF), // Invalid direction
		}
		_, err := cmd.Serialize()
		if err == nil {
			t.Fatal("expected error for invalid direction, got nil")
		}
		if err.Error() != "direction must be INBOUND_TABLE, OUTBOUND_TABLE or OUTBOUND_BROADCAST_TABLE" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})
}

// TestParseRdSecureDeviceV2ByIndexResponseOK verifies ParseRdSecureDeviceV2ByIndexResponseOK behavior.
func TestParseRdSecureDeviceV2ByIndexResponseOK(t *testing.T) {
	t.Run("parses secure device V2 by index response", func(t *testing.T) {
		// Response from Data: SecurityLevelFormat(1) + DeviceID(4) + PrivateKey(16) = 21 bytes
		// Response from OptData: RollingCode(4) + TeachInInfo(1) + PSK(16) = 21 bytes
		resp := response.Packet{
			Code: enums.ReturnCodeSUCCESS,
			Data: []byte{
				0x01,                   // SecurityLevelFormat
				0x12, 0x34, 0x56, 0x78, // DeviceID
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, // PrivateKey
			},
			OptData: []byte{
				0x00, 0x01, 0x00, 0x00, // RollingCode
				0x05, // TeachInInfo
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20, // PSK
			},
		}

		result, err := ParseRdSecureDeviceV2ByIndexResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.SecurityLevelFormat != 0x01 {
			t.Errorf("expected SecurityLevelFormat 1, got %d", result.SecurityLevelFormat)
		}

		if result.DeviceID != deviceid.DeviceID(0x12345678) {
			t.Errorf("expected DeviceID 0x12345678, got 0x%08x", result.DeviceID)
		}

		if result.RollingCode != 0x00010000 {
			t.Errorf("expected RollingCode 0x00010000, got 0x%08x", result.RollingCode)
		}

		if result.TeachInInfo != 0x05 {
			t.Errorf("expected TeachInInfo 5, got %d", result.TeachInInfo)
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x01, 0x12, 0x34, 0x56, 0x78},
			OptData: nil,
		}

		_, err := ParseRdSecureDeviceV2ByIndexResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for non-success return code, got nil")
		}

		if err.Error() != "invalid return code" {
			t.Errorf("expected error 'invalid return code', got '%s'", err.Error())
		}
	})

	t.Run("returns error for insufficient data", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{},
			OptData: []byte{},
		}

		_, err := ParseRdSecureDeviceV2ByIndexResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}

		if err.Error() != "failed to deserialize response" {
			t.Errorf("expected error 'failed to deserialize response', got '%s'", err.Error())
		}
	})
}

// TestNewWrSecureDeviceRemainCode verifies NewWrSecureDeviceRemainCode behavior.
func TestNewWrSecureDeviceRemainCode(t *testing.T) {
	t.Run("creates secure device reman key command", func(t *testing.T) {
		remanKey := [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
		cmd, err := NewWrSecureDeviceRemainCode(deviceid.DeviceID(0x12345678), remanKey, 0x05)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_SECUREDEVICE_REMAN_KEY {
			t.Errorf("expected CommandCode WR_SECUREDEVICE_REMAN_KEY, got 0x%02x", cmd.CommandCode)
		}

		if cmd.DeviceID != deviceid.DeviceID(0x12345678) {
			t.Errorf("expected DeviceID 0x12345678, got 0x%08x", cmd.DeviceID)
		}

		if cmd.RemanKeyNumber != 0x05 {
			t.Errorf("expected RemanKeyNumber 5, got %d", cmd.RemanKeyNumber)
		}
	})
}

// TestWrSecureDeviceRemanKey_Serialize verifies WrSecureDeviceRemanKey_Serialize behavior.
func TestWrSecureDeviceRemanKey_Serialize(t *testing.T) {
	t.Run("serializes secure device reman key command", func(t *testing.T) {
		remanKey := [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
		cmd, _ := NewWrSecureDeviceRemainCode(deviceid.DeviceID(0x12345678), remanKey, 0x05)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + DeviceID(4) + RemanKey(16) + RemanKeyNumber(1) = 22 bytes
		if len(telegram.Data) != 22 {
			t.Errorf("expected Data length 22, got %d", len(telegram.Data))
		}
	})
}

// TestNewRdSecureDeviceRemanKey verifies NewRdSecureDeviceRemanKey behavior.
func TestNewRdSecureDeviceRemanKey(t *testing.T) {
	t.Run("creates read secure device reman key command", func(t *testing.T) {
		cmd, err := NewRdSecureDeviceRemanKey(0x05)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandRD_SECUREDEVICE_REMAN_KEY {
			t.Errorf("expected CommandCode RD_SECUREDEVICE_REMAN_KEY, got 0x%02x", cmd.CommandCode)
		}

		if cmd.Index != 0x05 {
			t.Errorf("expected Index 5, got %d", cmd.Index)
		}
	})
}

// TestRdSecureDeviceRemanKey_Serialize verifies RdSecureDeviceRemanKey_Serialize behavior.
func TestRdSecureDeviceRemanKey_Serialize(t *testing.T) {
	t.Run("serializes read secure device reman key command", func(t *testing.T) {
		cmd, _ := NewRdSecureDeviceRemanKey(0x05)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + Index(1) = 2 bytes
		if len(telegram.Data) != 2 {
			t.Errorf("expected Data length 2, got %d", len(telegram.Data))
		}
	})
}

// TestParseRdSecureDeviceRemanKeyResponseOK verifies ParseRdSecureDeviceRemanKeyResponseOK behavior.
func TestParseRdSecureDeviceRemanKeyResponseOK(t *testing.T) {
	t.Run("parses secure device reman key response", func(t *testing.T) {
		// Response: Index(1) + DeviceID(4) + PrivateKey(16) + KeyNumber(1) +
		//           InboundRollingCode(4) + OutboundRollingCode(4) = 30 bytes
		resp := response.Packet{
			Code: enums.ReturnCodeSUCCESS,
			Data: []byte{
				0x05,                   // Index
				0x12, 0x34, 0x56, 0x78, // DeviceID
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, // PrivateKey
				0x05,                   // KeyNumber
				0x00, 0x01, 0x00, 0x00, // InboundRollingCode
				0x00, 0x02, 0x00, 0x00, // OutboundRollingCode
			},
			OptData: nil,
		}

		result, err := ParseRdSecureDeviceRemanKeyResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.Index != 0x05 {
			t.Errorf("expected Index 5, got %d", result.Index)
		}

		if result.DeviceID != deviceid.DeviceID(0x12345678) {
			t.Errorf("expected DeviceID 0x12345678, got 0x%08x", result.DeviceID)
		}

		if result.KeyNumber != 0x05 {
			t.Errorf("expected KeyNumber 5, got %d", result.KeyNumber)
		}

		if result.InboundRollingCode != 0x00010000 {
			t.Errorf("expected InboundRollingCode 0x00010000, got 0x%08x", result.InboundRollingCode)
		}

		if result.OutboundRollingCode != 0x00020000 {
			t.Errorf("expected OutboundRollingCode 0x00020000, got 0x%08x", result.OutboundRollingCode)
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x05, 0x12, 0x34, 0x56, 0x78},
			OptData: nil,
		}

		_, err := ParseRdSecureDeviceRemanKeyResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for non-success return code, got nil")
		}

		if err.Error() != "invalid return code" {
			t.Errorf("expected error 'invalid return code', got '%s'", err.Error())
		}
	})

	t.Run("returns error for insufficient data", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{},
			OptData: nil,
		}

		_, err := ParseRdSecureDeviceRemanKeyResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}

		if err.Error() != "failed to deserialize response" {
			t.Errorf("expected error 'failed to deserialize response', got '%s'", err.Error())
		}
	})
}

// TestWrSecureDeviceRemanKey_Serialize_InvalidKeyNumber verifies WrSecureDeviceRemanKey_Serialize_InvalidKeyNumber behavior.
func TestWrSecureDeviceRemanKey_Serialize_InvalidKeyNumber(t *testing.T) {
	t.Run("returns error for key number 0", func(t *testing.T) {
		cmd := WrSecureDeviceRemanKey{
			CommandCode:    enums.CommonCommandWR_SECUREDEVICE_REMAN_KEY,
			DeviceID:       0x12345678,
			RemanKey:       [16]byte{},
			RemanKeyNumber: 0x00,
		}
		_, err := cmd.Serialize()
		if err == nil {
			t.Fatal("expected error for key number 0, got nil")
		}
		if err.Error() != "reman key number must be between 1 and 15" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})

	t.Run("returns error for key number above 15", func(t *testing.T) {
		cmd := WrSecureDeviceRemanKey{
			CommandCode:    enums.CommonCommandWR_SECUREDEVICE_REMAN_KEY,
			DeviceID:       0x12345678,
			RemanKey:       [16]byte{},
			RemanKeyNumber: 0x10,
		}
		_, err := cmd.Serialize()
		if err == nil {
			t.Fatal("expected error for key number 0x10, got nil")
		}
		if err.Error() != "reman key number must be between 1 and 15" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})
}

// TestRdSecureDeviceRemanKey_Serialize_InvalidIndex verifies RdSecureDeviceRemanKey_Serialize_InvalidIndex behavior.
func TestRdSecureDeviceRemanKey_Serialize_InvalidIndex(t *testing.T) {
	t.Run("returns error for index 0", func(t *testing.T) {
		cmd := RdSecureDeviceRemanKey{
			CommandCode: enums.CommonCommandRD_SECUREDEVICE_REMAN_KEY,
			Index:       0x00,
		}
		_, err := cmd.Serialize()
		if err == nil {
			t.Fatal("expected error for index 0, got nil")
		}
		if err.Error() != "index must be between 1 and 15" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})

	t.Run("returns error for index above 15", func(t *testing.T) {
		cmd := RdSecureDeviceRemanKey{
			CommandCode: enums.CommonCommandRD_SECUREDEVICE_REMAN_KEY,
			Index:       0x10,
		}
		_, err := cmd.Serialize()
		if err == nil {
			t.Fatal("expected error for index 0x10, got nil")
		}
		if err.Error() != "index must be between 1 and 15" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})
}

package commoncommand

import (
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/deviceid"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

// TestNewWrIDBase verifies NewWrIDBase behavior.
func TestNewWrIDBase(t *testing.T) {
	t.Run("creates write ID base command with valid base ID", func(t *testing.T) {
		// Valid base ID: 0xff800000 (min), must be aligned to 128
		validBaseID := deviceid.DeviceID(0xff800000)
		cmd, err := NewWrIDBase(validBaseID)
		if err != nil {
			t.Fatalf("expected no error for valid base ID, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_IDBASE {
			t.Errorf("expected CommandCode WR_IDBASE, got 0x%02x", cmd.CommandCode)
		}

		if cmd.IDBase != validBaseID {
			t.Errorf("expected IDBase 0x%08x, got 0x%08x", validBaseID, cmd.IDBase)
		}
	})

	t.Run("creates write ID base command with max valid base ID", func(t *testing.T) {
		// Max valid base ID: 0xffffff80, must be aligned to 128
		maxBaseID := deviceid.DeviceID(0xffffff80)
		cmd, err := NewWrIDBase(maxBaseID)
		if err != nil {
			t.Fatalf("expected no error for max valid base ID, got: %v", err)
		}

		if cmd.IDBase != maxBaseID {
			t.Errorf("expected IDBase 0x%08x, got 0x%08x", maxBaseID, cmd.IDBase)
		}
	})

	t.Run("returns error for ID below minimum", func(t *testing.T) {
		// Below minimum: 0xff7fffff
		invalidBaseID := deviceid.DeviceID(0xff7fffff)
		_, err := NewWrIDBase(invalidBaseID)
		if err == nil {
			t.Fatal("expected error for ID below minimum, got nil")
		}

		if err.Error() != "device ID out of range" {
			t.Errorf("expected error 'device ID out of range', got '%s'", err.Error())
		}
	})

	t.Run("returns error for ID above maximum", func(t *testing.T) {
		// Above maximum: 0xffffff81
		invalidBaseID := deviceid.DeviceID(0xffffff81)
		_, err := NewWrIDBase(invalidBaseID)
		if err == nil {
			t.Fatal("expected error for ID above maximum, got nil")
		}

		if err.Error() != "device ID out of range" {
			t.Errorf("expected error 'device ID out of range', got '%s'", err.Error())
		}
	})

	t.Run("returns error for non-aligned base ID", func(t *testing.T) {
		// Valid range but not aligned to 128: 0xff800001
		invalidBaseID := deviceid.DeviceID(0xff800001)
		_, err := NewWrIDBase(invalidBaseID)
		if err == nil {
			t.Fatal("expected error for non-aligned base ID, got nil")
		}

		if err.Error() != "device ID is not a base ID" {
			t.Errorf("expected error 'device ID is not a base ID', got '%s'", err.Error())
		}
	})
}

// TestWrIDBase_Serialize verifies WrIDBase_Serialize behavior.
func TestWrIDBase_Serialize(t *testing.T) {
	t.Run("serializes write ID base command", func(t *testing.T) {
		validBaseID := deviceid.DeviceID(0xff800000)
		cmd, err := NewWrIDBase(validBaseID)
		if err != nil {
			t.Fatalf("expected no constructor error, got: %v", err)
		}
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + IDBase(4) = 5 bytes
		if len(telegram.Data) != 5 {
			t.Fatalf("expected Data length 5, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandWR_IDBASE) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandWR_IDBASE, telegram.Data[0])
		}
	})
}

// TestNewRdIDBase verifies NewRdIDBase behavior.
func TestNewRdIDBase(t *testing.T) {
	t.Run("creates read ID base command", func(t *testing.T) {
		cmd, err := NewRdIDBase()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandRD_IDBASE {
			t.Errorf("expected CommandCode RD_IDBASE, got 0x%02x", cmd.CommandCode)
		}
	})
}

// TestRdIDBase_Serialize verifies RdIDBase_Serialize behavior.
func TestRdIDBase_Serialize(t *testing.T) {
	t.Run("serializes read ID base command", func(t *testing.T) {
		cmd, err := NewRdIDBase()
		if err != nil {
			t.Fatalf("expected no constructor error, got: %v", err)
		}
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(telegram.Data) != 1 {
			t.Fatalf("expected Data length 1, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandRD_IDBASE) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandRD_IDBASE, telegram.Data[0])
		}
	})
}

// TestParseRdIDBaseResponseOK verifies ParseRdIDBaseResponseOK behavior.
func TestParseRdIDBaseResponseOK(t *testing.T) {
	t.Run("parses read ID base response", func(t *testing.T) {
		// Response: BaseID(4) from Data + RemainingWriteCount(1) from OptData = 5 bytes total
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0xff, 0x80, 0x00, 0x00},
			OptData: []byte{0x0A},
		}

		result, err := ParseRdIDBaseResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		expectedBaseID := deviceid.DeviceID(0xff800000)
		if result.BaseID != expectedBaseID {
			t.Errorf("expected BaseID 0x%08x, got 0x%08x", expectedBaseID, result.BaseID)
		}

		if result.RemainingWriteCount != 0x0A {
			t.Errorf("expected RemainingWriteCount 10, got %d", result.RemainingWriteCount)
		}
	})

	t.Run("parses read ID base response with maximum remaining writes", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0xff, 0xff, 0xff, 0x80},
			OptData: []byte{0xFF},
		}

		result, err := ParseRdIDBaseResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		expectedBaseID := deviceid.DeviceID(0xffffff80)
		if result.BaseID != expectedBaseID {
			t.Errorf("expected BaseID 0x%08x, got 0x%08x", expectedBaseID, result.BaseID)
		}

		if result.RemainingWriteCount != 0xFF {
			t.Errorf("expected RemainingWriteCount 255, got %d", result.RemainingWriteCount)
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0xff, 0x80, 0x00, 0x00},
			OptData: []byte{0x0A},
		}

		_, err := ParseRdIDBaseResponseOK(resp)
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

		_, err := ParseRdIDBaseResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}

		if err.Error() != "failed to deserialize response" {
			t.Errorf("expected error 'failed to deserialize response', got '%s'", err.Error())
		}
	})
}

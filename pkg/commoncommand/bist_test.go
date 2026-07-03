package commoncommand

import (
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

func TestNewBist(t *testing.T) {
	t.Run("creates BIST command successfully", func(t *testing.T) {
		cmd, err := NewBist()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_BIST {
			t.Errorf("expected CommandCode WR_BIST (0x%02x), got 0x%02x", enums.CommonCommandWR_BIST, cmd.CommandCode)
		}
	})
}

func TestWrBist_Serialize(t *testing.T) {
	t.Run("serializes BIST command successfully", func(t *testing.T) {
		cmd, _ := NewBist()
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Verify telegram type is COMMON_COMMAND
		if telegram.PacketType != enums.PacketTypeCOMMON_COMMAND {
			t.Errorf("expected PacketType COMMON_COMMAND, got %v", telegram.PacketType)
		}

		// Verify data contains command code (1 byte: 0x28)
		if len(telegram.Data) != 1 {
			t.Errorf("expected Data length 1, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandWR_BIST) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandWR_BIST, telegram.Data[0])
		}

		// OptData should be nil for this command
		if telegram.OptData != nil {
			t.Errorf("expected nil OptData, got %v", telegram.OptData)
		}
	})
}

func TestParseWrBistResponseOK(t *testing.T) {
	t.Run("parses successful BIST response", func(t *testing.T) {
		// BIST result: true (1 byte: 0x01)
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x01},
			OptData: nil,
		}

		result, err := ParseWrBistResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if !result.BistResult {
			t.Errorf("expected BistResult = true, got false")
		}
	})

	t.Run("parses failed BIST response", func(t *testing.T) {
		// BIST result: false (1 byte: 0x00)
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x00},
			OptData: nil,
		}

		result, err := ParseWrBistResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.BistResult {
			t.Errorf("expected BistResult = false, got true")
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x01},
			OptData: nil,
		}

		_, err := ParseWrBistResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for non-success return code, got nil")
		}

		if err.Error() != "invalid return code" {
			t.Errorf("expected error 'invalid return code', got '%s'", err.Error())
		}
	})

	t.Run("returns error for empty response data", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{},
			OptData: nil,
		}

		_, err := ParseWrBistResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for empty data, got nil")
		}
	})
}

func TestWrBist_Integration(t *testing.T) {
	t.Run("full roundtrip: create, serialize, and parse response", func(t *testing.T) {
		// Create command
		cmd, err := NewBist()
		if err != nil {
			t.Fatalf("failed to create command: %v", err)
		}

		// Serialize
		telegram, err := cmd.Serialize()
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Verify telegram structure
		if telegram.PacketType != enums.PacketTypeCOMMON_COMMAND {
			t.Error("wrong packet type")
		}

		// Simulate receiving response
		resp := response.Packet{
			Code: enums.ReturnCodeSUCCESS,
			// BIST result byte (0x01 = success)
			Data:    []byte{0x01},
			OptData: nil,
		}

		// Parse response
		result, err := ParseWrBistResponseOK(resp)
		if err != nil {
			t.Fatalf("failed to parse response: %v", err)
		}

		if !result.BistResult {
			t.Error("expected BIST to pass")
		}
	})
}

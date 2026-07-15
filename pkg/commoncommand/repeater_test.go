package commoncommand

import (
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

// TestNewWrRepeater verifies NewWrRepeater behavior.
func TestNewWrRepeater(t *testing.T) {
	t.Run("creates write repeater command", func(t *testing.T) {
		cmd, err := NewWrRepeater(enums.RepeaterModeON, enums.RepeaterLevel1_REPETITION)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_REPEATER {
			t.Errorf("expected CommandCode WR_REPEATER, got 0x%02x", cmd.CommandCode)
		}

		if cmd.RepeaterMode != enums.RepeaterModeON {
			t.Errorf("expected RepeaterMode ON, got %v", cmd.RepeaterMode)
		}

		if cmd.RepeaterLevel != enums.RepeaterLevel1_REPETITION {
			t.Errorf("expected RepeaterLevel 1_REPETITION, got %v", cmd.RepeaterLevel)
		}
	})

	t.Run("creates write repeater command with different mode and level", func(t *testing.T) {
		cmd, err := NewWrRepeater(enums.RepeaterModeOFF, enums.RepeaterLevel2_REPETITION)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.RepeaterMode != enums.RepeaterModeOFF {
			t.Errorf("expected RepeaterMode OFF, got %v", cmd.RepeaterMode)
		}

		if cmd.RepeaterLevel != enums.RepeaterLevel2_REPETITION {
			t.Errorf("expected RepeaterLevel 2_REPETITION, got %v", cmd.RepeaterLevel)
		}
	})
}

// TestWrRepeater_Serialize verifies WrRepeater_Serialize behavior.
func TestWrRepeater_Serialize(t *testing.T) {
	t.Run("serializes write repeater command", func(t *testing.T) {
		cmd, _ := NewWrRepeater(enums.RepeaterModeON, enums.RepeaterLevel1_REPETITION)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + RepeaterMode(1) + RepeaterLevel(1) = 3 bytes
		if len(telegram.Data) != 3 {
			t.Errorf("expected Data length 3, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandWR_REPEATER) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandWR_REPEATER, telegram.Data[0])
		}
	})
}

// TestParseRdRepeaterResponseOK verifies ParseRdRepeaterResponseOK behavior.
func TestParseRdRepeaterResponseOK(t *testing.T) {
	t.Run("parses read repeater response", func(t *testing.T) {
		// Response: RepeaterMode(1) + RepeaterLevel(1) = 2 bytes
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x01, 0x01},
			OptData: nil,
		}

		result, err := ParseRdRepeaterResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.RepeaterMode != enums.RepeaterModeON {
			t.Errorf("expected RepeaterMode ON, got %v", result.RepeaterMode)
		}

		if result.RepeaterLevel != enums.RepeaterLevel1_REPETITION {
			t.Errorf("expected RepeaterLevel 1_REPETITION, got %v", result.RepeaterLevel)
		}
	})

	t.Run("parses read repeater response with OFF mode", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x00, 0x02},
			OptData: nil,
		}

		result, err := ParseRdRepeaterResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.RepeaterMode != enums.RepeaterModeOFF {
			t.Errorf("expected RepeaterMode OFF, got %v", result.RepeaterMode)
		}

		if result.RepeaterLevel != enums.RepeaterLevel2_REPETITION {
			t.Errorf("expected RepeaterLevel 2_REPETITION, got %v", result.RepeaterLevel)
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x01, 0x01},
			OptData: nil,
		}

		_, err := ParseRdRepeaterResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for non-success return code, got nil")
		}

		if err.Error() != "invalid return code" {
			t.Errorf("expected error 'invalid return code', got '%s'", err.Error())
		}
	})

	t.Run("returns error for invalid repeater mode", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0xFF, 0x01}, // Invalid mode
			OptData: nil,
		}

		_, err := ParseRdRepeaterResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for invalid repeater mode, got nil")
		}

		if err.Error() != "invalid repeater mode" {
			t.Errorf("expected error 'invalid repeater mode', got '%s'", err.Error())
		}
	})

	t.Run("returns error for invalid repeater level", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x01, 0xFF}, // Invalid level
			OptData: nil,
		}

		_, err := ParseRdRepeaterResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for invalid repeater level, got nil")
		}

		if err.Error() != "invalid repeater level" {
			t.Errorf("expected error 'invalid repeater level', got '%s'", err.Error())
		}
	})

	t.Run("returns error for insufficient data", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{},
			OptData: nil,
		}

		_, err := ParseRdRepeaterResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}

		if err.Error() != "failed to deserialize response" {
			t.Errorf("expected error 'failed to deserialize response', got '%s'", err.Error())
		}
	})
}

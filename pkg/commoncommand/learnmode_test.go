package commoncommand

import (
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

// TestNewWrLearnMode verifies NewWrLearnMode behavior.
func TestNewWrLearnMode(t *testing.T) {
	t.Run("creates write learn mode command enabled", func(t *testing.T) {
		cmd, err := NewWrLearnMode(true, 30000, 0x01)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_LEARNMODE {
			t.Errorf("expected CommandCode WR_LEARNMODE, got 0x%02x", cmd.CommandCode)
		}

		if !cmd.EnableLearnMode {
			t.Errorf("expected EnableLearnMode = true, got false")
		}

		if cmd.Timeout != 30000 {
			t.Errorf("expected Timeout 30000, got %d", cmd.Timeout)
		}

		if cmd.Channel != 0x01 {
			t.Errorf("expected Channel 0x01, got 0x%02x", cmd.Channel)
		}
	})

	t.Run("creates write learn mode command disabled", func(t *testing.T) {
		cmd, err := NewWrLearnMode(false, 0, 0x00)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.EnableLearnMode {
			t.Errorf("expected EnableLearnMode = false, got true")
		}

		if cmd.Timeout != 0 {
			t.Errorf("expected Timeout 0, got %d", cmd.Timeout)
		}
	})
}

// TestWrLearnMode_Serialize verifies WrLearnMode_Serialize behavior.
func TestWrLearnMode_Serialize(t *testing.T) {
	t.Run("serializes write learn mode command", func(t *testing.T) {
		cmd, _ := NewWrLearnMode(true, 30000, 0x01)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + EnableLearnMode(1) + Timeout(4) + Channel(1) = 7 bytes
		if len(telegram.Data) != 7 {
			t.Errorf("expected Data length 7, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandWR_LEARNMODE) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandWR_LEARNMODE, telegram.Data[0])
		}
	})
}

// TestParseRdLearnModeResponseOK verifies ParseRdLearnModeResponseOK behavior.
func TestParseRdLearnModeResponseOK(t *testing.T) {
	t.Run("parses read learn mode response enabled", func(t *testing.T) {
		// Response: LearnModeStatus(1) + Channel(1) = 2 bytes
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x01},
			OptData: []byte{0x01},
		}

		result, err := ParseRdLearnModeResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if !result.LearnModeStatus {
			t.Errorf("expected LearnModeStatus = true, got false")
		}

		if result.Channel != 0x01 {
			t.Errorf("expected Channel 0x01, got 0x%02x", result.Channel)
		}
	})

	t.Run("parses read learn mode response disabled", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x00},
			OptData: []byte{0x00},
		}

		result, err := ParseRdLearnModeResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.LearnModeStatus {
			t.Errorf("expected LearnModeStatus = false, got true")
		}

		if result.Channel != 0x00 {
			t.Errorf("expected Channel 0x00, got 0x%02x", result.Channel)
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x01},
			OptData: []byte{0x01},
		}

		_, err := ParseRdLearnModeResponseOK(resp)
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

		_, err := ParseRdLearnModeResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}

		if err.Error() != "failed to deserialize response" {
			t.Errorf("expected error 'failed to deserialize response', got '%s'", err.Error())
		}
	})
}

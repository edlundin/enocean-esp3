package commoncommand

import (
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

func TestNewWrRemanCode(t *testing.T) {
	t.Run("creates write reman code command", func(t *testing.T) {
		cmd, err := NewWrRemanCode(0x12345678)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_REMAN_CODE {
			t.Errorf("expected CommandCode WR_REMAN_CODE, got 0x%02x", cmd.CommandCode)
		}

		if cmd.SecureCode != 0x12345678 {
			t.Errorf("expected SecureCode 0x12345678, got 0x%08x", cmd.SecureCode)
		}
	})

	t.Run("creates write reman code command with zero code", func(t *testing.T) {
		cmd, err := NewWrRemanCode(0x00000000)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.SecureCode != 0x00000000 {
			t.Errorf("expected SecureCode 0x00000000, got 0x%08x", cmd.SecureCode)
		}
	})
}

func TestWrRemanCode_Serialize(t *testing.T) {
	t.Run("serializes write reman code command", func(t *testing.T) {
		cmd, _ := NewWrRemanCode(0x12345678)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + SecureCode(4) = 5 bytes
		if len(telegram.Data) != 5 {
			t.Errorf("expected Data length 5, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandWR_REMAN_CODE) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandWR_REMAN_CODE, telegram.Data[0])
		}
	})
}

func TestNewWrRemanRepeating(t *testing.T) {
	t.Run("creates write reman repeating command with repetition enabled", func(t *testing.T) {
		cmd, err := NewWrRemanRepeating(true)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_REMAN_REPEATING {
			t.Errorf("expected CommandCode WR_REMAN_REPEATING, got 0x%02x", cmd.CommandCode)
		}

		if !cmd.SetRemanRepetition {
			t.Errorf("expected SetRemanRepetition = true, got false")
		}
	})

	t.Run("creates write reman repeating command with repetition disabled", func(t *testing.T) {
		cmd, err := NewWrRemanRepeating(false)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.SetRemanRepetition {
			t.Errorf("expected SetRemanRepetition = false, got true")
		}
	})
}

func TestWrRemanRepeating_Serialize(t *testing.T) {
	t.Run("serializes write reman repeating command", func(t *testing.T) {
		cmd, _ := NewWrRemanRepeating(true)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + SetRemanRepetition(1) = 2 bytes
		if len(telegram.Data) != 2 {
			t.Errorf("expected Data length 2, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandWR_REMAN_REPEATING) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandWR_REMAN_REPEATING, telegram.Data[0])
		}

		if telegram.Data[1] != 0x01 {
			t.Errorf("expected Data[1] = 0x01 (true), got 0x%02x", telegram.Data[1])
		}
	})

	t.Run("serializes write reman repeating command with false", func(t *testing.T) {
		cmd, _ := NewWrRemanRepeating(false)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if telegram.Data[1] != 0x00 {
			t.Errorf("expected Data[1] = 0x00 (false), got 0x%02x", telegram.Data[1])
		}
	})
}

func TestNewRdRemanRepeating(t *testing.T) {
	t.Run("creates read reman repeating command", func(t *testing.T) {
		cmd, err := NewRdRemanRepeating()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandRD_REMAN_REPEATING {
			t.Errorf("expected CommandCode RD_REMAN_REPEATING, got 0x%02x", cmd.CommandCode)
		}
	})
}

func TestRdRemanRepeating_Serialize(t *testing.T) {
	t.Run("serializes read reman repeating command", func(t *testing.T) {
		cmd, _ := NewRdRemanRepeating()
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(telegram.Data) != 1 {
			t.Errorf("expected Data length 1, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandRD_REMAN_REPEATING) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandRD_REMAN_REPEATING, telegram.Data[0])
		}
	})
}

func TestParseRdRemanRepeatingResponseOK(t *testing.T) {
	t.Run("parses reman repeating response enabled", func(t *testing.T) {
		// Response: RemanRepetitionEnabled(1) = 1 byte
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x01},
			OptData: nil,
		}

		result, err := ParseRdRemanRepeatingResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if !result.RemanRepetitionEnabled {
			t.Errorf("expected RemanRepetitionEnabled = true, got false")
		}
	})

	t.Run("parses reman repeating response disabled", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x00},
			OptData: nil,
		}

		result, err := ParseRdRemanRepeatingResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.RemanRepetitionEnabled {
			t.Errorf("expected RemanRepetitionEnabled = false, got true")
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x01},
			OptData: nil,
		}

		_, err := ParseRdRemanRepeatingResponseOK(resp)
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

		_, err := ParseRdRemanRepeatingResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}

		if err.Error() != "failed to deserialize response" {
			t.Errorf("expected error 'failed to deserialize response', got '%s'", err.Error())
		}
	})
}

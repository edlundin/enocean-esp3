package commoncommand

import (
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
)

func TestNewWrMode(t *testing.T) {
	t.Run("creates write mode command with ERP1 mode", func(t *testing.T) {
		cmd, err := NewWrMode(enums.RadioModeERP1)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_MODE {
			t.Errorf("expected CommandCode WR_MODE, got 0x%02x", cmd.CommandCode)
		}

		if cmd.Mode != enums.RadioModeERP1 {
			t.Errorf("expected Mode ERP1, got %v", cmd.Mode)
		}
	})

	t.Run("creates write mode command with ERP2 mode", func(t *testing.T) {
		cmd, err := NewWrMode(enums.RadioModeERP2)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.Mode != enums.RadioModeERP2 {
			t.Errorf("expected Mode ERP2, got %v", cmd.Mode)
		}
	})
}

func TestWrMode_Serialize(t *testing.T) {
	t.Run("serializes write mode command", func(t *testing.T) {
		cmd, _ := NewWrMode(enums.RadioModeERP1)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + Mode(1) = 2 bytes
		if len(telegram.Data) != 2 {
			t.Errorf("expected Data length 2, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandWR_MODE) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandWR_MODE, telegram.Data[0])
		}

		if telegram.Data[1] != byte(enums.RadioModeERP1) {
			t.Errorf("expected Data[1] = 0x%02x, got 0x%02x", enums.RadioModeERP1, telegram.Data[1])
		}
	})
}

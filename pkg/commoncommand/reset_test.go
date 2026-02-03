package commoncommand

import (
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
)

func TestNewWrReset(t *testing.T) {
	t.Run("creates reset command", func(t *testing.T) {
		cmd, err := NewWrReset()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_RESET {
			t.Errorf("expected CommandCode WR_RESET, got 0x%02x", cmd.CommandCode)
		}
	})
}

func TestWrReset_Serialize(t *testing.T) {
	t.Run("serializes reset command", func(t *testing.T) {
		cmd, _ := NewWrReset()
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) = 1 byte
		if len(telegram.Data) != 1 {
			t.Errorf("expected Data length 1, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandWR_RESET) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandWR_RESET, telegram.Data[0])
		}

		// OptData should be nil for this command
		if telegram.OptData != nil {
			t.Errorf("expected nil OptData, got %v", telegram.OptData)
		}
	})
}

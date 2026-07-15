package commoncommand

import (
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
)

// TestNewWrSubTel verifies NewWrSubTel behavior.
func TestNewWrSubTel(t *testing.T) {
	t.Run("creates write subtel command with toggle enabled", func(t *testing.T) {
		cmd, err := NewWrSubTel(true)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_SUBTEL {
			t.Errorf("expected CommandCode WR_SUBTEL, got 0x%02x", cmd.CommandCode)
		}

		if !cmd.Toggle {
			t.Errorf("expected Toggle = true, got false")
		}
	})

	t.Run("creates write subtel command with toggle disabled", func(t *testing.T) {
		cmd, err := NewWrSubTel(false)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.Toggle {
			t.Errorf("expected Toggle = false, got true")
		}
	})
}

// TestWrSubTel_Serialize verifies WrSubTel_Serialize behavior.
func TestWrSubTel_Serialize(t *testing.T) {
	t.Run("serializes write subtel command", func(t *testing.T) {
		cmd, _ := NewWrSubTel(true)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + Toggle(1) = 2 bytes
		if len(telegram.Data) != 2 {
			t.Errorf("expected Data length 2, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandWR_SUBTEL) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandWR_SUBTEL, telegram.Data[0])
		}

		if telegram.Data[1] != 0x01 {
			t.Errorf("expected Data[1] = 0x01 (true), got 0x%02x", telegram.Data[1])
		}
	})

	t.Run("serializes write subtel command with false", func(t *testing.T) {
		cmd, _ := NewWrSubTel(false)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if telegram.Data[1] != 0x00 {
			t.Errorf("expected Data[1] = 0x00 (false), got 0x%02x", telegram.Data[1])
		}
	})
}

package commoncommand

import (
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
)

// TestNewWrSleep verifies NewWrSleep behavior.
func TestNewWrSleep(t *testing.T) {
	t.Run("creates write sleep command", func(t *testing.T) {
		cmd, err := NewWrSleep(10000)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_SLEEP {
			t.Errorf("expected CommandCode WR_SLEEP, got 0x%02x", cmd.CommandCode)
		}

		if cmd.DeepSleepPeriod != 10000 {
			t.Errorf("expected DeepSleepPeriod 10000, got %d", cmd.DeepSleepPeriod)
		}
	})

	t.Run("creates write sleep command with period clamped to max", func(t *testing.T) {
		// Test that periods above max are clamped to 0xffffff
		cmd, err := NewWrSleep(0x1ffffff) // Larger than maxDeepSleepPeriod
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.DeepSleepPeriod != 0xffffff {
			t.Errorf("expected DeepSleepPeriod clamped to 0xffffff, got 0x%06x", cmd.DeepSleepPeriod)
		}
	})

	t.Run("creates write sleep command with zero period", func(t *testing.T) {
		cmd, err := NewWrSleep(0)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.DeepSleepPeriod != 0 {
			t.Errorf("expected DeepSleepPeriod 0, got %d", cmd.DeepSleepPeriod)
		}
	})
}

// TestWrSleep_Serialize verifies WrSleep_Serialize behavior.
func TestWrSleep_Serialize(t *testing.T) {
	t.Run("serializes write sleep command", func(t *testing.T) {
		cmd, _ := NewWrSleep(10000)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + DeepSleepPeriod(4) = 5 bytes
		if len(telegram.Data) != 5 {
			t.Errorf("expected Data length 5, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandWR_SLEEP) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandWR_SLEEP, telegram.Data[0])
		}
	})
}

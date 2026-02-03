package commoncommand

import (
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
)

func TestNewWrWaitMaturity(t *testing.T) {
	t.Run("creates write wait maturity command with forwarded immediately", func(t *testing.T) {
		cmd, err := NewWrWaitMaturity(enums.MaturityFORWARDED_IMMEDIATELY)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_WAIT_MATURITY {
			t.Errorf("expected CommandCode WR_WAIT_MATURITY, got 0x%02x", cmd.CommandCode)
		}

		if cmd.Maturity != enums.MaturityFORWARDED_IMMEDIATELY {
			t.Errorf("expected Maturity FORWARDED_IMMEDIATELY, got %v", cmd.Maturity)
		}
	})

	t.Run("creates write wait maturity command with forwarded on timeout", func(t *testing.T) {
		cmd, err := NewWrWaitMaturity(enums.MaturityFORWARDED_ON_TIMEOUT)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.Maturity != enums.MaturityFORWARDED_ON_TIMEOUT {
			t.Errorf("expected Maturity FORWARDED_ON_TIMEOUT, got %v", cmd.Maturity)
		}
	})
}

func TestWrWaitMaturity_Serialize(t *testing.T) {
	t.Run("serializes write wait maturity command", func(t *testing.T) {
		cmd, _ := NewWrWaitMaturity(enums.MaturityFORWARDED_IMMEDIATELY)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + Maturity(1) = 2 bytes
		if len(telegram.Data) != 2 {
			t.Errorf("expected Data length 2, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandWR_WAIT_MATURITY) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandWR_WAIT_MATURITY, telegram.Data[0])
		}

		if telegram.Data[1] != byte(enums.MaturityFORWARDED_IMMEDIATELY) {
			t.Errorf("expected Data[1] = 0x%02x, got 0x%02x", enums.MaturityFORWARDED_IMMEDIATELY, telegram.Data[1])
		}
	})
}

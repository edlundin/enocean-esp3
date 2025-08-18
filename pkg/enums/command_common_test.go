package enums

import (
	"testing"
)

func TestParseCommonCommandFromByte(t *testing.T) {
	t.Run("parses all valid common commands correctly", func(t *testing.T) {
		testCases := []struct {
			input    byte
			expected CommonCommand
		}{
			{0x01, CommonCommandWR_SLEEP},
			{0x02, CommonCommandWR_RESET},
			{0x03, CommonCommandRD_VERSION},
			{0x04, CommonCommandRD_SYS_LOG},
			{0x05, CommonCommandWR_SYS_LOG},
			{0x06, CommonCommandWR_BIST},
			{0x07, CommonCommandWR_IDBASE},
			{0x08, CommonCommandRD_IDBASE},
			{0x09, CommonCommandWR_REPEATER},
			{0x0a, CommonCommandRD_REPEATER},
			{0x0b, CommonCommandWR_FILTER_ADD},
			{0x0c, CommonCommandWR_FILTER_DEL},
			{0x0d, CommonCommandWR_FILTER_DEL_ALL},
			{0x0e, CommonCommandWR_FILTER_ENABLE},
			{0x0f, CommonCommandRD_FILTER},
			{0x10, CommonCommandWR_WAIT_MATURITY},
			{0x11, CommonCommandWR_SUBTEL},
			{0x12, CommonCommandWR_MEM},
			{0x13, CommonCommandRD_MEM},
			{0x14, CommonCommandRD_MEM_ADDRESS},
			{0x15, CommonCommandRD_SECURITY},
			{0x16, CommonCommandWR_SECURITY},
			{0x17, CommonCommandWR_LEARNMODE},
			{0x18, CommonCommandRD_LEARNMODE},
			{0x19, CommonCommandWR_SECUREDEVICE_ADD},
			{0x1a, CommonCommandWR_SECUREDEVICE_DEL},
			{0x1b, CommonCommandRD_SECUREDEVICE_BY_INDEX},
			{0x1c, CommonCommandWR_MODE},
			{0x1d, CommonCommandRD_NUMSECUREDEVICES},
			{0x1e, CommonCommandRD_SECUREDEVICE_BY_ID},
			{0x1f, CommonCommandWR_SECUREDEVICE_ADD_PSK},
			{0x20, CommonCommandWR_SECUREDEVICE_SENDTEACHIN},
			{0x21, CommonCommandWR_TEMPORARY_RLC_WINDOW},
			{0x22, CommonCommandRD_SECUREDEVICE_PSK},
			{0x23, CommonCommandRD_DUTYCYCLE_LIMIT},
			{0x24, CommonCommandSET_BAUDRATE},
			{0x25, CommonCommandGET_FREQUENCY_INFO},
			{0x27, CommonCommandGET_STEPCODE},
			{0x2e, CommonCommandWR_REMAN_CODE},
			{0x2f, CommonCommandWR_STARTUP_DELAY},
			{0x30, CommonCommandWR_REMAN_REPEATING},
			{0x31, CommonCommandRD_REMAN_REPEATING},
			{0x32, CommonCommandSET_NOISETHRESHOLD},
			{0x33, CommonCommandGET_NOISETHRESHOLD},
			{0x36, CommonCommandWR_RLC_SAVE_PERIOD},
			{0x37, CommonCommandWR_RLC_LEGACY_MODE},
			{0x38, CommonCommandWR_SECUREDEVICEV2_ADD},
			{0x39, CommonCommandRD_SECUREDEVICEV2_BY_INDEX},
			{0x3a, CommonCommandWR_RSSITEST_MODE},
			{0x3b, CommonCommandRD_RSSITEST_MODE},
			{0x3c, CommonCommandWR_SECUREDEVICE_MAINTENANCEKEY},
			{0x3d, CommonCommandRD_SECUREDEVICE_MAINTENANCEKEY},
			{0x3e, CommonCommandWR_TRANSPARENT_MODE},
			{0x3f, CommonCommandRD_TRANSPARENT_MODE},
			{0x40, CommonCommandWR_TX_ONLY_MODE},
			{0x41, CommonCommandRD_TX_ONLY_MODE},
		}

		for _, tc := range testCases {
			t.Run(tc.expected.String(), func(t *testing.T) {
				result, err := ParseCommonCommandFromByte(tc.input)
				if err != nil {
					t.Errorf("expected no error for input 0x%02x, got: %s", tc.input, err)
				}
				if result != tc.expected {
					t.Errorf("expected %v, got %v", tc.expected, result)
				}
			})
		}
	})

	t.Run("returns error for invalid common command", func(t *testing.T) {
		invalidInputs := []byte{0x00, 0x26, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x34, 0x35, 0x42, 0xff}

		for _, input := range invalidInputs {
			t.Run(t.Name(), func(t *testing.T) {
				result, err := ParseCommonCommandFromByte(input)
				if err == nil {
					t.Errorf("expected error for input 0x%02x, got nil", input)
				}
				if err.Error() != "invalid common command" {
					t.Errorf("expected error 'invalid common command', got '%s'", err.Error())
				}
				if result != 0 {
					t.Errorf("expected result 0, got %v", result)
				}
			})
		}
	})
}

func TestCommonCommand_String(t *testing.T) {
	t.Run("returns correct string for all valid common commands", func(t *testing.T) {
		testCases := []struct {
			input    CommonCommand
			expected string
		}{
			{CommonCommandWR_SLEEP, "WR_SLEEP"},
			{CommonCommandWR_RESET, "WR_RESET"},
			{CommonCommandRD_VERSION, "RD_VERSION"},
			{CommonCommandRD_SYS_LOG, "RD_SYS_LOG"},
			{CommonCommandWR_SYS_LOG, "WR_SYS_LOG"},
			{CommonCommandWR_BIST, "WR_BIST"},
			{CommonCommandWR_IDBASE, "WR_IDBASE"},
			{CommonCommandRD_IDBASE, "RD_IDBASE"},
			{CommonCommandWR_REPEATER, "WR_REPEATER"},
			{CommonCommandRD_REPEATER, "RD_REPEATER"},
			{CommonCommandWR_FILTER_ADD, "WR_FILTER_ADD"},
			{CommonCommandWR_FILTER_DEL, "WR_FILTER_DEL"},
			{CommonCommandWR_FILTER_DEL_ALL, "WR_FILTER_DEL_ALL"},
			{CommonCommandWR_FILTER_ENABLE, "WR_FILTER_ENABLE"},
			{CommonCommandRD_FILTER, "RD_FILTER"},
			{CommonCommandWR_WAIT_MATURITY, "WR_WAIT_MATURITY"},
			{CommonCommandWR_SUBTEL, "WR_SUBTEL"},
			{CommonCommandWR_MEM, "WR_MEM"},
			{CommonCommandRD_MEM, "RD_MEM"},
			{CommonCommandRD_MEM_ADDRESS, "RD_MEM_ADDRESS"},
			{CommonCommandRD_SECURITY, "RD_SECURITY"},
			{CommonCommandWR_SECURITY, "WR_SECURITY"},
			{CommonCommandWR_LEARNMODE, "WR_LEARNMODE"},
			{CommonCommandRD_LEARNMODE, "RD_LEARNMODE"},
			{CommonCommandWR_SECUREDEVICE_ADD, "WR_SECUREDEVICE_ADD"},
			{CommonCommandWR_SECUREDEVICE_DEL, "WR_SECUREDEVICE_DEL"},
			{CommonCommandRD_SECUREDEVICE_BY_INDEX, "RD_SECUREDEVICE_BY_INDEX"},
			{CommonCommandWR_MODE, "WR_MODE"},
			{CommonCommandRD_NUMSECUREDEVICES, "RD_NUMSECUREDEVICES"},
			{CommonCommandRD_SECUREDEVICE_BY_ID, "RD_SECUREDEVICE_BY_ID"},
			{CommonCommandWR_SECUREDEVICE_ADD_PSK, "WR_SECUREDEVICE_ADD_PSK"},
			{CommonCommandWR_SECUREDEVICE_SENDTEACHIN, "WR_SECUREDEVICE_SENDTEACHIN"},
			{CommonCommandWR_TEMPORARY_RLC_WINDOW, "WR_TEMPORARY_RLC_WINDOW"},
			{CommonCommandRD_SECUREDEVICE_PSK, "RD_SECUREDEVICE_PSK"},
			{CommonCommandRD_DUTYCYCLE_LIMIT, "RD_DUTYCYCLE_LIMIT"},
			{CommonCommandSET_BAUDRATE, "SET_BAUDRATE"},
			{CommonCommandGET_FREQUENCY_INFO, "GET_FREQUENCY_INFO"},
			{CommonCommandGET_STEPCODE, "GET_STEPCODE"},
			{CommonCommandWR_REMAN_CODE, "WR_REMAN_CODE"},
			{CommonCommandWR_STARTUP_DELAY, "WR_STARTUP_DELAY"},
			{CommonCommandWR_REMAN_REPEATING, "WR_REMAN_REPEATING"},
			{CommonCommandRD_REMAN_REPEATING, "RD_REMAN_REPEATING"},
			{CommonCommandSET_NOISETHRESHOLD, "SET_NOISETHRESHOLD"},
			{CommonCommandGET_NOISETHRESHOLD, "GET_NOISETHRESHOLD"},
			{CommonCommandWR_RLC_SAVE_PERIOD, "WR_RLC_SAVE_PERIOD"},
			{CommonCommandWR_RLC_LEGACY_MODE, "WR_RLC_LEGACY_MODE"},
			{CommonCommandWR_SECUREDEVICEV2_ADD, "WR_SECUREDEVICEV2_ADD"},
			{CommonCommandRD_SECUREDEVICEV2_BY_INDEX, "RD_SECUREDEVICEV2_BY_INDEX"},
			{CommonCommandWR_RSSITEST_MODE, "WR_RSSITEST_MODE"},
			{CommonCommandRD_RSSITEST_MODE, "RD_RSSITEST_MODE"},
			{CommonCommandWR_SECUREDEVICE_MAINTENANCEKEY, "WR_SECUREDEVICE_MAINTENANCEKEY"},
			{CommonCommandRD_SECUREDEVICE_MAINTENANCEKEY, "RD_SECUREDEVICE_MAINTENANCEKEY"},
			{CommonCommandWR_TRANSPARENT_MODE, "WR_TRANSPARENT_MODE"},
			{CommonCommandWR_TX_ONLY_MODE, "WR_TX_ONLY_MODE"},
			{CommonCommandRD_TX_ONLY_MODE, "RD_TX_ONLY_MODE"},
		}

		for _, tc := range testCases {
			t.Run(tc.expected, func(t *testing.T) {
				result := tc.input.String()
				if result != tc.expected {
					t.Errorf("expected '%s', got '%s'", tc.expected, result)
				}
			})
		}
	})

	t.Run("returns UNKNOWN for invalid common commands", func(t *testing.T) {
		invalidTypes := []CommonCommand{0x00, 0x26, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x34, 0x35, 0x42, 0xff}

		for _, input := range invalidTypes {
			t.Run(t.Name(), func(t *testing.T) {
				result := input.String()
				if result != "UNKNOWN" {
					t.Errorf("expected 'UNKNOWN' for input %v, got '%s'", input, result)
				}
			})
		}
	})

	t.Run("handles CommonCommandRD_TRANSPARENT_MODE correctly", func(t *testing.T) {
		// This command is handled by String() but not by ParseCommonCommandFromByte
		command := CommonCommand(0x3f) // CommonCommandRD_TRANSPARENT_MODE
		result := command.String()
		expected := "RD_TRANSPARENT_MODE"
		if result != expected {
			t.Errorf("expected '%s', got '%s'", expected, result)
		}
	})
}

func TestCommonCommandValid(t *testing.T) {
	t.Run("CommonCommand_Valid", func(t *testing.T) {
		// Test valid commands
		validCommands := []CommonCommand{
			CommonCommandWR_SLEEP,
			CommonCommandWR_RESET,
			CommonCommandRD_VERSION,
			CommonCommandWR_SYS_LOG,
			CommonCommandWR_BIST,
		}
		for _, cmd := range validCommands {
			if !cmd.Valid() {
				t.Errorf("CommonCommand %v should be valid", cmd)
			}
		}

		// Test invalid commands
		invalidCommands := []CommonCommand{0x42, 0x99, 0xFF}
		for _, cmd := range invalidCommands {
			if cmd.Valid() {
				t.Errorf("CommonCommand %v should not be valid", cmd)
			}
		}
	})
}

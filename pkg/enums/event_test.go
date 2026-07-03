package enums

import (
	"testing"
)

func TestEventCode(t *testing.T) {
	t.Run("ParseEventCodeFromByte", func(t *testing.T) {
		tests := []struct {
			name    string
			input   byte
			want    EventCode
			wantErr bool
		}{
			{"SA_RECLAIM_NOT_SUCCESSFUL", 0x01, EventCodeSA_RECLAIM_NOT_SUCCESSFUL, false},
			{"SA_CONFIRM_LEARN", 0x02, EventCodeSA_CONFIRM_LEARN, false},
			{"SA_LEARN_ACK", 0x03, EventCodeSA_LEARN_ACK, false},
			{"CO_READY", 0x04, EventCodeCO_READY, false},
			{"CO_EVENT_SECUREDEVICES", 0x05, EventCodeCO_EVENT_SECUREDEVICES, false},
			{"CO_DUTYCYCLE_LIMIT", 0x06, EventCodeCO_DUTYCYCLE_LIMIT, false},
			{"CO_TRANSMIT_FAILED", 0x07, EventCodeCO_TRANSMIT_FAILED, false},
			{"CO_TX_DONE", 0x08, EventCodeCO_TX_DONE, false},
			{"CO_LRN_MODE_DISABLED", 0x09, EventCodeCO_LRN_MODE_DISABLED, false},
			{"Invalid", 0x0A, 0, true},
			{"Invalid", 0xFF, 0, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := ParseEventCodeFromByte(tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("ParseEventCodeFromByte() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("ParseEventCodeFromByte() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("EventCode_String", func(t *testing.T) {
		tests := []struct {
			name string
			code EventCode
			want string
		}{
			{"SA_RECLAIM_NOT_SUCCESSFUL", EventCodeSA_RECLAIM_NOT_SUCCESSFUL, "SA_RECLAIM_NOT_SUCCESSFUL"},
			{"SA_CONFIRM_LEARN", EventCodeSA_CONFIRM_LEARN, "SA_CONFIRM_LEARN"},
			{"SA_LEARN_ACK", EventCodeSA_LEARN_ACK, "SA_LEARN_ACK"},
			{"CO_READY", EventCodeCO_READY, "CO_READY"},
			{"CO_EVENT_SECUREDEVICES", EventCodeCO_EVENT_SECUREDEVICES, "CO_EVENT_SECUREDEVICES"},
			{"CO_DUTYCYCLE_LIMIT", EventCodeCO_DUTYCYCLE_LIMIT, "CO_DUTYCYCLE_LIMIT"},
			{"CO_TRANSMIT_FAILED", EventCodeCO_TRANSMIT_FAILED, "CO_TRANSMIT_FAILED"},
			{"CO_TX_DONE", EventCodeCO_TX_DONE, "CO_TX_DONE"},
			{"CO_LRN_MODE_DISABLED", EventCodeCO_LRN_MODE_DISABLED, "CO_LRN_MODE_DISABLED"},
			{"Unknown", 0x0A, "UNKNOWN"},
			{"Unknown", 0xFF, "UNKNOWN"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.code.String(); got != tt.want {
					t.Errorf("EventCode.String() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("EventCode_Valid", func(t *testing.T) {
		// Test valid event codes
		validEvents := []EventCode{
			EventCodeSA_RECLAIM_NOT_SUCCESSFUL,
			EventCodeSA_CONFIRM_LEARN,
			EventCodeSA_LEARN_ACK,
			EventCodeCO_READY,
			EventCodeCO_EVENT_SECUREDEVICES,
			EventCodeCO_DUTYCYCLE_LIMIT,
			EventCodeCO_TRANSMIT_FAILED,
			EventCodeCO_TX_DONE,
			EventCodeCO_LRN_MODE_DISABLED,
		}
		for _, ec := range validEvents {
			if !ec.Valid() {
				t.Errorf("EventCode %v should be valid", ec)
			}
		}

		// Test invalid event codes
		invalidEvents := []EventCode{0x0A, 0x0B, 0x0C, 0xFF}
		for _, ec := range invalidEvents {
			if ec.Valid() {
				t.Errorf("EventCode %v should not be valid", ec)
			}
		}
	})
}

func TestLearnAckConfirmCode(t *testing.T) {
	t.Run("ParseLearnAckConfirmCodeFromByte", func(t *testing.T) {
		tests := []struct {
			name    string
			input   byte
			want    LearnAckConfirmCode
			wantErr bool
		}{
			{"LRN_IN", 0x00, LearnAckConfirmCodeLRN_IN, false},
			{"EEP_NOT_ACCEPTED", 0x11, LearnAckConfirmCodeEEP_NOT_ACCEPTED, false},
			{"NO_PLACE_IN_PM", 0x12, LearnAckConfirmCodeNO_PLACE_IN_PM, false},
			{"NO_PLACE_IN_CONTROLLER", 0x13, LearnAckConfirmCodeNO_PLACE_IN_CONTROLLER, false},
			{"RSSI_NOT_GOOD_ENOUGH", 0x14, LearnAckConfirmCodeRSSI_NOT_GOOD_ENOUGH, false},
			{"LRN_OUT", 0x20, LearnAckConfirmCodeLRN_OUT, false},
			{"FUNCTION_NOT_SUPPORTED", 0xff, LearnAckConfirmCodeFUNCTION_NOT_SUPPORTED, false},
			{"Invalid", 0x01, 0, true},
			{"Invalid", 0x10, 0, true},
			{"Invalid", 0x15, 0, true},
			{"Invalid", 0x21, 0, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := ParseLearnAckConfirmCodeFromByte(tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("ParseLearnAckConfirmCodeFromByte() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("ParseLearnAckConfirmCodeFromByte() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("LearnAckConfirmCode_String", func(t *testing.T) {
		tests := []struct {
			name string
			code LearnAckConfirmCode
			want string
		}{
			{"LRN_IN", LearnAckConfirmCodeLRN_IN, "LRN_IN"},
			{"EEP_NOT_ACCEPTED", LearnAckConfirmCodeEEP_NOT_ACCEPTED, "EEP_NOT_ACCEPTED"},
			{"NO_PLACE_IN_PM", LearnAckConfirmCodeNO_PLACE_IN_PM, "NO_PLACE_IN_PM"},
			{"NO_PLACE_IN_CONTROLLER", LearnAckConfirmCodeNO_PLACE_IN_CONTROLLER, "NO_PLACE_IN_CONTROLLER"},
			{"RSSI_NOT_GOOD_ENOUGH", LearnAckConfirmCodeRSSI_NOT_GOOD_ENOUGH, "RSSI_NOT_GOOD_ENOUGH"},
			{"LRN_OUT", LearnAckConfirmCodeLRN_OUT, "LRN_OUT"},
			{"FUNCTION_NOT_SUPPORTED", LearnAckConfirmCodeFUNCTION_NOT_SUPPORTED, "FUNCTION_NOT_SUPPORTED"},
			{"Unknown", 0x01, "UNKNOWN"},
			{"Unknown", 0x10, "UNKNOWN"},
			{"Unknown", 0x15, "UNKNOWN"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.code.String(); got != tt.want {
					t.Errorf("LearnAckConfirmCode.String() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("LearnAckConfirmCode_Valid", func(t *testing.T) {
		// Test valid codes
		validCodes := []LearnAckConfirmCode{
			LearnAckConfirmCodeLRN_IN,
			LearnAckConfirmCodeEEP_NOT_ACCEPTED,
			LearnAckConfirmCodeNO_PLACE_IN_PM,
			LearnAckConfirmCodeNO_PLACE_IN_CONTROLLER,
			LearnAckConfirmCodeRSSI_NOT_GOOD_ENOUGH,
			LearnAckConfirmCodeLRN_OUT,
			LearnAckConfirmCodeFUNCTION_NOT_SUPPORTED,
		}
		for _, code := range validCodes {
			if !code.Valid() {
				t.Errorf("LearnAckConfirmCode %v should be valid", code)
			}
		}

		// Test invalid codes
		invalidCodes := []LearnAckConfirmCode{0x01, 0x10, 0x15, 0x21, 0xfe}
		for _, code := range invalidCodes {
			if code.Valid() {
				t.Errorf("LearnAckConfirmCode %v should not be valid", code)
			}
		}
	})
}

func TestWakeUpCause(t *testing.T) {
	t.Run("ParseWakeUpCauseFromByte", func(t *testing.T) {
		tests := []struct {
			name    string
			input   byte
			want    WakeUpCause
			wantErr bool
		}{
			{"VOLTAGE_SUPPLY_DROP", 0x00, WakeUpCauseVOLTAGE_SUPPLY_DROP, false},
			{"RESET_BY_RESET_PIN", 0x01, WakeUpCauseRESET_BY_RESET_PIN, false},
			{"WATCHDOG_TIMEOUT", 0x02, WakeUpCauseWATCHDOG_TIMEOUT, false},
			{"FLYWHEEL_TIMEOUT", 0x03, WakeUpCauseFLYWHEEL_TIMEOUT, false},
			{"PARITY_ERROR", 0x04, WakeUpCausePARITY_ERROR, false},
			{"HARDWARE_PARITY_ERROR_IN_MEMORY", 0x05, WakeUpCauseHARDWARE_PARITY_ERROR_IN_MEMORY, false},
			{"REQUESTED_MEMORY_LOCATION_NOT_FOUND", 0x06, WakeUpCauseREQUESTED_MEMORY_LOCATION_NOT_FOUND, false},
			{"TRIGGER_PIN0", 0x07, WakeUpCauseTRIGGER_PIN0, false},
			{"TRIGGER_PIN1", 0x08, WakeUpCauseTRIGGER_PIN1, false},
			{"UNKNOWN_RESET_SOURCE", 0x09, WakeUpCauseUNKNOWN_RESET_SOURCE, false},
			{"UART_WAKE_UP", 0x0a, WakeUpCauseUART_WAKE_UP, false},
			{"Invalid", 0x0b, 0, true},
			{"Invalid", 0xFF, 0, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := ParseWakeUpCauseFromByte(tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("ParseWakeUpCauseFromByte() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("ParseWakeUpCauseFromByte() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("WakeUpCause_String", func(t *testing.T) {
		tests := []struct {
			name string
			code WakeUpCause
			want string
		}{
			{"VOLTAGE_SUPPLY_DROP", WakeUpCauseVOLTAGE_SUPPLY_DROP, "VOLTAGE_SUPPLY_DROP"},
			{"RESET_BY_RESET_PIN", WakeUpCauseRESET_BY_RESET_PIN, "RESET_BY_RESET_PIN"},
			{"WATCHDOG_TIMEOUT", WakeUpCauseWATCHDOG_TIMEOUT, "WATCHDOG_TIMEOUT"},
			{"FLYWHEEL_TIMEOUT", WakeUpCauseFLYWHEEL_TIMEOUT, "FLYWHEEL_TIMEOUT"},
			{"PARITY_ERROR", WakeUpCausePARITY_ERROR, "PARITY_ERROR"},
			{"HARDWARE_PARITY_ERROR_IN_MEMORY", WakeUpCauseHARDWARE_PARITY_ERROR_IN_MEMORY, "HARDWARE_PARITY_ERROR_IN_MEMORY"},
			{"REQUESTED_MEMORY_LOCATION_NOT_FOUND", WakeUpCauseREQUESTED_MEMORY_LOCATION_NOT_FOUND, "REQUESTED_MEMORY_LOCATION_NOT_FOUND"},
			{"TRIGGER_PIN0", WakeUpCauseTRIGGER_PIN0, "TRIGGER_PIN0"},
			{"TRIGGER_PIN1", WakeUpCauseTRIGGER_PIN1, "TRIGGER_PIN1"},
			{"UNKNOWN_RESET_SOURCE", WakeUpCauseUNKNOWN_RESET_SOURCE, "UNKNOWN_RESET_SOURCE"},
			{"UART_WAKE_UP", WakeUpCauseUART_WAKE_UP, "UART_WAKE_UP"},
			{"Unknown", 0x0b, "UNKNOWN"},
			{"Unknown", 0xFF, "UNKNOWN"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.code.String(); got != tt.want {
					t.Errorf("WakeUpCause.String() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("WakeUpCause_Valid", func(t *testing.T) {
		// Test valid codes
		validCodes := []WakeUpCause{
			WakeUpCauseVOLTAGE_SUPPLY_DROP,
			WakeUpCauseRESET_BY_RESET_PIN,
			WakeUpCauseWATCHDOG_TIMEOUT,
			WakeUpCauseFLYWHEEL_TIMEOUT,
			WakeUpCausePARITY_ERROR,
			WakeUpCauseHARDWARE_PARITY_ERROR_IN_MEMORY,
			WakeUpCauseREQUESTED_MEMORY_LOCATION_NOT_FOUND,
			WakeUpCauseTRIGGER_PIN0,
			WakeUpCauseTRIGGER_PIN1,
			WakeUpCauseUNKNOWN_RESET_SOURCE,
			WakeUpCauseUART_WAKE_UP,
		}
		for _, code := range validCodes {
			if !code.Valid() {
				t.Errorf("WakeUpCause %v should be valid", code)
			}
		}

		// Test invalid codes
		invalidCodes := []WakeUpCause{0x0b, 0x0c, 0xFF}
		for _, code := range invalidCodes {
			if code.Valid() {
				t.Errorf("WakeUpCause %v should not be valid", code)
			}
		}
	})
}

func TestWakeUpMode(t *testing.T) {
	t.Run("ParseWakeUpModeFromByte", func(t *testing.T) {
		tests := []struct {
			name    string
			input   byte
			want    WakeUpMode
			wantErr bool
		}{
			{"STANDARD_SECURITY", 0x00, WakeUpModeSTANDARD_SECURITY, false},
			{"EXTENDED_SECURITY", 0x01, WakeUpModeEXTENDED_SECURITY, false},
			{"Invalid", 0x02, 0, true},
			{"Invalid", 0xFF, 0, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := ParseWakeUpModeFromByte(tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("ParseWakeUpModeFromByte() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("ParseWakeUpModeFromByte() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("WakeUpMode_String", func(t *testing.T) {
		tests := []struct {
			name string
			code WakeUpMode
			want string
		}{
			{"STANDARD_SECURITY", WakeUpModeSTANDARD_SECURITY, "STANDARD_SECURITY"},
			{"EXTENDED_SECURITY", WakeUpModeEXTENDED_SECURITY, "EXTENDED_SECURITY"},
			{"Unknown", 0x02, "UNKNOWN"},
			{"Unknown", 0xFF, "UNKNOWN"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.code.String(); got != tt.want {
					t.Errorf("WakeUpMode.String() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("WakeUpMode_Valid", func(t *testing.T) {
		// Test valid codes
		validCodes := []WakeUpMode{
			WakeUpModeSTANDARD_SECURITY,
			WakeUpModeEXTENDED_SECURITY,
		}
		for _, code := range validCodes {
			if !code.Valid() {
				t.Errorf("WakeUpMode %v should be valid", code)
			}
		}

		// Test invalid codes
		invalidCodes := []WakeUpMode{0x02, 0x03, 0xFF}
		for _, code := range invalidCodes {
			if code.Valid() {
				t.Errorf("WakeUpMode %v should not be valid", code)
			}
		}
	})
}

func TestSecureDeviceEventCause(t *testing.T) {
	t.Run("ParseCOEventSecureFromByte", func(t *testing.T) {
		tests := []struct {
			name    string
			input   byte
			want    SecureDeviceEventCause
			wantErr bool
		}{
			{"SECURE_LINK_TABLE_FULL", 0x00, COEventSecureSECURE_LINK_TABLE_FULL, false},
			{"RESYNC_WRONG_PRIVATE_KEY", 0x01, COEventSecureRESYNC_WRONG_PRIVATE_KEY, false},
			{"WRONG_CMAC_TELEGRAM_THRESHOLD_HIT", 0x02, COEventSecureWRONG_CMAC_TELEGRAM_THRESHOLD_HIT, false},
			{"TEACH_IN_TELEGRAM_CORRUPTED", 0x03, COEventSecureTEACH_IN_TELEGRAM_CORRUPTED, false},
			{"TEACH_IN_PSK_FAILED_NO_PSK_SET", 0x04, COEventSecureTEACH_IN_PSK_FAILED_NO_PSK_SET, false},
			{"TEACH_IN_PSK_FAILED_NO_PSK_GIVEN", 0x05, COEventSecureTEACH_IN_PSK_FAILED_NO_PSK_GIVEN, false},
			{"INCORRECT_CMAC_OR_RKC", 0x06, COEventSecureINCORRECT_CMAC_OR_RKC, false},
			{"RECEIVED_TELEGRAM_FROM_DEVICE_IN_SECURE_LINK_TABLE", 0x07, COEventSecureRECEIVED_TELEGRAM_FROM_DEVICE_IN_SECURE_LINK_TABLE, false},
			{"TEACH_IN_SUCCESSFUL", 0x08, COEventSecureTEACH_IN_SUCCESSFUL, false},
			{"VALID_RLC_SYNC_RECEIVED_VIA_TEACH_IN", 0x09, COEventSecureVALID_RLC_SYNC_RECEIVED_VIA_TEACH_IN, false},
			{"Invalid", 0x0a, 0, true},
			{"Invalid", 0xFF, 0, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := ParseCOEventSecureFromByte(tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("ParseCOEventSecureFromByte() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("ParseCOEventSecureFromByte() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("SecureDeviceEventCause_String", func(t *testing.T) {
		tests := []struct {
			name string
			code SecureDeviceEventCause
			want string
		}{
			{"SECURE_LINK_TABLE_FULL", COEventSecureSECURE_LINK_TABLE_FULL, "SECURE_LINK_TABLE_FULL"},
			{"RESYNC_WRONG_PRIVATE_KEY", COEventSecureRESYNC_WRONG_PRIVATE_KEY, "RESYNC_WRONG_PRIVATE_KEY"},
			{"WRONG_CMAC_TELEGRAM_THRESHOLD_HIT", COEventSecureWRONG_CMAC_TELEGRAM_THRESHOLD_HIT, "WRONG_CMAC_TELEGRAM_THRESHOLD_HIT"},
			{"TEACH_IN_TELEGRAM_CORRUPTED", COEventSecureTEACH_IN_TELEGRAM_CORRUPTED, "TEACH_IN_TELEGRAM_CORRUPTED"},
			{"TEACH_IN_PSK_FAILED_NO_PSK_SET", COEventSecureTEACH_IN_PSK_FAILED_NO_PSK_SET, "TEACH_IN_PSK_FAILED_NO_PSK_SET"},
			{"TEACH_IN_PSK_FAILED_NO_PSK_GIVEN", COEventSecureTEACH_IN_PSK_FAILED_NO_PSK_GIVEN, "TEACH_IN_PSK_FAILED_NO_PSK_GIVEN"},
			{"INCORRECT_CMAC_OR_RKC", COEventSecureINCORRECT_CMAC_OR_RKC, "INCORRECT_CMAC_OR_RKC"},
			{"RECEIVED_TELEGRAM_FROM_DEVICE_IN_SECURE_LINK_TABLE", COEventSecureRECEIVED_TELEGRAM_FROM_DEVICE_IN_SECURE_LINK_TABLE, "RECEIVED_TELEGRAM_FROM_DEVICE_IN_SECURE_LINK_TABLE"},
			{"TEACH_IN_SUCCESSFUL", COEventSecureTEACH_IN_SUCCESSFUL, "TEACH_IN_SUCCESSFUL"},
			{"VALID_RLC_SYNC_RECEIVED_VIA_TEACH_IN", COEventSecureVALID_RLC_SYNC_RECEIVED_VIA_TEACH_IN, "VALID_RLC_SYNC_RECEIVED_VIA_TEACH_IN"},
			{"Unknown", 0x01, "UNKNOWN"},
			{"Unknown", 0xFF, "UNKNOWN"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.code.String(); got != tt.want {
					t.Errorf("SecureDeviceEventCause.String() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("SecureDeviceEventCause_Valid", func(t *testing.T) {
		// Test valid codes
		validCodes := []SecureDeviceEventCause{
			COEventSecureSECURE_LINK_TABLE_FULL,
			COEventSecureRESYNC_WRONG_PRIVATE_KEY,
			COEventSecureWRONG_CMAC_TELEGRAM_THRESHOLD_HIT,
			COEventSecureTEACH_IN_TELEGRAM_CORRUPTED,
			COEventSecureTEACH_IN_PSK_FAILED_NO_PSK_SET,
			COEventSecureTEACH_IN_PSK_FAILED_NO_PSK_GIVEN,
			COEventSecureINCORRECT_CMAC_OR_RKC,
			COEventSecureRECEIVED_TELEGRAM_FROM_DEVICE_IN_SECURE_LINK_TABLE,
			COEventSecureTEACH_IN_SUCCESSFUL,
			COEventSecureVALID_RLC_SYNC_RECEIVED_VIA_TEACH_IN,
		}
		for _, code := range validCodes {
			if !code.Valid() {
				t.Errorf("SecureDeviceEventCause %v should be valid", code)
			}
		}

		// Test invalid codes
		invalidCodes := []SecureDeviceEventCause{0x01, 0x0b, 0xFF}
		for _, code := range invalidCodes {
			if code.Valid() {
				t.Errorf("SecureDeviceEventCause %v should not be valid", code)
			}
		}
	})
}

func TestDutyCycleLimitCause(t *testing.T) {
	t.Run("ParseDutyCycleLimitCauseFromByte", func(t *testing.T) {
		tests := []struct {
			name    string
			input   byte
			want    DutyCycleLimitCause
			wantErr bool
		}{
			{"NOT_YET_REACHED", 0x00, DutyCycleLimitCauseNOT_YET_REACHED, false},
			{"REACHED", 0x01, DutyCycleLimitCauseREACHED, false},
			{"Invalid", 0x02, 0, true},
			{"Invalid", 0xFF, 0, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := ParseDutyCycleLimitCauseFromByte(tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("ParseDutyCycleLimitCauseFromByte() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("ParseDutyCycleLimitCauseFromByte() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("DutyCycleLimitCause_String", func(t *testing.T) {
		tests := []struct {
			name string
			code DutyCycleLimitCause
			want string
		}{
			{"NOT_YET_REACHED", DutyCycleLimitCauseNOT_YET_REACHED, "NOT_YET_REACHED"},
			{"REACHED", DutyCycleLimitCauseREACHED, "REACHED"},
			{"Unknown", 0x02, "UNKNOWN"},
			{"Unknown", 0xFF, "UNKNOWN"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.code.String(); got != tt.want {
					t.Errorf("DutyCycleLimitCause.String() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("DutyCycleLimitCause_Valid", func(t *testing.T) {
		// Test valid codes
		validCodes := []DutyCycleLimitCause{
			DutyCycleLimitCauseNOT_YET_REACHED,
			DutyCycleLimitCauseREACHED,
		}
		for _, code := range validCodes {
			if !code.Valid() {
				t.Errorf("DutyCycleLimitCause %v should be valid", code)
			}
		}

		// Test invalid codes
		invalidCodes := []DutyCycleLimitCause{0x02, 0x03, 0xFF}
		for _, code := range invalidCodes {
			if code.Valid() {
				t.Errorf("DutyCycleLimitCause %v should not be valid", code)
			}
		}
	})
}

func TestTransmitFailedCause(t *testing.T) {
	t.Run("ParseTransmitFailedCauseFromByte", func(t *testing.T) {
		tests := []struct {
			name    string
			input   byte
			want    TransmitFailedCause
			wantErr bool
		}{
			{"CSMA_FAILED_CHANNEL_NOT_FREE", 0x00, TransmitFailedCauseCSMA_FAILED_CHANNEL_NOT_FREE, false},
			{"NO_ACK_RECEIVED", 0x01, TransmitFailedCauseNO_ACK_RECEIVED, false},
			{"Invalid", 0x02, 0, true},
			{"Invalid", 0xFF, 0, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := ParseTransmitFailedCauseFromByte(tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("ParseTransmitFailedCauseFromByte() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("ParseTransmitFailedCauseFromByte() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("TransmitFailedCause_String", func(t *testing.T) {
		tests := []struct {
			name string
			code TransmitFailedCause
			want string
		}{
			{"CSMA_FAILED_CHANNEL_NOT_FREE", TransmitFailedCauseCSMA_FAILED_CHANNEL_NOT_FREE, "CSMA_FAILED_CHANNEL_NOT_FREE"},
			{"NO_ACK_RECEIVED", TransmitFailedCauseNO_ACK_RECEIVED, "NO_ACK_RECEIVED"},
			{"Unknown", 0x02, "UNKNOWN"},
			{"Unknown", 0xFF, "UNKNOWN"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.code.String(); got != tt.want {
					t.Errorf("TransmitFailedCause.String() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("TransmitFailedCause_Valid", func(t *testing.T) {
		// Test valid codes
		validCodes := []TransmitFailedCause{
			TransmitFailedCauseCSMA_FAILED_CHANNEL_NOT_FREE,
			TransmitFailedCauseNO_ACK_RECEIVED,
		}
		for _, code := range validCodes {
			if !code.Valid() {
				t.Errorf("TransmitFailedCause %v should be valid", code)
			}
		}

		// Test invalid codes
		invalidCodes := []TransmitFailedCause{0x02, 0x03, 0xFF}
		for _, code := range invalidCodes {
			if code.Valid() {
				t.Errorf("TransmitFailedCause %v should not be valid", code)
			}
		}
	})
}

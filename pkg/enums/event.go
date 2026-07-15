package enums

import "errors"

type EventCode byte

const (
	EventCodeSA_RECLAIM_NOT_SUCCESSFUL EventCode = iota + 1
	EventCodeSA_CONFIRM_LEARN
	EventCodeSA_LEARN_ACK
	EventCodeCO_READY
	EventCodeCO_EVENT_SECUREDEVICES
	EventCodeCO_DUTYCYCLE_LIMIT
	EventCodeCO_TRANSMIT_FAILED
	EventCodeCO_TX_DONE
	EventCodeCO_LRN_MODE_DISABLED
)

// ParseEventCodeFromByte parses a EventCode from a byte.
func ParseEventCodeFromByte(b byte) (EventCode, error) {
	code := EventCode(b)
	if !code.Valid() {
		return 0, errors.New("invalid event code")
	}
	return code, nil
}

// String returns the string representation of EventCode.
func (eventCode EventCode) String() string {
	switch eventCode {
	case EventCodeSA_RECLAIM_NOT_SUCCESSFUL:
		return "SA_RECLAIM_NOT_SUCCESSFUL"
	case EventCodeSA_CONFIRM_LEARN:
		return "SA_CONFIRM_LEARN"
	case EventCodeSA_LEARN_ACK:
		return "SA_LEARN_ACK"
	case EventCodeCO_READY:
		return "CO_READY"
	case EventCodeCO_EVENT_SECUREDEVICES:
		return "CO_EVENT_SECUREDEVICES"
	case EventCodeCO_DUTYCYCLE_LIMIT:
		return "CO_DUTYCYCLE_LIMIT"
	case EventCodeCO_TRANSMIT_FAILED:
		return "CO_TRANSMIT_FAILED"
	case EventCodeCO_TX_DONE:
		return "CO_TX_DONE"
	case EventCodeCO_LRN_MODE_DISABLED:
		return "CO_LRN_MODE_DISABLED"
	default:
		return "UNKNOWN"
	}
}

// Valid reports whether EventCode is valid.
func (eventCode EventCode) Valid() bool {
	switch eventCode {
	case EventCodeSA_RECLAIM_NOT_SUCCESSFUL,
		EventCodeSA_CONFIRM_LEARN,
		EventCodeSA_LEARN_ACK,
		EventCodeCO_READY,
		EventCodeCO_EVENT_SECUREDEVICES,
		EventCodeCO_DUTYCYCLE_LIMIT,
		EventCodeCO_TRANSMIT_FAILED,
		EventCodeCO_TX_DONE,
		EventCodeCO_LRN_MODE_DISABLED:
		return true
	default:
		return false
	}
}

type LearnAckConfirmCode byte

const (
	LearnAckConfirmCodeLRN_IN                 LearnAckConfirmCode = 0x00
	LearnAckConfirmCodeEEP_NOT_ACCEPTED       LearnAckConfirmCode = 0x11
	LearnAckConfirmCodeNO_PLACE_IN_PM         LearnAckConfirmCode = 0x12
	LearnAckConfirmCodeNO_PLACE_IN_CONTROLLER LearnAckConfirmCode = 0x13
	LearnAckConfirmCodeRSSI_NOT_GOOD_ENOUGH   LearnAckConfirmCode = 0x14
	LearnAckConfirmCodeLRN_OUT                LearnAckConfirmCode = 0x20
	LearnAckConfirmCodeFUNCTION_NOT_SUPPORTED LearnAckConfirmCode = 0xff
)

// ParseLearnAckConfirmCodeFromByte parses a LearnAckConfirmCode from a byte.
func ParseLearnAckConfirmCodeFromByte(b byte) (LearnAckConfirmCode, error) {
	code := LearnAckConfirmCode(b)
	if !code.Valid() {
		return 0, errors.New("invalid learn ack confirm code")
	}
	return code, nil
}

// String returns the string representation of LearnAckConfirmCode.
func (learnAckConfirmCode LearnAckConfirmCode) String() string {
	switch learnAckConfirmCode {
	case LearnAckConfirmCodeLRN_IN:
		return "LRN_IN"
	case LearnAckConfirmCodeEEP_NOT_ACCEPTED:
		return "EEP_NOT_ACCEPTED"
	case LearnAckConfirmCodeNO_PLACE_IN_PM:
		return "NO_PLACE_IN_PM"
	case LearnAckConfirmCodeNO_PLACE_IN_CONTROLLER:
		return "NO_PLACE_IN_CONTROLLER"
	case LearnAckConfirmCodeRSSI_NOT_GOOD_ENOUGH:
		return "RSSI_NOT_GOOD_ENOUGH"
	case LearnAckConfirmCodeLRN_OUT:
		return "LRN_OUT"
	case LearnAckConfirmCodeFUNCTION_NOT_SUPPORTED:
		return "FUNCTION_NOT_SUPPORTED"
	default:
		return "UNKNOWN"
	}
}

// Valid reports whether LearnAckConfirmCode is valid.
func (learnAckConfirmCode LearnAckConfirmCode) Valid() bool {
	switch learnAckConfirmCode {
	case LearnAckConfirmCodeLRN_IN,
		LearnAckConfirmCodeEEP_NOT_ACCEPTED,
		LearnAckConfirmCodeNO_PLACE_IN_PM,
		LearnAckConfirmCodeNO_PLACE_IN_CONTROLLER,
		LearnAckConfirmCodeRSSI_NOT_GOOD_ENOUGH,
		LearnAckConfirmCodeLRN_OUT,
		LearnAckConfirmCodeFUNCTION_NOT_SUPPORTED:
		return true
	default:
		return false
	}
}

type WakeUpCause byte

const (
	WakeUpCauseVOLTAGE_SUPPLY_DROP WakeUpCause = iota
	WakeUpCauseRESET_BY_RESET_PIN
	WakeUpCauseWATCHDOG_TIMEOUT
	WakeUpCauseFLYWHEEL_TIMEOUT
	WakeUpCausePARITY_ERROR
	WakeUpCauseHARDWARE_PARITY_ERROR_IN_MEMORY
	WakeUpCauseREQUESTED_MEMORY_LOCATION_NOT_FOUND
	WakeUpCauseTRIGGER_PIN0
	WakeUpCauseTRIGGER_PIN1
	WakeUpCauseUNKNOWN_RESET_SOURCE
	WakeUpCauseUART_WAKE_UP
)

// ParseWakeUpCauseFromByte parses a WakeUpCause from a byte.
func ParseWakeUpCauseFromByte(b byte) (WakeUpCause, error) {
	cause := WakeUpCause(b)
	if !cause.Valid() {
		return 0, errors.New("invalid wake up cause")
	}
	return cause, nil
}

// String returns the string representation of WakeUpCause.
func (wakeUpCause WakeUpCause) String() string {
	switch wakeUpCause {
	case WakeUpCauseVOLTAGE_SUPPLY_DROP:
		return "VOLTAGE_SUPPLY_DROP"
	case WakeUpCauseRESET_BY_RESET_PIN:
		return "RESET_BY_RESET_PIN"
	case WakeUpCauseWATCHDOG_TIMEOUT:
		return "WATCHDOG_TIMEOUT"
	case WakeUpCauseFLYWHEEL_TIMEOUT:
		return "FLYWHEEL_TIMEOUT"
	case WakeUpCausePARITY_ERROR:
		return "PARITY_ERROR"
	case WakeUpCauseHARDWARE_PARITY_ERROR_IN_MEMORY:
		return "HARDWARE_PARITY_ERROR_IN_MEMORY"
	case WakeUpCauseREQUESTED_MEMORY_LOCATION_NOT_FOUND:
		return "REQUESTED_MEMORY_LOCATION_NOT_FOUND"
	case WakeUpCauseTRIGGER_PIN0:
		return "TRIGGER_PIN0"
	case WakeUpCauseTRIGGER_PIN1:
		return "TRIGGER_PIN1"
	case WakeUpCauseUNKNOWN_RESET_SOURCE:
		return "UNKNOWN_RESET_SOURCE"
	case WakeUpCauseUART_WAKE_UP:
		return "UART_WAKE_UP"
	default:
		return "UNKNOWN"
	}
}

// Valid reports whether WakeUpCause is valid.
func (wakeUpCause WakeUpCause) Valid() bool {
	switch wakeUpCause {
	case WakeUpCauseVOLTAGE_SUPPLY_DROP,
		WakeUpCauseRESET_BY_RESET_PIN,
		WakeUpCauseWATCHDOG_TIMEOUT,
		WakeUpCauseFLYWHEEL_TIMEOUT,
		WakeUpCausePARITY_ERROR,
		WakeUpCauseHARDWARE_PARITY_ERROR_IN_MEMORY,
		WakeUpCauseREQUESTED_MEMORY_LOCATION_NOT_FOUND,
		WakeUpCauseTRIGGER_PIN0,
		WakeUpCauseTRIGGER_PIN1,
		WakeUpCauseUNKNOWN_RESET_SOURCE,
		WakeUpCauseUART_WAKE_UP:
		return true
	default:
		return false
	}
}

type WakeUpMode byte

const (
	WakeUpModeSTANDARD_SECURITY WakeUpMode = iota
	WakeUpModeEXTENDED_SECURITY
)

// ParseWakeUpModeFromByte parses a WakeUpMode from a byte.
func ParseWakeUpModeFromByte(b byte) (WakeUpMode, error) {
	mode := WakeUpMode(b)
	if !mode.Valid() {
		return 0, errors.New("invalid wake up mode")
	}
	return mode, nil
}

// String returns the string representation of WakeUpMode.
func (wakeUpMode WakeUpMode) String() string {
	switch wakeUpMode {
	case WakeUpModeSTANDARD_SECURITY:
		return "STANDARD_SECURITY"
	case WakeUpModeEXTENDED_SECURITY:
		return "EXTENDED_SECURITY"
	default:
		return "UNKNOWN"
	}
}

// Valid reports whether WakeUpMode is valid.
func (wakeUpMode WakeUpMode) Valid() bool {
	switch wakeUpMode {
	case WakeUpModeSTANDARD_SECURITY,
		WakeUpModeEXTENDED_SECURITY:
		return true
	default:
		return false
	}
}

type SecureDeviceEventCause byte

const (
	COEventSecureSECURE_LINK_TABLE_FULL   SecureDeviceEventCause = iota
	COEventSecureRESYNC_WRONG_PRIVATE_KEY SecureDeviceEventCause = iota + 1
	COEventSecureWRONG_CMAC_TELEGRAM_THRESHOLD_HIT
	COEventSecureTEACH_IN_TELEGRAM_CORRUPTED
	COEventSecureTEACH_IN_PSK_FAILED_NO_PSK_SET
	COEventSecureTEACH_IN_PSK_FAILED_NO_PSK_GIVEN
	COEventSecureINCORRECT_CMAC_OR_RKC
	COEventSecureRECEIVED_TELEGRAM_FROM_DEVICE_IN_SECURE_LINK_TABLE
	COEventSecureTEACH_IN_SUCCESSFUL
	COEventSecureVALID_RLC_SYNC_RECEIVED_VIA_TEACH_IN
)

// ParseCOEventSecureFromByte parses a SecureDeviceEventCause from a byte.
func ParseCOEventSecureFromByte(b byte) (SecureDeviceEventCause, error) {
	switch b {
	case 0x00:
		return COEventSecureSECURE_LINK_TABLE_FULL, nil
	case 0x01:
		return COEventSecureRESYNC_WRONG_PRIVATE_KEY, nil
	case 0x02:
		return COEventSecureWRONG_CMAC_TELEGRAM_THRESHOLD_HIT, nil
	case 0x03:
		return COEventSecureTEACH_IN_TELEGRAM_CORRUPTED, nil
	case 0x04:
		return COEventSecureTEACH_IN_PSK_FAILED_NO_PSK_SET, nil
	case 0x05:
		return COEventSecureTEACH_IN_PSK_FAILED_NO_PSK_GIVEN, nil
	case 0x06:
		return COEventSecureINCORRECT_CMAC_OR_RKC, nil
	case 0x07:
		return COEventSecureRECEIVED_TELEGRAM_FROM_DEVICE_IN_SECURE_LINK_TABLE, nil
	case 0x08:
		return COEventSecureTEACH_IN_SUCCESSFUL, nil
	case 0x09:
		return COEventSecureVALID_RLC_SYNC_RECEIVED_VIA_TEACH_IN, nil
	default:
		return 0, errors.New("invalid co event secure")
	}
}

// String returns the string representation of SecureDeviceEventCause.
func (coEventSecure SecureDeviceEventCause) String() string {
	switch coEventSecure {
	case COEventSecureSECURE_LINK_TABLE_FULL:
		return "SECURE_LINK_TABLE_FULL"
	case COEventSecureRESYNC_WRONG_PRIVATE_KEY:
		return "RESYNC_WRONG_PRIVATE_KEY"
	case COEventSecureWRONG_CMAC_TELEGRAM_THRESHOLD_HIT:
		return "WRONG_CMAC_TELEGRAM_THRESHOLD_HIT"
	case COEventSecureTEACH_IN_TELEGRAM_CORRUPTED:
		return "TEACH_IN_TELEGRAM_CORRUPTED"
	case COEventSecureTEACH_IN_PSK_FAILED_NO_PSK_SET:
		return "TEACH_IN_PSK_FAILED_NO_PSK_SET"
	case COEventSecureTEACH_IN_PSK_FAILED_NO_PSK_GIVEN:
		return "TEACH_IN_PSK_FAILED_NO_PSK_GIVEN"
	case COEventSecureINCORRECT_CMAC_OR_RKC:
		return "INCORRECT_CMAC_OR_RKC"
	case COEventSecureRECEIVED_TELEGRAM_FROM_DEVICE_IN_SECURE_LINK_TABLE:
		return "RECEIVED_TELEGRAM_FROM_DEVICE_IN_SECURE_LINK_TABLE"
	case COEventSecureTEACH_IN_SUCCESSFUL:
		return "TEACH_IN_SUCCESSFUL"
	case COEventSecureVALID_RLC_SYNC_RECEIVED_VIA_TEACH_IN:
		return "VALID_RLC_SYNC_RECEIVED_VIA_TEACH_IN"
	default:
		return "UNKNOWN"
	}
}

// Valid reports whether SecureDeviceEventCause is valid.
func (coEventSecure SecureDeviceEventCause) Valid() bool {
	switch coEventSecure {
	case COEventSecureSECURE_LINK_TABLE_FULL,
		COEventSecureRESYNC_WRONG_PRIVATE_KEY,
		COEventSecureWRONG_CMAC_TELEGRAM_THRESHOLD_HIT,
		COEventSecureTEACH_IN_TELEGRAM_CORRUPTED,
		COEventSecureTEACH_IN_PSK_FAILED_NO_PSK_SET,
		COEventSecureTEACH_IN_PSK_FAILED_NO_PSK_GIVEN,
		COEventSecureINCORRECT_CMAC_OR_RKC,
		COEventSecureRECEIVED_TELEGRAM_FROM_DEVICE_IN_SECURE_LINK_TABLE,
		COEventSecureTEACH_IN_SUCCESSFUL,
		COEventSecureVALID_RLC_SYNC_RECEIVED_VIA_TEACH_IN:
		return true
	default:
		return false
	}
}

type DutyCycleLimitCause byte

const (
	DutyCycleLimitCauseNOT_YET_REACHED DutyCycleLimitCause = iota
	DutyCycleLimitCauseREACHED
)

// ParseDutyCycleLimitCauseFromByte parses a DutyCycleLimitCause from a byte.
func ParseDutyCycleLimitCauseFromByte(b byte) (DutyCycleLimitCause, error) {
	cause := DutyCycleLimitCause(b)
	if !cause.Valid() {
		return 0, errors.New("invalid duty cycle limit cause")
	}
	return cause, nil
}

// String returns the string representation of DutyCycleLimitCause.
func (dutyCycleLimitCause DutyCycleLimitCause) String() string {
	switch dutyCycleLimitCause {
	case DutyCycleLimitCauseNOT_YET_REACHED:
		return "NOT_YET_REACHED"
	case DutyCycleLimitCauseREACHED:
		return "REACHED"
	default:
		return "UNKNOWN"
	}
}

// Valid reports whether DutyCycleLimitCause is valid.
func (dutyCycleLimitCause DutyCycleLimitCause) Valid() bool {
	switch dutyCycleLimitCause {
	case DutyCycleLimitCauseNOT_YET_REACHED,
		DutyCycleLimitCauseREACHED:
		return true
	default:
		return false
	}
}

type TransmitFailedCause byte

const (
	TransmitFailedCauseCSMA_FAILED_CHANNEL_NOT_FREE TransmitFailedCause = iota
	TransmitFailedCauseNO_ACK_RECEIVED
)

// ParseTransmitFailedCauseFromByte parses a TransmitFailedCause from a byte.
func ParseTransmitFailedCauseFromByte(b byte) (TransmitFailedCause, error) {
	cause := TransmitFailedCause(b)
	if !cause.Valid() {
		return 0, errors.New("invalid transmit failed cause")
	}
	return cause, nil
}

// String returns the string representation of TransmitFailedCause.
func (transmitFailedCause TransmitFailedCause) String() string {
	switch transmitFailedCause {
	case TransmitFailedCauseCSMA_FAILED_CHANNEL_NOT_FREE:
		return "CSMA_FAILED_CHANNEL_NOT_FREE"
	case TransmitFailedCauseNO_ACK_RECEIVED:
		return "NO_ACK_RECEIVED"
	default:
		return "UNKNOWN"
	}
}

// Valid reports whether TransmitFailedCause is valid.
func (transmitFailedCause TransmitFailedCause) Valid() bool {
	switch transmitFailedCause {
	case TransmitFailedCauseCSMA_FAILED_CHANNEL_NOT_FREE,
		TransmitFailedCauseNO_ACK_RECEIVED:
		return true
	default:
		return false
	}
}

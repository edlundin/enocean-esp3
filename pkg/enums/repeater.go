package enums

import "errors"

type RepeaterMode byte

const (
	RepeaterModeOFF RepeaterMode = iota
	RepeaterModeON
	RepeaterModeSELECTIVE
)

// ParseRepeaterModeFromByte parses a RepeaterMode from a byte.
func ParseRepeaterModeFromByte(b byte) (RepeaterMode, error) {
	switch b {
	case 0x00:
		return RepeaterModeOFF, nil
	case 0x01:
		return RepeaterModeON, nil
	case 0x02:
		return RepeaterModeSELECTIVE, nil
	default:
		return 0, errors.New("invalid repeater mode")
	}
}

// String returns the string representation of RepeaterMode.
func (repeaterMode RepeaterMode) String() string {
	switch repeaterMode {
	case RepeaterModeOFF:
		return "OFF"
	case RepeaterModeON:
		return "ON"
	case RepeaterModeSELECTIVE:
		return "SELECTIVE"
	default:
		return "UNKNOWN"
	}
}

// Valid reports whether RepeaterMode is valid.
func (repeaterMode RepeaterMode) Valid() bool {
	switch repeaterMode {
	case RepeaterModeOFF,
		RepeaterModeON,
		RepeaterModeSELECTIVE:
		return true
	default:
		return false
	}
}

type RepeaterLevel byte

const (
	RepeaterLevelNO_REPETITION RepeaterLevel = iota
	RepeaterLevel1_REPETITION
	RepeaterLevel2_REPETITION
)

// ParseRepeaterLevelFromByte parses a RepeaterLevel from a byte.
func ParseRepeaterLevelFromByte(b byte) (RepeaterLevel, error) {
	switch b {
	case 0x00:
		return RepeaterLevelNO_REPETITION, nil
	case 0x01:
		return RepeaterLevel1_REPETITION, nil
	case 0x02:
		return RepeaterLevel2_REPETITION, nil
	default:
		return 0, errors.New("invalid repeater mode")
	}
}

// String returns the string representation of RepeaterLevel.
func (repeaterLevel RepeaterLevel) String() string {
	switch repeaterLevel {
	case RepeaterLevelNO_REPETITION:
		return "NO_REPEATING"
	case RepeaterLevel1_REPETITION:
		return "1_REPEAT"
	case RepeaterLevel2_REPETITION:
		return "2_REPEAT"
	default:
		return "UNKNOWN"
	}
}

// Valid reports whether RepeaterLevel is valid.
func (repeaterLevel RepeaterLevel) Valid() bool {
	switch repeaterLevel {
	case RepeaterLevelNO_REPETITION,
		RepeaterLevel1_REPETITION,
		RepeaterLevel2_REPETITION:
		return true
	default:
		return false
	}
}

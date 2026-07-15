package enums

import "errors"

type RadioMode byte

const (
	RadioModeERP1 RadioMode = iota
	RadioModeERP2
)

// ParseRadioModeFromByte parses a RadioMode from a byte.
func ParseRadioModeFromByte(b byte) (RadioMode, error) {
	switch b {
	case 0x00:
		return RadioModeERP1, nil
	case 0x01:
		return RadioModeERP2, nil
	default:
		return 0, errors.New("invalid radio mode")
	}
}

// String returns the string representation of RadioMode.
func (radioMode RadioMode) String() string {
	switch radioMode {
	case RadioModeERP1:
		return "ERP1"
	case RadioModeERP2:
		return "ERP2"
	default:
		return "UNKNOWN"
	}
}

// Valid reports whether RadioMode is valid.
func (radioMode RadioMode) Valid() bool {
	switch radioMode {
	case RadioModeERP1, RadioModeERP2:
		return true
	default:
		return false
	}
}

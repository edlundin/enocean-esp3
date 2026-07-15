package eep

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/edlundin/enocean-esp3/pkg/enums"
)

const (
	minRorg = 0x00
	maxRorg = 0xff
	minFunc = 0x00
	maxFunc = 0xb0
	minType = 0x00
	maxType = 0x7f
)

type EEP struct {
	Rorg enums.Rorg
	Func byte
	Type byte
}

// FromTriplet constructs a value from Triplet.
func FromTriplet(eepRorg enums.Rorg, eepFunc byte, eepType byte) (EEP, error) {
	if eepFunc < minFunc || eepFunc > maxFunc {
		return EEP{}, errors.New("invalid FUNC: out of bounds")
	}

	if eepType < minType || eepType > maxType {
		return EEP{}, errors.New("invalid TYPE: out of bounds")
	}

	return EEP{
		Rorg: eepRorg,
		Func: eepFunc,
		Type: eepType,
	}, nil
}

// FromString constructs an EEP from an RR-FF-TT string, where each component
// is a one- or two-digit hexadecimal RORG, FUNC, or TYPE value.
func FromString(str string) (EEP, error) {
	const (
		strFieldLen = 3
		rorgIndex   = 0
		funcIndex   = 1
		typeIndex   = 2
	)

	strFields := strings.Split(str, "-")

	if len(strFields) != strFieldLen {
		return EEP{}, errors.New("invalid format (RR-FF-TT)")
	}

	eepRorg, err := strconv.ParseInt(strFields[rorgIndex], 16, 32)
	if err != nil {
		return EEP{}, errors.Join(errors.New("invalid RORG"), err)
	}

	if eepRorg < minRorg || eepRorg > maxRorg {
		return EEP{}, errors.New("invalid RORG: out of bounds")
	}

	eepFunc, err := strconv.ParseInt(strFields[funcIndex], 16, 32)
	if err != nil {
		return EEP{}, errors.Join(errors.New("invalid FUNC"), err)
	}

	if eepFunc < minFunc || eepFunc > maxFunc {
		return EEP{}, errors.New("invalid FUNC: out of bounds")
	}

	eepType, err := strconv.ParseInt(strFields[typeIndex], 16, 32)
	if err != nil {
		return EEP{}, errors.Join(errors.New("invalid TYPE"), err)
	}

	if eepType < minType || eepType > maxType {
		return EEP{}, errors.New("invalid TYPE: out of bounds")
	}

	return EEP{
		Rorg: enums.Rorg(eepRorg),
		Func: byte(eepFunc),
		Type: byte(eepType),
	}, nil
}

// String returns the string representation of EEP.
func (eep EEP) String() string {
	return fmt.Sprintf("%02X-%02X-%02X", byte(eep.Rorg), eep.Func, eep.Type)
}

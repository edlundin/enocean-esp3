package enums

import "errors"

type Maturity byte

const (
	MaturityFORWARDED_IMMEDIATELY Maturity = iota
	MaturityFORWARDED_ON_TIMEOUT
	MaturityFORWARD_SUBTELEGRAMS
)

func ParseMaturityFromByte(b byte) (Maturity, error) {
	switch b {
	case 0x00:
		return MaturityFORWARDED_IMMEDIATELY, nil
	case 0x01:
		return MaturityFORWARDED_ON_TIMEOUT, nil
	case 0x02:
		return MaturityFORWARD_SUBTELEGRAMS, nil
	default:
		return 0, errors.New("invalid maturity")
	}
}

func (maturity Maturity) String() string {
	switch maturity {
	case MaturityFORWARDED_IMMEDIATELY:
		return "FORWARDED_IMMEDIATELY"
	case MaturityFORWARDED_ON_TIMEOUT:
		return "FORWARDED_ON_TIMEOUT"
	case MaturityFORWARD_SUBTELEGRAMS:
		return "FORWARD_SUBTELEGRAMS"
	default:
		return "UNKNOWN"
	}
}

func (maturity Maturity) Valid() bool {
	switch maturity {
	case MaturityFORWARDED_IMMEDIATELY,
		MaturityFORWARDED_ON_TIMEOUT,
		MaturityFORWARD_SUBTELEGRAMS:
		return true
	default:
		return false
	}
}

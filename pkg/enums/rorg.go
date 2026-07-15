package enums

import "errors"

type Rorg byte

const (
	RorgRPS        Rorg = 0xf6
	Rorg1BS        Rorg = 0xd5
	Rorg4BS        Rorg = 0xa5
	RorgVLD        Rorg = 0xd2
	RorgMSC        Rorg = 0xd1
	RorgADT        Rorg = 0xa6
	RorgGP_TI      Rorg = 0xb0
	RorgGP_TR      Rorg = 0xb1
	RorgGP_CD      Rorg = 0xb2
	RorgGP_SD      Rorg = 0xb3
	RorgSM_LRN_REQ Rorg = 0xc6
	RorgSM_LRN_ANS Rorg = 0xc7
	RorgSM_REC     Rorg = 0xa7
	RorgSYS_EX     Rorg = 0xc5
	RorgSEC        Rorg = 0x30
	RorgSEC_R      Rorg = 0x31
	RorgSEC_ENCAPS Rorg = RorgSEC_R
	RorgSEC_D      Rorg = 0x32
	RorgSEC_CDM    Rorg = 0x33
	RorgSEC_MAN    Rorg = 0x34
	RorgSEC_TI     Rorg = 0x35
	RorgSIGNAL     Rorg = 0xd0
	RorgUTE        Rorg = 0xd4
)

// ParseRorgFromByte parses a Rorg from a byte.
func ParseRorgFromByte(b byte) (Rorg, error) {
	switch b {
	case 0xf6:
		return RorgRPS, nil
	case 0xd5:
		return Rorg1BS, nil
	case 0xa5:
		return Rorg4BS, nil
	case 0xd2:
		return RorgVLD, nil
	case 0xd1:
		return RorgMSC, nil
	case 0xa6:
		return RorgADT, nil
	case 0xb0:
		return RorgGP_TI, nil
	case 0xb1:
		return RorgGP_TR, nil
	case 0xb2:
		return RorgGP_CD, nil
	case 0xb3:
		return RorgGP_SD, nil
	case 0xc6:
		return RorgSM_LRN_REQ, nil
	case 0xc7:
		return RorgSM_LRN_ANS, nil
	case 0xa7:
		return RorgSM_REC, nil
	case 0xc5:
		return RorgSYS_EX, nil
	case 0x30:
		return RorgSEC, nil
	case 0x31:
		return RorgSEC_R, nil
	case 0x32:
		return RorgSEC_D, nil
	case 0x33:
		return RorgSEC_CDM, nil
	case 0x34:
		return RorgSEC_MAN, nil
	case 0x35:
		return RorgSEC_TI, nil
	case 0xd0:
		return RorgSIGNAL, nil
	case 0xd4:
		return RorgUTE, nil
	default:
		return 0, errors.New("invalid rorg")
	}
}

// String returns the string representation of Rorg.
func (rorg Rorg) String() string {
	switch rorg {
	case RorgRPS:
		return "RPS"
	case Rorg1BS:
		return "1BS"
	case Rorg4BS:
		return "4BS"
	case RorgVLD:
		return "VLD"
	case RorgMSC:
		return "MSC"
	case RorgADT:
		return "ADT"
	case RorgGP_TI:
		return "GP_TI"
	case RorgGP_TR:
		return "GP_TR"
	case RorgGP_CD:
		return "GP_CD"
	case RorgGP_SD:
		return "GP_SD"
	case RorgSM_LRN_REQ:
		return "SM_LRN_REQ"
	case RorgSM_LRN_ANS:
		return "SM_LRN_ANS"
	case RorgSM_REC:
		return "SM_REC"
	case RorgSYS_EX:
		return "SYS_EX"
	case RorgSEC:
		return "SEC"
	case RorgSEC_R:
		return "SEC_R"
	case RorgSEC_D:
		return "SEC_D"
	case RorgSEC_CDM:
		return "SEC_CDM"
	case RorgSEC_MAN:
		return "SEC_MAN"
	case RorgSEC_TI:
		return "SEC_TI"
	case RorgSIGNAL:
		return "SIGNAL"
	case RorgUTE:
		return "UTE"
	default:
		return "UNKNOWN"
	}
}

// Valid reports whether Rorg is valid.
func (rorg Rorg) Valid() bool {
	switch rorg {
	case RorgRPS,
		Rorg1BS,
		Rorg4BS,
		RorgVLD,
		RorgMSC,
		RorgADT,
		RorgGP_TI,
		RorgGP_TR,
		RorgGP_CD,
		RorgGP_SD,
		RorgSM_LRN_REQ,
		RorgSM_LRN_ANS,
		RorgSM_REC,
		RorgSYS_EX,
		RorgSEC,
		RorgSEC_R,
		RorgSEC_D,
		RorgSEC_CDM,
		RorgSEC_MAN,
		RorgSEC_TI,
		RorgSIGNAL,
		RorgUTE:
		return true
	default:
		return false
	}
}

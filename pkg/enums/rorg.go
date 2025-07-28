package enums

import "errors"

type Rorg uint8

const (
	RORG_RPS        Rorg = 0xf6
	RORG_1BS        Rorg = 0xd5
	RORG_4BS        Rorg = 0xa5
	RORG_VLD        Rorg = 0xd2
	RORG_MSC        Rorg = 0xd1
	RORG_ADT        Rorg = 0xa6
	RORG_SM_LRN_REQ Rorg = 0xc6
	RORG_SM_LRN_ANS Rorg = 0xc7
	RORG_SM_REC     Rorg = 0xa7
	RORG_SYS_EX     Rorg = 0xc5
	RORG_SEC        Rorg = 0x30
	RORG_SEC_ENCAPS Rorg = 0x31
	RORG_SEC_MAN    Rorg = 0x34
	RORG_SIGNAL     Rorg = 0xd0
	RORG_UTE        Rorg = 0xd4
)

func ParseRorgFromByte(byte uint8) (Rorg, error) {
	switch byte {
	case 0xf6:
		return RORG_RPS, nil
	case 0xd5:
		return RORG_1BS, nil
	case 0xa5:
		return RORG_4BS, nil
	case 0xd2:
		return RORG_VLD, nil
	case 0xd1:
		return RORG_MSC, nil
	case 0xa6:
		return RORG_ADT, nil
	case 0xc6:
		return RORG_SM_LRN_REQ, nil
	case 0xc7:
		return RORG_SM_LRN_ANS, nil
	case 0xa7:
		return RORG_SM_REC, nil
	case 0xc5:
		return RORG_SYS_EX, nil
	case 0x30:
		return RORG_SEC, nil
	case 0x31:
		return RORG_SEC_ENCAPS, nil
	case 0x34:
		return RORG_SEC_MAN, nil
	case 0xd0:
		return RORG_SIGNAL, nil
	case 0xd4:
		return RORG_UTE, nil
	default:
		return 0, errors.New("invalid rorg")
	}
}

func (rorg Rorg) String() string {
	switch rorg {
	case RORG_RPS:
		return "RPS"
	case RORG_1BS:
		return "1BS"
	case RORG_4BS:
		return "4BS"
	case RORG_VLD:
		return "VLD"
	case RORG_MSC:
		return "MSC"
	case RORG_ADT:
		return "ADT"
	case RORG_SM_LRN_REQ:
		return "SM_LRN_REQ"
	case RORG_SM_LRN_ANS:
		return "SM_LRN_ANS"
	case RORG_SM_REC:
		return "SM_REC"
	case RORG_SYS_EX:
		return "SYS_EX"
	case RORG_SEC:
		return "SEC"
	case RORG_SEC_ENCAPS:
		return "SEC_ENCAPS"
	case RORG_SEC_MAN:
		return "SEC_MAN"
	case RORG_SIGNAL:
		return "SIGNAL"
	case RORG_UTE:
		return "UTE"
	default:
		return "UNKNOWN"
	}
}

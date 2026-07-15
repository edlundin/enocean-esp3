package enums

import "errors"

type ReturnCode byte

const (
	ReturnCodeSUCCESS ReturnCode = iota
	ReturnCodeERROR
	ReturnCodeNOT_SUPPORTED
	ReturnCodeWRONG_ARGUMENT
	ReturnCodeOPERATION_DENIED
	ReturnCodeLOCK_SET
	ReturnCodeBUFFER_TO_SMALL
	ReturnCodeNO_FREE_BUFFER
	ReturnCodeBASEID_OUT_OF_RANGE ReturnCode = 0x90
	ReturnCodeBASEID_MAX_REACHED  ReturnCode = 0x91
)

// ParseReturnCodeFromByte parses a ReturnCode from a byte.
func ParseReturnCodeFromByte(b byte) (ReturnCode, error) {
	switch b {
	case 0x00:
		return ReturnCodeSUCCESS, nil
	case 0x01:
		return ReturnCodeERROR, nil
	case 0x02:
		return ReturnCodeNOT_SUPPORTED, nil
	case 0x03:
		return ReturnCodeWRONG_ARGUMENT, nil
	case 0x04:
		return ReturnCodeOPERATION_DENIED, nil
	case 0x05:
		return ReturnCodeLOCK_SET, nil
	case 0x06:
		return ReturnCodeBUFFER_TO_SMALL, nil
	case 0x07:
		return ReturnCodeNO_FREE_BUFFER, nil
	case 0x90:
		return ReturnCodeBASEID_OUT_OF_RANGE, nil
	case 0x91:
		return ReturnCodeBASEID_MAX_REACHED, nil
	default:
		return ReturnCodeERROR, errors.New("invalid return code")
	}
}

// String returns the string representation of ReturnCode.
func (returnCode ReturnCode) String() string {
	switch returnCode {
	case ReturnCodeSUCCESS:
		return "SUCCESS"
	case ReturnCodeERROR:
		return "ERROR"
	case ReturnCodeNOT_SUPPORTED:
		return "NOT_SUPPORTED"
	case ReturnCodeWRONG_ARGUMENT:
		return "WRONG_ARGUMENT"
	case ReturnCodeOPERATION_DENIED:
		return "OPERATION_DENIED"
	case ReturnCodeLOCK_SET:
		return "LOCK_SET"
	case ReturnCodeBUFFER_TO_SMALL:
		return "BUFFER_TO_SMALL"
	case ReturnCodeNO_FREE_BUFFER:
		return "NO_FREE_BUFFER"
	case ReturnCodeBASEID_OUT_OF_RANGE:
		return "BASEID_OUT_OF_RANGE"
	case ReturnCodeBASEID_MAX_REACHED:
		return "BASEID_MAX_REACHED"
	default:
		return "UNKNOWN"
	}
}

// Valid reports whether ReturnCode is valid.
func (returnCode ReturnCode) Valid() bool {
	switch returnCode {
	case ReturnCodeSUCCESS,
		ReturnCodeERROR,
		ReturnCodeNOT_SUPPORTED,
		ReturnCodeWRONG_ARGUMENT,
		ReturnCodeOPERATION_DENIED,
		ReturnCodeLOCK_SET,
		ReturnCodeBUFFER_TO_SMALL,
		ReturnCodeNO_FREE_BUFFER,
		ReturnCodeBASEID_OUT_OF_RANGE,
		ReturnCodeBASEID_MAX_REACHED:
		return true
	default:
		return false
	}
}

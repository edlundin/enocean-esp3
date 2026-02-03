package enums

import "errors"

type SmartAckCommand uint8

const (
	SmartAckCommandWR_LEARN_MODE SmartAckCommand = iota + 1
	SmartAckCommandRD_LEARN_MODE
	SmartAckCommandWR_LEARN_CONFIRM
	SmartAckCommandWR_CLIENT_LEARN_RQ
	SmartAckCommandWR_RESET
	SmartAckCommandWR_RD_LEARNED_CLIENTS
	SmartAckCommandWR_RECLAIMS
	SmartAckCommandWR_WR_POSTMASTER
)

func ParseSmartAckCommandFromByte(b byte) (SmartAckCommand, error) {
	switch b {
	case 0x01:
		return SmartAckCommandWR_LEARN_MODE, nil
	case 0x02:
		return SmartAckCommandRD_LEARN_MODE, nil
	case 0x03:
		return SmartAckCommandWR_LEARN_CONFIRM, nil
	case 0x04:
		return SmartAckCommandWR_CLIENT_LEARN_RQ, nil
	case 0x05:
		return SmartAckCommandWR_RESET, nil
	case 0x06:
		return SmartAckCommandWR_RD_LEARNED_CLIENTS, nil
	case 0x07:
		return SmartAckCommandWR_RECLAIMS, nil
	case 0x08:
		return SmartAckCommandWR_WR_POSTMASTER, nil
	default:
		return 0, errors.New("invalid smart ack command")
	}
}

func (command SmartAckCommand) String() string {
	switch command {
	case SmartAckCommandWR_LEARN_MODE:
		return "WR_LEARN_MODE"
	case SmartAckCommandRD_LEARN_MODE:
		return "RD_LEARN_MODE"
	case SmartAckCommandWR_LEARN_CONFIRM:
		return "WR_LEARN_CONFIRM"
	case SmartAckCommandWR_CLIENT_LEARN_RQ:
		return "WR_CLIENT_LEARN_RQ"
	case SmartAckCommandWR_RESET:
		return "WR_RESET"
	case SmartAckCommandWR_RD_LEARNED_CLIENTS:
		return "WR_RD_LEARNED_CLIENTS"
	case SmartAckCommandWR_RECLAIMS:
		return "WR_RECLAIMS"
	case SmartAckCommandWR_WR_POSTMASTER:
		return "WR_WR_POSTMASTER"
	default:
		return "UNKNOWN"
	}
}

func (command SmartAckCommand) Valid() bool {
	switch command {
	case SmartAckCommandWR_LEARN_MODE,
		SmartAckCommandRD_LEARN_MODE,
		SmartAckCommandWR_LEARN_CONFIRM,
		SmartAckCommandWR_CLIENT_LEARN_RQ,
		SmartAckCommandWR_RESET,
		SmartAckCommandWR_RD_LEARNED_CLIENTS,
		SmartAckCommandWR_RECLAIMS,
		SmartAckCommandWR_WR_POSTMASTER:
		return true
	default:
		return false
	}
}

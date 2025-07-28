package pkg

import (
	"errors"

	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
)

type Esp3ReturnCode uint8

const (
	Esp3ReturnCodeOK               Esp3ReturnCode = 0x00
	Esp3ReturnCodeERROR            Esp3ReturnCode = 0x01
	Esp3ReturnCodeNOT_SUPPORTED    Esp3ReturnCode = 0x02
	Esp3ReturnCodeWRONG_PARAM      Esp3ReturnCode = 0x03
	Esp3ReturnCodeOPERATION_DENIED Esp3ReturnCode = 0x04
	Esp3ReturnCodeLOCK_SET         Esp3ReturnCode = 0x05
	Esp3ReturnCodeBUFFER_TO_SMALL  Esp3ReturnCode = 0x06
	Esp3ReturnCodeNO_FREE_BUFFER   Esp3ReturnCode = 0x07
)

func (i Esp3ReturnCode) Valid() bool {
	switch i {
	case Esp3ReturnCodeOK,
		Esp3ReturnCodeERROR,
		Esp3ReturnCodeNOT_SUPPORTED,
		Esp3ReturnCodeWRONG_PARAM,
		Esp3ReturnCodeOPERATION_DENIED,
		Esp3ReturnCodeLOCK_SET,
		Esp3ReturnCodeBUFFER_TO_SMALL,
		Esp3ReturnCodeNO_FREE_BUFFER:
		return true
	default:
		return false
	}
}

func (i Esp3ReturnCode) String() string {
	switch i {
	case Esp3ReturnCodeOK:
		return "RET_OK"
	case Esp3ReturnCodeERROR:
		return "RET_ERROR"
	case Esp3ReturnCodeNOT_SUPPORTED:
		return "RET_NOT_SUPPORTED"
	case Esp3ReturnCodeWRONG_PARAM:
		return "RET_WRONG_PARAM"
	case Esp3ReturnCodeOPERATION_DENIED:
		return "RET_OPERATION_DENIED"
	case Esp3ReturnCodeLOCK_SET:
		return "RET_LOCK_SET"
	case Esp3ReturnCodeBUFFER_TO_SMALL:
		return "RET_BUFFER_TO_SMALL"
	case Esp3ReturnCodeNO_FREE_BUFFER:
		return "RET_NO_FREE_BUFFER"
	default:
		return "UNKNOWN"
	}
}

type ResponsePacket struct {
	ReturnCode   Esp3ReturnCode
	OptionalData []byte
}

type EventSAConfirmLearnResponse struct {
	ReturnCode   Esp3ReturnCode
	ResponseTime uint16
	ConfirmCode  uint8
}

func NewResponsePacketFromEsp3(telegram esp3.Esp3Telegram) (ResponsePacket, error) {
	if telegram.PacketType != enums.PACKET_TYPE_RESPONSE {
		return ResponsePacket{}, errors.New("invalid packet type")
	}

	if len(telegram.Data) != 1 {
		return ResponsePacket{}, errors.New("invalid data length")
	}

	const returnCodeOffset = 0

	returnCode := Esp3ReturnCode(telegram.Data[returnCodeOffset])

	if !returnCode.Valid() {
		return ResponsePacket{}, errors.New("invalid return code")
	}

	return ResponsePacket{
		ReturnCode:   returnCode,
		OptionalData: telegram.Data[returnCodeOffset+1:],
	}, nil
}

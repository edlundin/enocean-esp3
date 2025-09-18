package response

import (
	"errors"

	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
)

type Packet struct {
	Code    enums.ReturnCode
	Data    []byte
	OptData []byte
}

func NewPacketFromEsp3(telegram esp3.Telegram) (Packet, error) {
	if telegram.PacketType != enums.PacketTypeRESPONSE {
		return Packet{}, errors.New("invalid packet type")
	}

	const returnCodeOffset = 0

	returnCode, err := enums.ParseReturnCodeFromByte(telegram.Data[returnCodeOffset])
	if err != nil {
		return Packet{}, err
	}

	return Packet{
		Code:    returnCode,
		Data:    telegram.Data[returnCodeOffset+1:],
		OptData: telegram.OptData,
	}, nil
}

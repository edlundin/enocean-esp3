package commoncommand

import (
	"errors"

	"github.com/edlundin/enocean-esp3/internal/serializer"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

type WrMem struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	Type        enums.MemoryType    `enocean-esp3:"data"`
	Address     uint32              `enocean-esp3:"data"`
	Data        []byte              `enocean-esp3:"data"`
}

func (cmd *WrMem) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewWrMem only supported for TCM3xx and TCM4xx
func NewWrMem(memoryType enums.MemoryType, address uint32, data []byte) (WrMem, error) {
	return WrMem{
		CommandCode: enums.CommonCommandWR_MEM,
		Type:        memoryType,
		Address:     address,
		Data:        data,
	}, nil
}

type RdMem struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	Type        enums.MemoryType    `enocean-esp3:"data"`
	Address     uint32              `enocean-esp3:"data"`
	DataLength  uint16              `enocean-esp3:"data"`
}

func (cmd *RdMem) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewRdMem(memoryType enums.MemoryType, address uint32, dataLength uint16) (RdMem, error) {
	return RdMem{
		CommandCode: enums.CommonCommandRD_MEM,
		Type:        memoryType,
		Address:     address,
		DataLength:  dataLength,
	}, nil
}

type RdMemResponse struct {
	Data []byte
}

func ParseRdMemResponseOK(response response.Packet) (RdMemResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return RdMemResponse{}, errors.New("invalid return code")
	}

	return RdMemResponse{
		Data: response.Data,
	}, nil
}

type RdMemAddress struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	Area        enums.MemoryArea    `enocean-esp3:"data"`
}

func (cmd *RdMemAddress) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewRdMemAddress(area enums.MemoryArea) (RdMemAddress, error) {
	return RdMemAddress{
		CommandCode: enums.CommonCommandRD_MEM_ADDRESS,
		Area:        area,
	}, nil
}

type RdMemAddressResponse struct {
	Type    enums.MemoryType
	Address uint32
	Length  uint32
}

func ParseRdMemAddressResponseOK(response response.Packet) (RdMemAddressResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return RdMemAddressResponse{}, errors.New("invalid return code")
	}

	var result RdMemAddressResponse
	if err := serializer.BytesToStruct(response.Data, &result); err != nil {
		return RdMemAddressResponse{}, errors.New("failed to deserialize response")
	}

	if !result.Type.Valid() {
		return RdMemAddressResponse{}, errors.New("invalid memory type")
	}

	return result, nil
}

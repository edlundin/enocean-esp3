package commoncommand

import (
	"errors"

	"github.com/edlundin/enocean-esp3/internal/serializer"
	"github.com/edlundin/enocean-esp3/pkg/deviceid"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

const (
	minBaseID     = 0xff800000
	maxBaseID     = 0xffffff80
	baseIDGapSize = 128
)

type WrIDBase struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	IDBase      deviceid.DeviceID   `enocean-esp3:"data"`
}

// Serialize encodes WrIDBase into its wire representation.
func (cmd *WrIDBase) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewWrIDBase constructs WrIDBase.
func NewWrIDBase(deviceID deviceid.DeviceID) (WrIDBase, error) {
	if deviceID < minBaseID || deviceID > maxBaseID {
		return WrIDBase{}, errors.New("device ID out of range")
	}

	if deviceID%baseIDGapSize != 0 {
		return WrIDBase{}, errors.New("device ID is not a base ID")
	}

	return WrIDBase{
		CommandCode: enums.CommonCommandWR_IDBASE,
		IDBase:      deviceID,
	}, nil
}

type RdIDBase struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
}

// Serialize encodes RdIDBase into its wire representation.
func (cmd *RdIDBase) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewRdIDBase constructs RdIDBase.
func NewRdIDBase() (RdIDBase, error) {
	return RdIDBase{
		CommandCode: enums.CommonCommandRD_IDBASE,
	}, nil
}

type RdIDBaseResponse struct {
	BaseID              deviceid.DeviceID
	RemainingWriteCount uint8
}

// ParseRdIDBaseResponseOK parses RdIDBaseResponseOK.
func ParseRdIDBaseResponseOK(response response.Packet) (RdIDBaseResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return RdIDBaseResponse{}, errors.New("invalid return code")
	}

	mergedData := make([]byte, 0, len(response.Data)+len(response.OptData))
	mergedData = append(mergedData, response.Data...)
	mergedData = append(mergedData, response.OptData...)

	var rdIDBaseResponse RdIDBaseResponse
	if err := serializer.BytesToStruct(mergedData, &rdIDBaseResponse); err != nil {
		return RdIDBaseResponse{}, errors.New("failed to deserialize response")
	}

	return rdIDBaseResponse, nil
}

package commoncommand

import (
	"errors"

	"github.com/edlundin/enocean-esp3/internal/serializer"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

type WrRemanCode struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	SecureCode  uint32              `enocean-esp3:"data"`
}

// Serialize encodes WrRemanCode into its wire representation.
func (cmd *WrRemanCode) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewWrRemanCode constructs WrRemanCode.
func NewWrRemanCode(secureCode uint32) (WrRemanCode, error) {
	return WrRemanCode{
		CommandCode: enums.CommonCommandWR_REMAN_CODE,
		SecureCode:  secureCode,
	}, nil
}

type WrRemanRepeating struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`

	//0x00: Reman telegrams will not be repeated (STATUS will be set to 0x8F)
	//0x01: Reman answers will be repeated (STATUS will be set to 0x80)
	SetRemanRepetition bool `enocean-esp3:"data"`
}

// Serialize encodes WrRemanRepeating into its wire representation.
func (cmd *WrRemanRepeating) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewWrRemanRepeating constructs WrRemanRepeating.
func NewWrRemanRepeating(setRemanRepetition bool) (WrRemanRepeating, error) {
	return WrRemanRepeating{
		CommandCode:        enums.CommonCommandWR_REMAN_REPEATING,
		SetRemanRepetition: setRemanRepetition,
	}, nil
}

type RdRemanRepeating struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
}

// Serialize encodes RdRemanRepeating into its wire representation.
func (cmd *RdRemanRepeating) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewRdRemanRepeating constructs RdRemanRepeating.
func NewRdRemanRepeating() (RdRemanRepeating, error) {
	return RdRemanRepeating{
		CommandCode: enums.CommonCommandRD_REMAN_REPEATING,
	}, nil
}

type RdRemanRepeatingResponse struct {
	RemanRepetitionEnabled bool
}

// ParseRdRemanRepeatingResponseOK parses RdRemanRepeatingResponseOK.
func ParseRdRemanRepeatingResponseOK(response response.Packet) (RdRemanRepeatingResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return RdRemanRepeatingResponse{}, errors.New("invalid return code")
	}



	var result RdRemanRepeatingResponse
	if err := serializer.BytesToStruct(response.Data, &result); err != nil {
		return RdRemanRepeatingResponse{}, errors.New("failed to deserialize response")
	}

	return result, nil
}

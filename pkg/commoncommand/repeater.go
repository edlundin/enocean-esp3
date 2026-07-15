package commoncommand

import (
	"errors"

	"github.com/edlundin/enocean-esp3/internal/serializer"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

type WrRepeater struct {
	CommandCode   enums.CommonCommand `enocean-esp3:"data"`
	RepeaterMode  enums.RepeaterMode  `enocean-esp3:"data"`
	RepeaterLevel enums.RepeaterLevel `enocean-esp3:"data"`
}

// Serialize encodes WrRepeater into its wire representation.
func (cmd *WrRepeater) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewWrRepeater constructs WrRepeater.
func NewWrRepeater(repeaterMode enums.RepeaterMode, repeaterLevel enums.RepeaterLevel) (WrRepeater, error) {
	return WrRepeater{
		CommandCode:   enums.CommonCommandWR_REPEATER,
		RepeaterMode:  repeaterMode,
		RepeaterLevel: repeaterLevel,
	}, nil
}

type RdRepeaterResponse struct {
	RepeaterMode  enums.RepeaterMode
	RepeaterLevel enums.RepeaterLevel
}

// ParseRdRepeaterResponseOK parses RdRepeaterResponseOK.
func ParseRdRepeaterResponseOK(response response.Packet) (RdRepeaterResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return RdRepeaterResponse{}, errors.New("invalid return code")
	}



	var result RdRepeaterResponse
	if err := serializer.BytesToStruct(response.Data, &result); err != nil {
		return RdRepeaterResponse{}, errors.New("failed to deserialize response")
	}

	if !result.RepeaterMode.Valid() {
		return RdRepeaterResponse{}, errors.New("invalid repeater mode")
	}

	if !result.RepeaterLevel.Valid() {
		return RdRepeaterResponse{}, errors.New("invalid repeater level")
	}

	return result, nil
}

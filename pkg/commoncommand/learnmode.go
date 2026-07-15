package commoncommand

import (
	"errors"

	"github.com/edlundin/enocean-esp3/internal/serializer"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

type WrLearnMode struct {
	CommandCode     enums.CommonCommand `enocean-esp3:"data"`
	EnableLearnMode bool                `enocean-esp3:"data"`
	Timeout         uint32              `enocean-esp3:"data"`
	Channel         uint8               `enocean-esp3:"data"`
}

// Serialize encodes WrLearnMode into its wire representation.
func (cmd *WrLearnMode) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewWrLearnMode constructs WrLearnMode.
func NewWrLearnMode(enableLearnMode bool, timeout uint32, channel uint8) (WrLearnMode, error) {
	return WrLearnMode{
		CommandCode:     enums.CommonCommandWR_LEARNMODE,
		EnableLearnMode: enableLearnMode,
		Timeout:         timeout,
		Channel:         channel,
	}, nil
}

type RdLearnModeResponse struct {
	LearnModeStatus bool
	Channel         uint8
}

// ParseRdLearnModeResponseOK parses RdLearnModeResponseOK.
func ParseRdLearnModeResponseOK(response response.Packet) (RdLearnModeResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return RdLearnModeResponse{}, errors.New("invalid return code")
	}



	mergedData := make([]byte, 0, len(response.Data)+len(response.OptData))
	mergedData = append(mergedData, response.Data...)
	mergedData = append(mergedData, response.OptData...)

	var result RdLearnModeResponse
	if err := serializer.BytesToStruct(mergedData, &result); err != nil {
		return RdLearnModeResponse{}, errors.New("failed to deserialize response")
	}

	return result, nil
}

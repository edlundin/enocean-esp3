package commoncommand

import (
	"errors"
	"fmt"

	"github.com/edlundin/enocean-esp3/internal/serializer"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

// WrBist is a command to run the Built-in Self Test
type WrBist struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
}

func (cmd *WrBist) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewBist() (WrBist, error) {
	return WrBist{
		CommandCode: enums.CommonCommandWR_BIST,
	}, nil
}

type WrBistResponse struct {
	BistResult bool
}

func ParseWrBistResponseOK(response response.Packet) (WrBistResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return WrBistResponse{}, errors.New("invalid return code")
	}

	var result WrBistResponse
	if err := serializer.BytesToStruct(response.Data, &result); err != nil {
		return WrBistResponse{}, fmt.Errorf("failed to deserialize response: %w", err)
	}

	return result, nil
}

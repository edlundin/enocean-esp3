package commoncommand

import (
	"github.com/edlundin/enocean-esp3/internal/serializer"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
)

type WrMode struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	Mode        enums.RadioMode     `enocean-esp3:"data"`
}

func (cmd *WrMode) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewWrMode(mode enums.RadioMode) (WrMode, error) {
	return WrMode{
		CommandCode: enums.CommonCommandWR_MODE,
		Mode:        mode,
	}, nil
}

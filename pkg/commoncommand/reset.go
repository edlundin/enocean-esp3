package commoncommand

import (
	"github.com/edlundin/enocean-esp3/internal/serializer"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
)

type WrReset struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
}

func (c *WrReset) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(c)
}

func NewWrReset() (WrReset, error) {
	return WrReset{
		CommandCode: enums.CommonCommandWR_RESET,
	}, nil
}

package commoncommand

import (
	"github.com/edlundin/enocean-esp3/internal/serializer"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
)

type WrSubTel struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	Toggle      bool                `enocean-esp3:"data"`
}

func (cmd *WrSubTel) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewWrSubTel(toggle bool) (WrSubTel, error) {
	return WrSubTel{
		CommandCode: enums.CommonCommandWR_SUBTEL,
		Toggle:      toggle,
	}, nil
}

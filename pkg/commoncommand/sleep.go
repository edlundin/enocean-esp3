package commoncommand

import (
	"slices"

	"github.com/edlundin/enocean-esp3/internal/serializer"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
)

const maxDeepSleepPeriod = 0xffffff

type WrSleep struct {
	CommandCode     enums.CommonCommand `enocean-esp3:"data"`
	DeepSleepPeriod uint32              `enocean-esp3:"data"`
}

func (cmd *WrSleep) Serialize() (esp3.Telegram, error) {
	cmd.DeepSleepPeriod = slices.Min([]uint32{cmd.DeepSleepPeriod, maxDeepSleepPeriod})
	return serializer.CommandToTelegram(cmd)
}

func NewWrSleep(deepSleepPeriod uint32) (WrSleep, error) {
	return WrSleep{
		CommandCode:     enums.CommonCommandWR_SLEEP,
		DeepSleepPeriod: slices.Min([]uint32{deepSleepPeriod, maxDeepSleepPeriod}),
	}, nil
}

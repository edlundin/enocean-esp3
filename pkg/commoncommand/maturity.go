package commoncommand

import (
	"github.com/edlundin/enocean-esp3/internal/serializer"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
)

type WrWaitMaturity struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	Maturity    enums.Maturity      `enocean-esp3:"data"`
}

// Serialize encodes WrWaitMaturity into its wire representation.
func (cmd *WrWaitMaturity) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewWrWaitMaturity constructs WrWaitMaturity.
func NewWrWaitMaturity(maturity enums.Maturity) (WrWaitMaturity, error) {
	return WrWaitMaturity{
		CommandCode: enums.CommonCommandWR_WAIT_MATURITY,
		Maturity:    maturity,
	}, nil
}

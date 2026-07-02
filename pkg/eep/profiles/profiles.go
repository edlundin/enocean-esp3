package profiles

import (
	"errors"
	"github.com/edlundin/enocean-esp3/pkg/eep"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/erp1"
)

type Telegram interface {
	EEP() eep.EEP
	MarshalERP1UserData() ([]byte, byte, error)
	Format() string
}

func ParsePacket(p erp1.Packet, prof eep.EEP) (Telegram, error) {
	if p.Rorg != prof.Rorg {
		return nil, errors.New("packet RORG does not match EEP")
	}
	return ParseUserData(prof, p.UserData, p.Status)
}

func ParseUserData(prof eep.EEP, userData []byte, status byte) (Telegram, error) {
	switch prof {
	case mustEEP(enums.Rorg1BS, 0x00, 0x01):
		return parseD50001(userData, status)
	case mustEEP(enums.RorgRPS, 0x01, 0x01):
		return parseF60101(userData, status)
	case mustEEP(enums.Rorg4BS, 0x02, 0x01):
		return parseA50201(userData, status)
	default:
		return Decode(prof, userData, status)
	}
}

func mustEEP(r enums.Rorg, f, t byte) eep.EEP { e, _ := eep.FromTriplet(r, f, t); return e }

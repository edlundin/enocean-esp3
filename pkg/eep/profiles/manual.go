package profiles

import (
	"errors"
	"fmt"

	"github.com/edlundin/enocean-esp3/pkg/eep"
	"github.com/edlundin/enocean-esp3/pkg/enums"
)

type D50001 struct {
	ContactClosed bool
	LearnButton   bool
}

// parseD50001 parses D50001.
func parseD50001(data []byte, status byte) (Telegram, error) {
	if len(data) < 1 {
		return nil, errors.New("D5-00-01 user data too short")
	}
	return D50001{ContactClosed: getBits(data, 7, 1) == 1, LearnButton: getBits(data, 4, 1) == 0}, nil
}

// EEP returns the EEP associated with D50001.
func (d D50001) EEP() eep.EEP { return mustEEP(enums.Rorg1BS, 0, 1) }

// MarshalERP1UserData marshals ERP1UserData.
func (d D50001) MarshalERP1UserData() ([]byte, byte, error) {
	b := []byte{0}
	if d.ContactClosed {
		setBits(b, 7, 1, 1)
	}
	if !d.LearnButton {
		setBits(b, 4, 1, 1)
	}
	return b, 0, nil
}

// Format returns the formatted representation of D50001.
func (d D50001) Format() string {
	return fmt.Sprintf("D5-00-01 ContactClosed=%t LearnButton=%t", d.ContactClosed, d.LearnButton)
}

type F60101 struct{ Pressed bool }

// parseF60101 parses F60101.
func parseF60101(data []byte, status byte) (Telegram, error) {
	if len(data) < 1 {
		return nil, errors.New("F6-01-01 user data too short")
	}
	return F60101{Pressed: getBits(data, 3, 1) == 1}, nil
}

// EEP returns the EEP associated with F60101.
func (f F60101) EEP() eep.EEP { return mustEEP(enums.RorgRPS, 1, 1) }

// MarshalERP1UserData marshals ERP1UserData.
func (f F60101) MarshalERP1UserData() ([]byte, byte, error) {
	b := []byte{0}
	if f.Pressed {
		setBits(b, 3, 1, 1)
	}
	return b, 0, nil
}

// Format returns the formatted representation of F60101.
func (f F60101) Format() string { return fmt.Sprintf("F6-01-01 Pressed=%t", f.Pressed) }

type A50201 struct {
	TemperatureC   float64
	TemperatureRaw byte
	LearnButton    bool
}

// parseA50201 parses A50201.
func parseA50201(data []byte, status byte) (Telegram, error) {
	if len(data) < 4 {
		return nil, errors.New("A5-02-01 user data too short")
	}
	raw := byte(getBits(data, 16, 8))
	return A50201{TemperatureRaw: raw, TemperatureC: eep.ScaleRaw(uint64(raw), 255, 0, -40, 0), LearnButton: getBits(data, 28, 1) == 0}, nil
}

// EEP returns the EEP associated with A50201.
func (a A50201) EEP() eep.EEP { return mustEEP(enums.Rorg4BS, 2, 1) }

// MarshalERP1UserData marshals ERP1UserData.
func (a A50201) MarshalERP1UserData() ([]byte, byte, error) {
	b := make([]byte, 4)
	raw := a.TemperatureRaw
	if raw == 0 && a.TemperatureC != 0 {
		raw = byte(eep.UnscaleRaw(a.TemperatureC, 255, 0, -40, 0))
	}
	setBits(b, 16, 8, uint64(raw))
	if !a.LearnButton {
		setBits(b, 28, 1, 1)
	}
	return b, 0, nil
}

// Format returns the formatted representation of A50201.
func (a A50201) Format() string {
	return fmt.Sprintf("A5-02-01 Temperature=%.2f°C Raw=%d LearnButton=%t", a.TemperatureC, a.TemperatureRaw, a.LearnButton)
}

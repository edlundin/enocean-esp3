package commoncommand

import (
	"errors"
	"fmt"

	"github.com/edlundin/enocean-esp3/internal/serializer"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

type WrFilterAdd struct {
	CommandCode enums.CommonCommand   `enocean-esp3:"data"`
	Action      uint8                 `enocean-esp3:"data"`
	Criterion   enums.FilterCriterion `enocean-esp3:"data"`
	Value       uint32                `enocean-esp3:"data"`
}

// Serialize encodes WrFilterAdd into its wire representation.
func (cmd *WrFilterAdd) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewWrFilterAdd constructs WrFilterAdd.
func NewWrFilterAdd(criterion enums.FilterCriterion, value uint32, forward bool, repeat bool) (WrFilterAdd, error) {
	filterAction := byte(0)

	if forward {
		filterAction |= byte(enums.FilterActionFORWARD)
	} else {
		filterAction |= byte(enums.FilterActionNO_FORWARD)
	}

	if repeat {
		filterAction |= byte(enums.FilterActionREPEAT)
	} else {
		filterAction |= byte(enums.FilterActionNO_REPEAT)
	}

	return WrFilterAdd{
		CommandCode: enums.CommonCommandWR_FILTER_ADD,
		Action:      filterAction,
		Criterion:   criterion,
		Value:       value,
	}, nil
}

type WrFilterDel struct {
	CommandCode enums.CommonCommand   `enocean-esp3:"data"`
	Action      uint8                 `enocean-esp3:"data"`
	Criterion   enums.FilterCriterion `enocean-esp3:"data"`
	Value       uint32                `enocean-esp3:"data"`
}

// Serialize encodes WrFilterDel into its wire representation.
func (cmd *WrFilterDel) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewWrFilterDel constructs WrFilterDel.
func NewWrFilterDel(criterion enums.FilterCriterion, value uint32, forward bool, repeat bool) (WrFilterDel, error) {
	filterAction := byte(0)

	if forward {
		filterAction |= byte(enums.FilterActionFORWARD)
	} else {
		filterAction |= byte(enums.FilterActionNO_FORWARD)
	}

	if repeat {
		filterAction |= byte(enums.FilterActionREPEAT)
	} else {
		filterAction |= byte(enums.FilterActionNO_REPEAT)
	}

	return WrFilterDel{
		CommandCode: enums.CommonCommandWR_FILTER_DEL,
		Action:      filterAction,
		Criterion:   criterion,
		Value:       value,
	}, nil
}

type WrFilterDelAll struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
}

// Serialize encodes WrFilterDelAll into its wire representation.
func (cmd *WrFilterDelAll) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewWrFilterDelAll constructs WrFilterDelAll.
func NewWrFilterDelAll() (WrFilterDelAll, error) {
	return WrFilterDelAll{
		CommandCode: enums.CommonCommandWR_FILTER_DEL_ALL,
	}, nil
}

type WrFilterEnable struct {
	CommandCode   enums.CommonCommand `enocean-esp3:"data"`
	Toggle        bool                `enocean-esp3:"data"`
	FilerOperator enums.FilerOperator `enocean-esp3:"data"`
}

// Serialize encodes WrFilterEnable into its wire representation.
func (cmd *WrFilterEnable) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewWrFilterEnable constructs WrFilterEnable.
func NewWrFilterEnable(toggle bool, operator enums.FilerOperator) (WrFilterEnable, error) {
	return WrFilterEnable{
		CommandCode:   enums.CommonCommandWR_FILTER_ENABLE,
		Toggle:        toggle,
		FilerOperator: operator,
	}, nil
}

type RdFilter struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
}

// Serialize encodes RdFilter into its wire representation.
func (cmd *RdFilter) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewRdFilter constructs RdFilter.
func NewRdFilter() (RdFilter, error) {
	return RdFilter{
		CommandCode: enums.CommonCommandRD_FILTER,
	}, nil
}

type Filter struct {
	Criterion enums.FilterCriterion
	Value     uint32
}

type RdFilterResponse struct {
	Filters []Filter
}

// ParseRdFilterResponseOK parses RdFilterResponseOK.
func ParseRdFilterResponseOK(response response.Packet) (RdFilterResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return RdFilterResponse{}, errors.New("invalid return code")
	}

	var raw struct {
		Count   uint8
		Filters []Filter
	}

	if err := serializer.BytesToStruct(response.Data, &raw); err != nil {
		return RdFilterResponse{}, fmt.Errorf("failed to deserialize response: %w", err)
	}
	if len(raw.Filters) != int(raw.Count) {
		return RdFilterResponse{}, fmt.Errorf("filter count %d does not match %d records", raw.Count, len(raw.Filters))
	}

	return RdFilterResponse{
		Filters: raw.Filters,
	}, nil
}

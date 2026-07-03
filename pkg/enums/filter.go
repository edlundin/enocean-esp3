package enums

import "errors"

type FilterCriterion byte

const (
	FilterCriterionSENDER_ID FilterCriterion = iota
	FilterCriterionRORG
	FilterCriterionRSSI
	FilterCriterionDESTINATION_ID
)

func ParseFilterFromByte(b byte) (FilterCriterion, error) {
	switch b {
	case 0x01:
		return FilterCriterionSENDER_ID, nil
	case 0x02:
		return FilterCriterionRORG, nil
	case 0x03:
		return FilterCriterionRSSI, nil
	case 0x04:
		return FilterCriterionDESTINATION_ID, nil
	default:
		return 0, errors.New("invalid filter criterion")
	}
}

func (filterCriterion FilterCriterion) String() string {
	switch filterCriterion {
	case FilterCriterionSENDER_ID:
		return "SENDER_ID"
	case FilterCriterionRORG:
		return "RORG"
	case FilterCriterionRSSI:
		return "RSSI"
	case FilterCriterionDESTINATION_ID:
		return "DESTINATION_ID"
	default:
		return "UNKNOWN"
	}
}

func (filterCriterion FilterCriterion) Valid() bool {
	switch filterCriterion {
	case FilterCriterionSENDER_ID,
		FilterCriterionRORG,
		FilterCriterionRSSI,
		FilterCriterionDESTINATION_ID:
		return true
	default:
		return false
	}
}

type FilterActionMask byte

const (
	FilterActionNO_FORWARD FilterActionMask = 0x00
	FilterActionNO_REPEAT  FilterActionMask = 0x40
	FilterActionFORWARD    FilterActionMask = 0x80
	FilterActionREPEAT     FilterActionMask = 0xC0
)

func ParseFilterActionMaskFromByte(b byte) (FilterActionMask, error) {
	switch b {
	case 0x00:
		return FilterActionNO_FORWARD, nil
	case 0x40:
		return FilterActionNO_REPEAT, nil
	case 0x80:
		return FilterActionFORWARD, nil
	case 0xC0:
		return FilterActionREPEAT, nil
	default:
		return 0, errors.New("invalid filter action mask")
	}
}

func (filterActionMask FilterActionMask) String() string {
	switch filterActionMask {
	case FilterActionNO_FORWARD:
		return "NO_FORWARD"
	case FilterActionNO_REPEAT:
		return "NO_REPEAT"
	case FilterActionFORWARD:
		return "FORWARD"
	case FilterActionREPEAT:
		return "REPEAT"
	default:
		return "UNKNOWN"
	}
}

func (filterActionMask FilterActionMask) Valid() bool {
	switch filterActionMask {
	case FilterActionNO_FORWARD,
		FilterActionNO_REPEAT,
		FilterActionFORWARD,
		FilterActionREPEAT:
		return true
	default:
		return false
	}
}

type FilerOperator byte

const (
	FilerOperatorOR_ALL_FILTERS                FilerOperator = 0x00
	FilerOperatorAND_ALL_FILTERS               FilerOperator = 0x01
	FilerOperatorOR_FOR_RECEIVE_AND_FOR_REPEAT FilerOperator = 0x08
	FilerOperatorAND_FOR_RECEIVE_OR_FOR_REPEAT FilerOperator = 0x09
)

func ParseFilerOperatorFromByte(b byte) (FilerOperator, error) {
	switch b {
	case 0x00:
		return FilerOperatorOR_ALL_FILTERS, nil
	case 0x01:
		return FilerOperatorAND_ALL_FILTERS, nil
	case 0x08:
		return FilerOperatorOR_FOR_RECEIVE_AND_FOR_REPEAT, nil
	case 0x09:
		return FilerOperatorAND_FOR_RECEIVE_OR_FOR_REPEAT, nil
	default:
		return 0, errors.New("invalid filer operator")
	}
}

func (filerOperator FilerOperator) String() string {
	switch filerOperator {
	case FilerOperatorOR_ALL_FILTERS:
		return "OR_ALL_FILTERS"
	case FilerOperatorAND_ALL_FILTERS:
		return "AND_ALL_FILTERS"
	case FilerOperatorOR_FOR_RECEIVE_AND_FOR_REPEAT:
		return "OR_FOR_RECEIVE_AND_FOR_REPEAT"
	case FilerOperatorAND_FOR_RECEIVE_OR_FOR_REPEAT:
		return "AND_FOR_RECEIVE_OR_FOR_REPEAT"
	default:
		return "UNKNOWN"
	}
}

func (filerOperator FilerOperator) Valid() bool {
	switch filerOperator {
	case FilerOperatorOR_ALL_FILTERS,
		FilerOperatorAND_ALL_FILTERS,
		FilerOperatorOR_FOR_RECEIVE_AND_FOR_REPEAT,
		FilerOperatorAND_FOR_RECEIVE_OR_FOR_REPEAT:
		return true
	default:
		return false
	}
}

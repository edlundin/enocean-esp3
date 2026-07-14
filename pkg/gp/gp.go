package gp

import (
	"errors"
	"fmt"

	"github.com/edlundin/enocean-esp3/pkg/enums"
)

const MaxMessageLength = 512

type Purpose byte

const (
	PurposeTeachIn Purpose = iota
	PurposeTeachOut
	PurposeToggle
)

type Result byte

const (
	ResultRejected Result = iota
	ResultSuccess
	ResultTeachOut
	ResultRejectedChannels
)

type RequestHeader struct {
	ManufacturerID uint16
	Bidirectional  bool
	Purpose        Purpose
}

type ResponseHeader struct {
	ManufacturerID uint16
	Result         Result
}

func IsRorg(r enums.Rorg) bool {
	return r == enums.RorgGP_TI || r == enums.RorgGP_TR || r == enums.RorgGP_CD || r == enums.RorgGP_SD
}

func EncodeRequestHeader(h RequestHeader) ([]byte, error) {
	if h.ManufacturerID > 0x7ff || h.Purpose > PurposeToggle {
		return nil, errors.New("invalid GP teach-in request header")
	}
	word := uint64(h.ManufacturerID) << 5
	if h.Bidirectional {
		word |= 1 << 4
	}
	word |= uint64(h.Purpose) << 2
	out := make([]byte, 2)
	_ = writeUnsigned(out, 0, 16, word)
	return out, nil
}

func DecodeRequestHeader(data []byte) (RequestHeader, error) {
	word, err := readUnsigned(data, 0, 16)
	if err != nil {
		return RequestHeader{}, err
	}
	if word&0x3 != 0 {
		return RequestHeader{}, errors.New("invalid GP teach-in request reserved bits")
	}
	return RequestHeader{ManufacturerID: uint16(word >> 5), Bidirectional: word&(1<<4) != 0, Purpose: Purpose((word >> 2) & 0x3)}, nil
}

func EncodeResponseHeader(h ResponseHeader) ([]byte, error) {
	if h.ManufacturerID > 0x7ff || h.Result > ResultRejectedChannels {
		return nil, errors.New("invalid GP teach-in response header")
	}
	out := make([]byte, 2)
	_ = writeUnsigned(out, 0, 16, uint64(h.ManufacturerID)<<5|uint64(h.Result)<<3)
	return out, nil
}

func DecodeResponseHeader(data []byte) (ResponseHeader, error) {
	word, err := readUnsigned(data, 0, 16)
	if err != nil {
		return ResponseHeader{}, err
	}
	if word&0x7 != 0 {
		return ResponseHeader{}, errors.New("invalid GP teach-in response reserved bits")
	}
	return ResponseHeader{ManufacturerID: uint16(word >> 5), Result: Result((word >> 3) & 0x3)}, nil
}

type ChannelType byte

const (
	ChannelTeachInInformation ChannelType = iota
	ChannelData
	ChannelFlag
	ChannelEnumeration
)

type ValueType byte

const (
	ValueReserved ValueType = iota
	ValueCurrent
	ValueSetPointAbsolute
	ValueSetPointRelative
)

type Channel struct {
	Type           ChannelType
	SignalType     byte
	ValueType      ValueType
	ResolutionCode byte
	EngineeringMin byte // Two's-complement; use EngineeringRange for signed values.
	ScalingMin     byte
	EngineeringMax byte // Two's-complement; use EngineeringRange for signed values.
	ScalingMax     byte
}

func (c Channel) EngineeringRange() (int8, int8) {
	return int8(c.EngineeringMin), int8(c.EngineeringMax)
}

func ResolutionBits(code byte) (int, bool) {
	bits := [...]int{0, 2, 3, 4, 5, 6, 8, 10, 12, 16, 20, 24, 32}
	if int(code) >= len(bits) || bits[code] == 0 {
		return 0, false
	}
	return bits[code], true
}

func (c Channel) ValueBits() (int, error) {
	switch c.Type {
	case ChannelFlag:
		return 1, nil
	case ChannelData, ChannelEnumeration:
		if bits, ok := ResolutionBits(c.ResolutionCode); ok {
			return bits, nil
		}
		return 0, fmt.Errorf("invalid GP resolution code %d", c.ResolutionCode)
	case ChannelTeachInInformation:
		return 0, nil
	default:
		return 0, errors.New("invalid GP channel type")
	}
}

func EncodeChannelDefinition(c Channel) ([]byte, int, error) {
	bits, err := channelDefinitionBits(c)
	if err != nil {
		return nil, 0, err
	}
	out := bytesForBits(bits)
	if err := writeChannelDefinition(out, 0, c); err != nil {
		return nil, 0, err
	}
	return out, bits, nil
}

func DecodeChannelDefinition(data []byte, bitOffset int) (Channel, int, error) {
	ct, err := readUnsigned(data, bitOffset, 2)
	if err != nil {
		return Channel{}, 0, err
	}
	c := Channel{Type: ChannelType(ct)}
	bits, err := channelDefinitionBits(c)
	if err != nil {
		return Channel{}, 0, err
	}
	if _, err := readUnsigned(data, bitOffset, bits); err != nil {
		return Channel{}, 0, err
	}
	c.SignalType = byte(mustRead(data, bitOffset+2, 8))
	c.ValueType = ValueType(mustRead(data, bitOffset+10, 2))
	if c.Type == ChannelData {
		c.ResolutionCode = byte(mustRead(data, bitOffset+12, 4))
		c.EngineeringMin = byte(mustRead(data, bitOffset+16, 8))
		c.ScalingMin = byte(mustRead(data, bitOffset+24, 4))
		c.EngineeringMax = byte(mustRead(data, bitOffset+28, 8))
		c.ScalingMax = byte(mustRead(data, bitOffset+36, 4))
	} else if c.Type == ChannelEnumeration {
		c.ResolutionCode = byte(mustRead(data, bitOffset+12, 4))
	}
	return c, bits, nil
}

func channelDefinitionBits(c Channel) (int, error) {
	switch c.Type {
	case ChannelData:
		return 40, nil
	case ChannelFlag:
		return 12, nil
	case ChannelEnumeration:
		return 16, nil
	default:
		return 0, errors.New("unsupported GP channel definition")
	}
}

func writeChannelDefinition(out []byte, bitOffset int, c Channel) error {
	if err := writeUnsigned(out, bitOffset, 2, uint64(c.Type)); err != nil {
		return err
	}
	_ = writeUnsigned(out, bitOffset+2, 8, uint64(c.SignalType))
	_ = writeUnsigned(out, bitOffset+10, 2, uint64(c.ValueType))
	if c.Type == ChannelData || c.Type == ChannelEnumeration {
		_ = writeUnsigned(out, bitOffset+12, 4, uint64(c.ResolutionCode))
	}
	if c.Type == ChannelData {
		_ = writeUnsigned(out, bitOffset+16, 8, uint64(c.EngineeringMin))
		_ = writeUnsigned(out, bitOffset+24, 4, uint64(c.ScalingMin))
		_ = writeUnsigned(out, bitOffset+28, 8, uint64(c.EngineeringMax))
		_ = writeUnsigned(out, bitOffset+36, 4, uint64(c.ScalingMax))
	}
	return nil
}

func mustRead(data []byte, off, size int) uint64 {
	v, _ := readUnsigned(data, off, size)
	return v
}

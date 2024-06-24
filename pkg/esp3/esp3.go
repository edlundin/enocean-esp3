package esp3

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/edlundin/enocean-esp3/pkg/enocean"
)

const syncByte = 0x55

func crcTable() [256]byte {
	return [256]byte{
		0x00, 0x07, 0x0e, 0x09, 0x1c, 0x1b, 0x12, 0x15, 0x38, 0x3f, 0x36, 0x31, 0x24, 0x23, 0x2a, 0x2d,
		0x70, 0x77, 0x7e, 0x79, 0x6c, 0x6b, 0x62, 0x65, 0x48, 0x4f, 0x46, 0x41, 0x54, 0x53, 0x5a, 0x5d,
		0xe0, 0xe7, 0xee, 0xe9, 0xfc, 0xfb, 0xf2, 0xf5, 0xd8, 0xdf, 0xd6, 0xd1, 0xc4, 0xc3, 0xca, 0xcd,
		0x90, 0x97, 0x9e, 0x99, 0x8c, 0x8b, 0x82, 0x85, 0xa8, 0xaf, 0xa6, 0xa1, 0xb4, 0xb3, 0xba, 0xbd,
		0xc7, 0xc0, 0xc9, 0xce, 0xdb, 0xdc, 0xd5, 0xd2, 0xff, 0xf8, 0xf1, 0xf6, 0xe3, 0xe4, 0xed, 0xea,
		0xb7, 0xb0, 0xb9, 0xbe, 0xab, 0xac, 0xa5, 0xa2, 0x8f, 0x88, 0x81, 0x86, 0x93, 0x94, 0x9d, 0x9a,
		0x27, 0x20, 0x29, 0x2e, 0x3b, 0x3c, 0x35, 0x32, 0x1f, 0x18, 0x11, 0x16, 0x03, 0x04, 0x0d, 0x0a,
		0x57, 0x50, 0x59, 0x5e, 0x4b, 0x4c, 0x45, 0x42, 0x6f, 0x68, 0x61, 0x66, 0x73, 0x74, 0x7d, 0x7a,
		0x89, 0x8e, 0x87, 0x80, 0x95, 0x92, 0x9b, 0x9c, 0xb1, 0xb6, 0xbf, 0xb8, 0xad, 0xaa, 0xa3, 0xa4,
		0xf9, 0xfe, 0xf7, 0xf0, 0xe5, 0xe2, 0xeb, 0xec, 0xc1, 0xc6, 0xcf, 0xc8, 0xdd, 0xda, 0xd3, 0xd4,
		0x69, 0x6e, 0x67, 0x60, 0x75, 0x72, 0x7b, 0x7c, 0x51, 0x56, 0x5f, 0x58, 0x4d, 0x4a, 0x43, 0x44,
		0x19, 0x1e, 0x17, 0x10, 0x05, 0x02, 0x0b, 0x0c, 0x21, 0x26, 0x2f, 0x28, 0x3d, 0x3a, 0x33, 0x34,
		0x4e, 0x49, 0x40, 0x47, 0x52, 0x55, 0x5c, 0x5b, 0x76, 0x71, 0x78, 0x7f, 0x6A, 0x6d, 0x64, 0x63,
		0x3e, 0x39, 0x30, 0x37, 0x22, 0x25, 0x2c, 0x2b, 0x06, 0x01, 0x08, 0x0f, 0x1a, 0x1d, 0x14, 0x13,
		0xae, 0xa9, 0xa0, 0xa7, 0xb2, 0xb5, 0xbc, 0xbb, 0x96, 0x91, 0x98, 0x9f, 0x8a, 0x8D, 0x84, 0x83,
		0xde, 0xd9, 0xd0, 0xd7, 0xc2, 0xc5, 0xcc, 0xcb, 0xe6, 0xe1, 0xe8, 0xef, 0xfa, 0xfd, 0xf4, 0xf3,
	}
}

type Telegram struct {
	DataLen    uint16
	OptDataLen uint8
	PacketType enocean.Esp3PacketType
	Data       []byte
	OptData    []byte
}

func computeCrc8(byte byte, crc byte) byte {
	return crcTable()[byte^crc]
}

func computeCrcFromSlice(slice []byte) byte {
	crc := byte(0)

	for _, byte := range slice {
		crc = computeCrc8(byte, crc)
	}

	return crc
}

func (telegram Telegram) String() string {
	var stringBuilder strings.Builder

	stringBuilder.WriteString(fmt.Sprintf("PacketType: %s\n", telegram.PacketType))
	stringBuilder.WriteString(fmt.Sprintf("DataLen: %d\n", telegram.DataLen))
	stringBuilder.WriteString(fmt.Sprintf("OptDataLen: %d\n", telegram.OptDataLen))
	stringBuilder.WriteString(fmt.Sprintf("Data: %s\n", strings.ToUpper(hex.EncodeToString(telegram.Data))))
	stringBuilder.WriteString(fmt.Sprintf("OptData: %s", strings.ToUpper(hex.EncodeToString(telegram.OptData))))

	return stringBuilder.String()
}

func (telegram Telegram) Serialize() []byte {
	header := []byte{byte(telegram.DataLen >> 8), byte(telegram.DataLen), telegram.OptDataLen, byte(telegram.PacketType)}
	crc8h := computeCrcFromSlice(header)
	crc8d := computeCrcFromSlice(slices.Concat(telegram.Data, telegram.OptData))

	return slices.Concat([]byte{syncByte}, header, []byte{crc8h}, telegram.Data, telegram.OptData, []byte{crc8d})
}

func FromData(packetType enocean.Esp3PacketType, data []byte, optData []byte) Telegram {
	return Telegram{
		DataLen:    uint16(len(data)),
		OptDataLen: uint8(len(optData)),
		PacketType: packetType,
		Data:       data,
		OptData:    optData,
	}
}

func FromHexString(hexStr string) (Telegram, error) {
	const MIN_HEX_STRING_LENGTH = 7

	if len(hexStr)%2 != 0 {
		hexStr = strings.Repeat("0", 1) + hexStr
	}

	telegram := Telegram{}
	bytes, err := hex.DecodeString(hexStr)

	if err != nil {
		return Telegram{}, fmt.Errorf("invalid hex string: %s", err)
	}

	if len(bytes) < MIN_HEX_STRING_LENGTH {
		return Telegram{}, errors.New("hex string too short")
	}

	if bytes[0] != 0x55 {
		return Telegram{}, errors.New("sync byte missing")
	}

	crc8h := bytes[5]
	crc := byte(0)

	for _, byte := range bytes[1:5] {
		crc = computeCrc8(byte, crc)
	}

	if crc != crc8h {
		return Telegram{}, fmt.Errorf("invalid CRC8H (got:0x%x, valid:0x%x)", crc8h, crc)
	}

	crc8dIndex := len(bytes) - 1
	crc8d := bytes[crc8dIndex]
	crc = 0

	for _, byte := range bytes[6:crc8dIndex] {
		crc = computeCrc8(byte, crc)
	}

	if crc != crc8d {
		return Telegram{}, fmt.Errorf("invalid CRC8D (got:0x%x, valid:0x%x)", crc8d, crc)
	}

	telegram.DataLen = binary.BigEndian.Uint16(bytes[1:3])
	telegram.OptDataLen = bytes[3]

	packetType, err := enocean.ParseEsp3PacketFromByte(bytes[4])

	if err != nil {
		return Telegram{}, err
	}

	telegram.PacketType = packetType
	telegram.Data = bytes[6 : telegram.DataLen+6]
	telegram.OptData = bytes[telegram.DataLen+6 : crc8dIndex]

	return telegram, nil
}

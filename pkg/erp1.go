package pkg

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type Erp1Packet struct {
	DestinationID DeviceID
	Rorg          Rorg
	Rssi          byte
	SecurityLevel byte
	Status        byte
	SubTelNum     byte
	SenderID      DeviceID
	UserData      []byte
}

func NewErp1PacketFromEsp3(telegram Esp3Telegram) (Erp1Packet, error) {
	if telegram.PacketType != PACKET_TYPE_RADIO_ERP1 {
		return Erp1Packet{}, errors.New("invalid packet type")
	}

	const rorgOffset = 0
	const rorgLen = 1
	const deviceIdLen = 4
	const statusLen = 1
	const dataOffset = 1
	const destinationIdOffset = 1
	const subTelNumOffset = 0
	const rssiOffset = destinationIdOffset + deviceIdLen
	const securityLevelOffset = rssiOffset + 1

	dataLen := telegram.DataLen - rorgLen - deviceIdLen - statusLen
	senderIdOffset := dataOffset + dataLen
	statusOffset := senderIdOffset + deviceIdLen

	destinationId, err := DeviceIdFromByteArray(telegram.OptData[destinationIdOffset : destinationIdOffset+deviceIdLen])

	if err != nil {
		return Erp1Packet{}, err
	}

	senderId, err := DeviceIdFromByteArray(telegram.Data[senderIdOffset:deviceIdLen])

	if err != nil {
		return Erp1Packet{}, err
	}

	return Erp1Packet{
		DestinationID: destinationId,
		Rorg:          Rorg(telegram.Data[rorgOffset]),
		Rssi:          telegram.Data[rssiOffset],
		SecurityLevel: telegram.Data[securityLevelOffset],
		Status:        telegram.Data[statusOffset],
		SubTelNum:     telegram.Data[subTelNumOffset],
		SenderID:      senderId,
		UserData:      telegram.Data[dataOffset : dataOffset+dataLen],
	}, nil
}

func (p Erp1Packet) ToEsp3() Esp3Telegram {
	return Esp3Telegram{
		PacketType: PACKET_TYPE_RADIO_ERP1,
		Data:       p.UserData,
		OptData: []byte{
			byte(p.DestinationID),
		},
	}
}

func (p Erp1Packet) Serialize() []byte {
	return p.ToEsp3().Serialize()
}

func (p Erp1Packet) String() string {
	var stringBuilder strings.Builder

	stringBuilder.WriteString(fmt.Sprintf("DestinationID: %s\n", p.DestinationID))
	stringBuilder.WriteString(fmt.Sprintf("Rorg: %s\n", p.Rorg))
	stringBuilder.WriteString(fmt.Sprintf("Rssi: %d\n", p.Rssi))
	stringBuilder.WriteString(fmt.Sprintf("SecurityLevel: %d\n", p.SecurityLevel))
	stringBuilder.WriteString(fmt.Sprintf("Status: %d\n", p.Status))
	stringBuilder.WriteString(fmt.Sprintf("SubTelNum: %d\n", p.SubTelNum))
	stringBuilder.WriteString(fmt.Sprintf("SenderID: %s\n", p.SenderID))
	stringBuilder.WriteString(fmt.Sprintf("UserData: %s", hex.EncodeToString(p.UserData)))

	return stringBuilder.String()
}

package pkg

import (
	"errors"
)

type Packet struct {
	DestinationId DeviceId
	Rorg          Rorg
	Rssi          byte
	SecurityLevel byte
	Status        byte
	SubTelNum     byte
	SenderId      DeviceId
	UserData      []byte
}

func FromEsp3(esp3Telegram Telegram) (Packet, error) {
	if esp3Telegram.PacketType != PACKET_TYPE_RADIO_ERP1 {
		return Packet{}, errors.New("invalid packet type")
	}

	const SIZE_RORG = 1
	const SIZE_DEVICE_ID = 4
	const SIZE_STATUS = 1
	const OFFSET_RORG = 0
	const OFFSET_DATA = 1
	const OFFSET_DESTINATION_ID = 1
	const OFFSET_SUB_TEL_NUM = 0
	const OFFSET_RSSI = OFFSET_DESTINATION_ID + SIZE_DEVICE_ID
	const OFFSET_SECURITY_LEVEL = OFFSET_RSSI + 1

	sizeData := esp3Telegram.DataLen - SIZE_RORG - SIZE_DEVICE_ID - SIZE_STATUS
	offsetSenderId := OFFSET_DATA + sizeData
	offsetStatus := offsetSenderId + SIZE_DEVICE_ID
	destinationId, err := DeviceIdFromByteArray(esp3Telegram.OptData[OFFSET_DESTINATION_ID : OFFSET_DESTINATION_ID+SIZE_DEVICE_ID])

	if err != nil {
		return Packet{}, err
	}

	senderId, err := DeviceIdFromByteArray(esp3Telegram.Data[offsetSenderId:SIZE_DEVICE_ID])

	if err != nil {
		return Packet{}, err
	}

	return Packet{
		DestinationId: destinationId,
		Rorg:          Rorg(esp3Telegram.Data[OFFSET_RORG]),
		Rssi:          esp3Telegram.Data[OFFSET_RSSI],
		SecurityLevel: esp3Telegram.Data[OFFSET_SECURITY_LEVEL],
		Status:        esp3Telegram.Data[offsetStatus],
		SubTelNum:     esp3Telegram.Data[OFFSET_SUB_TEL_NUM],
		SenderId:      senderId,
		UserData:      esp3Telegram.Data[OFFSET_DATA : OFFSET_DATA+sizeData],
	}, nil
}

func (p Packet) ToEsp3() Telegram {
	return Telegram{
		PacketType: PACKET_TYPE_RADIO_ERP1,
		Data:       p.UserData,
		OptData: []byte{
			byte(p.DestinationId),
		},
	}
}

func (p Packet) Serialize() []byte {
	return p.ToEsp3().Serialize()
}

func (p Packet) String() string {
	return "not implemented yet"
}

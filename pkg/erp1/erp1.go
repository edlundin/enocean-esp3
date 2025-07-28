package erp1

import (
	"errors"

	device_id "github.com/edlundin/enocean-esp3/pkg/device-id"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
)

type Erp1Packet struct {
	DestinationID device_id.DeviceID
	Rorg          enums.Rorg
	Rssi          byte
	SecurityLevel byte
	Status        byte
	SubTelNum     byte
	SenderID      device_id.DeviceID
	UserData      []byte
}

func NewErp1PacketFromEsp3(telegram esp3.Esp3Telegram) (Erp1Packet, error) {
	const minDataLen = 6    // 1 rorg + 4 sender ID + 1 status
	const minOptDataLen = 7 // 1 subTelNum + 4 destination ID + 1 rssi + 1 security level
	const destinationIdOffset = 1
	const rorgOffset = 0
	const rssiOffset = 5
	const securityLevelOffset = 6
	const subTelNumOffset = 0
	const userDataOffset = 1

	statusOffset := len(telegram.Data) - 1
	senderIdOffset := statusOffset - device_id.DeviceIDSize

	if telegram.PacketType != enums.PACKET_TYPE_RADIO_ERP1 {
		return Erp1Packet{}, errors.New("invalid packet type")
	}

	if len(telegram.Data) < minDataLen {
		return Erp1Packet{}, errors.New("data too short")
	}

	if len(telegram.OptData) < minOptDataLen {
		return Erp1Packet{}, errors.New("optData too short for destination ID")
	}

	destinationId, _ := device_id.FromByteArray(telegram.OptData[destinationIdOffset : destinationIdOffset+device_id.DeviceIDSize])
	senderId, _ := device_id.FromByteArray(telegram.Data[senderIdOffset : senderIdOffset+device_id.DeviceIDSize])

	rorg := enums.Rorg(telegram.Data[rorgOffset])
	rssi := telegram.OptData[rssiOffset]
	securityLevel := telegram.OptData[securityLevelOffset]
	status := telegram.Data[statusOffset]
	subTelNum := telegram.OptData[subTelNumOffset]
	userData := telegram.Data[userDataOffset:senderIdOffset]

	return Erp1Packet{
		DestinationID: destinationId,
		Rorg:          rorg,
		Rssi:          rssi,
		SecurityLevel: securityLevel,
		Status:        status,
		SubTelNum:     subTelNum,
		SenderID:      senderId,
		UserData:      userData,
	}, nil
}

func (p Erp1Packet) ToEsp3() esp3.Esp3Telegram {
	senderID := p.SenderID.ToArray()
	destinationID := p.DestinationID.ToArray()

	data := make([]byte, 0, 1+len(p.UserData)+device_id.DeviceIDSize+1)
	data = append(data, byte(p.Rorg))
	data = append(data, p.UserData...)
	data = append(data, senderID[:]...)
	data = append(data, p.Status)

	optData := make([]byte, 0, 3+device_id.DeviceIDSize)
	optData = append(optData, p.SubTelNum)
	optData = append(optData, destinationID[:]...)
	optData = append(optData, 0xff)
	optData = append(optData, 0x03)

	return esp3.Esp3Telegram{
		PacketType: enums.PACKET_TYPE_RADIO_ERP1,
		Data:       data,
		OptData:    optData,
	}
}

func (p Erp1Packet) Serialize() []byte {
	return p.ToEsp3().Serialize()
}

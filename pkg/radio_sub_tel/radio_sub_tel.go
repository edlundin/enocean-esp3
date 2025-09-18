package radiosubtel

import (
	"encoding/binary"
	"errors"

	device_id "github.com/edlundin/enocean-esp3/pkg/device_id"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
)

type SubTel struct {
	Tick   byte
	Rssi   byte
	Status byte
}

type Packet struct {
	DestinationID device_id.DeviceID
	Rorg          enums.Rorg
	Rssi          byte
	SecurityLevel byte
	Status        byte
	SubTelNum     byte
	SenderID      device_id.DeviceID
	SubTels       []SubTel
	Timestamp     uint16
	UserData      []byte
}

func NewPacketFromEsp3(telegram esp3.Telegram) (Packet, error) {
	const minDataLen = 6    // 1 rorg + 4 sender ID + 1 status
	const minOptDataLen = 9 // 1 subTelNum + 4 destination ID + 1 rssi + 1 security level + 2 timestamp
	const subTelSize = 3    // 1 tick + 1 rssi + 1 status
	const destinationIdOffset = 1
	const rorgOffset = 0
	const rssiOffset = 5
	const securityLevelOffset = 6
	const subTelNumOffset = 0
	const timestampOffset = 7
	const timestampSize = 2
	const userDataOffset = 1

	statusOffset := len(telegram.Data) - 1
	senderIdOffset := statusOffset - device_id.DeviceIDSize

	if telegram.PacketType != enums.PacketTypeRADIO_ERP1 {
		return Packet{}, errors.New("invalid packet type")
	}

	if len(telegram.Data) < minDataLen {
		return Packet{}, errors.New("data too short")
	}

	if len(telegram.OptData) < minOptDataLen {
		return Packet{}, errors.New("optData too short for destination ID")
	}

	destinationId, _ := device_id.FromByteArray(telegram.OptData[destinationIdOffset : destinationIdOffset+device_id.DeviceIDSize])
	senderId, _ := device_id.FromByteArray(telegram.Data[senderIdOffset : senderIdOffset+device_id.DeviceIDSize])

	rorg := enums.Rorg(telegram.Data[rorgOffset])
	rssi := telegram.OptData[rssiOffset]
	securityLevel := telegram.OptData[securityLevelOffset]
	status := telegram.Data[statusOffset]
	subTelNum := telegram.OptData[subTelNumOffset]
	userData := telegram.Data[userDataOffset:senderIdOffset]

	subTelsOffset := len(telegram.OptData) - (timestampOffset + timestampSize)
	subTels := make([]SubTel, subTelsOffset/subTelSize)
	for i := subTelsOffset; i < len(telegram.OptData); i += subTelSize {
		subTels = append(subTels, SubTel{
			Tick:   telegram.OptData[i],
			Rssi:   telegram.OptData[i+1],
			Status: telegram.OptData[i+2],
		})
	}

	return Packet{
		DestinationID: destinationId,
		Rorg:          rorg,
		Rssi:          rssi,
		SecurityLevel: securityLevel,
		Status:        status,
		SubTelNum:     subTelNum,
		SenderID:      senderId,
		SubTels:       subTels,
		Timestamp:     binary.BigEndian.Uint16(telegram.OptData[timestampOffset : timestampOffset+2]),
		UserData:      userData,
	}, nil
}

func (p Packet) ToEsp3() esp3.Telegram {
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
	optData = append(optData, byte(p.Timestamp&0xFF00>>8), byte(p.Timestamp&0x00FF))

	for _, subTel := range p.SubTels {
		optData = append(optData, subTel.Tick, subTel.Rssi, subTel.Status)
	}

	return esp3.Telegram{
		PacketType: enums.PacketTypeRADIO_ERP1,
		Data:       data,
		OptData:    optData,
	}
}

func (p Packet) Serialize() []byte {
	return p.ToEsp3().Serialize()
}

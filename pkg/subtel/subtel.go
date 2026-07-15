package subtel

import (
	"encoding/binary"
	"errors"

	"github.com/edlundin/enocean-esp3/pkg/deviceid"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
)

type SubTel struct {
	Tick   byte
	Rssi   byte
	Status byte
}

type Packet struct {
	DestinationID deviceid.DeviceID
	Rorg          enums.Rorg
	Rssi          byte
	SecurityLevel byte
	Status        byte
	SubTelNum     byte
	SenderID      deviceid.DeviceID
	SubTels       []SubTel
	Timestamp     uint16
	UserData      []byte
}

// NewPacketFromEsp3 parses a subtelegram packet from an ESP3 telegram.
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
	senderIdOffset := statusOffset - deviceid.DeviceIDSize

	if telegram.PacketType != enums.PacketTypeRADIO_ERP1 {
		return Packet{}, errors.New("invalid packet type")
	}

	if len(telegram.Data) < minDataLen {
		return Packet{}, errors.New("data too short")
	}

	if len(telegram.OptData) < minOptDataLen {
		return Packet{}, errors.New("optData too short for destination ID")
	}

	destinationId, _ := deviceid.FromByteArray(telegram.OptData[destinationIdOffset : destinationIdOffset+deviceid.DeviceIDSize])
	senderId, _ := deviceid.FromByteArray(telegram.Data[senderIdOffset : senderIdOffset+deviceid.DeviceIDSize])

	rorg := enums.Rorg(telegram.Data[rorgOffset])
	rssi := telegram.OptData[rssiOffset]
	securityLevel := telegram.OptData[securityLevelOffset]
	status := telegram.Data[statusOffset]
	subTelNum := telegram.OptData[subTelNumOffset]
	userData := telegram.Data[userDataOffset:senderIdOffset]

	subTelsOffset := timestampOffset + timestampSize
	subTelsData := telegram.OptData[subTelsOffset:]
	subTels := make([]SubTel, 0, len(subTelsData)/subTelSize)
	for i := 0; i < len(subTelsData); i += subTelSize {
		if i+2 < len(subTelsData) {
			subTels = append(subTels, SubTel{
				Tick:   subTelsData[i],
				Rssi:   subTelsData[i+1],
				Status: subTelsData[i+2],
			})
		}
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

// ToEsp3 converts the packet to an ESP3 telegram.
func (p Packet) ToEsp3() esp3.Telegram {
	senderID := p.SenderID.ToArray()
	destinationID := p.DestinationID.ToArray()

	data := make([]byte, 0, 1+len(p.UserData)+deviceid.DeviceIDSize+1)
	data = append(data, byte(p.Rorg))
	data = append(data, p.UserData...)
	data = append(data, senderID[:]...)
	data = append(data, p.Status)

	optData := make([]byte, 0, 3+deviceid.DeviceIDSize)
	optData = append(optData, p.SubTelNum)
	optData = append(optData, destinationID[:]...)
	optData = append(optData, 0xff) // RSSI is unknown when transmitting.
	optData = append(optData, 0x03)
	optData = append(optData, byte(p.Timestamp>>8), byte(p.Timestamp&0xFF))

	for _, subTel := range p.SubTels {
		optData = append(optData, subTel.Tick, subTel.Rssi, subTel.Status)
	}

	return esp3.Telegram{
		PacketType: enums.PacketTypeRADIO_ERP1,
		Data:       data,
		OptData:    optData,
	}
}

// Serialize encodes Packet into its wire representation.
func (p Packet) Serialize() []byte {
	return p.ToEsp3().Serialize()
}

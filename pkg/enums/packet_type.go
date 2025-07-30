package enums

import "errors"

type PacketType uint8

const (
	PacketTypeRADIO_ERP1         PacketType = 0x01
	PacketTypeRESPONSE           PacketType = 0x02
	PacketTypeRADIO_SUB_TEL      PacketType = 0x03
	PacketTypeEVENT              PacketType = 0x04
	PacketTypeCOMMON_COMMAND     PacketType = 0x05
	PacketTypeSMART_ACK_COMMAND  PacketType = 0x06
	PacketTypeREMOTE_MAN_COMMAND PacketType = 0x07
	PacketTypeRADIO_MESSAGE      PacketType = 0x09
	PacketTypeRADIO_ERP2         PacketType = 0x0a
	PacketTypeCONFIG_COMMAND     PacketType = 0x0b
	PacketTypeCOMMAND_ACCEPTED   PacketType = 0x0c
	PacketTypeRADIO_802_15_4     PacketType = 0x10
	PacketTypeCOMMAND_2_4        PacketType = 0x11
)

func ParsePacketTypeFromByte(byte uint8) (PacketType, error) {
	switch byte {
	case 0x01:
		return PacketTypeRADIO_ERP1, nil
	case 0x02:
		return PacketTypeRESPONSE, nil
	case 0x03:
		return PacketTypeRADIO_SUB_TEL, nil
	case 0x04:
		return PacketTypeEVENT, nil
	case 0x05:
		return PacketTypeCOMMON_COMMAND, nil
	case 0x06:
		return PacketTypeSMART_ACK_COMMAND, nil
	case 0x07:
		return PacketTypeREMOTE_MAN_COMMAND, nil
	case 0x09:
		return PacketTypeRADIO_MESSAGE, nil
	case 0x0a:
		return PacketTypeRADIO_ERP2, nil
	case 0x0b:
		return PacketTypeCONFIG_COMMAND, nil
	case 0x0c:
		return PacketTypeCOMMAND_ACCEPTED, nil
	case 0x10:
		return PacketTypeRADIO_802_15_4, nil
	case 0x11:
		return PacketTypeCOMMAND_2_4, nil
	default:
		return 0, errors.New("invalid packet type")
	}
}

func (packetType PacketType) String() string {
	switch packetType {
	case PacketTypeRADIO_ERP1:
		return "RADIO_ERP1"
	case PacketTypeRESPONSE:
		return "RESPONSE"
	case PacketTypeRADIO_SUB_TEL:
		return "RADIO_SUB_TEL"
	case PacketTypeEVENT:
		return "EVENT"
	case PacketTypeCOMMON_COMMAND:
		return "COMMON_COMMAND"
	case PacketTypeSMART_ACK_COMMAND:
		return "SMART_ACK_COMMAND"
	case PacketTypeREMOTE_MAN_COMMAND:
		return "REMOTE_MAN_COMMAND"
	case PacketTypeRADIO_MESSAGE:
		return "RADIO_MESSAGE"
	case PacketTypeRADIO_ERP2:
		return "RADIO_ERP2"
	case PacketTypeCONFIG_COMMAND:
		return "CONFIG_COMMAND"
	case PacketTypeCOMMAND_ACCEPTED:
		return "COMMAND_ACCEPTED"
	case PacketTypeRADIO_802_15_4:
		return "RADIO_802_15_4"
	case PacketTypeCOMMAND_2_4:
		return "COMMAND_2_4"
	default:
		return "UNKNOWN"
	}
}

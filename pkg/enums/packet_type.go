package enums

import "errors"

type PacketType uint8

const (
	PACKET_TYPE_RADIO_ERP1         PacketType = 0x01
	PACKET_TYPE_RESPONSE           PacketType = 0x02
	PACKET_TYPE_RADIO_SUB_TEL      PacketType = 0x03
	PACKET_TYPE_EVENT              PacketType = 0x04
	PACKET_TYPE_COMMON_COMMAND     PacketType = 0x05
	PACKET_TYPE_SMART_ACK_COMMAND  PacketType = 0x06
	PACKET_TYPE_REMOTE_MAN_COMMAND PacketType = 0x07
	PACKET_TYPE_RADIO_MESSAGE      PacketType = 0x09
	PACKET_TYPE_RADIO_ERP2         PacketType = 0x0a
	PACKET_TYPE_CONFIG_COMMAND     PacketType = 0x0b
	PACKET_TYPE_COMMAND_ACCEPTED   PacketType = 0x0c
	PACKET_TYPE_RADIO_802_15_4     PacketType = 0x10
	PACKET_TYPE_COMMAND_2_4        PacketType = 0x11
)

func ParsePacketTypeFromByte(byte uint8) (PacketType, error) {
	switch byte {
	case 0x01:
		return PACKET_TYPE_RADIO_ERP1, nil
	case 0x02:
		return PACKET_TYPE_RESPONSE, nil
	case 0x03:
		return PACKET_TYPE_RADIO_SUB_TEL, nil
	case 0x04:
		return PACKET_TYPE_EVENT, nil
	case 0x05:
		return PACKET_TYPE_COMMON_COMMAND, nil
	case 0x06:
		return PACKET_TYPE_SMART_ACK_COMMAND, nil
	case 0x07:
		return PACKET_TYPE_REMOTE_MAN_COMMAND, nil
	case 0x09:
		return PACKET_TYPE_RADIO_MESSAGE, nil
	case 0x0a:
		return PACKET_TYPE_RADIO_ERP2, nil
	case 0x0b:
		return PACKET_TYPE_CONFIG_COMMAND, nil
	case 0x0c:
		return PACKET_TYPE_COMMAND_ACCEPTED, nil
	case 0x10:
		return PACKET_TYPE_RADIO_802_15_4, nil
	case 0x11:
		return PACKET_TYPE_COMMAND_2_4, nil
	default:
		return 0, errors.New("invalid packet type")
	}
}

func (packetType PacketType) String() string {
	switch packetType {
	case PACKET_TYPE_RADIO_ERP1:
		return "RADIO_ERP1"
	case PACKET_TYPE_RESPONSE:
		return "RESPONSE"
	case PACKET_TYPE_RADIO_SUB_TEL:
		return "RADIO_SUB_TEL"
	case PACKET_TYPE_EVENT:
		return "EVENT"
	case PACKET_TYPE_COMMON_COMMAND:
		return "COMMON_COMMAND"
	case PACKET_TYPE_SMART_ACK_COMMAND:
		return "SMART_ACK_COMMAND"
	case PACKET_TYPE_REMOTE_MAN_COMMAND:
		return "REMOTE_MAN_COMMAND"
	case PACKET_TYPE_RADIO_MESSAGE:
		return "RADIO_MESSAGE"
	case PACKET_TYPE_RADIO_ERP2:
		return "RADIO_ERP2"
	case PACKET_TYPE_CONFIG_COMMAND:
		return "CONFIG_COMMAND"
	case PACKET_TYPE_COMMAND_ACCEPTED:
		return "COMMAND_ACCEPTED"
	case PACKET_TYPE_RADIO_802_15_4:
		return "RADIO_802_15_4"
	case PACKET_TYPE_COMMAND_2_4:
		return "COMMAND_2_4"
	default:
		return "UNKNOWN"
	}
}

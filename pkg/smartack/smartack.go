package smartack

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/edlundin/enocean-esp3/pkg/deviceid"
	"github.com/edlundin/enocean-esp3/pkg/eep"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/erp1"
)

const (
	LearnReplyIndex       byte = 0x01
	LearnAcknowledgeIndex byte = 0x02
	SignalMailboxEmpty    byte = 0x01
	SignalMailboxMissing  byte = 0x02
	SignalReset           byte = 0x03
)

type RequestCode byte

const (
	RequestCandidateMailboxSpace   RequestCode = 0x00
	RequestCandidateNoMailboxSpace RequestCode = 0x01
	RequestNoCandidateMailboxSpace RequestCode = 0x02
	RequestNoCandidateNoSpace      RequestCode = 0x03
	RequestDefaultSensor           RequestCode = 0x1f
)

type AckCode byte

func (a AckCode) Class() string {
	switch {
	case a == 0x00:
		return "first learn-in"
	case a >= 0x01 && a <= 0x0f:
		return "repeated learn-in"
	case a >= 0x10 && a <= 0x1f:
		return "failed learn-in"
	case a == 0x20:
		return "complete learn-out"
	case a >= 0x21 && a <= 0x2f:
		return "partial learn-out"
	default:
		return "application-specific"
	}
}

type Message interface {
	ERP1(sender deviceid.DeviceID) erp1.Packet
}

type LearnRequest struct {
	RequestCode    RequestCode
	ManufacturerID uint16 // 11 bits
	EEP            eep.EEP
	RSSI           byte
	RepeaterID     deviceid.DeviceID
}

type LearnReply struct {
	ResponseTime uint16
	AckCode      AckCode
	SensorID     deviceid.DeviceID
}

type LearnAcknowledge struct {
	ResponseTime uint16
	AckCode      AckCode
	MailboxIndex byte
}

type LearnReclaim struct{ Data bool }
type DataReclaim struct{ MailboxIndex byte }
type Signal struct{ Index byte }

func Parse(p erp1.Packet) (Message, error) {
	switch p.Rorg {
	case enums.RorgSM_LRN_REQ:
		return parseLearnRequest(p.UserData)
	case enums.RorgSM_LRN_ANS:
		return parseLearnAnswer(p.UserData)
	case enums.RorgSM_REC:
		return parseReclaim(p.UserData)
	case enums.RorgSIGNAL:
		return parseSignal(p.UserData)
	default:
		return nil, fmt.Errorf("not a Smart Ack RORG: %s", p.Rorg)
	}
}

func parseLearnRequest(b []byte) (LearnRequest, error) {
	if len(b) != 10 {
		return LearnRequest{}, fmt.Errorf("learn request length %d, want 10", len(b))
	}
	head := binary.BigEndian.Uint16(b[:2])
	rep, _ := deviceid.FromByteArray(b[6:10])
	return LearnRequest{RequestCode: RequestCode(head >> 11), ManufacturerID: head & 0x07ff, EEP: eep.EEP{Rorg: enums.Rorg(b[2]), Func: b[3], Type: b[4]}, RSSI: b[5], RepeaterID: rep}, nil
}

func parseLearnAnswer(b []byte) (Message, error) {
	if len(b) == 0 {
		return nil, errors.New("empty learn answer")
	}
	switch b[0] {
	case LearnReplyIndex:
		if len(b) != 8 {
			return nil, fmt.Errorf("learn reply length %d, want 8", len(b))
		}
		id, _ := deviceid.FromByteArray(b[4:8])
		return LearnReply{ResponseTime: binary.BigEndian.Uint16(b[1:3]), AckCode: AckCode(b[3]), SensorID: id}, nil
	case LearnAcknowledgeIndex:
		if len(b) != 5 {
			return nil, fmt.Errorf("learn acknowledge length %d, want 5", len(b))
		}
		return LearnAcknowledge{ResponseTime: binary.BigEndian.Uint16(b[1:3]), AckCode: AckCode(b[3]), MailboxIndex: b[4]}, nil
	default:
		return nil, fmt.Errorf("unknown learn answer index 0x%02x", b[0])
	}
}

func parseReclaim(b []byte) (Message, error) {
	if len(b) != 1 {
		return nil, fmt.Errorf("reclaim length %d, want 1", len(b))
	}
	if b[0]&0x80 == 0 {
		return LearnReclaim{Data: b[0]&1 == 1}, nil
	}
	return DataReclaim{MailboxIndex: b[0] & 0x7f}, nil
}

func parseSignal(b []byte) (Signal, error) {
	if len(b) != 1 {
		return Signal{}, fmt.Errorf("signal length %d, want 1", len(b))
	}
	if b[0] < SignalMailboxEmpty || b[0] > SignalReset {
		return Signal{}, fmt.Errorf("unknown signal index 0x%02x", b[0])
	}
	return Signal{Index: b[0]}, nil
}

func (m LearnRequest) ERP1(sender deviceid.DeviceID) erp1.Packet {
	head := uint16(m.RequestCode&0x1f)<<11 | (m.ManufacturerID & 0x07ff)
	b := make([]byte, 10)
	binary.BigEndian.PutUint16(b[:2], head)
	b[2], b[3], b[4], b[5] = byte(m.EEP.Rorg), m.EEP.Func, m.EEP.Type, m.RSSI
	rep := m.RepeaterID.ToArray()
	copy(b[6:], rep[:])
	return packet(enums.RorgSM_LRN_REQ, b, sender, 3)
}

func (m LearnReply) ERP1(sender deviceid.DeviceID) erp1.Packet {
	b := make([]byte, 8)
	b[0] = LearnReplyIndex
	binary.BigEndian.PutUint16(b[1:3], m.ResponseTime)
	b[3] = byte(m.AckCode)
	id := m.SensorID.ToArray()
	copy(b[4:], id[:])
	return packet(enums.RorgSM_LRN_ANS, b, sender, 3)
}

func (m LearnAcknowledge) ERP1(sender deviceid.DeviceID) erp1.Packet {
	b := []byte{LearnAcknowledgeIndex, 0, 0, byte(m.AckCode), m.MailboxIndex}
	binary.BigEndian.PutUint16(b[1:3], m.ResponseTime)
	return packet(enums.RorgSM_LRN_ANS, b, sender, 1)
}

func (m LearnReclaim) ERP1(sender deviceid.DeviceID) erp1.Packet {
	var v byte
	if m.Data {
		v = 1
	}
	return packet(enums.RorgSM_REC, []byte{v}, sender, 1)
}

func (m DataReclaim) ERP1(sender deviceid.DeviceID) erp1.Packet {
	return packet(enums.RorgSM_REC, []byte{0x80 | (m.MailboxIndex & 0x7f)}, sender, 1)
}

func (m Signal) ERP1(sender deviceid.DeviceID) erp1.Packet {
	return packet(enums.RorgSIGNAL, []byte{m.Index}, sender, 1)
}

func packet(r enums.Rorg, data []byte, sender deviceid.DeviceID, subTel byte) erp1.Packet {
	return erp1.Packet{Rorg: r, UserData: data, SenderID: sender, SubTelNum: subTel, SecurityLevel: 3, Rssi: 0xff, DestinationID: deviceid.BroadcastId()}
}

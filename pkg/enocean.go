package pkg

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/erp1"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
	"github.com/edlundin/enocean-esp3/pkg/event"
	"github.com/edlundin/enocean-esp3/pkg/gp"
	"github.com/edlundin/enocean-esp3/pkg/reman"
	"github.com/edlundin/enocean-esp3/pkg/response"
	"github.com/edlundin/enocean-esp3/pkg/smartack"
	"go.bug.st/serial"
)

// GetSerialPortList returns the available serial port names.
func GetSerialPortList() ([]string, error) { return serial.GetPortsList() }

type Message struct {
	Kind string
	ESP3 esp3.Telegram
	ERP1 *erp1.Packet
	Data any
	Err  error
}

// Channels contains independent, buffered, best-effort event streams.
// A message is dropped from an individual stream when that stream's buffer is
// full; slow consumers do not block serial parsing or other streams.
type Channels struct {
	All        <-chan Message
	ESP3       <-chan esp3.Telegram
	ERP1       <-chan erp1.Packet
	Response   <-chan response.Packet
	Event      <-chan event.Event
	SmartAck   <-chan smartack.Message
	ReMan      <-chan reman.Message
	ReManPart  <-chan reman.Part
	GPHeader   <-chan any
	Unparsed   <-chan Message
	ParseError <-chan Message
}

type channelSet struct {
	all        chan Message
	esp3       chan esp3.Telegram
	erp1       chan erp1.Packet
	response   chan response.Packet
	event      chan event.Event
	smartAck   chan smartack.Message
	reman      chan reman.Message
	remanPart  chan reman.Part
	gpHeader   chan any
	unparsed   chan Message
	parseError chan Message
}

// newChannelSet constructs ChannelSet.
func newChannelSet(size int) (*channelSet, *Channels) {
	set := &channelSet{
		all:        make(chan Message, size),
		esp3:       make(chan esp3.Telegram, size),
		erp1:       make(chan erp1.Packet, size),
		response:   make(chan response.Packet, size),
		event:      make(chan event.Event, size),
		smartAck:   make(chan smartack.Message, size),
		reman:      make(chan reman.Message, size),
		remanPart:  make(chan reman.Part, size),
		gpHeader:   make(chan any, size),
		unparsed:   make(chan Message, size),
		parseError: make(chan Message, size),
	}
	return set, &Channels{All: set.all, ESP3: set.esp3, ERP1: set.erp1, Response: set.response, Event: set.event, SmartAck: set.smartAck, ReMan: set.reman, ReManPart: set.remanPart, GPHeader: set.gpHeader, Unparsed: set.unparsed, ParseError: set.parseError}
}

// close closes all parser output channels.
func (c *channelSet) close() {
	close(c.all)
	close(c.esp3)
	close(c.erp1)
	close(c.response)
	close(c.event)
	close(c.smartAck)
	close(c.reman)
	close(c.remanPart)
	close(c.gpHeader)
	close(c.unparsed)
	close(c.parseError)
}

var serialOpen = serial.Open

// OpenSerialPort opens a serial port and starts ESP3 parsing.
func OpenSerialPort(ctx context.Context, portPath string) (serial.Port, *Channels, error) {
	portSettings := &serial.Mode{
		BaudRate: 57600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	port, err := serialOpen(portPath, portSettings)

	if err != nil {
		return nil, nil, err
	}

	err = port.SetReadTimeout(time.Second * 2)

	if err != nil {
		_ = port.Close()
		return nil, nil, err
	}

	set, channels := newChannelSet(64)
	go parser(ctx, port, set)

	return port, channels, nil
}

// parser parses r.
func parser(ctx context.Context, serialPort serial.Port, channels *channelSet) {
	defer channels.close()
	type ParserState uint8

	const (
		ParserStateWaitingForSyncByte ParserState = iota
		ParserStateWaitingForHeader
		ParserStateWaitingForCrc8H
		ParserStateWaitingForData
		ParserStateWaitingForCrc8D
	)

	const interByteTimeout = time.Millisecond * 100
	const syncByte = 0x55
	const dataLengthOffset = 0
	const dataLengthLen = 2
	const optDataLengthOffset = dataLengthOffset + dataLengthLen
	const packetTypeOffset = 3
	const headerLen = 4

	lastByteReceivedTime := time.Now()
	parserState := ParserStateWaitingForSyncByte
	readBuffer := make([]uint8, 64)

	parserBuffer := make([]uint8, 0)
	parserCrc := uint8(0)
	parserDataLen := 0
	parserOptDataLen := 0
	parserPacketType := uint8(0)
	remanMessages := newReManAssembler(remanChainPeriod)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			byteReceived, err := serialPort.Read(readBuffer)

			if err != nil {
				return
			}

			if byteReceived == 0 {
				continue
			}

			if time.Since(lastByteReceivedTime) >= interByteTimeout {
				parserState = ParserStateWaitingForSyncByte
			}

			for i := 0; i < byteReceived; i++ {
				parserByte := readBuffer[i]

				switch parserState {
				case ParserStateWaitingForSyncByte:
					if parserByte == syncByte {
						parserState = ParserStateWaitingForHeader
						parserBuffer = make([]uint8, 0)
						parserCrc = 0
					}
				case ParserStateWaitingForHeader:
					parserBuffer = append(parserBuffer, parserByte)
					parserCrc = esp3.ComputeCrc8(parserByte, parserCrc)

					if len(parserBuffer) == headerLen { // Header received
						parserState = ParserStateWaitingForCrc8H
					}
				case ParserStateWaitingForCrc8H:
					// CRC8H invalid
					if parserCrc != parserByte {
						syncByteIdx := bytes.IndexByte(parserBuffer, syncByte)

						// Header and CRC8H does not contain the sync code, wait for new packet to start
						if syncByteIdx < 0 && parserByte != syncByte {
							parserState = ParserStateWaitingForSyncByte
							break
						}

						// Header does not have sync code but CRC8H does, reset state, this is a new packet
						if syncByteIdx < 0 {
							parserState = ParserStateWaitingForHeader
							parserBuffer = make([]uint8, 0)
							parserCrc = 0
							break
						}

						parserBuffer = append(parserBuffer[:0], parserBuffer[syncByteIdx+1:]...)
						parserBuffer = append(parserBuffer, parserByte)
						parserCrc = esp3.ComputeCrcSlice(parserBuffer)

						if len(parserBuffer) < headerLen {
							parserState = ParserStateWaitingForHeader
							break
						}

						break
					}

					parserDataLen = int(binary.BigEndian.Uint16(parserBuffer[dataLengthOffset : dataLengthOffset+dataLengthLen]))
					parserOptDataLen = int(parserBuffer[optDataLengthOffset])
					parserPacketType = parserBuffer[packetTypeOffset]

					parserState = ParserStateWaitingForData
					if parserDataLen+parserOptDataLen == 0 {
						parserState = ParserStateWaitingForCrc8D
					}
					parserBuffer = make([]uint8, 0)
					parserCrc = 0
				case ParserStateWaitingForData:
					parserBuffer = append(parserBuffer, parserByte)
					parserCrc = esp3.ComputeCrc8(parserByte, parserCrc)

					if len(parserBuffer) == parserDataLen+parserOptDataLen {
						parserState = ParserStateWaitingForCrc8D
					}
				case ParserStateWaitingForCrc8D:
					parserState = ParserStateWaitingForSyncByte
					if parserByte != parserCrc {
						if parserByte == syncByte {
							parserState = ParserStateWaitingForHeader
							parserBuffer = make([]uint8, 0)
							parserCrc = 0
						}
						break
					}

					packetType, err := enums.ParsePacketTypeFromByte(parserPacketType)
					if err != nil {
						break
					}

					if !publish(ctx, channels, parseTelegram(remanMessages, esp3.NewTelegramFromData(packetType, parserBuffer[:parserDataLen], parserBuffer[parserDataLen:]))) {
						return
					}
				default:
					parserState = ParserStateWaitingForSyncByte
				}
			}

			lastByteReceivedTime = time.Now()
		}
	}
}

const remanChainPeriod = time.Second

type remanKey struct {
	seq          byte
	source, dest uint32
}

type remanBuffer struct {
	parts   []reman.Part
	updated time.Time
}

type remanAssembler struct {
	ttl     time.Duration
	buffers map[remanKey]remanBuffer
}

// newReManAssembler constructs ReManAssembler.
func newReManAssembler(ttl time.Duration) *remanAssembler {
	return &remanAssembler{ttl: ttl, buffers: map[remanKey]remanBuffer{}}
}

// add adds a ReMan chain part to the assembler.
func (a *remanAssembler) add(part reman.Part) ([]Message, error) {
	now := time.Now()
	a.expire(now)
	key := remanKey{seq: part.Seq, source: uint32(part.SourceID), dest: uint32(part.DestinationID)}
	buf := a.buffers[key]
	for _, existing := range buf.parts {
		if existing.Index == part.Index {
			delete(a.buffers, key)
			return nil, fmt.Errorf("duplicate ReMan part index %d", part.Index)
		}
	}
	buf.parts = append(buf.parts, part)
	buf.updated = now

	msg, ok, err := reman.Merge(append([]reman.Part(nil), buf.parts...))
	if err != nil {
		delete(a.buffers, key)
		return nil, err
	}
	if !ok {
		a.buffers[key] = buf
		return nil, nil
	}
	delete(a.buffers, key)
	return []Message{{Kind: "reman", Data: msg}}, nil
}

// expire removes stale ReMan chains.
func (a *remanAssembler) expire(now time.Time) {
	for key, buf := range a.buffers {
		if now.Sub(buf.updated) >= a.ttl {
			delete(a.buffers, key)
		}
	}
}

// parseTelegram parses Telegram.
func parseTelegram(remanMessages *remanAssembler, t esp3.Telegram) []Message {
	messages := []Message{{Kind: "esp3", ESP3: t, Data: t}}

	switch t.PacketType {
	case enums.PacketTypeRADIO_ERP1:
		p, err := erp1.NewPacketFromEsp3(t)
		if err != nil {
			return append(messages, Message{Kind: "parse_error", ESP3: t, Err: err})
		}
		messages = append(messages, Message{Kind: "erp1", ESP3: t, ERP1: &p, Data: p})
		messages = append(messages, parseERP1(remanMessages, t, p)...)
	case enums.PacketTypeRESPONSE:
		p, err := response.NewPacketFromEsp3(t)
		if err != nil {
			return append(messages, Message{Kind: "parse_error", ESP3: t, Err: err})
		}
		messages = append(messages, Message{Kind: "response", ESP3: t, Data: p})
	case enums.PacketTypeEVENT:
		p, err := event.NewPacketFromEsp3(t)
		if err != nil {
			return append(messages, Message{Kind: "parse_error", ESP3: t, Err: err})
		}
		messages = append(messages, Message{Kind: "event", ESP3: t, Data: p})
	default:
		messages = append(messages, Message{Kind: "unparsed", ESP3: t, Data: t})
	}

	return messages
}

// parseERP1 parses ERP1.
func parseERP1(remanMessages *remanAssembler, t esp3.Telegram, p erp1.Packet) []Message {
	switch {
	case p.Rorg == enums.RorgSYS_EX:
		part, err := reman.ParsePacket(p)
		if err != nil {
			return []Message{{Kind: "parse_error", ESP3: t, ERP1: &p, Err: err}}
		}
		out := []Message{{Kind: "reman_part", ESP3: t, ERP1: &p, Data: part}}
		messages, err := remanMessages.add(part)
		if err != nil {
			return append(out, Message{Kind: "parse_error", ESP3: t, ERP1: &p, Err: err})
		}
		for i := range messages {
			messages[i].ESP3 = t
			messages[i].ERP1 = &p
		}
		return append(out, messages...)
	case gp.IsRorg(p.Rorg):
		header, err := parseGPHeader(p)
		if err != nil {
			return []Message{{Kind: "parse_error", ESP3: t, ERP1: &p, Err: err}}
		}
		return []Message{{Kind: "gp_header", ESP3: t, ERP1: &p, Data: header}}
	default:
		msg, err := smartack.Parse(p)
		if err == nil {
			return []Message{{Kind: "smart_ack", ESP3: t, ERP1: &p, Data: msg}}
		}
		return nil
	}
}

// parseGPHeader parses GPHeader.
func parseGPHeader(p erp1.Packet) (any, error) {
	switch p.Rorg {
	case enums.RorgGP_TI:
		return gp.DecodeRequestHeader(p.UserData)
	case enums.RorgGP_TR:
		return gp.DecodeResponseHeader(p.UserData)
	default:
		return p, nil
	}
}

// publish dispatches parsed messages to output channels.
func publish(ctx context.Context, channels *channelSet, messages []Message) bool {
	for _, msg := range messages {
		if !send(ctx, channels.all, msg) {
			return false
		}
		switch msg.Kind {
		case "esp3":
			if !send(ctx, channels.esp3, msg.Data.(esp3.Telegram)) {
				return false
			}
		case "erp1":
			if !send(ctx, channels.erp1, msg.Data.(erp1.Packet)) {
				return false
			}
		case "response":
			if !send(ctx, channels.response, msg.Data.(response.Packet)) {
				return false
			}
		case "event":
			if !send(ctx, channels.event, msg.Data.(event.Event)) {
				return false
			}
		case "smart_ack":
			if !send(ctx, channels.smartAck, msg.Data.(smartack.Message)) {
				return false
			}
		case "reman":
			if !send(ctx, channels.reman, msg.Data.(reman.Message)) {
				return false
			}
		case "reman_part":
			if !send(ctx, channels.remanPart, msg.Data.(reman.Part)) {
				return false
			}
		case "gp_header":
			if !send(ctx, channels.gpHeader, msg.Data) {
				return false
			}
		case "unparsed":
			if !send(ctx, channels.unparsed, msg) {
				return false
			}
		case "parse_error":
			if !send(ctx, channels.parseError, msg) {
				return false
			}
		}
	}
	return true
}

// send delivers a value without blocking the parser.
func send[T any](ctx context.Context, ch chan<- T, v T) bool {
	select {
	case ch <- v:
		return true
	case <-ctx.Done():
		return false
	default:
		return true
	}
}

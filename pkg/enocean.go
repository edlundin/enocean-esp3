package pkg

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
	"go.bug.st/serial"
)

func GetSerialPortList() ([]string, error) {
	ports, err := serial.GetPortsList()

	if err != nil {
		return nil, err
	}

	return ports, nil
}

func OpenSerialPort(ctx context.Context, portPath string) (serial.Port, error) {
	portSettings := &serial.Mode{
		BaudRate: 57600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(portPath, portSettings)

	if err != nil {
		return nil, err
	}

	err = port.SetReadTimeout(time.Second * 2)

	if err != nil {
		return nil, err
	}

	go parser(ctx, port)

	//TODO handle channel for esp3 telegrams and add cancel with context

	return port, nil
}

func parser(ctx context.Context, serialPort serial.Port) {
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
	parserDataLen := uint16(0)
	parserOptDataLen := uint8(0)
	parserPacketType := uint8(0)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			byteReceived, err := serialPort.Read(readBuffer)

			if err != nil {
				fmt.Println(fmt.Errorf("error reading from serial port: %w", err))
				continue
			}

			if byteReceived == 0 {
				fmt.Println("no bytes received")
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
					const syncByteIdxInit = -1

					// CRC8H invalid
					if parserCrc != parserByte {
						syncByteIdx := syncByteIdxInit

						for i := 0; i < headerLen; i++ {
							// Header have a sync byte, indicates the start of the new packet
							if parserBuffer[i] == syncByte {
								syncByteIdx = i + 1
								break
							}
						}

						// Header and CRC8H does not contain the sync code, wait for new packet to start
						if syncByteIdx == syncByteIdxInit && parserByte != syncByte {
							parserState = ParserStateWaitingForSyncByte
							break
						}

						// Header does not have sync code but CRC8H does, reset state, this is a new packet
						if syncByteIdx == syncByteIdxInit && parserByte == syncByte {
							parserState = ParserStateWaitingForHeader
							parserBuffer = make([]uint8, 0)
							parserCrc = 0
							break
						}

						parserCrc = 0
						tmpBuffer := make([]uint8, 0)

						for i := 0; i < headerLen-syncByteIdx; i++ {
							tmpBuffer = append(tmpBuffer, parserBuffer[syncByteIdx+i])
							parserCrc = esp3.ComputeCrc8(parserBuffer[i], parserCrc)
						}

						parserBuffer = append(tmpBuffer, parserByte)
						parserCrc = esp3.ComputeCrc8(parserByte, parserCrc)

						if len(parserBuffer) < headerLen {
							parserState = ParserStateWaitingForHeader
							break
						}

						break
					}

					parserDataLen = binary.BigEndian.Uint16(parserBuffer[dataLengthOffset : dataLengthOffset+dataLengthLen])
					parserOptDataLen = parserBuffer[optDataLengthOffset]
					parserPacketType = parserBuffer[packetTypeOffset]

					// Data length fields are invalid
					if parserDataLen+uint16(parserOptDataLen) == 0 {
						if parserByte == syncByte { // Sync already received
							parserState = ParserStateWaitingForHeader
							parserBuffer = make([]uint8, 0)
							parserCrc = 0
							break
						}

						parserState = ParserStateWaitingForSyncByte
						break
					}

					parserState = ParserStateWaitingForData
					parserBuffer = make([]uint8, 0)
					parserCrc = 0
				case ParserStateWaitingForData:
					parserBuffer = append(parserBuffer, parserByte)
					parserCrc = esp3.ComputeCrc8(parserByte, parserCrc)

					if uint16(len(parserBuffer)) == parserDataLen+uint16(parserOptDataLen) {
						parserState = ParserStateWaitingForCrc8D
					}
				case ParserStateWaitingForCrc8D:
					// Parsing done, packet invalid, sync byte already received
					if parserByte == syncByte {
						parserState = ParserStateWaitingForHeader
						parserBuffer = make([]uint8, 0)
						parserCrc = 0

						break
					}

					parserState = ParserStateWaitingForSyncByte

					// Parsing done, packet valid, calling callback
					if parserByte == parserCrc {
						telegram :=
							esp3.NewEsp3TelegramFromData(enums.PacketType(parserPacketType), parserBuffer[:parserDataLen], parserBuffer[parserDataLen:]) //TODO: check packet type

						fmt.Println(telegram)
					}
				default:
					parserState = ParserStateWaitingForSyncByte
				}
			}

			lastByteReceivedTime = time.Now()
		}
	}
}

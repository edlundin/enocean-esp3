package pkg

import (
	"context"
	"errors"
	"io"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/edlundin/enocean-esp3/pkg/deviceid"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/erp1"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
	"github.com/edlundin/enocean-esp3/pkg/gp"
	"github.com/edlundin/enocean-esp3/pkg/reman"
	"github.com/edlundin/enocean-esp3/pkg/response"
	"github.com/edlundin/enocean-esp3/pkg/smartack"
	"go.bug.st/serial"
)

type fakePort struct {
	reads             [][]byte
	setReadTimeoutErr error
	closed            bool
}

// Read reads the value.
func (p *fakePort) Read(b []byte) (int, error) {
	if len(p.reads) == 0 {
		time.Sleep(time.Millisecond)
		return 0, nil
	}
	n := copy(b, p.reads[0])
	p.reads = p.reads[1:]
	return n, nil
}

// Write writes the value.
func (p *fakePort) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }

// SetMode updates Mode.
func (p *fakePort) SetMode(*serial.Mode) error { return nil }

// Drain waits for pending serial writes to complete.
func (p *fakePort) Drain() error { return nil }

// ResetInputBuffer clears buffered serial input.
func (p *fakePort) ResetInputBuffer() error { return nil }

// ResetOutputBuffer clears buffered serial output.
func (p *fakePort) ResetOutputBuffer() error { return nil }

// SetDTR updates DTR.
func (p *fakePort) SetDTR(bool) error { return nil }

// SetRTS updates RTS.
func (p *fakePort) SetRTS(bool) error { return nil }

// GetModemStatusBits returns ModemStatusBits.
func (p *fakePort) GetModemStatusBits() (*serial.ModemStatusBits, error) { return nil, nil }

// SetReadTimeout updates ReadTimeout.
func (p *fakePort) SetReadTimeout(time.Duration) error { return p.setReadTimeoutErr }

// Close closes the serial port.
func (p *fakePort) Close() error {
	p.closed = true
	return nil
}

// Break requests a serial break for the supplied duration.
func (p *fakePort) Break(time.Duration) error { return nil }

type permanentErrorPort struct {
	fakePort
	readCalls atomic.Int32
	unblock   chan struct{}
}

// Read reads the value.
func (p *permanentErrorPort) Read([]byte) (int, error) {
	if p.readCalls.Add(1) == 1 {
		return 0, io.ErrUnexpectedEOF
	}
	<-p.unblock
	return 0, nil
}

// TestOpenSerialPortClosesPortWhenReadTimeoutFails verifies OpenSerialPortClosesPortWhenReadTimeoutFails behavior.
func TestOpenSerialPortClosesPortWhenReadTimeoutFails(t *testing.T) {
	wantErr := errors.New("timeout setup failed")
	port := &fakePort{setReadTimeoutErr: wantErr}
	oldOpen := serialOpen
	serialOpen = func(string, *serial.Mode) (serial.Port, error) { return port, nil }
	t.Cleanup(func() { serialOpen = oldOpen })

	gotPort, channels, err := OpenSerialPort(context.Background(), "fake")
	if !errors.Is(err, wantErr) || gotPort != nil || channels != nil || !port.closed {
		t.Fatalf("port=%v channels=%v closed=%v err=%v", gotPort, channels, port.closed, err)
	}
}

// TestParserStopsAfterPermanentReadError verifies ParserStopsAfterPermanentReadError behavior.
func TestParserStopsAfterPermanentReadError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	port := &permanentErrorPort{unblock: make(chan struct{})}
	defer func() {
		cancel()
		close(port.unblock)
	}()
	set, channels := newChannelSet(1)
	go parser(ctx, port, set)
	select {
	case _, ok := <-channels.All:
		if ok {
			t.Fatal("All channel remained open")
		}
	case <-time.After(time.Second):
		t.Fatal("parser did not stop after read error")
	}
	if got := port.readCalls.Load(); got != 1 {
		t.Fatalf("Read called %d times, want 1", got)
	}
}

// TestParserDispatchesTelegramsAndClosesOnCancel verifies ParserDispatchesTelegramsAndClosesOnCancel behavior.
func TestParserDispatchesTelegramsAndClosesOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	want := esp3.NewTelegramFromData(enums.PacketTypeRADIO_ERP1, []byte{0xd2}, []byte{0x00})
	set, channels := newChannelSet(4)

	go parser(ctx, &fakePort{reads: [][]byte{want.Serialize()}}, set)

	select {
	case got := <-channels.ESP3:
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("telegram mismatch: got %+v want %+v", got, want)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for telegram")
	}

	cancel()
	select {
	case _, ok := <-channels.ESP3:
		if ok {
			t.Fatal("telegram channel still open after cancel")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for parser cancellation")
	}
}

// TestParserDispatchesParsedERP1 verifies ParserDispatchesParsedERP1 behavior.
func TestParserDispatchesParsedERP1(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	telegram := esp3.NewTelegramFromData(enums.PacketTypeRADIO_ERP1, []byte{0xd2, 0x01, 0, 0, 0, 1, 0}, []byte{1, 0xff, 0xff, 0xff, 0xff, 0x40, 3})
	set, channels := newChannelSet(4)

	go parser(ctx, &fakePort{reads: [][]byte{telegram.Serialize()}}, set)

	select {
	case got := <-channels.ERP1:
		if got.Rorg != enums.RorgVLD || !reflect.DeepEqual(got.UserData, []byte{0x01}) {
			t.Fatalf("bad ERP1: %+v", got)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for ERP1 packet")
	}
}

// TestParserDispatchesMergedReMan verifies ParserDispatchesMergedReMan behavior.
func TestParserDispatchesMergedReMan(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	want := reman.Message{Seq: 1, ManufacturerID: reman.ManufacturerID, Function: reman.FuncSetCode, Payload: []byte{1, 2, 3, 4, 5}, SourceID: deviceid.DeviceID(1), DestinationID: deviceid.BroadcastId()}
	packets, err := want.Packets()
	if err != nil {
		t.Fatal(err)
	}
	reads := make([][]byte, len(packets))
	for i, packet := range packets {
		reads[i] = packet.Serialize()
	}
	set, channels := newChannelSet(8)

	go parser(ctx, &fakePort{reads: reads}, set)

	select {
	case got := <-channels.ReMan:
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("ReMan mismatch: got %+v want %+v", got, want)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for ReMan message")
	}
}

// TestParserSkipsInvalidPacketType verifies ParserSkipsInvalidPacketType behavior.
func TestParserSkipsInvalidPacketType(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bad := esp3.NewTelegramFromData(enums.PacketType(0xff), []byte{0x01}, nil).Serialize()
	set, channels := newChannelSet(4)

	go parser(ctx, &fakePort{reads: [][]byte{bad}}, set)

	select {
	case got := <-channels.ESP3:
		t.Fatalf("unexpected telegram: %+v", got)
	case <-time.After(20 * time.Millisecond):
	}
}

// TestParserAcceptsCRC8DEqualToSyncByte verifies ParserAcceptsCRC8DEqualToSyncByte behavior.
func TestParserAcceptsCRC8DEqualToSyncByte(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	want := esp3.NewTelegramFromData(enums.PacketTypeCOMMON_COMMAND, []byte{0xc5}, nil)
	serialized := want.Serialize()
	if got := serialized[len(serialized)-1]; got != 0x55 {
		t.Fatalf("test vector CRC8D = %#x, want 0x55", got)
	}
	set, channels := newChannelSet(4)
	go parser(ctx, &fakePort{reads: [][]byte{serialized}}, set)
	select {
	case got := <-channels.ESP3:
		if got.PacketType != want.PacketType || !reflect.DeepEqual(got.Data, want.Data) || len(got.OptData) != 0 {
			t.Fatalf("telegram mismatch: got %+v want %+v", got, want)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for telegram")
	}
}

// TestParserStreamsZeroPayloadFrame verifies ParserStreamsZeroPayloadFrame behavior.
func TestParserStreamsZeroPayloadFrame(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	zero := esp3.NewTelegramFromData(enums.PacketTypeCOMMON_COMMAND, []byte{}, []byte{})
	next := esp3.NewTelegramFromData(enums.PacketTypeCOMMON_COMMAND, []byte{1}, []byte{})
	set, channels := newChannelSet(4)
	go parser(ctx, &fakePort{reads: [][]byte{append(zero.Serialize(), next.Serialize()...)}}, set)

	for _, want := range []esp3.Telegram{zero, next} {
		select {
		case got := <-channels.ESP3:
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("telegram mismatch: got %+v want %+v", got, want)
			}
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for telegram")
		}
	}
}

// TestParserDropsBadCRC8HThenResyncs verifies ParserDropsBadCRC8HThenResyncs behavior.
func TestParserDropsBadCRC8HThenResyncs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bad := esp3.NewTelegramFromData(enums.PacketTypeRESPONSE, []byte{byte(enums.ReturnCodeSUCCESS)}, nil).Serialize()
	bad[5] ^= 0xff
	want := esp3.NewTelegramFromData(enums.PacketTypeRESPONSE, []byte{byte(enums.ReturnCodeERROR)}, nil)
	set, channels := newChannelSet(4)

	go parser(ctx, &fakePort{reads: [][]byte{append(bad, want.Serialize()...)}}, set)

	select {
	case got := <-channels.Response:
		if got.Code != enums.ReturnCodeERROR {
			t.Fatalf("got response code %s", got.Code)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for parser resync")
	}
}

// TestParserResyncsFromSyncByteInsideBadHeader verifies ParserResyncsFromSyncByteInsideBadHeader behavior.
func TestParserResyncsFromSyncByteInsideBadHeader(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	want := esp3.NewTelegramFromData(enums.PacketTypeRESPONSE, []byte{byte(enums.ReturnCodeERROR)}, nil)
	packet := want.Serialize()
	stream := append([]byte{0x55, 0x00}, packet[:4]...)
	stream = append(stream, packet[4:]...)
	set, channels := newChannelSet(4)
	go parser(ctx, &fakePort{reads: [][]byte{stream}}, set)
	select {
	case got := <-channels.Response:
		if got.Code != enums.ReturnCodeERROR {
			t.Fatalf("got response code %s", got.Code)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for parser resync")
	}
}

// TestParseTelegramResponseAndUnparsed verifies ParseTelegramResponseAndUnparsed behavior.
func TestParseTelegramResponseAndUnparsed(t *testing.T) {
	remanMessages := newReManAssembler(time.Second)
	resp := esp3.NewTelegramFromData(enums.PacketTypeRESPONSE, []byte{byte(enums.ReturnCodeSUCCESS), 1, 2}, []byte{3})
	msgs := parseTelegram(remanMessages, resp)
	if len(msgs) != 2 || msgs[0].Kind != "esp3" || msgs[1].Kind != "response" {
		t.Fatalf("messages = %#v", msgs)
	}
	if got := msgs[1].Data.(response.Packet); got.Code != enums.ReturnCodeSUCCESS || !reflect.DeepEqual(got.Data, []byte{1, 2}) || !reflect.DeepEqual(got.OptData, []byte{3}) {
		t.Fatalf("response = %#v", got)
	}

	unknown := esp3.NewTelegramFromData(enums.PacketTypeCOMMON_COMMAND, []byte{1}, nil)
	msgs = parseTelegram(remanMessages, unknown)
	if len(msgs) != 2 || msgs[1].Kind != "unparsed" {
		t.Fatalf("unparsed messages = %#v", msgs)
	}
}

// TestParseTelegramReportsParseErrors verifies ParseTelegramReportsParseErrors behavior.
func TestParseTelegramReportsParseErrors(t *testing.T) {
	cases := []esp3.Telegram{
		esp3.NewTelegramFromData(enums.PacketTypeRESPONSE, nil, nil),
		esp3.NewTelegramFromData(enums.PacketTypeRADIO_ERP1, []byte{0xd2}, nil),
	}
	for _, tc := range cases {
		msgs := parseTelegram(newReManAssembler(time.Second), tc)
		if msgs[len(msgs)-1].Kind != "parse_error" || msgs[len(msgs)-1].Err == nil {
			t.Fatalf("messages = %#v", msgs)
		}
	}
}

// TestParseERP1SmartAckAndGPHeader verifies ParseERP1SmartAckAndGPHeader behavior.
func TestParseERP1SmartAckAndGPHeader(t *testing.T) {
	sa := smartack.DataReclaim{MailboxIndex: 3}.ERP1(deviceid.DeviceID(1))
	msgs := parseERP1(newReManAssembler(time.Second), sa.ToEsp3(), sa)
	if len(msgs) != 1 || msgs[0].Kind != "smart_ack" || msgs[0].Data.(smartack.DataReclaim).MailboxIndex != 3 {
		t.Fatalf("smart ack messages = %#v", msgs)
	}

	header, err := gp.EncodeRequestHeader(gp.RequestHeader{ManufacturerID: 0x123, Bidirectional: true, Purpose: gp.PurposeTeachIn})
	if err != nil {
		t.Fatal(err)
	}
	gpPacket := erp1.Packet{Rorg: enums.RorgGP_TI, UserData: header}
	msgs = parseERP1(newReManAssembler(time.Second), gpPacket.ToEsp3(), gpPacket)
	if len(msgs) != 1 || msgs[0].Kind != "gp_header" {
		t.Fatalf("GP messages = %#v", msgs)
	}
	if got := msgs[0].Data.(gp.RequestHeader); got.ManufacturerID != 0x123 || !got.Bidirectional || got.Purpose != gp.PurposeTeachIn {
		t.Fatalf("GP header = %#v", got)
	}
}

// TestPublishDispatchesTypedChannels verifies PublishDispatchesTypedChannels behavior.
func TestPublishDispatchesTypedChannels(t *testing.T) {
	set, channels := newChannelSet(2)
	msg := Message{Kind: "response", Data: response.Packet{Code: enums.ReturnCodeSUCCESS}}
	if !publish(context.Background(), set, []Message{msg}) {
		t.Fatal("publish returned false")
	}
	if (<-channels.All).Kind != "response" {
		t.Fatal("all channel missed message")
	}
	if (<-channels.Response).Code != enums.ReturnCodeSUCCESS {
		t.Fatal("response channel missed message")
	}
}

// TestPublishFullAllStillDispatchesTypedChannel verifies PublishFullAllStillDispatchesTypedChannel behavior.
func TestPublishFullAllStillDispatchesTypedChannel(t *testing.T) {
	set, channels := newChannelSet(1)
	set.all <- Message{Kind: "occupied"}
	msg := Message{Kind: "response", Data: response.Packet{Code: enums.ReturnCodeSUCCESS}}
	if !publish(context.Background(), set, []Message{msg}) {
		t.Fatal("publish returned false")
	}
	if got := <-channels.Response; got.Code != enums.ReturnCodeSUCCESS {
		t.Fatalf("response code = %s", got.Code)
	}
}

// TestPublishStopsWhenContextCancelled verifies PublishStopsWhenContextCancelled behavior.
func TestPublishStopsWhenContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	set, _ := newChannelSet(0)
	if publish(ctx, set, []Message{{Kind: "esp3", Data: esp3.Telegram{}}}) {
		t.Fatal("publish should stop on canceled context")
	}
}

// TestReManChainPeriod verifies ReManChainPeriod behavior.
func TestReManChainPeriod(t *testing.T) {
	if remanChainPeriod != time.Second {
		t.Fatalf("ReMan chain period = %v", remanChainPeriod)
	}
	base := time.Unix(0, 0)
	key := remanKey{seq: 1}
	assembler := newReManAssembler(remanChainPeriod)
	assembler.buffers[key] = remanBuffer{updated: base}
	assembler.expire(base.Add(remanChainPeriod - time.Nanosecond))
	if _, ok := assembler.buffers[key]; !ok {
		t.Fatal("ReMan chain expired before one second")
	}
	assembler.expire(base.Add(remanChainPeriod))
	if _, ok := assembler.buffers[key]; ok {
		t.Fatal("ReMan chain did not expire at one second")
	}
}

// TestSendDropsWhenChannelFull verifies SendDropsWhenChannelFull behavior.
func TestSendDropsWhenChannelFull(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 1
	if !send(context.Background(), ch, 2) {
		t.Fatal("send should drop instead of blocking")
	}
	if got := <-ch; got != 1 {
		t.Fatalf("got %d", got)
	}
}

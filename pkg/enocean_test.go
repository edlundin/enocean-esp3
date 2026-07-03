package pkg

import (
	"context"
	"io"
	"reflect"
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
	reads [][]byte
}

func (p *fakePort) Read(b []byte) (int, error) {
	if len(p.reads) == 0 {
		time.Sleep(time.Millisecond)
		return 0, nil
	}
	n := copy(b, p.reads[0])
	p.reads = p.reads[1:]
	return n, nil
}

func (p *fakePort) Write([]byte) (int, error)                            { return 0, io.ErrClosedPipe }
func (p *fakePort) SetMode(*serial.Mode) error                           { return nil }
func (p *fakePort) Drain() error                                         { return nil }
func (p *fakePort) ResetInputBuffer() error                              { return nil }
func (p *fakePort) ResetOutputBuffer() error                             { return nil }
func (p *fakePort) SetDTR(bool) error                                    { return nil }
func (p *fakePort) SetRTS(bool) error                                    { return nil }
func (p *fakePort) GetModemStatusBits() (*serial.ModemStatusBits, error) { return nil, nil }
func (p *fakePort) SetReadTimeout(time.Duration) error                   { return nil }
func (p *fakePort) Close() error                                         { return nil }
func (p *fakePort) Break(time.Duration) error                            { return nil }

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

func TestPublishStopsWhenContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	set, _ := newChannelSet(0)
	if publish(ctx, set, []Message{{Kind: "esp3", Data: esp3.Telegram{}}}) {
		t.Fatal("publish should stop on canceled context")
	}
}

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

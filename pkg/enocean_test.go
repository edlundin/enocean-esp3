package pkg

import (
	"context"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/edlundin/enocean-esp3/pkg/deviceid"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
	"github.com/edlundin/enocean-esp3/pkg/reman"
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

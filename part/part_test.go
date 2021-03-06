package part

import (
	"github.com/cyrilix/robocar-base/testtools"
	"github.com/cyrilix/robocar-protobuf/go/events"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	"testing"
	"time"
)

type fakeLed struct {
	red, green, blue int
	blink            bool
}

func (f *fakeLed) SetBlink(freq float64) {
	if freq > 0 {
		f.blink = true
	} else {
		f.blink = false
	}
}

func (f *fakeLed) SetRed(value int) {
	f.red = value
}

func (f *fakeLed) SetGreen(value int) {
	f.green = value
}

func (f *fakeLed) SetBlue(value int) {
	f.blue = value
}

func TestLedPart_OnDriveMode(t *testing.T) {
	led := fakeLed{}
	p := LedPart{led: &led}

	cases := []struct {
		msg              mqtt.Message
		red, green, blue int
	}{
		{testtools.NewFakeMessageFromProtobuf("drive", &events.DriveModeMessage{DriveMode: events.DriveMode_USER}), 0, 255, 0},
		{testtools.NewFakeMessageFromProtobuf("drive", &events.DriveModeMessage{DriveMode: events.DriveMode_PILOT}), 0, 0, 255},
		{testtools.NewFakeMessageFromProtobuf("drive", &events.DriveModeMessage{DriveMode: events.DriveMode_INVALID}), 0, 0, 255},
	}

	for _, c := range cases {
		p.onDriveMode(nil, c.msg)
		time.Sleep(1 * time.Millisecond)
		var msg events.DriveModeMessage
		err := proto.Unmarshal(c.msg.Payload(), &msg)
		if err != nil {
			t.Errorf("unable to unmarshal drive mode message: %v", err)
		}
		value := msg.DriveMode
		if led.red != c.red {
			t.Errorf("driveMode(%v)=invalid value for red channel: %v, wants %v", value, led.red, c.red)
		}
		if led.green != c.green {
			if err != nil {
				t.Errorf("payload isn't a led value: %v", err)
			}
			t.Errorf("driveMode(%v)=invalid value for green channel: %v, wants %v", value, led.green, c.green)
		}
		if led.blue != c.blue {
			if err != nil {
				t.Errorf("payload isn't a led value: %v", err)
			}
			t.Errorf("driveMode(%v)=invalid value for blue channel: %v, wants %v", value, led.blue, c.blue)
		}
	}
}
func TestLedPart_OnRecord(t *testing.T) {
	led := fakeLed{}
	p := LedPart{led: &led}

	cases := []struct {
		msg    mqtt.Message
		record bool
		blink  bool
	}{
		{testtools.NewFakeMessageFromProtobuf("record", &events.SwitchRecordMessage{Enabled: false}), true, false},
		{testtools.NewFakeMessageFromProtobuf("record", &events.SwitchRecordMessage{Enabled: true}), false, true},
		{testtools.NewFakeMessageFromProtobuf("record", &events.SwitchRecordMessage{Enabled: false}), true, false},
		{testtools.NewFakeMessageFromProtobuf("record", &events.SwitchRecordMessage{Enabled: true}), false, true},
	}

	for _, c := range cases {
		p.onRecord(nil, c.msg)
		if led.blink != c.blink {
			var msg events.SwitchRecordMessage
			err := proto.Unmarshal(c.msg.Payload(), &msg)
			if err != nil {
				t.Errorf("unable to unmarshal %T message: %v", msg, err)
			}

			value := msg.Enabled
			t.Errorf("onRecord(%v): %v, wants %v", value, c.record, led.blink)
		}
	}
}

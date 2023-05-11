package part

import (
	"github.com/cyrilix/robocar-base/testtools"
	"github.com/cyrilix/robocar-led/pkg/led"
	"github.com/cyrilix/robocar-protobuf/go/events"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

type fakeLed struct {
	color led.Color
	blink bool
}

func (f *fakeLed) SetColor(color led.Color) {
	f.color = color
}

func (f *fakeLed) SetBlink(freq float64) {
	if freq > 0 {
		f.blink = true
	} else {
		f.blink = false
	}
}

func TestLedPart_OnDriveMode(t *testing.T) {
	l := fakeLed{}
	p := LedPart{led: &l, speedZone: events.SpeedZone_FAST}

	cases := []struct {
		msg   mqtt.Message
		color led.Color
	}{
		{testtools.NewFakeMessageFromProtobuf("drive", &events.DriveModeMessage{DriveMode: events.DriveMode_USER}), led.ColorGreen},
		{testtools.NewFakeMessageFromProtobuf("drive", &events.DriveModeMessage{DriveMode: events.DriveMode_PILOT}), led.ColorBlue},
		{testtools.NewFakeMessageFromProtobuf("drive", &events.DriveModeMessage{DriveMode: events.DriveMode_INVALID}), led.ColorBlue},
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
		if l.color != c.color {
			t.Errorf("driveMode(%v)=invalid value for expectedColor: %v, wants %v", value, l.color, c.color)
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

func TestLedPart_OnSpeedZone(t *testing.T) {
	l := fakeLed{}
	p := LedPart{led: &l, mode: LedModeSpeedZone, driveMode: events.DriveMode_PILOT}

	cases := []struct {
		msg   mqtt.Message
		color led.Color
	}{
		{testtools.NewFakeMessageFromProtobuf("speedzone", &events.SpeedZoneMessage{SpeedZone: events.SpeedZone_SLOW}), led.ColorRed},
		{testtools.NewFakeMessageFromProtobuf("speedzone", &events.SpeedZoneMessage{SpeedZone: events.SpeedZone_NORMAL}), led.ColorYellow},
		{testtools.NewFakeMessageFromProtobuf("speedzone", &events.SpeedZoneMessage{SpeedZone: events.SpeedZone_FAST}), led.ColorBlue},
		{testtools.NewFakeMessageFromProtobuf("speedzone", &events.SpeedZoneMessage{SpeedZone: events.SpeedZone_UNKNOWN}), led.ColorWhite},
	}

	for _, c := range cases {
		p.onSpeedZone(nil, c.msg)
		time.Sleep(1 * time.Millisecond)
		var msg events.SpeedZoneMessage
		err := proto.Unmarshal(c.msg.Payload(), &msg)
		if err != nil {
			t.Errorf("unable to unmarshal drive mode message: %v", err)
		}
		value := msg.GetSpeedZone()
		if l.color != c.color {
			t.Errorf("driveMode(%v)=invalid value for expectedColor: %v, wants %v", value, l.color, c.color)
		}
	}
}

func TestLedPart_OnThrottle(t *testing.T) {

	cases := []struct {
		name          string
		msg           mqtt.Message
		expectedColor led.Color
	}{
		{"throttle stop",
			testtools.NewFakeMessageFromProtobuf("throttle", &events.ThrottleMessage{Throttle: 0.}),
			led.ColorBlue,
		},
		{
			"throttle normal",
			testtools.NewFakeMessageFromProtobuf("throttle", &events.ThrottleMessage{Throttle: 0.5}),
			led.ColorBlue,
		},
		{
			"near zero",
			testtools.NewFakeMessageFromProtobuf("throttle", &events.ThrottleMessage{Throttle: -0.01}),
			led.ColorBlue,
		},
		{
			"slow brake",
			testtools.NewFakeMessageFromProtobuf("throttle", &events.ThrottleMessage{Throttle: -0.06}),
			led.ColorWhite,
		},
		{
			"normal brake",
			testtools.NewFakeMessageFromProtobuf("throttle", &events.ThrottleMessage{Throttle: -0.4}),
			led.ColorYellow,
		},
		{
			"normal high brake",
			testtools.NewFakeMessageFromProtobuf("throttle", &events.ThrottleMessage{Throttle: -0.6}),
			led.ColorRed,
		},
		{
			"high brake",
			testtools.NewFakeMessageFromProtobuf("throttle", &events.ThrottleMessage{Throttle: -1.}),
			led.ColorPurple,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := fakeLed{}
			p := LedPart{led: &l, mode: LedModeBrake, driveMode: events.DriveMode_PILOT}

			p.onThrottle(nil, c.msg)
			time.Sleep(1 * time.Millisecond)
			var msg events.ThrottleMessage
			err := proto.Unmarshal(c.msg.Payload(), &msg)
			if err != nil {
				t.Errorf("unable to unmarshal drive mode message: %v", err)
			}
			value := msg.GetThrottle()
			if l.color != c.expectedColor {
				t.Errorf("driveMode(%v)=invalid value for expectedColor: %v, wants %v", value, l.color, c.expectedColor)
			}
		})
	}
}

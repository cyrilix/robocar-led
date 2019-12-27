package part

import (
	"github.com/cyrilix/robocar-base/mqttdevice"
	"github.com/cyrilix/robocar-base/testtools"
	"github.com/cyrilix/robocar-base/types"
	mqtt "github.com/eclipse/paho.mqtt.golang"
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
		{testtools.NewFakeMessage("drive", mqttdevice.NewMqttValue(types.DriveModeUser)), 0, 255, 0},
		{testtools.NewFakeMessage("drive", mqttdevice.NewMqttValue(types.DriveModePilot)), 0, 0, 255},
		{testtools.NewFakeMessage("drive", mqttdevice.NewMqttValue(types.DriveModeInvalid)), 0, 0, 255},
	}

	for _, c := range cases {
		p.onDriveMode(nil, c.msg)
		time.Sleep(1 * time.Millisecond)
		if led.red != c.red {
			payload := mqttdevice.NewMqttValue(c.msg.Payload())
			value, err := payload.IntValue()
			if err != nil {
				t.Errorf("payload isn't a led value: %v", err)
			}
			t.Errorf("driveMode(%v)=invalid value for red channel: %v, wants %v", value, led.red, c.red)
		}
		if led.green != c.green {
			payload := mqttdevice.NewMqttValue(c.msg.Payload())
			value, err := payload.IntValue()
			if err != nil {
				t.Errorf("payload isn't a led value: %v", err)
			}
			t.Errorf("driveMode(%v)=invalid value for green channel: %v, wants %v", value, led.green, c.green)
		}
		if led.blue != c.blue {
			payload := mqttdevice.NewMqttValue(c.msg.Payload())
			value, err := payload.IntValue()
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
		{testtools.NewFakeMessage("record", mqttdevice.NewMqttValue(false)), true, false},
		{testtools.NewFakeMessage("record", mqttdevice.NewMqttValue(true)), false, true},
		{testtools.NewFakeMessage("record", mqttdevice.NewMqttValue(false)), true, false},
		{testtools.NewFakeMessage("record", mqttdevice.NewMqttValue(true)), false, true},
	}

	for _, c := range cases {
		p.onRecord(nil, c.msg)
		if led.blink != c.blink {
			payload := mqttdevice.NewMqttValue(c.msg.Payload())
			value, err := payload.BoolValue()
			if err != nil {
				t.Errorf("payload isn't a bool value: %v", err)
			}
			t.Errorf("onRecord(%v): %v, wants %v", value, c.record, led.blink)
		}
	}
}

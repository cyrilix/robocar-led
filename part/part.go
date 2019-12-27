package part

import (
	"fmt"
	"github.com/cyrilix/robocar-base/mqttdevice"
	"github.com/cyrilix/robocar-base/service"
	"github.com/cyrilix/robocar-base/types"
	"github.com/cyrilix/robocar-led/led"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"sync"
	"time"
)

func NewPart(client mqtt.Client, driveModeTopic, recordTopic string) *LedPart {
	return &LedPart{
		led:              led.New(),
		client:           client,
		onDriveModeTopic: driveModeTopic,
		onRecordTopic:    recordTopic,
		muDriveMode:      sync.Mutex{},
		m:                types.DriveModeInvalid,
		muRecord:         sync.Mutex{},
		recordEnabled:    false,
	}

}

type LedPart struct {
	led              led.ColoredLed
	client           mqtt.Client
	onDriveModeTopic string
	onRecordTopic    string

	muDriveMode   sync.Mutex
	m             types.DriveMode
	muRecord      sync.Mutex
	recordEnabled bool
}

func (p *LedPart) Start() error {
	if err := p.registerCallbacks(); err != nil {
		return fmt.Errorf("unable to start service: %v", err)
	}
	for {
		time.Sleep(1 * time.Hour)
	}
}

func (p *LedPart) Stop() {
	defer p.led.SetBlink(0)
	defer p.led.SetGreen(0)
	defer p.led.SetBlue(0)
	defer p.led.SetRed(0)
	service.StopService("led", p.client, p.onDriveModeTopic, p.onRecordTopic)
}

func (p *LedPart) onDriveMode(_ mqtt.Client, message mqtt.Message) {
	payload := message.Payload()
	value := mqttdevice.NewMqttValue(payload)
	m, err := value.DriveModeValue()
	if err != nil {
		log.Printf("invalid drive mode: %v", err)
		return
	}
	switch m {
	case types.DriveModeUser:
		p.led.SetRed(0)
		p.led.SetGreen(255)
		p.led.SetBlue(0)
	case types.DriveModePilot:
		p.led.SetRed(0)
		p.led.SetGreen(0)
		p.led.SetBlue(255)
	}
}

func (p *LedPart) onRecord(client mqtt.Client, message mqtt.Message) {
	mqttValue := mqttdevice.NewMqttValue(message.Payload())
	rec, err := mqttValue.BoolValue()
	if err != nil {
		log.Printf("unable to convert message payload '%v' to bool: %v", message.Payload(), err)
		return
	}
	if rec {
		p.led.SetBlink(2)
	} else {
		p.led.SetBlink(0)
	}
}

func (p *LedPart) registerCallbacks() error {
	err := service.RegisterCallback(p.client, p.onDriveModeTopic, p.onDriveMode)
	if err != nil {
		return err
	}

	err = service.RegisterCallback(p.client, p.onRecordTopic, p.onRecord)
	if err != nil {
		return err
	}
	return nil
}


package part

import (
	"fmt"
	"github.com/cyrilix/robocar-base/mode"
	"github.com/cyrilix/robocar-base/mqttdevice"
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
		m:                mode.DriveModeInvalid,
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
	m             mode.DriveMode
	muRecord      sync.Mutex
	recordEnabled bool
}

func (p *LedPart) Start() error {
	if err := p.registerCallbacks(); err != nil {
		return fmt.Errorf("unable to start service: %v", err)
	}
	defer p.Stop()
	for {
		time.Sleep(1 * time.Hour)
	}
}

func (p *LedPart) Stop() {
	StopService("led", p.client, p.onDriveModeTopic, p.onRecordTopic)
}

func (p *LedPart) onDriveMode(_ mqtt.Client, message mqtt.Message) {
	mqttValue := mqttdevice.NewMqttValue(message.Payload())
	m, err := mqttValue.IntValue()
	if err != nil {
		log.Printf("unable to convert message payload '%v' to DriveMode: %v", message.Payload(), err)
		return
	}
	switch m {
	case mode.DriveModeUser:
		p.led.SetRed(0)
		p.led.SetGreen(255)
		p.led.SetBlue(0)
	case mode.DriveModePilot:
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
	err := RegisterCallback(p.client, p.onDriveModeTopic, p.onDriveMode)
	if err != nil {
		return err
	}

	err = RegisterCallback(p.client, p.onRecordTopic, p.onRecord)
	if err != nil {
		return err
	}
	return nil
}

func StopService(name string, client mqtt.Client, topics ...string) {
	log.Printf("Stop %s service", name)
	token := client.Unsubscribe(topics...)
	token.Wait()
	if token.Error() != nil {
		log.Printf("unable to unsubscribe service: %v", token.Error())
	}
	client.Disconnect(50)
}

func RegisterCallback(client mqtt.Client, topic string, callback mqtt.MessageHandler) error {
	log.Printf("Register callback on topic %v", topic)
	token := client.Subscribe(topic, 0, callback)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("unable to register callback on topic %s: %v", topic, token.Error())
	}
	return nil
}

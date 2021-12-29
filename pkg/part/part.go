package part

import (
	"fmt"
	"github.com/cyrilix/robocar-base/service"
	"github.com/cyrilix/robocar-led/pkg/led"
	"github.com/cyrilix/robocar-protobuf/go/events"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
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
		m:                events.DriveMode_INVALID,
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
	m             events.DriveMode
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
	defer p.led.SetColor(led.ColorBlack)
	service.StopService("led", p.client, p.onDriveModeTopic, p.onRecordTopic)
}

func (p *LedPart) onDriveMode(_ mqtt.Client, message mqtt.Message) {
	var driveModeMessage events.DriveModeMessage
	err := proto.Unmarshal(message.Payload(), &driveModeMessage)
	if err != nil {
		zap.S().Errorf("unable to unmarshal %T message: %v", driveModeMessage, err)
		return
	}
	switch driveModeMessage.GetDriveMode() {
	case events.DriveMode_USER:
		p.led.SetColor(led.ColorGreen)
	case events.DriveMode_PILOT:
		p.led.SetColor(led.ColorBlue)
	}
}

func (p *LedPart) onRecord(client mqtt.Client, message mqtt.Message) {
	var switchRecord events.SwitchRecordMessage
	err := proto.Unmarshal(message.Payload(), &switchRecord)
	if err != nil {
		zap.S().Errorf("unable to unmarchal %T message: %v", switchRecord, err)
		return
	}

	p.muRecord.Lock()
	defer p.muRecord.Unlock()
	if p.recordEnabled == switchRecord.GetEnabled() {
		return
	}
	p.recordEnabled = switchRecord.GetEnabled()

	if switchRecord.GetEnabled() {
		zap.S().Info("record mode enabled")
		p.led.SetBlink(2)
	} else {
		zap.S().Info("record mode disabled")
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

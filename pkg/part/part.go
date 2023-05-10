package part

import (
	"fmt"
	"github.com/cyrilix/robocar-base/service"
	"github.com/cyrilix/robocar-led/pkg/led"
	"github.com/cyrilix/robocar-protobuf/go/events"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

const (
	LedModeBrake LedMode = iota
	LedModeSpeedZone
)

type LedMode int

func NewPart(client mqtt.Client, driveModeTopic, recordTopic, speedZoneTopic, throttleTopic string, ledMode LedMode) *LedPart {
	return &LedPart{
		led:              led.New(),
		mode:             ledMode,
		client:           client,
		onDriveModeTopic: driveModeTopic,
		onRecordTopic:    recordTopic,
		onSpeedZoneTopic: speedZoneTopic,
		onThrottleTopic:  throttleTopic,
		muDriveMode:      sync.Mutex{},
		driveMode:        events.DriveMode_INVALID,
		muRecord:         sync.Mutex{},
		recordEnabled:    false,
		muSpeedZone:      sync.Mutex{},
		speedZone:        events.SpeedZone_UNKNOWN,
		muThrottle:       sync.Mutex{},
	}

}

type LedPart struct {
	led              led.ColoredLed
	mode             LedMode
	client           mqtt.Client
	onDriveModeTopic string
	onRecordTopic    string
	onSpeedZoneTopic string
	onThrottleTopic  string

	muDriveMode   sync.Mutex
	driveMode     events.DriveMode
	muRecord      sync.Mutex
	recordEnabled bool

	muSpeedZone sync.Mutex
	speedZone   events.SpeedZone

	muThrottle sync.Mutex
	throttle   float32
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
	service.StopService("led", p.client, p.onDriveModeTopic, p.onRecordTopic, p.onSpeedZoneTopic, p.onThrottleTopic)
}

func (p *LedPart) setDriveMode(m events.DriveMode) {
	p.muDriveMode.Lock()
	defer p.muDriveMode.Unlock()
	p.driveMode = m
}

func (p *LedPart) onDriveMode(_ mqtt.Client, message mqtt.Message) {
	var driveModeMessage events.DriveModeMessage
	err := proto.Unmarshal(message.Payload(), &driveModeMessage)
	if err != nil {
		zap.S().Errorf("unable to unmarshal %T message: %v", driveModeMessage, err)
		return
	}
	p.setDriveMode(driveModeMessage.GetDriveMode())
	p.updateColor()
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

func (p *LedPart) setSpeedZone(sz events.SpeedZone) {
	p.muSpeedZone.Lock()
	defer p.muSpeedZone.Unlock()
	p.speedZone = sz
}

func (p *LedPart) onSpeedZone(_ mqtt.Client, message mqtt.Message) {
	var speedZoneMessage events.SpeedZoneMessage
	err := proto.Unmarshal(message.Payload(), &speedZoneMessage)
	if err != nil {
		zap.S().Errorf("unable to unmarshal %T message: %v", speedZoneMessage, err)
		return
	}

	p.setSpeedZone(speedZoneMessage.GetSpeedZone())
	p.updateColor()
}

func (p *LedPart) setThrottle(throttle float32) {
	p.muThrottle.Lock()
	defer p.muThrottle.Unlock()
	p.throttle = throttle
}

func (p *LedPart) onThrottle(_ mqtt.Client, message mqtt.Message) {
	var throttleMessage events.ThrottleMessage
	err := proto.Unmarshal(message.Payload(), &throttleMessage)
	if err != nil {
		zap.S().Errorf("unable to unmarshal %T message: %v", throttleMessage, err)
		return
	}

	p.setThrottle(throttleMessage.GetThrottle())
	p.updateColor()
}

func (p *LedPart) updateColor() {
	p.muSpeedZone.Lock()
	defer p.muSpeedZone.Unlock()
	p.muDriveMode.Lock()
	defer p.muDriveMode.Unlock()
	p.muThrottle.Lock()
	defer p.muThrottle.Unlock()

	if p.throttle <= -0.05 {
		p.led.SetColor(led.Color{Red: int(p.throttle * -255)})
		return
	}

	switch p.mode {
	case LedModeBrake:
		p.updateBrakeColor()
	case LedModeSpeedZone:
		p.updateSpeedZoneColor()
	}
}

func (p *LedPart) updateSpeedZoneColor() {
	switch p.driveMode {
	case events.DriveMode_USER:
		p.led.SetColor(led.ColorGreen)
	case events.DriveMode_PILOT:
		switch p.speedZone {
		case events.SpeedZone_UNKNOWN:
			p.led.SetColor(led.ColorWhite)
		case events.SpeedZone_SLOW:
			p.led.SetColor(led.ColorRed)
		case events.SpeedZone_NORMAL:
			p.led.SetColor(led.ColorYellow)
		case events.SpeedZone_FAST:
			p.led.SetColor(led.ColorBlue)
		}
	}
}

func (p *LedPart) updateBrakeColor() {

	switch p.driveMode {
	case events.DriveMode_USER:
		p.led.SetColor(led.ColorGreen)
	case events.DriveMode_PILOT:
		p.led.SetColor(led.ColorBlue)
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

	err = service.RegisterCallback(p.client, p.onSpeedZoneTopic, p.onSpeedZone)
	if err != nil {
		return err
	}

	err = service.RegisterCallback(p.client, p.onThrottleTopic, p.onThrottle)
	if err != nil {
		return err
	}

	return nil
}

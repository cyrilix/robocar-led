package led

import (
	log "github.com/sirupsen/logrus"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
	"sync"
	"time"
)

func init() {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

}

func New() *PiColorLed {
	led := PiColorLed{
		pinRed:          rpi.P1_16,
		pinGreen:        rpi.P1_18,
		pinBlue:         rpi.P1_22,
		redValue:        0,
		greenValue:      0,
		blueValue:       0,
		cancelBlinkChan: make(chan interface{}),
		blinkEnabled:    false,
	}
	return &led
}

type Led interface {
	SetBlink(freq float64)
}

type ColoredLed interface {
	Led
	SetRed(value int)
	SetGreen(value int)
	SetBlue(value int)
}

type PiColorLed struct {
	muPinRed, muPinGreen, muPinBlue sync.Mutex
	pinRed                          gpio.PinIO
	pinGreen                        gpio.PinIO
	pinBlue                         gpio.PinIO

	muColorValue                    sync.RWMutex
	redValue, greenValue, blueValue int

	cancelBlinkChan chan interface{}

	muBlink      sync.Mutex
	blinkEnabled bool
}

func (l *PiColorLed) SetBlink(freq float64) {
	l.muBlink.Lock()
	defer l.muBlink.Unlock()
	if freq > 0 {
		if !l.blinkEnabled {
			l.blinkEnabled = true
			go l.blink(freq)
		}
	} else {
		if l.blinkEnabled {
			l.blinkEnabled = false
			l.cancelBlinkChan <- struct{}{}
		}
	}
}

func (l *PiColorLed) blink(freq float64) {
	factor := 0
	ticker := time.NewTicker(time.Duration(float64(time.Second) / freq))
	red := l.Red()
	green := l.Green()
	blue := l.Blue()
	var tmpR, tmpG, tmpB int
	for {
		select {
		case <-ticker.C:
		case <-l.cancelBlinkChan:
			// Restore values
			l.SetRed(red)
			l.SetGreen(green)
			l.SetBlue(blue)
			return
		}

		tmpR = l.Red()
		tmpG = l.Green()
		tmpB = l.Blue()
		if factor == 1 {
			// Led is off
			if tmpR > 0 {
				red = tmpR
			}
			if tmpG > 0 {
				green = tmpG
			}
			if tmpB > 0 {
				blue = tmpB
			}
		} else {
			// Led on: get updated value
			red = tmpR
			green = tmpG
			blue = tmpB
		}
		log.Infof("factor: %v", factor)
		l.SetRed(red * factor)
		l.SetGreen(green * factor)
		l.SetBlue(blue * factor)

		if factor == 0 {
			factor = 1
		} else {
			factor = 0
		}
	}
}

var setLed = func(v int, led gpio.PinIO, mutex *sync.Mutex) {
	mutex.Lock()
	defer mutex.Unlock()
	lvl := gpio.High
	if v == 0 {
		lvl = gpio.Low
	}
	err := led.Out(lvl)
	if err != nil {
		log.Errorf("unable to sed pin to %v: %v", lvl, err)
	}
}

func (l *PiColorLed) Red() int {
	l.muColorValue.RLock()
	defer l.muColorValue.RUnlock()
	return l.redValue
}
func (l *PiColorLed) SetRed(v int) {
	setLed(v, l.pinRed, &l.muPinRed)
	l.muColorValue.Lock()
	defer l.muColorValue.Unlock()
	l.redValue = v
}

func (l *PiColorLed) Green() int {
	l.muColorValue.RLock()
	defer l.muColorValue.RUnlock()
	return l.greenValue
}
func (l *PiColorLed) SetGreen(v int) {
	setLed(v, l.pinGreen, &l.muPinGreen)
	l.muColorValue.Lock()
	defer l.muColorValue.Unlock()
	l.greenValue = v
}

func (l *PiColorLed) Blue() int {
	l.muColorValue.RLock()
	defer l.muColorValue.RUnlock()
	return l.blueValue
}
func (l *PiColorLed) SetBlue(v int) {
	setLed(v, l.pinBlue, &l.muPinBlue)
	l.muColorValue.Lock()
	defer l.muColorValue.Unlock()
	l.blueValue = v
}

package led

import (
	"go.uber.org/zap"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
	"sync"
	"time"
)

func init() {
	zap.S().Info("init pin")
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		zap.S().Fatalf("unable to init host driver: %v", err)
	}

}

var (
	ColorBlack = Color{0, 0, 0}
	ColorRed   = Color{255, 0, 0}
	ColorGreen = Color{0, 255, 0}
	ColorBlue  = Color{0, 0, 255}
)

func New() *PiColorLed {
	led := PiColorLed{
		pinRed:          rpi.P1_16,
		pinGreen:        rpi.P1_18,
		pinBlue:         rpi.P1_22,
		currentColor:    ColorBlack,
		cancelBlinkChan: make(chan interface{}),
		blinkEnabled:    false,
	}

	return &led
}

type Color struct {
	Red   int
	Green int
	Blue  int
}

type Led interface {
	SetBlink(freq float64)
}

type ColoredLed interface {
	Led
	SetColor(color Color)
}

type PiColorLed struct {
	muPinRed, muPinGreen, muPinBlue sync.Mutex
	pinRed                          gpio.PinIO
	pinGreen                        gpio.PinIO
	pinBlue                         gpio.PinIO

	muColorValue sync.RWMutex
	currentColor Color

	cancelBlinkChan chan interface{}

	muBlink      sync.Mutex
	blinkEnabled bool
}

func (l *PiColorLed) SetColor(color Color) {
	l.muColorValue.Lock()
	defer l.muColorValue.Unlock()
	if color == l.currentColor {
		return
	}
	l.currentColor = color
	setLed(color.Red, l.pinRed, &l.muPinRed)
	setLed(color.Green, l.pinGreen, &l.muPinGreen)
	setLed(color.Blue, l.pinBlue, &l.muPinBlue)
}

func (l *PiColorLed) on() {
	l.muColorValue.RLock()
	defer l.muColorValue.RUnlock()

	setLed(l.currentColor.Red, l.pinRed, &l.muPinRed)
	setLed(l.currentColor.Green, l.pinGreen, &l.muPinGreen)
	setLed(l.currentColor.Blue, l.pinBlue, &l.muPinBlue)
}
func (l *PiColorLed) off() {
	l.muColorValue.RLock()
	defer l.muColorValue.RUnlock()

	setLed(0, l.pinRed, &l.muPinRed)
	setLed(0, l.pinGreen, &l.muPinGreen)
	setLed(0, l.pinBlue, &l.muPinBlue)
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
	log := zap.S().With("func", "blink")
	ticker := time.NewTicker(time.Duration(float64(time.Second) / freq))

	// Restore values
	defer l.SetColor(l.Color())

	for {
		select {
		case <-ticker.C:
		case <-l.cancelBlinkChan:
			return
		}
		log.Debugf("off with color %v", ColorBlack)
		l.off()

		select {
		case <-ticker.C:
		case <-l.cancelBlinkChan:
			return
		}
		log.Debugf("on with color %v", l.Color())
		l.on()
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
		zap.S().Errorf("unable to sed pin to %v: %v", lvl, err)
	}
}

func (l *PiColorLed) Red() int {
	l.muColorValue.RLock()
	defer l.muColorValue.RUnlock()
	return l.currentColor.Red
}

func (l *PiColorLed) Green() int {
	l.muColorValue.RLock()
	defer l.muColorValue.RUnlock()
	return l.currentColor.Green
}

func (l *PiColorLed) Blue() int {
	l.muColorValue.RLock()
	defer l.muColorValue.RUnlock()
	return l.currentColor.Blue
}

func (l *PiColorLed) Color() Color {
	l.muColorValue.RLock()
	defer l.muColorValue.RUnlock()
	return l.currentColor
}

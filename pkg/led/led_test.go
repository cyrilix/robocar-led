package led

import (
	"periph.io/x/conn/v3/gpio"
	"sync"
	"testing"
	"time"
)

func TestColorLed_Red(t *testing.T) {
	setLedBackup := setLed
	defer func() { setLed = setLedBackup }()

	ledColors := make(map[gpio.PinIO]int)
	setLed = func(v int, led gpio.PinIO, mutex *sync.Mutex) {
		mutex.Lock()
		defer mutex.Unlock()
		ledColors[led] = v
	}

	l := New()
	if l.Red() != 0 {
		t.Errorf("%T.Red(): %v, wants %v", l, l.Red(), 0)
	}
	if ledColors[l.pinRed] != 0 {
		t.Errorf("colorValue: %v, wants %v", ledColors[l.pinRed], 0)
	}

	l.SetColor(ColorRed)
	if l.Red() != 255 {
		t.Errorf("%T.Red(): %v, wants %v", l, l.Red(), 255)
	}
	if ledColors[l.pinRed] != 255 {
		t.Errorf("colorValue: %v, wants %v", ledColors[l.pinRed], 255)
	}
}

func TestColorLed_Green(t *testing.T) {
	setLedBackup := setLed
	defer func() { setLed = setLedBackup }()

	ledColors := make(map[gpio.PinIO]int)
	setLed = func(v int, led gpio.PinIO, mutex *sync.Mutex) {
		mutex.Lock()
		defer mutex.Unlock()
		ledColors[led] = v
	}

	l := New()
	if l.Green() != 0 {
		t.Errorf("%T.Green(): %v, wants %v", l, l.Green(), 0)
	}
	if ledColors[l.pinGreen] != 0 {
		t.Errorf("colorValue: %v, wants %v", ledColors[l.pinGreen], 0)
	}

	l.SetColor(ColorGreen)
	if l.Green() != 255 {
		t.Errorf("%T.Green(): %v, wants %v", l, l.Green(), 255)
	}
	if ledColors[l.pinGreen] != 255 {
		t.Errorf("colorValue: %v, wants %v", ledColors[l.pinGreen], 255)
	}
}

func TestColorLed_Blue(t *testing.T) {
	setLedBackup := setLed
	defer func() { setLed = setLedBackup }()

	ledColors := make(map[gpio.PinIO]int)
	setLed = func(v int, led gpio.PinIO, mutex *sync.Mutex) {
		mutex.Lock()
		defer mutex.Unlock()
		ledColors[led] = v
	}

	l := New()
	if l.Blue() != 0 {
		t.Errorf("%T.Blue(): %v, wants %v", l, l.Blue(), 0)
	}
	if ledColors[l.pinBlue] != 0 {
		t.Errorf("colorValue: %v, wants %v", ledColors[l.pinBlue], 0)
	}

	l.SetColor(ColorBlue)
	if l.Blue() != 255 {
		t.Errorf("%T.Blue(): %v, wants %v", l, l.Blue(), 255)
	}
	if ledColors[l.pinBlue] != 255 {
		t.Errorf("colorValue: %v, wants %v", ledColors[l.pinBlue], 255)
	}
}

func TestColorLed_SetBlink(t *testing.T) {
	setLedBackup := setLed
	defer func() { setLed = setLedBackup }()

	var muFakeValue sync.Mutex
	ledColors := make(map[gpio.PinIO]int)
	setLed = func(v int, led gpio.PinIO, mutex *sync.Mutex) {
		mutex.Lock()
		defer mutex.Unlock()
		muFakeValue.Lock()
		defer muFakeValue.Unlock()
		ledColors[led] = v
	}
	readValue := func(p gpio.PinIO) int {
		muFakeValue.Lock()
		defer muFakeValue.Unlock()
		return ledColors[p]
	}

	l := New()
	l.SetColor(ColorBlue)
	v := ledColors[l.pinBlue]
	if v != 255 {
		t.Errorf("colorValue: %v, wants %v", v, 255)
	}
	l.SetBlink(100)
	v = readValue(l.pinBlue)
	if v != 255 {
		t.Errorf("colorValue: %v, wants %v", v, 255)
	}
	time.Sleep(12 * time.Millisecond)
	v = readValue(l.pinBlue)
	if v != 0 {
		t.Errorf("colorValue: %v, wants %v", v, 0)
	}
	time.Sleep(12 * time.Millisecond)
	v = readValue(l.pinBlue)
	if v != 255 {
		t.Errorf("colorValue: %v, wants %v", v, 255)
	}
	time.Sleep(12 * time.Millisecond)
	v = readValue(l.pinBlue)
	if v != 0 {
		t.Errorf("colorValue: %v, wants %v", v, 0)
	}

	// Stop blink
	l.SetBlink(0)
	time.Sleep(5 * time.Millisecond)
	v = readValue(l.pinBlue)
	if v != 255 {
		t.Errorf("colorValue: %v, wants %v", v, 255)
	}
	time.Sleep(12 * time.Millisecond)
	v = readValue(l.pinBlue)
	if v != 255 {
		t.Errorf("colorValue: %v, wants %v", v, 255)
	}
}

func TestColorLed_SetBlinkAndUpdadeColor(t *testing.T) {
	setLedBackup := setLed
	defer func() { setLed = setLedBackup }()

	var muFakeValue sync.Mutex
	ledColors := make(map[gpio.PinIO]int)
	setLed = func(v int, led gpio.PinIO, mutex *sync.Mutex) {
		mutex.Lock()
		defer mutex.Unlock()
		muFakeValue.Lock()
		defer muFakeValue.Unlock()
		ledColors[led] = v
	}
	readValue := func(p gpio.PinIO) int {
		muFakeValue.Lock()
		defer muFakeValue.Unlock()
		return ledColors[p]
	}

	l := New()
	l.SetColor(ColorBlue)
	l.SetBlink(100)
	v := readValue(l.pinBlue)
	if v != 255 {
		t.Errorf("colorValue: %v, wants %v", v, 255)
	}
	time.Sleep(6 * time.Millisecond)
	l.SetColor(ColorBlue)

	v = readValue(l.pinBlue)
	if v != 255 {
		t.Errorf("colorValue: %v, wants %v", v, 128)
	}
	time.Sleep(6 * time.Millisecond)

	time.Sleep(12 * time.Millisecond)

	v = readValue(l.pinBlue)
	if v != 255 {
		t.Errorf("colorValue: %v, wants %v", v, 128)
	}
	time.Sleep(12 * time.Millisecond)

	v = readValue(l.pinBlue)
	if v != 0 {
		t.Errorf("colorValue: %v, wants %v", v, 0)
	}

	// Stop blink
	l.SetBlink(0)
	time.Sleep(5 * time.Millisecond)
	v = readValue(l.pinBlue)
	if v != 255 {
		t.Errorf("colorValue: %v, wants %v", v, 128)
	}
	time.Sleep(12 * time.Millisecond)
	v = readValue(l.pinBlue)
	if v != 255 {
		t.Errorf("colorValue: %v, wants %v", v, 128)
	}
}

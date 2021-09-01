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

	var colorValue int
	setLed = func(v int, led gpio.PinIO, mutex *sync.Mutex) {
		mutex.Lock()
		defer mutex.Unlock()
		colorValue = v
	}

	l := New()
	if l.Red() != 0 {
		t.Errorf("%T.Red(): %v, wants %v", l, l.Red(), 0)
	}
	if colorValue != 0 {
		t.Errorf("colorValue: %v, wants %v", colorValue, 0)
	}

	l.SetRed(255)
	if l.Red() != 255 {
		t.Errorf("%T.Red(): %v, wants %v", l, l.Red(), 255)
	}
	if colorValue != 255 {
		t.Errorf("colorValue: %v, wants %v", colorValue, 255)
	}
}

func TestColorLed_Green(t *testing.T) {
	setLedBackup := setLed
	defer func() { setLed = setLedBackup }()

	var colorValue int
	setLed = func(v int, led gpio.PinIO, mutex *sync.Mutex) {
		mutex.Lock()
		defer mutex.Unlock()
		colorValue = v
	}

	l := New()
	if l.Green() != 0 {
		t.Errorf("%T.Green(): %v, wants %v", l, l.Green(), 0)
	}
	if colorValue != 0 {
		t.Errorf("colorValue: %v, wants %v", colorValue, 0)
	}

	l.SetGreen(255)
	if l.Green() != 255 {
		t.Errorf("%T.Green(): %v, wants %v", l, l.Green(), 255)
	}
	if colorValue != 255 {
		t.Errorf("colorValue: %v, wants %v", colorValue, 255)
	}
}

func TestColorLed_Blue(t *testing.T) {
	setLedBackup := setLed
	defer func() { setLed = setLedBackup }()

	var colorValue int
	setLed = func(v int, led gpio.PinIO, mutex *sync.Mutex) {
		mutex.Lock()
		defer mutex.Unlock()
		colorValue = v
	}

	l := New()
	if l.Blue() != 0 {
		t.Errorf("%T.Blue(): %v, wants %v", l, l.Blue(), 0)
	}
	if colorValue != 0 {
		t.Errorf("colorValue: %v, wants %v", colorValue, 0)
	}

	l.SetBlue(255)
	if l.Blue() != 255 {
		t.Errorf("%T.Blue(): %v, wants %v", l, l.Blue(), 255)
	}
	if colorValue != 255 {
		t.Errorf("colorValue: %v, wants %v", colorValue, 255)
	}
}

func TestColorLed_SetBlink(t *testing.T) {
	setLedBackup := setLed
	defer func() { setLed = setLedBackup }()

	var muFakeValue sync.Mutex
	var colorValue int
	setLed = func(v int, led gpio.PinIO, mutex *sync.Mutex) {
		mutex.Lock()
		defer mutex.Unlock()
		muFakeValue.Lock()
		defer muFakeValue.Unlock()
		colorValue = v
	}
	readValue := func() int {
		muFakeValue.Lock()
		defer muFakeValue.Unlock()
		return colorValue
	}

	l := New()
	l.SetBlue(255)
	v := readValue()
	if v != 255 {
		t.Errorf("colorValue: %v, wants %v", v, 255)
	}
	l.SetBlink(100)
	v = readValue()
	if v != 255 {
		t.Errorf("colorValue: %v, wants %v", v, 255)
	}
	time.Sleep(12 * time.Millisecond)
	v = readValue()
	if v != 0 {
		t.Errorf("colorValue: %v, wants %v", v, 0)
	}
	time.Sleep(12 * time.Millisecond)
	v = readValue()
	if v != 255 {
		t.Errorf("colorValue: %v, wants %v", v, 255)
	}
	time.Sleep(12 * time.Millisecond)
	v = readValue()
	if v != 0 {
		t.Errorf("colorValue: %v, wants %v", v, 0)
	}

	// Stop blink
	l.SetBlink(0)
	time.Sleep(5 * time.Millisecond)
	v = readValue()
	if v != 255 {
		t.Errorf("colorValue: %v, wants %v", v, 255)
	}
	time.Sleep(12 * time.Millisecond)
	v = readValue()
	if v != 255 {
		t.Errorf("colorValue: %v, wants %v", v, 255)
	}
}

func TestColorLed_SetBlinkAndUpdadeColor(t *testing.T) {
	setLedBackup := setLed
	defer func() { setLed = setLedBackup }()

	var muFakeValue sync.Mutex
	var colorValue int
	setLed = func(v int, led gpio.PinIO, mutex *sync.Mutex) {
		mutex.Lock()
		defer mutex.Unlock()
		muFakeValue.Lock()
		defer muFakeValue.Unlock()
		colorValue = v
	}
	readValue := func() int {
		muFakeValue.Lock()
		defer muFakeValue.Unlock()
		return colorValue
	}

	l := New()
	l.SetBlue(255)
	l.SetBlink(100)
	v := readValue()
	if v != 255 {
		t.Errorf("colorValue: %v, wants %v", v, 255)
	}
	time.Sleep(6 * time.Millisecond)
	l.SetBlue(128)

	v = readValue()
	if v != 128 {
		t.Errorf("colorValue: %v, wants %v", v, 128)
	}
	time.Sleep(6 * time.Millisecond)

	time.Sleep(12 * time.Millisecond)

	v = readValue()
	if v != 128 {
		t.Errorf("colorValue: %v, wants %v", v, 128)
	}
	time.Sleep(12 * time.Millisecond)

	v = readValue()
	if v != 0 {
		t.Errorf("colorValue: %v, wants %v", v, 0)
	}

	// Stop blink
	l.SetBlink(0)
	time.Sleep(5 * time.Millisecond)
	v = readValue()
	if v != 128 {
		t.Errorf("colorValue: %v, wants %v", v, 128)
	}
	time.Sleep(12 * time.Millisecond)
	v = readValue()
	if v != 128 {
		t.Errorf("colorValue: %v, wants %v", v, 128)
	}
}

package gpio

import (
	"sync"

	"github.com/dlefevre/go.ventilation-service/config"
	"github.com/stianeikeland/go-rpio/v4"
)

var once sync.Once

// GPIORPiAdapter is an adapter for the Raspberry Pi GPIO pins.
type GPIORPiAdapter struct {
	speed1Pin rpio.Pin
	speed2Pin rpio.Pin
	speed3Pin rpio.Pin
	awayPin   rpio.Pin
	autoPin   rpio.Pin
	timerPin  rpio.Pin
}

// NewGPIORPiAdapter creates a new GPIORPiAdapter.
func NewGPIORPiAdapter() *GPIORPiAdapter {
	once.Do(func() {
		rpio.Open()
	})

	adapter := &GPIORPiAdapter{
		speed1Pin: rpio.Pin(config.GetGPIOSpeed1Pin()),
		speed2Pin: rpio.Pin(config.GetGPIOSpeed2Pin()),
		speed3Pin: rpio.Pin(config.GetGPIOSpeed3Pin()),
		awayPin:   rpio.Pin(config.GetGPIOAwayPin()),
		autoPin:   rpio.Pin(config.GetGPIOAutoPin()),
		timerPin:  rpio.Pin(config.GetGPIOTimerPin()),
	}
	adapter.speed1Pin.Output()
	adapter.speed2Pin.Output()
	adapter.speed3Pin.Output()
	adapter.awayPin.Output()
	adapter.autoPin.Output()
	adapter.timerPin.Output()

	return adapter
}

func (g *GPIORPiAdapter) writePin(pin rpio.Pin, value bool) {
	if value {
		pin.High()
	} else {
		pin.Low()
	}
}

// WriteSpeed1Pin writes a value to the first speed pin.
func (g *GPIORPiAdapter) WriteSpeed1Pin(value bool) {
	g.writePin(g.speed1Pin, value)
}

// WriteSpeed2Pin writes a value to the second speed pin.
func (g *GPIORPiAdapter) WriteSpeed2Pin(value bool) {
	g.writePin(g.speed2Pin, value)
}

// WriteSpeed3Pin writes a value to the third speed pin.
func (g *GPIORPiAdapter) WriteSpeed3Pin(value bool) {
	g.writePin(g.speed3Pin, value)
}

// WriteAwayPin writes a value to the away pin.
func (g *GPIORPiAdapter) WriteAwayPin(value bool) {
	g.writePin(g.awayPin, value)
}

// WriteAutoPin writes a value to the auto pin.
func (g *GPIORPiAdapter) WriteAutoPin(value bool) {
	g.writePin(g.autoPin, value)
}

// WriteTimerPin writes a value to the timer pin.
func (g *GPIORPiAdapter) WriteTimerPin(value bool) {
	g.writePin(g.timerPin, value)
}

package gpio

import "github.com/dlefevre/go.ventilation-service/config"

// GPIOAdapter specifies the interface for GPIO operations.
type GPIOAdapter interface {
	WriteSpeed1Pin(value bool)
	WriteSpeed2Pin(value bool)
	WriteSpeed3Pin(value bool)
	WriteAwayPin(value bool)
	WriteAutoPin(value bool)
	WriteTimerPin(value bool)
}

// GetGPIOAdapter returns the GPIO adapter based on the current mode.
func GetGPIOAdapter() GPIOAdapter {
	switch config.GetMode() {
	case "production":
		return NewGPIORPiAdapter()
	case "development":
		return NewGPIOMockAdapter()
	default:
		panic("Unknown mode")
	}
}

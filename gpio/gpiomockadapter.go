package gpio

import (
	"github.com/dlefevre/go.ventilation-service/config"
	"github.com/rs/zerolog/log"
)

// GPIOMockAdapter is a mock GPIO adapter, which:
// - mimicks the behavior of the garage door, without the delays of a physical door and motor.
// - reports all actions to the log.
type GPIOMockAdapter struct {
	speed1Pin int
	speed2Pin int
	speed3Pin int
	awayPin   int
	autoPin   int
	timerPin  int
}

// NewGPIOMockAdapter creates a new GPIOMockAdapter.
func NewGPIOMockAdapter() *GPIOMockAdapter {
	log.Info().Msg("Mock GPIO: Creating mock GPIO adapter")
	return &GPIOMockAdapter{
		speed1Pin: config.GetGPIOSpeed1Pin(),
		speed2Pin: config.GetGPIOSpeed2Pin(),
		speed3Pin: config.GetGPIOSpeed3Pin(),
		awayPin:   config.GetGPIOAwayPin(),
		autoPin:   config.GetGPIOAutoPin(),
		timerPin:  config.GetGPIOTimerPin(),
	}
}

// WriteSpeed1Pin writes a value to the first speed pin.
func (g *GPIOMockAdapter) WriteSpeed1Pin(value bool) {
	log.Info().Bool("value", value).Msgf("Mock GPIO: Writing to speed 1 pin: %d", g.speed1Pin)
}

// WriteSpeed2Pin writes a value to the second speed pin.
func (g *GPIOMockAdapter) WriteSpeed2Pin(value bool) {
	log.Info().Bool("value", value).Msgf("Mock GPIO: Writing to speed 2 pin: %d", g.speed2Pin)
}

// WriteSpeed3Pin writes a value to the third speed pin.
func (g *GPIOMockAdapter) WriteSpeed3Pin(value bool) {
	log.Info().Bool("value", value).Msgf("Mock GPIO: Writing to speed 3 pin: %d", g.speed3Pin)
}

// WriteAwayPin writes a value to the away pin.
func (g *GPIOMockAdapter) WriteAwayPin(value bool) {
	log.Info().Bool("value", value).Msgf("Mock GPIO: Writing to away pin: %d", g.awayPin)
}

// WriteAutoPin writes a value to the auto pin.
func (g *GPIOMockAdapter) WriteAutoPin(value bool) {
	log.Info().Bool("value", value).Msgf("Mock GPIO: Writing to auto pin: %d", g.autoPin)
}

// WriteTimerPin writes a value to the timer pin.
func (g *GPIOMockAdapter) WriteTimerPin(value bool) {
	log.Info().Bool("value", value).Msgf("Mock GPIO: Writing to timer pin: %d", g.timerPin)
}

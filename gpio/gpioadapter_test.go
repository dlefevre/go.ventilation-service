package gpio

import (
	"os"
	"testing"
)

func init() {
	os.Setenv("VENTILATIONSERVICE_CONFIG_PATH", "..")
	//zerolog.SetGlobalLevel(zerolog.ErrorLevel)
}

func TestWritePins(t *testing.T) {
	adapter := NewGPIOMockAdapter()
	adapter.WriteSpeed1Pin(true)
	adapter.WriteSpeed2Pin(true)
	adapter.WriteSpeed3Pin(true)
	adapter.WriteAwayPin(true)
	adapter.WriteAutoPin(true)
	adapter.WriteTimerPin(true)
}

package controller

import (
	"os"
	"testing"
	"time"
)

func init() {
	// Set the environment variable for the configuration path
	os.Setenv("VENTILATIONSERVICE_CONFIG_PATH", "..")
	//zerolog.SetGlobalLevel(zerolog.ErrorLevel)
}

func TestCreateVentilationController(t *testing.T) {
	controller := GetVentilationControllerService()
	if controller == nil {
		t.Fatalf("Expected controller to be created")
	}
}

func TestStartStop(t *testing.T) {
	controller := GetVentilationControllerService()
	controller.Start()
	time.Sleep(1 * time.Second)
	controller.Stop()
}

func TestRestart(t *testing.T) {
	controller := GetVentilationControllerService()
	controller.Start()
	time.Sleep(1 * time.Second)
	controller.Stop()
	controller.Start()
	time.Sleep(1 * time.Second)
	controller.Stop()
}

func TestSpeedToggles(t *testing.T) {
	controller := GetVentilationControllerService()
	controller.Start()

	controller.SendCommand(CmdSpeed1)
	controller.SendCommand(CmdSpeed2)
	controller.SendCommand(CmdSpeed3)
	time.Sleep(10 * time.Second)

	controller.Stop()
}

func TestOtherToggles(t *testing.T) {
	controller := GetVentilationControllerService()
	controller.Start()

	controller.SendCommand(CmdAway)
	controller.SendCommand(CmdAuto)
	time.Sleep(10 * time.Second)

	controller.Stop()
}

func TestTimerToggles(t *testing.T) {
	controller := GetVentilationControllerService()
	controller.Start()

	controller.SendCommand(CmdTimer15)
	controller.SendCommand(CmdTimer30)
	controller.SendCommand(CmdTimer60)
	time.Sleep(10 * time.Second)

	controller.Stop()
}

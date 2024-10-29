package config

import (
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func init() {
	// Set the environment variable for the configuration path
	os.Setenv("VENTILATIONSERVICE_CONFIG_PATH", "..")
}

func TestVerify(t *testing.T) {
	if err := Verify(); err != nil {
		t.Fatalf("Error verifying configuration: %v", err)
	}
}

func TestMode(t *testing.T) {
	if GetMode() != "development" {
		t.Fatalf("Expected mode to be development, got %s", GetMode())
	}
}

func TestBindPortAndHost(t *testing.T) {
	if GetBindHost() != "127.0.0.1" {
		t.Fatalf("Expected bind  host to be 127.0.0.1, got %s", GetBindHost())
	}
	if GetBindPort() != 8000 {
		t.Fatalf("Expected bind port to be 8080, got %d", GetBindPort())
	}
}

func TestGPIO(t *testing.T) {
	if GetGPIOBackoff() != 3000 {
		t.Fatalf("Expected GPIO backoff to be 3000, got %d", GetGPIOBackoff())
	}
	if GetGPIOSpeed1Pin() != 11 {
		t.Fatalf("Expected GPIO speed 1 pin to be 11, got %d", GetGPIOSpeed1Pin())
	}
	if GetGPIOSpeed2Pin() != 12 {
		t.Fatalf("Expected GPIO speed 2 pin to be 12, got %d", GetGPIOSpeed2Pin())
	}
	if GetGPIOSpeed3Pin() != 13 {
		t.Fatalf("Expected GPIO speed 3 pin to be 13, got %d", GetGPIOSpeed3Pin())
	}
	if GetGPIOAwayPin() != 20 {
		t.Fatalf("Expected GPIO away pin to be 20, got %d", GetGPIOAwayPin())
	}
	if GetGPIOAutoPin() != 21 {
		t.Fatalf("Expected GPIO auto pin to be 21, got %d", GetGPIOAutoPin())
	}
	if GetGPIOTimerPin() != 27 {
		t.Fatalf("Expected GPIO timer pin to be 27, got %d", GetGPIOTimerPin())
	}
}

func TestAPIKeys(t *testing.T) {
	keys := GetAPIKeys()
	if len(keys) != 1 {
		t.Fatalf("Expected 1 API key, got %d", len(keys))
	}
	if err := bcrypt.CompareHashAndPassword([]byte(keys[0]), []byte("test")); err != nil {
		t.Fatalf("Expected API key to be bcrypt digest of 'test'")
	}
}

func TestMQTT(t *testing.T) {
	if !GetMQTTEnabled() {
		t.Fatalf("Expected MQTT to be enabled")
	}
	if GetMQTTURL() != "mqtt://127.0.0.1:1883" {
		t.Fatalf("Expected MQTT URL to be mqtt://127.0.0.1:1883, got %s", GetMQTTURL())
	}
	if GetMQTTClientID() != "ventilation" {
		t.Fatalf("Expected MQTT client ID to be ventilation, got %s", GetMQTTClientID())
	}
	if GetMQTTDiscoveryPrefix() != "homeassistant" {
		t.Fatalf("Expected MQTT topic prefix to be 'homeassistant', got %s", GetMQTTDiscoveryPrefix())
	}
	if GetMQTTID() != "vent01" {
		t.Fatalf("Expected MQTT object ID to be 'vent01', got %s", GetMQTTID())
	}
}

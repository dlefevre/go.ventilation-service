package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

var (
	// All known configuration properties, and weither they are mandatory or not
	knownKeys = map[string]bool{
		"mode":                  true,
		"bind.port":             true,
		"bind.host":             true,
		"gpio.backoff":          true,
		"gpio.pins.speed_1":     true,
		"gpio.pins.speed_2":     true,
		"gpio.pins.speed_3":     true,
		"gpio.pins.away":        true,
		"gpio.pins.auto":        true,
		"gpio.pins.timer":       true,
		"api_keys":              true,
		"mqtt.enabled":          true,
		"mqtt.url":              false,
		"mqtt.username":         false,
		"mqtt.password":         false,
		"mqtt.client_id":        false,
		"mqtt.discovery_prefix": false,
		"mqtt.id":               false,
	}

	viperInst *viper.Viper
	once      sync.Once
)

// Create a new Viper instance and load the configuration file.
func loadConfig() {
	viperInst = viper.New()

	viperInst.SetConfigName("config")
	viperInst.SetConfigType("yaml")
	configPath := os.Getenv("VENTILATIONSERVICE_CONFIG_PATH")
	if configPath != "" {
		viperInst.AddConfigPath(configPath)
	}
	viperInst.AddConfigPath(".")

	if err := viperInst.ReadInConfig(); err != nil {
		panic(fmt.Errorf("config: fatal error while parsing config file: %s", err))
	}
}

// Verifies that all mandatory keys are set in the configuration file,
// and that no unknown keys are present.
func verifyKeys() error {
	for key, mandatory := range knownKeys {
		if mandatory && !viperInst.IsSet(key) {
			return fmt.Errorf("config: configuration property %s is mandatory", key)
		}
	}
	for _, key := range viperInst.AllKeys() {
		if _, found := knownKeys[key]; !found {
			return fmt.Errorf("config: configuration property %s is unknown", key)
		}
	}
	return nil
}

// Verify that the configuration is valid.
func Verify() error {
	once.Do(loadConfig)
	if err := verifyKeys(); err != nil {
		return err
	}
	mode := viperInst.GetString("mode")
	if mode != "development" && mode != "production" {
		return fmt.Errorf("config: mode must be either 'development' or 'production'")
	}
	port := viperInst.GetInt("bind.port")
	if port < 0 || port > 65535 {
		return fmt.Errorf("config: bind.port must be a valid port number")
	}
	gpioBackoff := viperInst.GetInt("gpio.backoff")
	if gpioBackoff < 0 {
		return fmt.Errorf("config: gpio.backoff must be a positive integer")
	}
	for _, pin := range []string{"speed_1", "speed_2", "speed_3", "away", "auto", "timer"} {
		pinNum := viperInst.GetInt(fmt.Sprintf("gpio.pins.%s", pin))
		if pinNum < 2 || pinNum > 27 {
			return fmt.Errorf("config: gpio.pins.%s must be a valid pin number", pin)
		}
	}
	apiKeys := viperInst.GetStringSlice("api_keys")
	if len(apiKeys) == 0 {
		return fmt.Errorf("config: api_keys must contain at least one key")
	}
	mqttEnabled := viperInst.GetBool("mqtt.enabled")
	if mqttEnabled {
		mqttURL := viperInst.GetString("mqtt.url")
		if mqttURL == "" {
			return fmt.Errorf("config: mqtt.url must be set when mqtt.enabled is true")
		}
		if viperInst.GetString("mqtt.client_id") == "" {
			return fmt.Errorf("config: mqtt.client_id must be set when mqtt.enabled is true")
		}
		if viperInst.GetString("mqtt.discovery_prefix") == "" {
			return fmt.Errorf("config: mqtt.discovery_prefix must be set when mqtt.enabled is true")
		}
		if viperInst.GetString("mqtt.username") != "" && viperInst.GetString("mqtt.password") == "" {
			return fmt.Errorf("config: mqtt.password must be set when mqtt.username is set")
		}
		if viperInst.GetString("mqtt.id") == "" {
			return fmt.Errorf("config: mqtt.id must be set when mqtt.enabled is true")
		}
	}

	return nil
}

// GetMode returns the current mode.
func GetMode() string {
	once.Do(loadConfig)
	return viperInst.GetString("mode")
}

// GetBindPort returns the port to bind the web server to.
func GetBindPort() int {
	once.Do(loadConfig)
	return viperInst.GetInt("bind.port")
}

// GetBindHost returns the host to bind the web server to.
func GetBindHost() string {
	once.Do(loadConfig)
	return viperInst.GetString("bind.host")
}

// GetGPIOBackoff returns the backoff time for the GPIO pins.
func GetGPIOBackoff() int {
	once.Do(loadConfig)
	return viperInst.GetInt("gpio.backoff")
}

// GetGPIOSpeed1Pin returns the pin number for the first speed button.
func GetGPIOSpeed1Pin() int {
	once.Do(loadConfig)
	return viperInst.GetInt("gpio.pins.speed_1")
}

// GetGPIOSpeed2Pin returns the pin number for the second speed button.
func GetGPIOSpeed2Pin() int {
	once.Do(loadConfig)
	return viperInst.GetInt("gpio.pins.speed_2")
}

// GetGPIOSpeed3Pin returns the pin number for the third speed button.
func GetGPIOSpeed3Pin() int {
	once.Do(loadConfig)
	return viperInst.GetInt("gpio.pins.speed_3")
}

// GetGPIOAwayPin returns the pin number for the away button.
func GetGPIOAwayPin() int {
	once.Do(loadConfig)
	return viperInst.GetInt("gpio.pins.away")
}

// GetGPIOAutoPin returns the pin number for the auto button.
func GetGPIOAutoPin() int {
	once.Do(loadConfig)
	return viperInst.GetInt("gpio.pins.auto")
}

// GetGPIOTimerPin returns the pin number for the timer button.
func GetGPIOTimerPin() int {
	once.Do(loadConfig)
	return viperInst.GetInt("gpio.pins.timer")
}

// GetAPIKeys returns the list of API keys.
func GetAPIKeys() []string {
	once.Do(loadConfig)
	return viperInst.GetStringSlice("api_keys")
}

// GetMQTTEnabled returns whether MQTT is enabled.
func GetMQTTEnabled() bool {
	once.Do(loadConfig)
	return viperInst.GetBool("mqtt.enabled")
}

// GetMQTTURL returns the URL for the MQTT broker.
func GetMQTTURL() string {
	once.Do(loadConfig)
	if !viperInst.IsSet("mqtt.url") {
		return ""
	}
	return viperInst.GetString("mqtt.url")
}

// GetMQTTUsername returns the client ID for the MQTT client.
func GetMQTTUsername() string {
	once.Do(loadConfig)
	if !viperInst.IsSet("mqtt.username") {
		return ""
	}
	return viperInst.GetString("mqtt.username")
}

// GetMQTTPassword returns the password for the MQTT client.
func GetMQTTPassword() string {
	once.Do(loadConfig)
	if !viperInst.IsSet("mqtt.password") {
		return ""
	}
	return viperInst.GetString("mqtt.password")
}

// GetMQTTClientID returns the client ID for the MQTT client.
func GetMQTTClientID() string {
	once.Do(loadConfig)
	if !viperInst.IsSet("mqtt.client_id") {
		return ""
	}
	return viperInst.GetString("mqtt.client_id")
}

// GetMQTTDiscoveryPrefix returns the discovery prefix for the MQTT client.
func GetMQTTDiscoveryPrefix() string {
	once.Do(loadConfig)
	if !viperInst.IsSet("mqtt.discovery_prefix") {
		return ""
	}
	return viperInst.GetString("mqtt.discovery_prefix")
}

// GetMQTTNodeID returns the object ID for the MQTT client.
func GetMQTTID() string {
	once.Do(loadConfig)
	if !viperInst.IsSet("mqtt.id") {
		return ""
	}
	return viperInst.GetString("mqtt.id")
}

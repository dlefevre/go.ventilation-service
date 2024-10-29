package mqtt

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/dlefevre/go.ventilation-service/controller"
	mochi_mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/hooks/debug"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/rs/zerolog/log"
)

var (
	done   = make(chan bool, 1)
	server = mochi_mqtt.New(&mochi_mqtt.Options{
		InlineClient: true,
	})
)

func init() {
	os.Setenv("VENTILATIONSERVICE_CONFIG_PATH", "..")
	//zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	// Allow all connections.
	_ = server.AddHook(new(auth.AllowHook), &auth.Options{
		Ledger: &auth.Ledger{
			Auth: auth.AuthRules{ // Auth disallows all by default
				{Username: "test", Password: "test", Allow: true},
				{Remote: "127.0.0.1:*", Allow: true},
				{Remote: "localhost:*", Allow: true},
			},
			ACL: auth.ACLRules{ // ACL allows all by default
				{Remote: "127.0.0.1:*"},
				{
					Username: "test", Filters: auth.Filters{
						"homeassistent/#":   auth.ReadWrite,
						"homeassistent/+/+": auth.ReadWrite,
					},
				},
			},
		},
	})
	_ = server.AddHook(new(debug.Hook), nil)

	// Create a TCP listener on a standard port.
	tcp := listeners.NewTCP(listeners.Config{ID: "t1", Address: ":1883"})
	if err := server.AddListener(tcp); err != nil {
		log.Fatal().Msgf("Could not add Listener: %v", err)
	}
}

func startBroker() {
	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal().Msgf("Could not start server: %v", err)
		}
	}()

	// Run server until interrupted
	<-done
}

func TestSingleton(t *testing.T) {
	mqttService := GetMQTTService()
	if mqttService == nil {
		t.Fatalf("Expected MQTTService to be non-nil")
	}
	if GetMQTTService() != mqttService {
		t.Fatalf("Expected GetMQTTService to return the same instance")
	}
}

func TestConnect(t *testing.T) {
	mqttService := GetMQTTService()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := mqttService.Connect(ctx); err != nil {
		t.Fatalf("Error connecting to MQTT broker: %v", err)
	}
}

func TestPublish(t *testing.T) {
	go startBroker()
	time.Sleep(1 * time.Second)

	mqttService := GetMQTTService()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dc := controller.GetVentilationControllerService()
	dc.Start()

	if err := mqttService.Connect(ctx); err != nil {
		t.Fatalf("Error connecting to MQTT broker: %v", err)
	}
	time.Sleep(1 * time.Second)

	log.Info().Msg("Publishing message to action topic")
	if err := server.Publish("homeassistant/button/vent01/action", []byte("speed1"), false, 1); err != nil {
		t.Fatalf("Error publishing message: %v", err)
	}
	time.Sleep(3 * time.Second)
	log.Info().Msg("Publishing message to action topic")
	if err := server.Publish("homeassistant/button/vent01/action", []byte("timer30"), false, 1); err != nil {
		t.Fatalf("Error publishing message: %v", err)
	}

	time.Sleep(1 * time.Second)
	done <- true
}

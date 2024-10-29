package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"github.com/dlefevre/go.ventilation-service/config"
	"github.com/dlefevre/go.ventilation-service/controller"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/rs/zerolog/log"
)

var (
	instance *MQTTManager
	once     sync.Once
)

type Enum uint

var (
	entities = map[string]controller.Enum{
		"speed1":  controller.CmdSpeed1,
		"speed2":  controller.CmdSpeed2,
		"speed3":  controller.CmdSpeed3,
		"away":    controller.CmdAway,
		"auto":    controller.CmdAuto,
		"timer15": controller.CmdTimer15,
		"timer30": controller.CmdTimer30,
		"timer60": controller.CmdTimer60,
	}
)

// MQTTManager is a singleton that encapsulates the MQTT client and .
type MQTTManager struct {
	actionTopic       string
	mqttCfg           autopaho.ClientConfig
	connectionManager *autopaho.ConnectionManager
}

// GetMQTTService returns the one and only MQTTService instance.
func GetMQTTService() *MQTTManager {
	once.Do(func() {
		instance = newMQTTService()
	})
	return instance
}

// Creates a new MQTTService object.
func newMQTTService() *MQTTManager {
	u, err := url.Parse(config.GetMQTTURL())
	if err != nil {
		panic(err)
	}

	mqttService := &MQTTManager{
		actionTopic: fmt.Sprintf("%s/button/%s/action", config.GetMQTTDiscoveryPrefix(), config.GetMQTTID()),
	}
	mqttCfg := autopaho.ClientConfig{
		ServerUrls:                    []*url.URL{u},
		CleanStartOnInitialConnection: false,
		KeepAlive:                     30,
		SessionExpiryInterval:         0,
		OnConnectionUp:                mqttService.connectHandler,
		OnConnectError:                mqttService.connectErrorHandler,
		ClientConfig: paho.ClientConfig{
			ClientID:           config.GetMQTTClientID(),
			OnPublishReceived:  []func(paho.PublishReceived) (bool, error){mqttService.publishHandler},
			OnClientError:      mqttService.clientErrorHandler,
			OnServerDisconnect: mqttService.disconnectHandler,
		},
	}
	if config.GetMQTTUsername() != "" {
		mqttCfg.ConnectUsername = config.GetMQTTUsername()
		mqttCfg.ConnectPassword = []byte(config.GetMQTTPassword())
	}
	mqttService.mqttCfg = mqttCfg
	return mqttService
}

func (s *MQTTManager) Connect(ctx context.Context) error {
	cm, err := autopaho.NewConnection(ctx, s.mqttCfg)
	if err != nil {
		return fmt.Errorf("failed to create connection manager: %v", err)
	}
	s.connectionManager = cm
	return nil
}

func (s *MQTTManager) connectHandler(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
	log.Info().Msgf("connected to MQTT broker: %s", connAck.String())

	// Subscribe to the action topic.
	if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
		Subscriptions: []paho.SubscribeOptions{
			{
				Topic: s.actionTopic,
				QoS:   1,
			},
		},
	}); err != nil {
		log.Error().Msgf("failed to subscribe (%s). This is likely to mean no messages will be received.", err)
	}
	log.Info().Msgf("subscribed to MQTT topic: %s", s.actionTopic)

	s.sendHomeAssistantAutodiscoveryPayload()
}

func (s *MQTTManager) connectErrorHandler(err error) {
	log.Error().Msgf("mqtt connection error: %v", err)
}

func (s *MQTTManager) publishHandler(pr paho.PublishReceived) (bool, error) {
	dc := controller.GetVentilationControllerService()
	command := string(pr.Packet.Payload)
	if cmd, ok := entities[command]; ok {
		dc.SendCommand(cmd)
	} else {
		log.Error().Msgf("received unknown command on action topic: %s", command)
		return false, fmt.Errorf("unknown command: %s", command)
	}

	log.Trace().Msgf("received command '%s' to VentilationControllerService", command)
	return true, nil
}

func (s *MQTTManager) clientErrorHandler(err error) {
	log.Error().Msgf("mqtt client error: %v", err)
}

func (s *MQTTManager) disconnectHandler(d *paho.Disconnect) {
	if d.Properties != nil {
		log.Info().Msgf("server requested disconnect: %s\n", d.Properties.ReasonString)
	} else {
		log.Info().Msgf("server requested disconnect; reason code: %d\n", d.ReasonCode)
	}
}

func (s *MQTTManager) entityPayload(entity string) map[string]interface{} {
	unique_id := fmt.Sprintf("%s%s", config.GetMQTTID(), entity)
	payload := map[string]interface{}{
		"unique_id":        unique_id,
		"command_topic":    s.actionTopic,
		"command_template": entity,
		"name":             unique_id,
		"retain":           true,
		"device": map[string]interface{}{
			"identifiers": []string{
				config.GetMQTTID(),
			},
			"name": config.GetMQTTID(),
		},
	}
	return payload
}

func (s *MQTTManager) publishEntityDiscoveryPayload(entity string) {
	payload := s.entityPayload(entity)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Error().Msgf("failed to marshal discovery payload: %v", err)
		return
	}

	topic := fmt.Sprintf("%s/button/%s%s/config", config.GetMQTTDiscoveryPrefix(), config.GetMQTTID(), entity)
	message := &paho.Publish{
		Topic:   topic,
		Payload: payloadBytes,
		QoS:     1,
		Retain:  true,
	}
	if _, err := s.connectionManager.Publish(context.Background(), message); err != nil {
		log.Error().Msgf("failed to publish discovery payload: %v", err)
	} else {
		log.Info().Msgf("published discovery payload to MQTT topic: %s", topic)
	}
}

func (s *MQTTManager) sendHomeAssistantAutodiscoveryPayload() {
	for entity := range entities {
		s.publishEntityDiscoveryPayload(entity)
	}
}

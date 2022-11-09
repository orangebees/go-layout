package application

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
)

var mqttClient *mqtt.Client

func GetMQTT() *mqtt.Client {
	if mqttClient != nil {
		return mqttClient
	}
	return createMQTT()
}

func createMQTT() *mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://127.0.0.1:1883")
	if v := os.Getenv("MQTT_BROKER"); v != "" {
		opts.AddBroker(v)
	}
	if v := os.Getenv("MQTT_CLIENTID"); v != "" {
		opts.SetClientID(v)
	}
	if v := os.Getenv("MQTT_USERNAME"); v != "" {
		opts.SetUsername(v)
	}
	if v := os.Getenv("MQTT_PASSWORD"); v != "" {
		opts.SetPassword(v)
	}
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	GetLogger().Info().Msg("mqtt connect success")
	return &client
}

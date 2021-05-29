package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"bitbucket.org/aravtechnologies/vibrationtemperature2cloud/src/configdef"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

/*
func TestChecker(t *testing.T) {
	a := checker{QuantityName: "Sample quantity", HighThreshold: 10, LowThreshold: 5}
	a.checkValue(6)
	a.checkValue(12)
	a.checkValue(3)
	a.checkValue(7)
}
*/
func TestMQTT(t *testing.T) {
	const configFilePath = "VT2Cconf.yaml"
	var config configdef.Config
	log.Println("Reading configuration from", configFilePath)
	if checkError(config.ReadFromYAML(configFilePath)) {
		log.Println("Couldn't read configuration.")
		t.Fail()
		os.Exit(1)
	}
	log.Println(config)

	cert, err := tls.LoadX509KeyPair(config.MQTT.CertificatePath, config.MQTT.KeyPath)
	if checkError(err) {
		log.Println("Failed to load certificates.")
		t.Fail()
	}

	clientOpts := &mqtt.ClientOptions{
		ClientID:             config.MQTT.ClientName,
		CleanSession:         true,
		AutoReconnect:        true,
		MaxReconnectInterval: time.Second,
		KeepAlive:            int64(30 * time.Second),
		TLSConfig:            &tls.Config{Certificates: []tls.Certificate{cert}},
	}

	brokerURL := fmt.Sprintf("tcps://%s:%d%s", config.MQTT.Host, config.MQTT.Port, config.MQTT.Path)
	clientOpts.AddBroker(brokerURL)

	mqttClient := mqtt.NewClient(clientOpts)
	log.Println("Connecting to MQTT broker")
	token := mqttClient.Connect()
	token.Wait()
	if checkError(token.Error()) {
		log.Println("Connection to MQTT broker failed")
		t.Fail()
	}

	mqttClient.Publish("Sensor1", 0, false, `{ "message": "Do you see me?" }`)

	mqttClient.Disconnect(250)
}

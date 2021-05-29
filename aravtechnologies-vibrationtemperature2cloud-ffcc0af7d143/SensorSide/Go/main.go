package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"bitbucket.org/aravtechnologies/vibrationtemperature2cloud/src/configdef"
	"bitbucket.org/aravtechnologies/vibrationtemperature2cloud/src/sensorQM30VT2"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	exe, _ := os.Executable()
	var configFilePath = filepath.Dir(exe) + "/VT2Cconf.yaml"
	var config configdef.Config
	log.Println("Reading configuration from", configFilePath)
	if checkError(config.ReadFromYAML(configFilePath)) {
		log.Println("Couldn't read configuration.")
		os.Exit(1)
	}

	var sensor = sensorQM30VT2.New(config.Sensor.Port, config.Sensor.Baud, config.Sensor.Parity, byte(config.Sensor.SlaveID))
	log.Println("Connecting to sensor")
	if !checkError(sensor.Connect()) {
		log.Println("Connected to sensor")
	}

	cert, err := tls.LoadX509KeyPair(config.MQTT.CertificatePath, config.MQTT.KeyPath)
	if checkError(err) {
		log.Println("Failed to load certificates.")
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
	} else {
		log.Println("Connected to MQTT broker")
	}

	var terminator = make(chan os.Signal, 1)
	signal.Notify(terminator, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-terminator
		log.Println("Disconnecting from sensor")
		sensor.Close()
		log.Println("Disconnected from sensor")
		log.Println("Disconnecting from MQTT broker")
		mqttClient.Disconnect(250)
		log.Println("Disconnected from MQTT broker")
		log.Println("Exiting cleanly")
		os.Exit(0)
	}()

	for {
		for i := 0; i < len(config.Quantities); i++ {
			var err error
			switch config.Quantities[i].Name {
			case "Temperature":
				config.Quantities[i].Value, err = sensor.Temperature()
			case "XVibrationVelocity":
				config.Quantities[i].Value, err = sensor.XVelocity()
			case "ZVibrationVelocity":
				config.Quantities[i].Value, err = sensor.ZVelocity()
			case "XPeakVibrationVelocity":
				config.Quantities[i].Value, err = sensor.XPeakVelocity()
			case "ZPeakVibrationVelocity":
				config.Quantities[i].Value, err = sensor.ZPeakVelocity()
			case "XVibrationFrequency":
				config.Quantities[i].Value, err = sensor.XFrequency()
			case "ZVibrationFrequency":
				config.Quantities[i].Value, err = sensor.ZFrequency()
			}
			if checkError(err) {
				log.Printf("Error reading %s from sensor", config.Quantities[i].Name)
				config.Quantities[i].Alert = "Error"
				config.Quantities[i].Value = 0
			} else {
				switch config.Quantities[i].OutOfBounds() {
				case 0:
					log.Printf("%s %f too low!", config.Quantities[i].Name, config.Quantities[i].Value)
					config.Quantities[i].Alert = "Low"
				case 2:
					log.Printf("%s %f too high!", config.Quantities[i].Name, config.Quantities[i].Value)
					config.Quantities[i].Alert = "High"
				default:
					log.Printf("%s %f okay.", config.Quantities[i].Name, config.Quantities[i].Value)
					config.Quantities[i].Alert = "None"
				}
			}
		}

		var pubMsg = publishMsg{SensorID: config.MQTT.ClientName, Quantities: config.Quantities}
		jsonBytes, err := json.Marshal(pubMsg)
		if checkError(err) {
			log.Println("Failed to marshal JSON")
		} else {
			log.Println("Publishing values")
			mqttClient.Publish("values", 0, false, string(jsonBytes))
		}

		time.Sleep(time.Millisecond * time.Duration(config.Sensor.MillisBetweenPolls))
	}
}

type publishMsg struct {
	SensorID   string               `json:"SensorID"`
	Quantities []configdef.Quantity `json:"Quantities"`
}

func checkError(e error) bool {
	var errorExists bool = e != nil
	if errorExists {
		log.Println(e.Error())
	}
	return errorExists
}

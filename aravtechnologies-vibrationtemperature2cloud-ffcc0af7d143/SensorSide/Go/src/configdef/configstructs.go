package configdef

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	MQTT       MQTTConfig   `yaml:"MQTT"`
	Sensor     SensorConfig `yaml:"Sensor"`
	Quantities []Quantity   `yaml:"Quantities"`
}

type SensorConfig struct {
	Baud               int    `yaml:"Baud"`
	SlaveID            int    `yaml:"SlaveID"`
	Port               string `yaml:"Port"`
	Parity             string `yaml:"Parity"`
	MillisBetweenPolls int    `yaml:"MillisBetweenPolls"`
}

type Quantity struct {
	Name  string  `yaml:"Name" json:"Name"`
	Max   float64 `yaml:"Max"`
	Min   float64 `yaml:"Min"`
	Value float64 `json:"Value"`
	Alert string  `json:"Alert"`
}

type MQTTConfig struct {
	CertificatePath string `yaml:"CertificatePath"`
	KeyPath         string `yaml:"KeyPath"`
	ClientName      string `yaml:"ClientName"`
	Host            string `yaml:"Host"`
	Port            int    `yaml:"Port"`
	Path            string `yaml:"Path"`
}

func (c *Config) ReadFromYAML(filepath string) (err error) {
	configfileBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(configfileBytes, c)
	return
}

// OutOfBounds returns 2 if Quantity.Value is too high, 0 if too low, and 1 if okay.
func (qty Quantity) OutOfBounds() uint {
	if qty.Value > qty.Max {
		return 2
	} else if qty.Value < qty.Min {
		return 0
	} else {
		return 1
	}
}

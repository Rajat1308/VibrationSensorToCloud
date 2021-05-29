package configdef

import (
	"log"
	"testing"
)

/*
func TestYAMLParse(t *testing.T) {
	var config Config
	err := config.ReadFromYAML("/home/prithvi/Projects/vibrationtemperature2cloud/GoCodeOnPC/VT2Cconf.yaml")
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(config)
	}
}
*/

func TestQuantity(t *testing.T) {
	qty := Quantity{Name: "Sample", Min: 3, Max: 6}
	qty.Value = 1
	log.Println(qty.OutOfBounds())
	qty.Value = 5
	log.Println(qty.OutOfBounds())
	qty.Value = 9
	log.Println(qty.OutOfBounds())
}

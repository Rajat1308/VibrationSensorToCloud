package sensorQM30VT2

import (
	"github.com/goburrow/modbus"
)

type QM30VT2 struct {
	mbhandler *modbus.RTUClientHandler
}

func New(port string, baud int, parity string, slaveid byte) QM30VT2 {
	handler := modbus.NewRTUClientHandler(port)
	handler.BaudRate = baud
	handler.Parity = parity
	handler.SlaveId = slaveid
	return QM30VT2{mbhandler: handler}
}

// Connect connects to QM30VT2 sensor
func (qm QM30VT2) Connect() error {
	return qm.mbhandler.Connect()
}

// Close closes connection with QM30VT2 sensor
func (qm QM30VT2) Close() error {
	return qm.mbhandler.Close()
}

func (qm QM30VT2) readFromeRegistersAndDivide(firstRegister uint16, numberOfRegisters uint16, divisor float64) (values []float64, err error) {
	client := modbus.NewClient(qm.mbhandler)
	incoming, err := client.ReadHoldingRegisters(firstRegister, numberOfRegisters)
	for i := 0; i < len(incoming); i += 2 {
		v := int16(incoming[i])<<8 | int16(incoming[i+1])
		values = append(values, float64(v)/divisor)
	}
	return
}

func (qm QM30VT2) checkedResult(val []float64, err error) (float64, error) {
	if err == nil && len(val) != 0 {
		return val[0], err
	}
	return 0, err
}

// TemperatureCelsius gets temperature in celsius
func (qm QM30VT2) Temperature() (temp float64, err error) {
	return qm.checkedResult(qm.readFromeRegistersAndDivide(5203, 1, 100))
}

// XPeakVelocity gets X-axis peak RMS velocity in millimeters per second
func (qm QM30VT2) XPeakVelocity() (XVelocity float64, err error) {
	return qm.checkedResult(qm.readFromeRegistersAndDivide(5219, 1, 1000))
}

// ZPeakVelocity gets Z-axis peak RMS velocity in millimeters per second
func (qm QM30VT2) ZPeakVelocity() (ZVelocity float64, err error) {
	return qm.checkedResult(qm.readFromeRegistersAndDivide(5216, 1, 1000))
}

// XFrequency gets X-axis peak velocity frequency component in Hertz
func (qm QM30VT2) XFrequency() (XFrequency float64, err error) {
	return qm.checkedResult(qm.readFromeRegistersAndDivide(5209, 1, 10))
}

// ZFrequency gets Z-axis peak velocity frequency component in Hertz
func (qm QM30VT2) ZFrequency() (ZFrequency float64, err error) {
	return qm.checkedResult(qm.readFromeRegistersAndDivide(5208, 1, 10))
}

// XVelocity gets X-axis RMS velocity in millimeters per second
func (qm QM30VT2) XVelocity() (XVelocity float64, err error) {
	return qm.checkedResult(qm.readFromeRegistersAndDivide(5205, 1, 1000))
}

// ZVelocity gets Z-axis RMS velocity in millimeters per second
func (qm QM30VT2) ZVelocity() (ZVelocity float64, err error) {
	return qm.checkedResult(qm.readFromeRegistersAndDivide(5201, 1, 1000))
}

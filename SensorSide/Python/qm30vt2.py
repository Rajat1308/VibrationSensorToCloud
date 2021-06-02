from pymodbus.client.sync import ModbusSerialClient
from pymodbus.constants import Defaults

class QM30VT2:
    def __init__(self, port, baud, parity, slaveid):
        self.slaveid = slaveid
        self.client = ModbusSerialClient(method="rtu", port=port, stopbits=1, parity=parity, baudrate=baud)
    
    def _read_from_register(self, register):
        response = self.client.read_holding_registers(register, count=1, unit=self.slaveid)
        return response.registers[0]
    
    def _read_scaled_value(self, register, divisor):
        return self._read_from_register(register) / divisor
    
    def Temperature(self):
        return self._read_scaled_value(5203, 100)
    
    def XPeakVelocity(self):
        return self._read_scaled_value(5219, 1000)
    
    def ZPeakVelocity(self):
        return self._read_scaled_value(5216, 1000)
    
    def XFrequency(self):
        return self._read_scaled_value(5209, 10)
    
    def ZFrequency(self):
        return self._read_scaled_value(5208, 10)
    
    def XVelocity(self):
        return self._read_scaled_value(5205, 1000)
    
    def ZVelocity(self):
        return self._read_scaled_value(5201, 1000)
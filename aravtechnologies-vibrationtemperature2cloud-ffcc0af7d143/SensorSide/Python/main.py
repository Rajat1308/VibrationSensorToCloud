#!/usr/bin/python3

from mqtt_client import MQTTClient
from quantity_def import Quantity
from qm30vt2 import QM30VT2
from time import sleep
import yaml
import os
import signal
import logging

logging.basicConfig(level=logging.INFO)

config_file_path = os.path.join(os.path.dirname(os.path.realpath(__file__)), "VT2Cconf.yaml")
logging.info("Reading configuration from {}".format(config_file_path))
try:
    with open(config_file_path, 'r') as yamlfile:
        config = yaml.safe_load(yamlfile)
    logging.info("Read configuration successfully.")
except:
    logging.critical("Failed to read configuration!")
    exit(1)

sensor = QM30VT2(
    config["Sensor"]["Port"],
    config["Sensor"]["Baud"],
    config["Sensor"]["Parity"],
    config["Sensor"]["SlaveID"]
)

logging.info("Connecting to MQTT broker")
mqtt = MQTTClient(config["MQTT"]["ClientName"])
mqtt.set_tls(
    config["MQTT"]["CAPath"],
    config["MQTT"]["CertificatePath"],
    config["MQTT"]["KeyPath"]
)
mqtt.connect(
    config["MQTT"]["Host"],
    config["MQTT"]["Port"]
)

function_for_qty_name = {
    "Temperature": sensor.Temperature,
    "XVibrationVelocity": sensor.XVelocity,
    "ZVibrationVelocity": sensor.ZVelocity,
    "XPeakVibrationVelocity": sensor.XPeakVelocity,
    "ZPeakVibrationVelocity": sensor.ZPeakVelocity,
    "XVibrationFrequency": sensor.XFrequency,
    "ZVibrationFrequency": sensor.ZFrequency
}

quantities = []
for q in config["Quantities"]:
    if q["Name"] in function_for_qty_name.keys():
        quantities.append(Quantity(q["Name"], q["Min"], q["Max"]))
    else:
        logging.warning("Invalid quantity name {}, ignoring".format(q.name))

class GracefulKiller:
    def __init__(self):
        self.killNow = False
        signal.signal(signal.SIGINT, self.exitGracefully)
    
    def exitGracefully(self, signum, frame):
        self.killNow = True


def assemble_message_dict(sensorID, qtys):
    msg_dict = {
        "SensorID": sensorID,
        "Quantities": []
    }
    for quantity in qtys:
        q = dict()
        q["Name"] = quantity.name
        q["Min"] = quantity.min
        q["Max"] = quantity.max
        q["Value"] = quantity.value
        q["Alert"] = quantity.alert
        msg_dict["Quantities"].append(q)
    return msg_dict

killer = GracefulKiller()

mqtt.start_loop()
while not killer.killNow:
    for quantity in quantities:
        try:
            quantity.value = function_for_qty_name[quantity.name]()
            value_state = quantity.out_of_bounds()
            if value_state == Quantity.TOO_HIGH:
                logging.warning("{} {} too high!".format(quantity.name, quantity.value))
                quantity.alert = "High"
            elif value_state == Quantity.TOO_LOW:
                logging.warning("{} {} too low!".format(quantity.name, quantity.value))
                quantity.alert = "Low"
            else:
                quantity.alert = "None"
                logging.info("{} {} okay.".format(quantity.name, quantity.value))
        except:
            logging.error("Reading from sensor failed!")
            quantity.alert = "Error"
    
    
    mqtt.publish_as_json(
        "values",
        assemble_message_dict(
            config["MQTT"]["ClientName"],
            quantities
        )
    )
    
    sleep(config["Sensor"]["MillisBetweenPolls"] / 1000)

mqtt.stop_loop()
logging.info("Exiting cleanly")
logging.info("Disconnecting from MQTT broker")
mqtt.disconnect()
import paho.mqtt.client as mqtt
import json
import logging

class MQTTClient():
    def __init__(self, clientID):
        self._client = mqtt.Client(client_id=clientID)
        self._client.on_connect = self._on_connect
        self._client.on_disconnect = self._on_disconnect
    
    def set_tls(self, ca_cert, client_cert, client_key):
        self._client.tls_set(ca_cert, client_cert, client_key)
    
    def connect(self, host, port):
        self._client.connect(host, port)
    
    def disconnect(self):
        self._client.disconnect()

    def loop_once(self):
        self._client.loop()
    
    def start_loop(self):
        self._client.loop_start()
    
    def stop_loop(self):
        self._client.loop_stop()

    def publish_as_json(self, topic, payload_dict):
        if not self._client.publish(topic, json.dumps(payload_dict), 0, False).is_published():
            logging.error("Publish failed!")
    
    def _on_connect(self, client, userdata, flags, returncode):
        if returncode == 0:
            logging.info("Connected to MQTT broker")
        else:
            logging.error("Connection failed with return code {}".format(returncode))
    
    def _on_disconnect(self, client, userdata, returncode):
        if returncode == 0:
            logging.info("Disconnected from MQTT broker as intended")
        else:
            logging.error("Disconnected unexpectedly with return code {}".format(returncode))
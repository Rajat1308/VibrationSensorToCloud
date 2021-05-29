# VibrationTemperature2Cloud

## Aims
- Have a computer sample **temperature** and **vibration** data from a **QM30VT2** sensor
- **Check** this data on the spot for whether it is within acceptable limits, and **alert** if not
- Store this data in a **cloud database**
- Have the data **analyzed** in the cloud
- Present the data in a **web-based dashboard**

---

## QM30VT2 Sensor

This sensor uses the **modbus** protocol over **RS-485**. Relevant information, including the modbus registers for specific data, can be found in its [datasheet](https://info.bannerengineering.com/cs/groups/public/documents/literature/210732.pdf). It connects to a computer through an **RS-485 to USB converter**.

---

## Program on the Sensor Side

There are two impementations of the code on the sensor side, one in [Golang](golang.org) and the other in [Python](python.org).

The Go implementation uses [go modbus](https://github.com/goburrow/modbus), the [Eclipse Paho](github.com/eclipse/paho.mqtt.golang) MQTT library, and [go YAML](gopkg.in/yaml.v2).

The Python implementation uses [pymodbus](https://pypi.org/project/pymodbus/), the [Eclipse Paho](https://pypi.org/project/paho-mqtt/) MQTT library, and [pyYAML](https://pypi.org/project/PyYAML/).

---

#### Relevant Files

The executable and [`VT2Cconf.yaml`](SensorSide/VT2Cconf.yaml) must reside in the same folder. A client certificate and key is required for TLS-encrypted MQTT. The [`systemd` unit file](SensorSide/vt2c.service) must be placed in `/etc/systemd/system` on a Linux system.

---

#### Configuration

Configuration is done through a single YAML configuration file, called [`VT2Cconf.yaml`](SensorSide/VT2Cconf.yaml) in the same directory as the program executable.

It must have three sections, `MQTT`, `Sensor`, and `Quantities`.

In the `MQTT` section, ensure that `CertificatePath` and `KeyPath` for TLS are absolute paths. `Path` is the URL path appended to `Host`.

The `Sensor` section is quite self-explanatory.

Each Quantity has a `Name`, `Max`, and `Min`. Only the quantities specified will be sampled.

`Name` can be one of the following:
- `Temperature` for temperature in Celsius
- `XVibrationVelocity` for RMS vibration velocity in the X-axis.
- `ZVibrationVelocity` for RMS vibration velocity in the Z-axis.
- `XVibrationFrequency` for the frequency component of vibration in the X-axis.
- `ZVibrationFrequency` for the frequency component of vibration in the Z-axis.
- `XPeakVibrationVelocity` for the peak velocity in the X-axis.
- `ZPeakVibrationVelocity` for the peak velocity in the Z-axis.

`Max` and `Min` are floating-point numbers that define the upper and lower thresholds of the quantity. If the value reported by the sensor goes above `Max` or below `Min`, the program will alert accordingly.

**All fields are mandatory.**

---

#### More on MQTT

All publishes are made to a topic called `values`. Each publish is JSON, in the following format:
```
{
    "SensorID": "SameAsClientName",
    "Quantities": [
        {
            "Name": "Temperature",
            "Min": 0.5,
            "Max": 10.4,
            "Value": 4.8,
            "Alert": "None"
        },
        ...
    ]
}
```
As of now, `Min` and `Max` are sent with every publish, though later, they could be published to a retained topic.

---

#### Alerts

Alerts are printed out to the console, as well as sent over MQTT. Alerts can have the following values:
- `None`
- `High`: Above `Max`
- `Low`: Below `Min`
- `Error`: Error while reading from sensor, `Value` is invalid

---

## AWS

Amazon Web Services was the cloud platform of choice for this project. They had some ready [articles](https://aws.amazon.com/iot/solutions/industrial-iot/) and [infographics](https://d1.awsstatic.com/IoT/Predictive%20Maintenance%20Infographic.pdf) that were suited to the requirements of this project.

[This tutorial](https://aws.amazon.com/blogs/big-data/build-a-visualization-and-monitoring-dashboard-for-iot-data-with-amazon-kinesis-analytics-and-amazon-quicksight/) goes through most of the requirements of this project, and is a good starting point.

Before going through the Kinesis setup, make sure AWS is receiving MQTT data from the sensor in the format above, to avoid having to edit the Kinesis schema later. With that said, follow all steps as instructed, till the Amazon QuickSight setup.

#### Viewing Data on QuickSight

First, navigate to the Amazon Athena dashboard, and create a new data source. You will be redirected to Amazon Glue. There, make a new Crawler to query your S3 database with aggregate values, and dump them to a database within Amazon Athena. Make this crawler run every five minutes, and have it import only new rows.

In QuickSight, create a new dataset from Amazon Athena, and choose the database you just created. Create your analyses in QuickSight, and publish them as dashboards.

---

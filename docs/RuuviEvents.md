## NewMeasurementData

`NewMeasurementData` is a struct type in Go that represents a new measurement data event. It has the following fields:

- `DataFormat`: An integer representing the data format.
- `Temperature`: A float representing the temperature measurement.
- `Humidity`: A float representing the humidity measurement.
- `Pressure`: An integer representing the pressure measurement.
- `Acceleration`: A struct representing the acceleration measurement. It has three integer fields: `X`, `Y`, and `Z`.
- `Battery`: An integer representing the battery level.
- `TXPower`: An integer representing the transmission power.
- `Movement`: An integer representing the movement measurement.
- `Sequence`: An integer representing the sequence number.
- `MAC`: A string representing the MAC address.
- `RSSI`: An integer representing the Received Signal Strength Indicator (RSSI).
- `Address`: A string representing the address.
- `LocalName`: A string representing the local name.

This type is used to unmarshal the JSON data into a Go object.

```json
{
    "DataFormat": 5,
    "Temperature": 22.34,
    "Humidity": 42.975,
    "Pressure": 97465,
    "Acceleration": {
        "X": -8,
        "Y": -20,
        "Z": 1056
    },
    "Battery": 2857,
    "TXPower": 4,
    "Movement": 75,
    "Sequence": 6256,
    "MAC": "E8:D3:AD:C4:6E:18",
    "RSSI": -75,
    "Address": "E8:D3:AD:C4:6E:18",
    "LocalName": ""
}
```

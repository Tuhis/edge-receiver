package events

import "encoding/json"

type RuuviEventTypes string

const (
	NewMeasurement RuuviEventTypes = "new_measurement"
)

type RuuviEvent struct {
	Type       RuuviEventTypes `json:"type"`
	Data       json.RawMessage `json:"data"`
	SourceUuid string          `json:"source_uuid"`
}

type RuuviKafkaEvent struct {
	Type       RuuviEventTypes `json:"type"`
	Data       interface{}     `json:"data"`
	SourceUuid string          `json:"source_uuid"`
}

type NewMeasurementData struct {
	DataFormat   *int     `json:"DataFormat"`
	Temperature  *float64 `json:"Temperature"`
	Humidity     *float64 `json:"Humidity"`
	Pressure     *int     `json:"Pressure"`
	Acceleration *struct {
		X *int `json:"X"`
		Y *int `json:"Y"`
		Z *int `json:"Z"`
	} `json:"Acceleration"`
	Battery   *int    `json:"Battery"`
	TXPower   *int    `json:"TXPower"`
	Movement  *int    `json:"Movement"`
	Sequence  *int    `json:"Sequence"`
	MAC       *string `json:"MAC"`
	RSSI      *int    `json:"RSSI"`
	Address   *string `json:"Address"`
	LocalName *string `json:"LocalName"`
}

type NewMeasurementDataNoPtr struct {
	DataFormat   int     `json:"DataFormat"`
	Temperature  float64 `json:"Temperature"`
	Humidity     float64 `json:"Humidity"`
	Pressure     int     `json:"Pressure"`
	Acceleration struct {
		X int `json:"X"`
		Y int `json:"Y"`
		Z int `json:"Z"`
	} `json:"Acceleration"`
	Battery   int    `json:"Battery"`
	TXPower   int    `json:"TXPower"`
	Movement  int    `json:"Movement"`
	Sequence  int    `json:"Sequence"`
	MAC       string `json:"MAC"`
	RSSI      int    `json:"RSSI"`
	Address   string `json:"Address"`
	LocalName string `json:"LocalName"`
}

// IsValid checks if all fields in the NewMeasurementData struct are not nil.
// It returns true if all fields are not nil, and false otherwise.
//
// This method is useful for validating that a NewMeasurementData object
// has all required fields before using it. For example, it can be used
// after unmarshalling a JSON object to ensure that all fields were present
// in the JSON data.
//
// Note that this method does not check if the fields have valid values,
// it only checks if they are not nil. Additional validation may be needed
// depending on the use case.
func (n *NewMeasurementData) IsValid() bool {
	return n.DataFormat != nil &&
		n.Temperature != nil &&
		n.Humidity != nil &&
		n.Pressure != nil &&
		n.Acceleration != nil &&
		n.Acceleration.X != nil &&
		n.Acceleration.Y != nil &&
		n.Acceleration.Z != nil &&
		n.Battery != nil &&
		n.TXPower != nil &&
		n.Movement != nil &&
		n.Sequence != nil &&
		n.MAC != nil &&
		n.RSSI != nil &&
		n.Address != nil &&
		n.LocalName != nil
}

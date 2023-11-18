package events_test

import (
	"testing"

	"github.com/Tuhis/edge-receiver/pkg/events"
)

func TestNewMeasurementData_IsValid(t *testing.T) {
	tests := []struct {
		name string
		data events.NewMeasurementData
		want bool
	}{
		{
			name: "All fields present",
			data: events.NewMeasurementData{
				DataFormat:  intPtr(1),
				Temperature: float64Ptr(1.0),
				Humidity:    float64Ptr(1.0),
				Pressure:    intPtr(1),
				Acceleration: structPtr(struct {
					X *int `json:"X"`
					Y *int `json:"Y"`
					Z *int `json:"Z"`
				}{X: intPtr(1), Y: intPtr(1), Z: intPtr(1)}),
				Battery:   intPtr(1),
				TXPower:   intPtr(1),
				Movement:  intPtr(1),
				Sequence:  intPtr(1),
				MAC:       stringPtr("MAC"),
				RSSI:      intPtr(1),
				Address:   stringPtr("Address"),
				LocalName: stringPtr("LocalName"),
			},
			want: true,
		},
		{
			name: "Missing field",
			data: events.NewMeasurementData{
				DataFormat:  intPtr(1),
				Temperature: float64Ptr(1.0),
				Humidity:    float64Ptr(1.0),
				Pressure:    intPtr(1),
				Acceleration: structPtr(struct {
					X *int `json:"X"`
					Y *int `json:"Y"`
					Z *int `json:"Z"`
				}{X: intPtr(1), Y: intPtr(1), Z: intPtr(1)}),
				Battery:   intPtr(1),
				TXPower:   intPtr(1),
				Movement:  intPtr(1),
				Sequence:  intPtr(1),
				MAC:       stringPtr("MAC"),
				RSSI:      intPtr(1),
				Address:   nil, // Missing field
				LocalName: stringPtr("LocalName"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.data.IsValid(); got != tt.want {
				t.Errorf("NewMeasurementData.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func stringPtr(s string) *string {
	return &s
}

func structPtr(s struct {
	X *int `json:"X"`
	Y *int `json:"Y"`
	Z *int `json:"Z"`
}) *struct {
	X *int `json:"X"`
	Y *int `json:"Y"`
	Z *int `json:"Z"`
} {
	return &s
}

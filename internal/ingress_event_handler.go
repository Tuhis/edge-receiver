package internal

import (
	// "encoding/json"
	// "fmt"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Tuhis/edge-receiver/pkg/apiresponse"
	"github.com/Tuhis/edge-receiver/pkg/events"
)

func CreateIncomingEventHandler(messageChan chan<- string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log incoming request
		log.Printf("Received a request: %s %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)

		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path != "/event" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(apiresponse.ApiResponse{Message: apiresponse.NotFound})
			return
		}

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(apiresponse.ApiResponse{Message: apiresponse.InvalidRequest})
			return
		}

		// Create a new RuuviEvent object
		var event events.RuuviEvent

		// Try to decode the request body into the RuuviEvent object
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(apiresponse.ApiResponse{Message: apiresponse.InvalidRequest})
			return
		}

		if event.SourceUuid == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(apiresponse.ApiResponse{Message: apiresponse.InvalidRequest})
			return
		}

		// Handle different types of events
		switch event.Type {
		case events.NewMeasurement:
			// Create a new NewMeasurementData object
			var data events.NewMeasurementData

			// Try to unmarshal the Data field into the NewMeasurementData object
			err := json.Unmarshal(event.Data, &data)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(apiresponse.ApiResponse{Message: apiresponse.InvalidRequest})
				return
			}

			// Check if all fields are present
			if !data.IsValid() {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(apiresponse.ApiResponse{Message: apiresponse.InvalidRequest})
				return
			}

			// Handle the NewMeasurementData here
			// For example, print it to the console
			fmt.Printf("Received new measurement data: %+v\n", data)

			// Send the data to the Kafka messaging channel
			kafkaEvent := events.RuuviKafkaEvent{
				Type:       event.Type,
				Data:       data,
				SourceUuid: event.SourceUuid,
			}
			kafkaEventJson, err := json.Marshal(kafkaEvent)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(apiresponse.ApiResponse{Message: apiresponse.InvalidRequest})
				return
			}
			messageChan <- string(kafkaEventJson)

			// Return a response to the client
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(apiresponse.ApiResponse{Message: apiresponse.Ok})

			return

		default:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(apiresponse.ApiResponse{Message: apiresponse.UnknownEvent})
			return
		}
	}
}

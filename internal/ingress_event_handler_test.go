package internal_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Tuhis/edge-receiver/internal"
)

func TestHandleIncomingEvent(t *testing.T) {
	tests := []struct {
		name                string
		method              string
		url                 string
		body                string
		expectedStatus      int
		expectedChanMessage string
	}{
		{
			name:                "valid request",
			method:              "POST",
			url:                 "/event",
			body:                `{"type":"new_measurement","source_uuid":"7d01818b-0332-4adf-99c1-13f833e59c6b","data":{"DataFormat":5,"Temperature":22.34,"Humidity":42.975,"Pressure":97465,"Acceleration":{"X":-8,"Y":-20,"Z":1056},"Battery":2857,"TXPower":4,"Movement":75,"Sequence":6256,"MAC":"E8:D3:AD:C4:6E:18","RSSI":-75,"Address":"E8:D3:AD:C4:6E:18","LocalName":""}}`,
			expectedStatus:      http.StatusOK,
			expectedChanMessage: "{\"type\":\"new_measurement\",\"data\":{\"DataFormat\":5,\"Temperature\":22.34,\"Humidity\":42.975,\"Pressure\":97465,\"Acceleration\":{\"X\":-8,\"Y\":-20,\"Z\":1056},\"Battery\":2857,\"TXPower\":4,\"Movement\":75,\"Sequence\":6256,\"MAC\":\"E8:D3:AD:C4:6E:18\",\"RSSI\":-75,\"Address\":\"E8:D3:AD:C4:6E:18\",\"LocalName\":\"\"},\"source_uuid\":\"7d01818b-0332-4adf-99c1-13f833e59c6b\"}",
		},
		{
			name:                "invalid method",
			method:              "GET",
			url:                 "/event",
			body:                "",
			expectedStatus:      http.StatusMethodNotAllowed,
			expectedChanMessage: "",
		},
		{
			name:                "invalid url",
			method:              "POST",
			url:                 "/invalid",
			body:                "",
			expectedStatus:      http.StatusNotFound,
			expectedChanMessage: "",
		},
		{
			name:                "invalid body",
			method:              "POST",
			url:                 "/event",
			body:                `{"type":"new_measurement","data":{}}`,
			expectedStatus:      http.StatusBadRequest,
			expectedChanMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			messageChan := make(chan string, 1)
			handler := http.HandlerFunc(internal.CreateIncomingEventHandler(messageChan))

			handler.ServeHTTP(rr, req)

			if tt.expectedChanMessage != "" {
				receivedMsg := <-messageChan
				if receivedMsg != tt.expectedChanMessage {
					t.Errorf("handler did not send correct message to channel: got %v want %v",
						receivedMsg, tt.expectedChanMessage)
				}
			}

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

package kafkawrapper

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap/zaptest"
)

func TestConfigFromEnvironment(t *testing.T) {
	tests := []struct {
		name          string
		brokers       string
		statusTopic   string
		ownName       string
		authMechanism string
		username      string
		password      string
		wantErr       bool
	}{
		{
			name:          "All environment variables set, PLAIN auth",
			brokers:       "localhost:9092",
			statusTopic:   "status",
			ownName:       "test",
			authMechanism: "PLAIN",
			wantErr:       false,
		},
		{
			name:          "All environment variables set, SCRAM auth",
			brokers:       "localhost:9092",
			statusTopic:   "status",
			ownName:       "test",
			authMechanism: "SCRAM-SHA-256",
			username:      "user",
			password:      "pass",
			wantErr:       false,
		},
		{
			name:          "Missing KAFKA_USERNAME and KAFKA_PASSWORD for SCRAM auth",
			brokers:       "localhost:9092",
			statusTopic:   "status",
			ownName:       "test",
			authMechanism: "SCRAM-SHA-256",
			wantErr:       true,
		},
		{
			name:        "All environment variables set",
			brokers:     "localhost:9092",
			statusTopic: "status",
			ownName:     "test",
			wantErr:     false,
		},
		{
			name:        "Missing KAFKA_BROKERS",
			statusTopic: "status",
			ownName:     "test",
			wantErr:     true,
		},
		{
			name:    "Missing KAFKA_STATUS_TOPIC",
			brokers: "localhost:9092",
			ownName: "test",
			wantErr: true,
		},
		{
			name:        "Missing OWN_NAME",
			brokers:     "localhost:9092",
			statusTopic: "status",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("KAFKA_BROKERS", tt.brokers)
			os.Setenv("KAFKA_STATUS_TOPIC", tt.statusTopic)
			os.Setenv("OWN_NAME", tt.ownName)
			os.Setenv("KAFKA_AUTH_MECHANISM", tt.authMechanism)
			os.Setenv("KAFKA_USERNAME", tt.username)
			os.Setenv("KAFKA_PASSWORD", tt.password)

			_, err := configFromEnvironment()

			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigFromEnvironment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type MockWriter struct {
	Messages []kafka.Message
}

func (mw *MockWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	mw.Messages = append(mw.Messages, msgs...)
	return nil
}

func (mw *MockWriter) Close() error {
	return nil
}

func TestKafkaProducer(t *testing.T) {

	// Create a new logger
	logger := zaptest.NewLogger(t)

	// Create a new Kafka producer with a mock writer
	writer := &MockWriter{}
	producer := &KafkaProducer{
		w:           writer,
		log:         logger.Sugar(),
		statusTopic: "status",
		ownName:     "test",
	}

	// Test ProduceMessage
	err := producer.ProduceMessage("test message", "test topic")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test ProduceRawMessage
	err = producer.ProduceRawMessage([]byte("test raw message"), "test topic")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test ProduceMessagesFromChan
	messageChan := make(chan string, 1)
	messageChan <- "test message from chan"
	close(messageChan)
	producer.ProduceMessagesFromChan(messageChan, "test topic")

	// Test ProduceMessagesFromRawChan
	rawMessageChan := make(chan []byte, 1)
	rawMessageChan <- []byte("test raw message from chan")
	close(rawMessageChan)
	producer.ProduceMessagesFromRawChan(rawMessageChan, "test topic")

	// Check that the messages were written correctly
	expectedMessages := []kafka.Message{
		{Topic: "test topic", Value: []byte("test message")},
		{Topic: "test topic", Value: []byte("test raw message")},
		{Topic: "test topic", Value: []byte("test message from chan")},
		{Topic: "test topic", Value: []byte("test raw message from chan")},
	}
	if !reflect.DeepEqual(writer.Messages, expectedMessages) {
		t.Errorf("Expected messages %v, got %v", expectedMessages, writer.Messages)
	}
}

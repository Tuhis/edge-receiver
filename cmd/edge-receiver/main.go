package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Tuhis/edge-receiver/internal"
	"github.com/Tuhis/edge-receiver/pkg/kafkawrapper"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

func main() {

	// Init logger
	logger, _ := zap.NewProduction()

	// Read KAFKA_INGRESS_TOPIC from environment
	kafkaIngressTopic := os.Getenv("KAFKA_INGRESS_TOPIC")
	if kafkaIngressTopic == "" {
		kafkaIngressTopic = "ruuvi-event-ingress"
	}

	// Initialize Kafka producer
	kf := kafkawrapper.NewKafkaProducer(logger.Sugar())

	// Create Kafka messaging channel
	messageChannel := make(chan string, 100)

	go kf.ProduceMessagesFromChan(messageChannel, kafkaIngressTopic)

	http.HandleFunc("/event", internal.CreateIncomingEventHandler(messageChannel))

	// Add k8s health and liveliness check endpoints.
	// TODO: Add more sophisticated checks.
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	fmt.Println("Starting server on port 8088")
	if err := http.ListenAndServe(":8088", nil); err != nil {
		panic(err)
	}
}

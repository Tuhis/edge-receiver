package main

import (
	"fmt"
	"net/http"

	"github.com/Tuhis/edge-receiver/internal"
	"github.com/Tuhis/edge-receiver/pkg/kafkawrapper"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

func main() {

	// Init logger
	logger, _ := zap.NewProduction()

	// Initialize Kafka producer
	kf := kafkawrapper.NewKafkaProducer(logger.Sugar())

	// Create Kafka messaging channel
	messageChannel := make(chan string, 100)

	go kf.ProduceMessagesFromChan(messageChannel, "test")

	http.HandleFunc("/event", internal.CreateIncomingEventHandler(messageChannel))

	fmt.Println("Starting server on port 8088")
	if err := http.ListenAndServe(":8088", nil); err != nil {
		panic(err)
	}
}

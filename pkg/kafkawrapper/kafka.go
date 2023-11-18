package kafkawrapper

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type IWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type KafkaProducer struct {
	w           IWriter
	log         *zap.SugaredLogger
	statusTopic string
	ownName     string
}

type IKafkaProducer interface {
	Close() error
	ProduceMessage(message string, topic string) error
	ProduceMessagesFromChan(messages <-chan string, topic string)
}

type envConfig struct {
	Brokers     string
	StatusTopic string
	OwnName     string
}

func NewKafkaProducer(logger *zap.SugaredLogger) IKafkaProducer {
	logger.Infow("Initializing Kafka producer")

	config, err := configFromEnvironment()
	if err != nil {
		logger.Errorf("Failed to read configuration from environment: %v", err)
		panic(err)
	}

	kafkaWriter := &kafka.Writer{
		Addr:                   kafka.TCP(config.Brokers),
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           kafka.RequireOne,
		BatchTimeout:           1 * time.Second,
		Logger:                 kafka.LoggerFunc(logger.Infof),
		ErrorLogger:            kafka.LoggerFunc(logger.Errorf),
		AllowAutoTopicCreation: false,
		Async:                  true,
		Completion:             completionHandlerWithLogger(logger),
	}

	// Test the writer
	err = kafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Topic: config.StatusTopic,
			Value: []byte(config.OwnName + " is alive"),
		},
	)

	if err != nil {
		logger.Errorf("Failed to write message to Kafka: %v", err)
		panic(err)
	}

	return &KafkaProducer{
		w:           kafkaWriter,
		log:         logger,
		statusTopic: config.StatusTopic,
		ownName:     config.OwnName,
	}
}

func configFromEnvironment() (*envConfig, error) {
	brokers := os.Getenv("KAFKA_BROKERS")
	statusTopic := os.Getenv("KAFKA_STATUS_TOPIC")
	ownName := os.Getenv("OWN_NAME")

	if brokers == "" || statusTopic == "" || ownName == "" {
		return nil, errors.New("KAFKA_BROKERS, KAFKA_STATUS_TOPIC and OWN_NAME must be set")
	}

	return &envConfig{
		Brokers:     brokers,
		StatusTopic: statusTopic,
		OwnName:     ownName,
	}, nil
}

func completionHandlerWithLogger(logger *zap.SugaredLogger) func([]kafka.Message, error) {
	return func(messages []kafka.Message, err error) {
		logger.Info(strconv.Itoa(len(messages)) + " messages produced")
		if err != nil {
			panic(err)
		}
	}
}

func (k *KafkaProducer) Close() error {
	return k.w.Close()
}

func (k *KafkaProducer) ProduceMessage(message string, topic string) error {
	return k.w.WriteMessages(context.Background(),
		kafka.Message{
			Topic: topic,
			Value: []byte(message),
		},
	)
}

func (k *KafkaProducer) ProduceRawMessage(message []byte, topic string) error {
	return k.w.WriteMessages(context.Background(),
		kafka.Message{
			Topic: topic,
			Value: message,
		},
	)
}

func (k *KafkaProducer) ProduceMessagesFromChan(messages <-chan string, topic string) {
	for message := range messages {
		err := k.ProduceMessage(message, topic)
		if err != nil {
			k.log.Errorf("Failed to produce message: %v", err)
		}
	}
}

func (k *KafkaProducer) ProduceMessagesFromRawChan(messages <-chan []byte, topic string) {
	for message := range messages {
		err := k.ProduceRawMessage(message, topic)
		if err != nil {
			k.log.Errorf("Failed to produce raw message: %v", err)
		}
	}
}

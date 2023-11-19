package kafkawrapper

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
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

type AuthMechanism string

const (
	AuthMechanismPlain AuthMechanism = "PLAIN"
	AuthMechanismScram AuthMechanism = "SCRAM-SHA-512"
)

type envConfig struct {
	Brokers       string
	StatusTopic   string
	OwnName       string
	AuthMechanism AuthMechanism
	Username      string
	Password      string
}

func NewKafkaProducer(logger *zap.SugaredLogger) IKafkaProducer {
	logger.Infow("Initializing Kafka producer")

	config, err := configFromEnvironment()
	if err != nil {
		logger.Errorf("Failed to read configuration from environment: %v", err)
		panic(err)
	}

	transport := &kafka.Transport{}

	if config.AuthMechanism == AuthMechanismScram {
		mechanism, err := scram.Mechanism(scram.SHA512, config.Username, config.Password)
		if err != nil {
			panic(err)
		}

		transport.SASL = mechanism
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
		Transport:              transport, // Use the custom transport
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
	authMechanismStr := os.Getenv("KAFKA_AUTH_MECHANISM")
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")

	if brokers == "" || statusTopic == "" || ownName == "" {
		return nil, errors.New("KAFKA_BROKERS, KAFKA_STATUS_TOPIC and OWN_NAME must be set")
	}

	authMechanism, err := parseAuthMechanism(authMechanismStr)
	if err != nil {
		return nil, err
	}

	if authMechanism == AuthMechanismScram {
		if username == "" || password == "" {
			return nil, errors.New("KAFKA_USERNAME and KAFKA_PASSWORD must be set for SCRAM authentication")
		}
	}

	return &envConfig{
		Brokers:       brokers,
		StatusTopic:   statusTopic,
		OwnName:       ownName,
		AuthMechanism: authMechanism,
		Username:      username,
		Password:      password,
	}, nil
}

func parseAuthMechanism(s string) (AuthMechanism, error) {
	switch s {
	case string(AuthMechanismPlain), "":
		return AuthMechanismPlain, nil
	case string(AuthMechanismScram):
		return AuthMechanismScram, nil
	default:
		return "", errors.New("invalid auth mechanism: " + s)
	}
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

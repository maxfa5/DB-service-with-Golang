package consumer

import (
	"fmt"
	"kafka_with_go/internal/config"
	"log/slog"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ConsumerService struct {
	logger   *slog.Logger
	cfg      config.Consumer
	consumer *kafka.Consumer
	stop     chan bool
}

func NewConsumerService(logger *slog.Logger, cfg config.Consumer) (*ConsumerService, error) {

	// Проверка подключения
	consumer, err := connToKafka(cfg)
	if err != nil {
		logger.Error("Ошибка подключения к Kafka", slog.String("error", err.Error()))
	} else {
		defer consumer.Close()
		logger.Info("Consumer успешно подключен к Kafka. Ожидание сообщений...")
	}

	return &ConsumerService{
		logger:   logger,
		cfg:      cfg,
		consumer: consumer,
		stop:     make(chan bool),
	}, err

}

func connToKafka(cfg config.Consumer) (*kafka.Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": cfg.Brokers,
		"group.id":          cfg.GroupId,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, fmt.Errorf("error creating kafka consumer: %w", err)
	}

	err = consumer.SubscribeTopics([]string{cfg.Topic}, nil)
	if err != nil {
		consumer.Close()
		return nil, fmt.Errorf("error subscribing to topic: %w", err)
	}

	return consumer, nil
}

func loopGetMsg(c *ConsumerService) (msg *kafka.Message, err error) {
	const op = "consumer.loopGetMsg"
	c.logger = slog.With(
		slog.String("op", op),
	)

	for {
		select {
		case <-c.stop:
			c.logger.Warn("consumer stopped")
			return
		default:
			msg, err := c.consumer.ReadMessage(time.Second)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() == kafka.ErrTimedOut {
					c.logger.Debug("No message received yet, continuing loop", slog.String("error", err.Error()))
					continue
				}
				err = fmt.Errorf("error while reading from kafka: %w", err)
				c.logger.Error("Error while reading from kafka", slog.String("error", err.Error()))
				return nil, err
			}
			c.logger.Info("Message received", slog.String("topic", *msg.TopicPartition.Topic),
				slog.Int("partition", int(msg.TopicPartition.Partition)),
				slog.Any("offset", msg.TopicPartition.Offset),
				slog.String("value", string(msg.Value)))
			return msg, nil // TODO добавть добавление в db
		}

	}
}

func connToPosgre() {}

func SendMessToDB() {}

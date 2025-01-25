package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// Укажите адреса брокеров
	brokers := "localhost:9092,localhost:9093,localhost:9091"
	topic := "orders"
	groupId := "test_group"
	// Проверка подключения
	consumer, err := connToKafka(brokers, topic, groupId)
	defer consumer.Close()
	if err != nil {
		logger.Error("Ошибка подключения к Kafka", slog.String("error", err.Error()))
	} else {
		logger.Info("Consumer успешно подключен к Kafka. Ожидание сообщений...")
	}
	msg, err := loopCathMsg(logger, consumer, 10000*time.Millisecond)
	if err == nil {
		logger.Info("Получено сообщение:", slog.Any("", string(msg.Value)))
	}

	logger.Info("Программа заверена")
}

func connToKafka(brokers string, topic string, groupId string) (c *kafka.Consumer, err error) {
	c, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	// Попробуйте подписаться на тему
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		c.Close()
		return nil, err
	}
	return c, nil
}

func loopCathMsg(logger *slog.Logger, consumer *kafka.Consumer, timeout time.Duration) (msg *kafka.Message, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			err := fmt.Errorf("timeout while waiting for message: %w", ctx.Err())
			logger.Warn("Timeout while waiting for message", slog.String("error", err.Error()))
			return nil, err
		default:
			msg, err := consumer.ReadMessage(100 * time.Millisecond) // Non-blocking call with a timeout
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() == kafka.ErrTimedOut {
					logger.Debug("No message received yet, continuing loop", slog.String("error", err.Error()))
					continue
				}
				err = fmt.Errorf("error while reading from kafka: %w", err)
				logger.Error("Error while reading from kafka", slog.String("error", err.Error())) //Логируем на уровне где это имеет смысл
				return nil, err
			}
			// Если сообщение получено, возвращаем его
			logger.Debug("Message received", slog.Any("key", msg.Key), slog.Any("headers", msg.Headers))
			return msg, nil
		}

	}
}

func connToPosgre() {}

func SendMessToDB() {}

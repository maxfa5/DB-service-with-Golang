package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// Укажите адреса брокеров
	brokers := "localhost:9092,localhost:9093,localhost:9091"
	topic := "orders"
	groupId := "test_group"
	// Проверка подключения
	consumer, err := connToKafka(brokers, topic, groupId)
	if err != nil {
		logger.Error("Ошибка подключения к Kafka: %v", err)
	}
	defer consumer.Close()

	logger.Info("Consumer успешно подключен к Kafka. Ожидание сообщений...")

	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			// Если произошла ошибка, выводим ее
			logger.Error("Ошибка при чтении сообщения: %v\n", err)
			break
		}
		// Если сообщение получено, выводим его
		logger.Info("Получено сообщение: %s\n", string(msg.Value))
	}

	// Закрываем Consumer
	consumer.Close()
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

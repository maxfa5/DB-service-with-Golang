package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	brokers := "localhost:9092,localhost:9093,localhost:9091"

	// Создаем продюсера
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers, // Замените на адрес вашего Kafka брокера
	})
	if err != nil {
		log.Fatalf("Ошибка создания продюсера: %s\n", err)
	}
	defer producer.Close()

	// Отправляем сообщение
	topic := "orders" // Замените на имя вашей темы
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte("Привет, Kafka!"),
	}

	// Отправляем сообщение асинхронно
	err = producer.Produce(message, nil)
	if err != nil {
		log.Fatalf("Ошибка отправки сообщения: %s\n", err)
		return
	}

	// Ждем, пока все сообщения будут отправлены
	producer.Flush(15 * 1000)

	fmt.Println("Сообщение успешно отправлено!")
}

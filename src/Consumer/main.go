package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

func main() {
	// Укажите адреса брокеров
	brokers := "localhost:9092,localhost:9093,localhost:9091"
	topic := "orders"
	group_id := "test_group"
	// Проверка подключения
	consumer, err := connToKafka(brokers, topic, group_id)
	if err != nil {
		log.Fatalf("Ошибка подключения к Kafka: %v", err)
	}
	defer consumer.Close()

	fmt.Println("Consumer успешно подключен к Kafka. Ожидание сообщений...")

	// Получение сообщений (можно просто ждать, чтобы проверить подключение)
	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			// Если произошла ошибка, выводим ее
			fmt.Printf("Ошибка при чтении сообщения: %v\n", err)
			break
		}
		// Если сообщение получено, выводим его
		fmt.Printf("Получено сообщение: %s\n", string(msg.Value))
	}

	// Закрываем Consumer
	consumer.Close()
}

func connToKafka(brokers string, topic string, group_id string) (c *kafka.Consumer, err error) {
	c, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          group_id,
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

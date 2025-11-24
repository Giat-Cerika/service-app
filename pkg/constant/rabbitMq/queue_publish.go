package rabbitmq

import (
	"encoding/json"
	"fmt"
	"giat-cerika-service/configs"
	"log"

	"github.com/streadway/amqp"
)

func PublishToQueue(exchangeName string, queueName string, payload interface{}) error {
	conn := configs.GetRabbitConn() // NEW
	if conn == nil {
		return fmt.Errorf("failed to get active RabbitMQ connection")
	}

	ch, err := conn.Channel() // NEW → channel per publish
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = ch.Publish(
		exchangeName,
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("📨 Published message to queue %s", queueName)
	return nil
}

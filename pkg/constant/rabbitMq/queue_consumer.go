package rabbitmq

import (
	"giat-cerika-service/configs"
	"log"

	"github.com/streadway/amqp"
)

func ConsumeQueueManual(queueName string, handler func(amqp.Delivery)) error {
	_, err := configs.RabbitChannel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare queue %s: %v", queueName, err)
	}

	msgs, err := configs.RabbitChannel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			handler(msg)
		}
	}()

	log.Printf("🕐 Waiting for messages in queue %s ...", queueName)
	return nil
}

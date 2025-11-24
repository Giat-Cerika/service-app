package rabbitmq

import (
	"giat-cerika-service/configs"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func ConsumeQueueManual(queueName string, handler func(amqp.Delivery)) error {
	go func() {
		for {
			err := startConsumer(queueName, handler)
			if err != nil {
				log.Printf("❌ Consumer stopped for %s: %v. Reconnecting in 5s...", queueName, err)
				time.Sleep(5 * time.Second)
			}
		}
	}()

	return nil
}

func startConsumer(queueName string, handler func(amqp.Delivery)) error {
	conn := configs.GetRabbitConn()
	if conn == nil {
		return amqp.ErrClosed
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	log.Printf("🕐 Waiting for messages in queue %s ...", queueName)

	for msg := range msgs {
		handler(msg)
	}

	return amqp.ErrClosed
}

package configs

import (
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel

func InitRabbitMQ() {
	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672/"
	}

	var err error
	for i := 0; i < 5; i++ {
		RabbitConn, err = amqp.Dial(amqpURL)
		if err == nil {
			RabbitChannel, err = RabbitConn.Channel()
			if err == nil {
				log.Println("RabbitMQ connected successfully")
				return
			}
		}
		log.Printf("Failed to connect to RabbitMQ, retrying in %d seconds: %v", i+1, err)
		time.Sleep(time.Second * time.Duration(i+1))
	}
	log.Fatalf("Fatal: Could not connect to RabbitMQ after multiple retries")
}

func GetRabbitChannel() *amqp.Channel {
	if RabbitConn == nil || RabbitConn.IsClosed() {
		log.Println("RabbitMQ connection is closed. Reconnecting...")
		InitRabbitMQ()
	}

	if RabbitChannel == nil {
		log.Println("RabbitMQ channel is not available. Opening a new channel...")
		var err error
		RabbitChannel, err = RabbitConn.Channel()
		if err != nil {
			log.Printf("Failed to open a new channel: %v", err)
			return nil
		}
	}
	return RabbitChannel
}

func CloseConnections() {
	if RabbitChannel != nil {
		_ = RabbitChannel.Close()
	}
	if RabbitConn != nil {
		_ = RabbitConn.Close()
	}
}

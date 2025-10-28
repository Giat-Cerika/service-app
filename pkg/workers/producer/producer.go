package producer

import (
	rabbitmq "giat-cerika-service/pkg/constant/rabbitMq"
	"giat-cerika-service/pkg/workers/consumer"
	handlerconsumer "giat-cerika-service/pkg/workers/handler_consumer"
	"giat-cerika-service/pkg/workers/payload"
)

func StartWorker() {
	go consumer.StartImageConsumer(rabbitmq.SendImageProfileStudentQueueName, &handlerconsumer.StudentImageHandler{}, func() any { return &payload.ImageUploadPayload{} })
	select {}
}

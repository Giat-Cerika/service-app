package producer

import (
	rabbitmq "giat-cerika-service/pkg/constant/rabbitMq"
	"giat-cerika-service/pkg/workers/consumer"
	handlerconsumer "giat-cerika-service/pkg/workers/handler_consumer"
	"giat-cerika-service/pkg/workers/payload"
)

func StartWorker() {
	materiHandler := handlerconsumer.NewMateriHandler()
	go consumer.StartImageConsumer(rabbitmq.SendImageProfileStudentQueueName, &handlerconsumer.StudentImageHandler{}, func() any { return &payload.ImageUploadPayload{} })
	go consumer.StartImageConsumer(rabbitmq.SendImageProfileAdminQueueName, &handlerconsumer.AdminPhotoHandler{}, func() any { return &payload.ImageUploadPayload{} })
	go consumer.StartImageConsumer(
		rabbitmq.SendImageMateriQueueName,
		materiHandler,
		func() any { return &payload.ImageUploadPayload{} },
	)
	select {}
}

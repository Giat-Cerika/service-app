package producer

import (
	rabbitmq "giat-cerika-service/pkg/constant/rabbitMq"
	"giat-cerika-service/pkg/workers/consumer"
	handlerconsumer "giat-cerika-service/pkg/workers/handler_consumer"
	"giat-cerika-service/pkg/workers/payload"
)

func StartWorker() {
	materiHandler := handlerconsumer.NewMateriHandler()
	studentImageHandler := handlerconsumer.NewStudentImageHandler()
	adminPhotoHandler := handlerconsumer.NewAdminPhotoHandler()
	questionHandler := handlerconsumer.NewQuestionHandler()
	go consumer.StartImageConsumer(rabbitmq.SendImageProfileStudentQueueName, studentImageHandler, func() any { return &payload.ImageUploadPayload{} })
	go consumer.StartImageConsumer(rabbitmq.SendImageProfileAdminQueueName, adminPhotoHandler, func() any { return &payload.ImageUploadPayload{} })
	go consumer.StartImageConsumer(
		rabbitmq.SendImageMateriQueueName,
		materiHandler,
		func() any { return &payload.ImageUploadPayload{} },
	)
	go consumer.StartImageConsumer(
		rabbitmq.SendImageQuestionQueueName,
		questionHandler,
		func() any { return &payload.ImageUploadPayload{} },
	)
	select {}
}

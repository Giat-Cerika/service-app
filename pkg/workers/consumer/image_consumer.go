package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	datasources "giat-cerika-service/internal/dataSources"
	"giat-cerika-service/internal/models"
	rabbitmq "giat-cerika-service/pkg/constant/rabbitMq"
	"log"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// Interface handler untuk masing-masing entity (shop, event, dll)
type ImageHandler interface {
	HandleSingle(ctx context.Context, imageURL string, payload any) error
	HandleMany(ctx context.Context, image *models.Image, payload any) error
}

// Interface umum untuk payload upload
type uploadable interface {
	GetFileBytes() []byte
	GetFolder() string
	GetFilename() string
	GetType() string
}

// Consumer utama
func StartImageConsumer(queueName string, handler ImageHandler, payloadFactory func() any) {
	err := rabbitmq.ConsumeQueueManual(queueName, func(msg amqp.Delivery) {
		ctx := context.Background()
		logPrefix := "[image-consumer]"

		// Helper untuk NACK
		nack := func(requeue bool, reason string, err error) {
			if err != nil {
				log.Printf("%s %s: %v", logPrefix, reason, err)
			} else {
				log.Printf("%s %s", logPrefix, reason)
			}
			msg.Nack(false, requeue)
		}

		// 1️⃣ Unmarshal payload
		payload := payloadFactory()
		if err := json.Unmarshal(msg.Body, payload); err != nil {
			nack(false, "failed to unmarshal payload", err)
			return
		}

		// 2️⃣ Pastikan implement interface uploadable
		u, ok := payload.(uploadable)
		if !ok {
			nack(false, "invalid payload type (not uploadable)", nil)
			return
		}

		// 3️⃣ Upload ke Cloudinary
		cld, err := datasources.NewCloudinaryService()
		if err != nil {
			nack(true, "cloudinary init error", err)
			return
		}

		upload, err := cld.UploadImageBytes(ctx,
			bytes.NewReader(u.GetFileBytes()), u.GetFolder(), u.GetFilename())
		if err != nil {
			nack(true, "cloudinary upload error", err)
			return
		}

		// 4️⃣ Proses berdasarkan type
		switch u.GetType() {
		case "single":
			// Tidak buat ImageEntity, langsung simpan URL
			if err := handler.HandleSingle(ctx, upload.URL, payload); err != nil {
				nack(true, "handler single error", err)
				return
			}

		case "many":
			// Buat entity image karena digunakan di relasi many-to-many
			img := &models.Image{
				ID:        uuid.New(),
				ImagePath: upload.URL,
			}
			if err := handler.HandleMany(ctx, img, payload); err != nil {
				nack(true, "handler many error", err)
				return
			}

		default:
			nack(false, "unknown upload type: "+u.GetType(), nil)
			return
		}

		// 5️⃣ Sukses
		msg.Ack(false)
		log.Printf("%s processed message successfully from queue %s", logPrefix, queueName)
	})

	if err != nil {
		log.Fatalf("Failed to start consumer for %s: %v", queueName, err)
	}

	log.Printf("✅ Consumer for %s started. Waiting for messages...", queueName)
	select {}
}

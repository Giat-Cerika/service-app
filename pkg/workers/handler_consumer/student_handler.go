package handlerconsumer

import (
	"context"
	"giat-cerika-service/internal/models"
	"giat-cerika-service/pkg/workers/payload"
)

type StudentImageHandler struct{}

func (h *StudentImageHandler) HandleSingle(ctx context.Context, imageURL string, payloads any) error {
	_, ok := payloads.(*payload.ImageUploadPayload)
	if !ok {
		return nil
	}

	// repo := repositories.NewShopRepositoryImpl(configs.DB)
	// return repo.UpdateShopCover(ctx, p.ID, imageURL)
	return nil
}

func (h *StudentImageHandler) HandleMany(ctx context.Context, image *models.Image, payloads any) error {
	_, ok := payloads.(*payload.ImageUploadPayload)
	if !ok {
		return nil
	}

	// repo := repositories.NewShopRepositoryImpl(configs.DB)
	// if err := repo.CreateImage(ctx, image); err != nil {
	// 	return err
	// }

	// return repo.CreateGallery(ctx, p.ID, image.ID, p.Filename)
	return nil
}

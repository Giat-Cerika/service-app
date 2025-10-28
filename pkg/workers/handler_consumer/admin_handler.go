package handlerconsumer

import (
	"context"
	"giat-cerika-service/configs"
	"giat-cerika-service/internal/models"
	adminrepo "giat-cerika-service/internal/repositories/admin_repo"
	"giat-cerika-service/pkg/workers/payload"
)

type AdminPhotoHandler struct{}

func (h *AdminPhotoHandler) HandleSingle(ctx context.Context, photoUrl string, payloads any) error {
	p, ok := payloads.(*payload.ImageUploadPayload)
	if !ok {
		return nil
	}

	repo := adminrepo.NewAdminRepositoryImpl(configs.DB)
	return repo.UpdatePhotoAdmin(ctx, p.ID, photoUrl)
}

func (h *AdminPhotoHandler) HandleMany(ctx context.Context, image *models.Image, payloads any) error {
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

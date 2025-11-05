package handlerconsumer

import (
	"context"
	"giat-cerika-service/configs"
	"giat-cerika-service/internal/models"
	studentrepo "giat-cerika-service/internal/repositories/student_repo"
	"giat-cerika-service/pkg/workers/payload"
)

type StudentImageHandler struct{}

func (h *StudentImageHandler) HandleSingle(ctx context.Context, imageURL string, payloads any) error {
	p, ok := payloads.(*payload.ImageUploadPayload)
	if !ok {
		return nil
	}

	repo := studentrepo.NewStudentRepositoryImpl(configs.DB)
	return repo.UpdatePhotoStudent(ctx, p.ID, imageURL)
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

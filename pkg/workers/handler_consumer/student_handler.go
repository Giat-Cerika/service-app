package handlerconsumer

import (
	"context"
	"fmt"
	"giat-cerika-service/configs"
	datasources "giat-cerika-service/internal/dataSources"
	"giat-cerika-service/internal/models"
	studentrepo "giat-cerika-service/internal/repositories/student_repo"
	"giat-cerika-service/pkg/utils"
	"giat-cerika-service/pkg/workers/payload"
)

type StudentImageHandler struct {
	repo studentrepo.IStudentRepository
	cld  datasources.CloudinaryService
}

func NewStudentImageHandler() *StudentImageHandler {
	cloudSv, err := datasources.NewCloudinaryService()
	if err != nil {
		panic("Failed to initialize Cloudinary service: " + err.Error())
	}
	return &StudentImageHandler{
		repo: studentrepo.NewStudentRepositoryImpl(configs.DB),
		cld:  cloudSv,
	}
}

func (h *StudentImageHandler) HandleSingle(ctx context.Context, imageURL string, payloads any) error {
	p, ok := payloads.(*payload.ImageUploadPayload)
	if !ok {
		return nil
	}

	newURL := imageURL

	if p.OldPhotoURL != "" {
		publicID := utils.ExtractPublicIDFromCloudinaryURL(p.OldPhotoURL)
		if publicID != "" {
			// aman → tidak ganggu proses meski error
			if err := h.cld.DestroyImage(ctx, publicID); err != nil {
				// bisa log error saja
				fmt.Println("failed delete old photo:", err)
			}
		}
	}

	return h.repo.UpdatePhotoStudent(ctx, p.ID, newURL)
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

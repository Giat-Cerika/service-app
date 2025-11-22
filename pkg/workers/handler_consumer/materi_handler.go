package handlerconsumer

import (
	"context"
	"giat-cerika-service/configs"
	"giat-cerika-service/internal/models"
	materialrepo "giat-cerika-service/internal/repositories/material_repo"
	"giat-cerika-service/pkg/workers/payload"

	"github.com/redis/go-redis/v9"
)

type MateriHandler struct {
	repo materialrepo.IMaterialRepository
	rdb  *redis.Client
}

func NewMateriHandler() *MateriHandler {
	return &MateriHandler{
		repo: materialrepo.NewMaterialRepositoryImpl(configs.DB),
		rdb:  configs.RDB,
	}
}
func (h *MateriHandler) HandleSingle(ctx context.Context, photoUrl string, payloads any) error {
	p, ok := payloads.(*payload.ImageUploadPayload)
	if !ok {
		return nil
	}

	err := h.repo.UpdateCoverMateri(ctx, p.ID, photoUrl)
	if err != nil {
		return err
	}

	h.deleteCache(ctx, p.ID.String())

	return nil
}

func (h *MateriHandler) HandleMany(ctx context.Context, image *models.Image, payloads any) error {
	p, ok := payloads.(*payload.ImageUploadPayload)
	if !ok {
		return nil
	}

	if err := h.repo.CreateImage(ctx, image); err != nil {
		return err
	}

	if err := h.repo.CreateGallery(ctx, p.ID, image.ID, p.Filename); err != nil {
		return err
	}

	h.deleteCache(ctx, p.ID.String())

	return nil
}

func (h *MateriHandler) deleteCache(ctx context.Context, materiID string) {
	h.rdb.Del(ctx, "material:"+materiID)

	iter := h.rdb.Scan(ctx, 0, "materiales:*", 0).Iterator()
	for iter.Next(ctx) {
		h.rdb.Del(ctx, iter.Val())
	}
}

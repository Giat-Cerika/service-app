package videoservice

import (
	"context"
	videorequest "giat-cerika-service/internal/dto/request/video_request"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IVideoService interface {
	CreateVideo(ctx context.Context, req videorequest.CreateVideoRequest, creatorID uuid.UUID) error
	GetAllVideo(ctx context.Context, page, limit int, search string) ([]*models.Video, int, error)
	GetByIdVideo(ctx context.Context, videoId uuid.UUID) (*models.Video, error)
	UpdateVideo(ctx context.Context, videoId uuid.UUID, req videorequest.UpdateVideoRequest) error
	DeleteVideo(ctx context.Context, videoId uuid.UUID) error
	GetAllLatestVideo(ctx context.Context) ([]*models.Video, error)
	GetAllPublicVideo(ctx context.Context, page, limit int, search string) ([]*models.Video, int, error)
	GetByIdPublicVideo(ctx context.Context, videoId uuid.UUID) (*models.Video, error)
}

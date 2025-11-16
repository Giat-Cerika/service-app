package videoservice

import (
	"context"
	videorequest "giat-cerika-service/internal/dto/request/video_request"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IVideoService interface {
	CreateVideo(ctx context.Context, req videorequest.CreateVideoRequest) error
	GetAllVideo(ctx context.Context, page, limit int, search string) ([]*models.Video, int, error)
	GetByIdVideo(ctx context.Context, videoId uuid.UUID) (*models.Video, error)
	UpdateVideo(ctx context.Context, videoId uuid.UUID, req videorequest.UpdateVideoRequest) error
	DeleteVideo(ctx context.Context, videoId uuid.UUID) error
}

package videoresponse

import (
	"giat-cerika-service/internal/models"
	"giat-cerika-service/pkg/utils"

	"github.com/google/uuid"
)

type VideoResponse struct {
	ID          uuid.UUID `json:"id"`
	VideoPath   string    `json:"video_path"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

func ToVideoResponse(video models.Video) VideoResponse {
	return VideoResponse{
		ID:          video.ID,
		VideoPath:   video.VideoPath,
		Title:       video.Title,
		Description: video.Description,
		CreatedAt:   utils.FormatDate(video.CreatedAt),
		UpdatedAt:   utils.FormatDate(video.UpdatedAt),
	}
}

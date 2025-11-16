package materialresponse

import (
	"giat-cerika-service/internal/models"
	"giat-cerika-service/pkg/utils"

	"github.com/google/uuid"
)

type MaterialResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

func ToMaterialResponse(class models.Materials) MaterialResponse {
	return MaterialResponse{
		ID:          class.ID,
		Title:       class.Title,
		Description: class.Description,
		CreatedAt:   utils.FormatDate(class.CreatedAt),
		UpdatedAt:   utils.FormatDate(class.UpdatedAt),
	}
}

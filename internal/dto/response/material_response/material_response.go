package materialresponse

import (
	"giat-cerika-service/internal/models"
	"giat-cerika-service/pkg/utils"

	"github.com/google/uuid"
)

type MaterialResponse struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Cover          string    `json:"cover"`
	MaterialImages []string  `json:"material_images"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
}

func ToMaterialResponse(material models.Materials) MaterialResponse {
	materialImages := []string{}
	for _, materialImage := range material.MaterialImages {
		materialImages = append(materialImages, materialImage.Image.ImagePath)
	}
	return MaterialResponse{
		ID:             material.ID,
		Title:          material.Title,
		Description:    material.Description,
		Cover:          material.Cover,
		MaterialImages: materialImages,
		CreatedAt:      utils.FormatDate(material.CreatedAt),
		UpdatedAt:      utils.FormatDate(material.UpdatedAt),
	}
}

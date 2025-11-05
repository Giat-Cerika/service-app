package classresponse

import (
	"giat-cerika-service/internal/models"
	"giat-cerika-service/pkg/utils"

	"github.com/google/uuid"
)

type ClassResponse struct {
	ID        uuid.UUID `json:"id"`
	NameClass string    `json:"name_class"`
	Grade     string    `json:"grade"`
	Teacher   string    `json:"teacher"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func ToClassResponse(class models.Class) ClassResponse {
	return ClassResponse{
		ID:        class.ID,
		NameClass: class.NameClass,
		Grade:     class.Grade,
		Teacher:   class.Teacher,
		CreatedAt: utils.FormatDate(class.CreatedAt),
		UpdatedAt: utils.FormatDate(class.UpdatedAt),
	}
}

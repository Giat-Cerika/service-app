package toothbrushresponse

import (
	"giat-cerika-service/internal/models"
	"giat-cerika-service/pkg/utils"

	"github.com/google/uuid"
)

type ToothBrushResponse struct {
	ID        uuid.UUID `json:"id"`
	User      string    `json:"user"`
	TimeType  string    `json:"time_type"`
	LogDate   string    `json:"log_date"`
	LogTime   string    `json:"log_time"`
	CreatedAt string    `json:"created_at"`
}

func ToToothBrushResponse(toothBrush models.ToootBrushLog) ToothBrushResponse {
	return ToothBrushResponse{
		ID:        toothBrush.ID,
		User:      *toothBrush.User.Name,
		TimeType:  toothBrush.TimeType,
		LogDate:   utils.FormatLogDate(toothBrush.LogDate),
		LogTime:   utils.FormatTime(toothBrush.LogTime),
		CreatedAt: utils.FormatDate(toothBrush.CreatedAt),
	}
}

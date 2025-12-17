package questionresponse

import (
	"giat-cerika-service/internal/models"
	"giat-cerika-service/pkg/utils"

	"github.com/google/uuid"
)

type QuestionResponse struct {
	ID            uuid.UUID `json:"id"`
	Quiz          string    `json:"quiz"`
	QuestionText  string    `json:"question_text"`
	QuestionImage string    `json:"question_image"`
	Answers       []any     `json:"answers"`
	CratedAt      string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}

func ToQuestionResponse(question models.Question) QuestionResponse {
	ans := []any{}
	for _, answer := range question.Answers {
		ans = append(ans, map[string]any{
			"answer_id":   answer.ID,
			"answer_text": answer.AnswerText,
			"score_value": answer.ScoreValue,
		})
	}
	return QuestionResponse{
		ID:            question.ID,
		Quiz:          question.Quiz.Title,
		QuestionText:  question.QuestionText,
		QuestionImage: question.QuestionImage,
		Answers:       ans,
		CratedAt:      utils.FormatDate(question.CreatedAt),
		UpdatedAt:     utils.FormatDate(question.UpdatedAt),
	}
}

package quizhistoryresponse

import (
	"giat-cerika-service/internal/models"
	"giat-cerika-service/pkg/utils"

	"github.com/google/uuid"
)

type QuizHistoryResponse struct {
	ID              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	StartDate       string    `json:"start_date"`
	EndDate         string    `json:"end_date"`
	AmountQuestions int       `json:"amount_questions"`
	AmountAssigned  int       `json:"amount_assigned"`
	Score           int       `json:"score"`
	MaxScore        int       `json:"max_score"`
	Percentage      float64   `json:"percentage"`
	Status          string    `json:"status"`
	StartedAt       string    `json:"started_at"`
	CompletedAt     string    `json:"completed_at"`
	StatusCategory  int       `json:"status_category"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
}

func ToQuizHistoryResponse(qh models.QuizHistory) QuizHistoryResponse {
	return QuizHistoryResponse{
		ID:              qh.ID,
		Title:           qh.Title,
		Description:     qh.Description,
		StartDate:       utils.FormatDateTime(qh.StartDate),
		EndDate:         utils.FormatDateTime(qh.EndDate),
		AmountQuestions: qh.AmountQuestions,
		AmountAssigned:  qh.AmountAssigned,
		Score:           qh.Score,
		MaxScore:        qh.MaxScore,
		Percentage:      qh.Percentage,
		Status:          string(qh.Status),
		StartedAt:       utils.FormatDateTime(qh.StartedAt),
		CompletedAt:     utils.FormatDateTime(qh.CompletedAt),
		StatusCategory:  qh.StatusCategory,
		CreatedAt:       utils.FormatDate(qh.CreatedAt),
		UpdatedAt:       utils.FormatDate(qh.UpdatedAt),
	}
}

type QuestionHistory struct {
	ID              uuid.UUID `json:"id"`
	QuizHistoryID   uuid.UUID `json:"quiz_history_id"`
	QuestionID      uuid.UUID `json:"question_id"`
	QuestionText    string    `json:"question_text"`
	QuestionImage   string    `json:"question_image"`
	AnswerHistories any       `json:"answer_histories"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
}

func ToQuestionHistory(qh models.QuestionHistory) QuestionHistory {
	ansHistories := []any{}
	for _, ans := range qh.AnswerHistory {
		ansHistories = append(ansHistories, map[string]any{
			"answer_id":    ans.AnswerID,
			"answer_text":  ans.AnswerText,
			"score_value":  ans.ScoreValue,
			"score_earned": ans.ScoreEarned,
		})
	}
	return QuestionHistory{
		ID:              qh.ID,
		QuizHistoryID:   qh.QuizHistoryID,
		QuestionID:      qh.QuestionID,
		QuestionText:    qh.QuestionText,
		QuestionImage:   qh.QuestionImage,
		AnswerHistories: ansHistories,
		CreatedAt:       utils.FormatDate(qh.CreatedAt),
		UpdatedAt:       utils.FormatDate(qh.UpdatedAt),
	}
}

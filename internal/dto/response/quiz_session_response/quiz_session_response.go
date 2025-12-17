package quizsessionresponse

import (
	"giat-cerika-service/internal/models"
	"giat-cerika-service/pkg/utils"
	"time"

	"github.com/google/uuid"
)

type QuizSessionStartResponse struct {
	QuizSessionID   uuid.UUID  `json:"quiz_session_id"`
	QuizID          uuid.UUID  `json:"quiz_id"`
	DurationSeconds int64      `json:"duration_seconds"`
	IsUnlimited     bool       `json:"is_unlimited"`
	EndTime         *time.Time `json:"end_time"`
}

type QuizSessionDurationResponse struct {
	RemainingSeconds int64 `json:"remaining_seconds"` // Sisa waktu dalam detik
	IsExpired        bool  `json:"is_expired"`        // True jika waktu habis
	IsUnlimited      bool  `json:"is_unlimited"`      // True jika tidak ada batas waktu
}

type AnswerResponse struct {
	ID         uuid.UUID `json:"id"`
	AnswerText string    `json:"answer_text"`
	ScoreValue int       `json:"score_value"`
}

type QuestionDetailResponse struct {
	ID            uuid.UUID        `json:"id"`
	QuestionText  string           `json:"question_text"`
	QuestionImage *string          `json:"question_image"`
	Answers       []AnswerResponse `json:"answers"`
}

type OrderedQuizQuestionsResponse struct {
	QuizID    uuid.UUID                `json:"quiz_id"`
	Questions []QuestionDetailResponse `json:"questions"`
}

type DetailQuizSession struct {
	ID          uuid.UUID                `json:"id"`
	Student     string                   `json:"student"`
	Score       int                      `json:"score"`
	MaxScore    int                      `json:"max_score"`
	Status      models.QuizSessionStatus `json:"status"`
	StartedAt   string                   `json:"started_at"`
	CompletedAt string                   `json:"completed_at"`
}

type ListQuestionSessionResponse struct {
	Quiz              string              `json:"quiz"`
	DetailQuizSession []DetailQuizSession `json:"detail_quiz_session"`
	CreatedAt         string              `json:"created_at"`
}

// Helper untuk membuat detail satuan
func ToDetailQuizSession(qs models.QuizSession) DetailQuizSession {
	studentName := ""
	if qs.User.Name != nil {
		studentName = *qs.User.Name
	}

	return DetailQuizSession{
		ID:          qs.ID,
		Student:     studentName,
		Score:       qs.Score,
		MaxScore:    qs.MaxScore,
		Status:      qs.Status,
		StartedAt:   utils.FormatDateTime(qs.StartedAt),
		CompletedAt: utils.FormatDateTime(qs.CompletedAt),
	}
}

// Fungsi ini digunakan saat inisialisasi pertama kali di map
func ToListQuestionSessionResponse(qs models.QuizSession) ListQuestionSessionResponse {
	return ListQuestionSessionResponse{
		Quiz:              qs.Quiz.Title,
		CreatedAt:         utils.FormatDateTime(&qs.Quiz.CreatedAt),
		DetailQuizSession: []DetailQuizSession{ToDetailQuizSession(qs)},
	}
}

package quizrequest

import "github.com/google/uuid"

type SubmitAnswerRequest struct {
	QuestionID uuid.UUID `json:"question_id" binding:"required"`
	AnswerID   uuid.UUID `json:"answer_id" binding:"required"`
}

type SubmitQuizRequest struct {
	Answers []SubmitAnswerRequest `json:"answers" binding:"required"`
}

package questionrequest

import (
	answerrequest "giat-cerika-service/internal/dto/request/answer_request"
	"mime/multipart"

	"github.com/google/uuid"
)

type CreateQuestionRequest struct {
	QuizId        uuid.UUID                           `json:"quiz_id" binding:"required,uuid"`
	QuestionText  string                              `json:"question_text" binding:"required"`
	QuestionImage *multipart.FileHeader               `form:"question_image" swaggerignore:"true"`
	Answers       []answerrequest.CreateAnswerRequest `json:"answers" binding:"required,dive,required"`
}

type UpdateQuestionRequest struct {
	QuizId        uuid.UUID                           `json:"quiz_id" binding:"required,uuid"`
	QuestionText  string                              `json:"question_text" binding:"required"`
	QuestionImage *multipart.FileHeader               `form:"question_image" swaggerignore:"true"`
	Answers       []answerrequest.CreateAnswerRequest `json:"answers" binding:"required,dive,required"`
}

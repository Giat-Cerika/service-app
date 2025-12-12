package answerrequest

type CreateAnswerRequest struct {
	AnswerText string `json:"answer_text" binding:"required"`
	ScoreValue int    `json:"score_value" binding:"required"`
}

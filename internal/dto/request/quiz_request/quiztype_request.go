package quizrequest

type CreateQuizTypeRequest struct {
	Name        string `form:"name" json:"name"`
	Description string `form:"description" json:"description"`
}

type UpdateQuizTypeRequest struct {
	Name        string `form:"name" json:"name"`
	Description string `form:"description" json:"description"`
}

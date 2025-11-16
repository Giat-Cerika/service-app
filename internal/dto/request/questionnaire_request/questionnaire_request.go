package questionnairerequest

type CreateQuestionnaireRequest struct {
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
	Amount      string    `form:"amount" json:"amount"`
	Code        string    `form:"code" json:"code"`
	Status      string    `form:"status" json:"status"`
	Type        string `form:"type" json:"type"`
	Duration    string `form:"duration" json:"duration"`
}

type UpdateQuestionnaireRequest struct {
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
	Amount      string    `form:"amount" json:"amount"`
	Code        string    `form:"code" json:"code"`
	Status      string    `form:"status" json:"status"`
	Type        string `form:"type" json:"type"`
	Duration    string `form:"duration" json:"duration"`
}

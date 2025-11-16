package materialrequest

type CreateMaterialRequest struct {
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
}

type UpdateMaterialRequest struct {
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
}

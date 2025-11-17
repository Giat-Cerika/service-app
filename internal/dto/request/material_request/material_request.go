package materialrequest

import "mime/multipart"

type CreateMaterialRequest struct {
	Title       string                  `form:"title" json:"title"`
	Description string                  `form:"description" json:"description"`
	Cover       *multipart.FileHeader   `form:"cover" swaggerignore:"true"`
	Gallery     []*multipart.FileHeader `form:"gallery" swaggerignore:"true"`
}

type UpdateMaterialRequest struct {
	Title       string                  `form:"title" json:"title"`
	Description string                  `form:"description" json:"description"`
	Cover       *multipart.FileHeader   `form:"cover" swaggerignore:"true"`
	Gallery     []*multipart.FileHeader `form:"gallery" swaggerignore:"true"`
}

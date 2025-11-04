package adminrequest

import "mime/multipart"

type RegisterAdminRequest struct {
	Username string                `form:"username" json:"username"`
	Password string                `form:"password" json:"password"`
	Photo    *multipart.FileHeader `form:"photo" json:"photo"`
}

type LoginAdminRequest struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

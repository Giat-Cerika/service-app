package rolerequest

type CreateRoleRequest struct {
	Name string `form:"name" json:"name"`
}

type UpdateRoleRequest struct {
	Name string `form:"name" json:"name"`
}

package request

type UpdateAdminRequest struct {
    Name  string `json:"name" validate:"omitempty,min=3"`
    Phone string `json:"phone" validate:"omitempty"`
}
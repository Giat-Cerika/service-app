package request

type RegisterStudentRequest struct {
    Name     string `json:"name" validate:"required,min=3"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
    Phone    string `json:"phone" validate:"omitempty"`
    NIM      string `json:"nim" validate:"required"`
    Major    string `json:"major" validate:"omitempty"`
}

type UpdateStudentRequest struct {
    Name  string `json:"name" validate:"omitempty,min=3"`
    Phone string `json:"phone" validate:"omitempty"`
    Major string `json:"major" validate:"omitempty"`
}
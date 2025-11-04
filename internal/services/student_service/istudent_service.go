package studentservice

import (
	"context"
	studentrequest "giat-cerika-service/internal/dto/request/student_request"
)

type IStudentService interface {
	Register(ctx context.Context, req studentrequest.RegisterStudentRequest) error
}

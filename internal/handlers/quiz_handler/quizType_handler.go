package quizhandler

import (
	quizrequest "giat-cerika-service/internal/dto/request/quiz_request"
	quizresponse "giat-cerika-service/internal/dto/response/quiz_response"
	quizservice "giat-cerika-service/internal/services/quiz_service"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"giat-cerika-service/pkg/constant/response"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type QuizTypeHandler struct {
	quizService quizservice.IQuizTypeService
}

func NewQuizTypeHandler(quizService quizservice.IQuizTypeService) *QuizTypeHandler {
	return &QuizTypeHandler{quizService: quizService}
}

func (q *QuizTypeHandler) CreateQuizType(c echo.Context) error {
	var req quizrequest.CreateQuizTypeRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	err := q.quizService.CreateQt(c.Request().Context(), req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to create quiz type")
	}

	return response.Success(c, http.StatusOK, "Quiz Type Created Succssfully", nil)
}

func (q *QuizTypeHandler) GetAllQuizType(c echo.Context) error {
	quizTypes, err := q.quizService.GetAllQt(c.Request().Context())
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get quiz types")
	}

	data := make([]quizresponse.QuizTypeResponse, len(quizTypes))
	for i, qt := range quizTypes {
		data[i] = quizresponse.ToQuizTypeResponse(*qt)
	}
	return response.Success(c, http.StatusOK, "Get All Quiz Types Successfully", quizTypes)
}

func (q *QuizTypeHandler) GetQuizTypeByID(c echo.Context) error {
	quizTypeid, err := uuid.Parse(c.Param("quizTypeid"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", "invalid quiz type id")
	}
	quizType, err := q.quizService.GetByIdQt(c.Request().Context(), quizTypeid)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get quiz type by id")
	}

	res := quizresponse.ToQuizTypeResponse(*quizType)
	return response.Success(c, http.StatusOK, "Get Quiz Type By ID Successfully", res)
}

func (q *QuizTypeHandler) UpdateQuizType(c echo.Context) error {
	quizTypeid, err := uuid.Parse(c.Param("quizTypeid"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", "invalid quiz type id")
	}
	var req quizrequest.UpdateQuizTypeRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}
	err = q.quizService.UpdateQt(c.Request().Context(), quizTypeid, req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to update quiz type")
	}
	return response.Success(c, http.StatusOK, "Quiz Type Updated Successfully", nil)
}

func (q *QuizTypeHandler) DeleteQuizType(c echo.Context) error {
	quizTypeid, err := uuid.Parse(c.Param("quizTypeid"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", "invalid quiz type id")
	}
	err = q.quizService.DeleteQt(c.Request().Context(), quizTypeid)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to delete quiz type")
	}
	return response.Success(c, http.StatusOK, "Quiz Type Deleted Successfully", nil)
}

package quizhandler

import (
	quizrequest "giat-cerika-service/internal/dto/request/quiz_request"
	quizresponse "giat-cerika-service/internal/dto/response/quiz_response"
	quizservice "giat-cerika-service/internal/services/quiz_service"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"giat-cerika-service/pkg/constant/response"
	"giat-cerika-service/pkg/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type QuizHandler struct {
	quizService quizservice.IQuizService
}

func NewQuizHandler(service quizservice.IQuizService) *QuizHandler {
	return &QuizHandler{quizService: service}
}

func (q *QuizHandler) CreateQuiz(c echo.Context) error {
	var req quizrequest.CreateQuizRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}
	err := q.quizService.CreateQuiz(c.Request().Context(), req)
	if err != nil {
		if cutomErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, cutomErr.Status, cutomErr.Msg, cutomErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, "failed to create quiz", err.Error())
	}
	return response.Success(c, http.StatusCreated, "quiz created successfully", nil)
}

func (q *QuizHandler) GetQuizAll(c echo.Context) error {
	pageInt, limitInt := utils.ParsePaginationParams(c, 10)
	search := c.QueryParam("search")

	quizzes, total, err := q.quizService.GetAllQuiz(c.Request().Context(), pageInt, limitInt, search)
	if err != nil {
		if cutomErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, cutomErr.Status, cutomErr.Msg, cutomErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, "failed to get quizzes", err.Error())
	}
	meta := utils.BuildPaginationMeta(c, pageInt, limitInt, total)
	data := make([]quizresponse.QuizResponse, len(quizzes))
	for i, q := range quizzes {
		data[i] = quizresponse.ToQuizResponse(*q)
	}

	return response.PaginatedSuccess(c, http.StatusOK, "Get All Quizzes Successfully", data, meta)
}

func (q *QuizHandler) GetQuizByID(c echo.Context) error {
	quizId, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}
	quiz, err := q.quizService.GetQuizById(c.Request().Context(), quizId)
	if err != nil {
		if cutomErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, cutomErr.Status, cutomErr.Msg, cutomErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, "failed to get quiz", err.Error())
	}
	data := quizresponse.ToQuizResponse(*quiz)
	return response.Success(c, http.StatusOK, "Get Quiz By ID Successfully", data)
}

func (q *QuizHandler) UpdateQuiz(c echo.Context) error {
	quizId, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}
	var req quizrequest.UpdateQuizRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}
	err = q.quizService.UpdateQuiz(c.Request().Context(), quizId, req)
	if err != nil {
		if cutomErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, cutomErr.Status, cutomErr.Msg, cutomErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, "failed to update quiz", err.Error())
	}
	return response.Success(c, http.StatusOK, "Quiz Updated Successfully", nil)
}

func (q *QuizHandler) DeleteQuiz(c echo.Context) error {
	quizId, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}
	err = q.quizService.DeleteQuiz(c.Request().Context(), quizId)
	if err != nil {
		if cutomErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, cutomErr.Status, cutomErr.Msg, cutomErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, "failed to delete quiz", err.Error())
	}
	return response.Success(c, http.StatusOK, "Quiz Deleted Successfully", nil)
}

func (q *QuizHandler) UpdateStatusQuiz(c echo.Context) error {
	quizId, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}
	var req quizrequest.UpdateStatusQuizRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}
	err = q.quizService.UpdateStatusQuiz(c.Request().Context(), quizId, req)
	if err != nil {
		if cutomErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, cutomErr.Status, cutomErr.Msg, cutomErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, "failed to update quiz status", err.Error())
	}

	return response.Success(c, http.StatusOK, "Quiz Status Updated Successfully", nil)
}

func (q *QuizHandler) UpdateQuestionOrderMode(c echo.Context) error {
	quizId, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}
	var req quizrequest.UpdateQuestionOrderModeRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}
	err = q.quizService.UpdateQuestionOrderMode(c.Request().Context(), quizId, req)
	if err != nil {
		if cutomErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, cutomErr.Status, cutomErr.Msg, cutomErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, "failed to update question order mode", err.Error())
	}

	return response.Success(c, http.StatusOK, "Quiz Question Order Mode Updated Successfully", nil)
}

func (q *QuizHandler) GetAllQuizAvailable(c echo.Context) error {
	search := c.QueryParam("search")

	items, err := q.quizService.GetAllQuizAvailable(c.Request().Context(), search)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err)
		}
		return response.Error(c, http.StatusInternalServerError, "failed to get quiz available", 500)
	}

	data := make([]quizresponse.QuizResponse, len(items))
	for i, quiz := range items {
		data[i] = quizresponse.ToQuizResponse(*quiz)
	}

	return response.Success(c, http.StatusOK, "Get Quiz Available Successfully", data)
}

func (q *QuizHandler) GetQuizAvailableById(c echo.Context) error {
	quizId, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	quiz, err := q.quizService.GetQuizAvailableById(c.Request().Context(), quizId)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err)
		}
		return response.Error(c, http.StatusInternalServerError, "failed to get detail quiz", err.Error())
	}

	data := quizresponse.ToQuizResponse(*quiz)

	return response.Success(c, http.StatusOK, "Get Detail Quiz Successfully", data)
}

package quizhistoryhandler

import (
	quizhistoryresponse "giat-cerika-service/internal/dto/response/quiz_history_response"
	quizhistoryservice "giat-cerika-service/internal/services/quiz_history_service"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"giat-cerika-service/pkg/constant/response"
	"giat-cerika-service/pkg/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type QuizHistoryHandler struct {
	qhService quizhistoryservice.IQuizHistoryService
}

func NewQuizHistoryHandler(qhService quizhistoryservice.IQuizHistoryService) *QuizHistoryHandler {
	return &QuizHistoryHandler{qhService: qhService}
}

func (qh *QuizHistoryHandler) GetHistoryQuizStudent(c echo.Context) error {
	search := c.QueryParam("search")

	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: "+err.Error(), nil)
	}
	studentId := claims.UserID

	items, err := qh.qhService.GetHistoryQuizStudent(c.Request().Context(), uuid.MustParse(studentId), search)
	if err != nil {
		if cutomErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, cutomErr.Status, cutomErr.Msg, cutomErr.Err)
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get quiz history student")
	}

	data := make([]quizhistoryresponse.QuizHistoryResponse, len(items))
	for i, qh := range items {
		data[i] = quizhistoryresponse.ToQuizHistoryResponse(*qh)
	}

	return response.Success(c, http.StatusOK, "Get Quiz HIstory Student Successfully", data)
}

func (qh *QuizHistoryHandler) GetAllQuestionHistory(c echo.Context) error {
	quizHistoryId, err := uuid.Parse(c.Param("quizHistoryId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	questionHistory, err := qh.qhService.GetAllHistoryQuestionByQuizHistory(c.Request().Context(), quizHistoryId)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err)
		}
		return response.Error(c, http.StatusInternalServerError, "failed to get question history", err.Error())
	}

	data := make([]quizhistoryresponse.QuestionHistory, len(questionHistory))
	for i, question := range questionHistory {
		data[i] = quizhistoryresponse.ToQuestionHistory(*question)
	}

	return response.Success(c, http.StatusOK, "Get all question history successfully", data)
}

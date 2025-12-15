package quizsessionhandler

import (
	quizsessionservice "giat-cerika-service/internal/services/quiz_session_service"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"giat-cerika-service/pkg/constant/response"
	"giat-cerika-service/pkg/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type QuizSessionHandler struct {
	qsService quizsessionservice.IQuizSessionService
}

func NewQuizSessionHandler(qsService quizsessionservice.IQuizSessionService) *QuizSessionHandler {
	return &QuizSessionHandler{qsService: qsService}
}

func (qs *QuizSessionHandler) AssignCodeQuiz(c echo.Context) error {
	type InputCode struct {
		Code string
	}
	var req InputCode
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: "+err.Error(), nil)
	}
	studentId := claims.UserID

	quizId, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	data, err := qs.qsService.AssignCodeQuiz(c.Request().Context(), uuid.MustParse(studentId), quizId, req.Code)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err)
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get code access quiz")
	}

	return response.Success(c, http.StatusOK, "code access quiz succeefully", data)
}

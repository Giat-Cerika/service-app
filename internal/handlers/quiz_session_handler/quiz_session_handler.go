package quizsessionhandler

import (
	quizrequest "giat-cerika-service/internal/dto/request/quiz_request"
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

	return response.Success(c, http.StatusOK, "code access quiz succeefully", map[string]string{
		"qustion_session_id": data.ID.String(),
		"quiz_id":            data.QuizID.String(),
		"status":             string(data.Status),
	})
}

func (qs *QuizSessionHandler) StartedQuiz(c echo.Context) error {
	quizSessionId, err := uuid.Parse(c.Param("quizSessionId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: "+err.Error(), nil)
	}
	studentId := claims.UserID

	data, err := qs.qsService.StartQuizSession(c.Request().Context(), uuid.MustParse(studentId), quizSessionId)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err)
		}
		return response.Error(c, http.StatusInternalServerError, "failed to start quiz", err.Error())
	}

	return response.Success(c, http.StatusOK, "Start Quiz Successfully", data)
}

func (qs *QuizSessionHandler) GetDuration(c echo.Context) error {
	quizSessionId, err := uuid.Parse(c.Param("quizSessionId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: "+err.Error(), nil)
	}
	studentId := claims.UserID

	duration, err := qs.qsService.GetQuizSessionDuration(c.Request().Context(), uuid.MustParse(studentId), quizSessionId)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err)
		}
		return response.Error(c, http.StatusInternalServerError, "failed to get duration session quiz", err.Error())
	}

	return response.Success(c, http.StatusOK, "Get Duration Success", duration)
}

func (qs QuizSessionHandler) SubmitQuizSession(c echo.Context) error {
	quizSessionId, err := uuid.Parse(c.Param("quizSessionId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: "+err.Error(), nil)
	}
	studentId := claims.UserID

	var req quizrequest.SubmitQuizRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	err = qs.qsService.SubmtiQuizSession(c.Request().Context(), uuid.MustParse(studentId), quizSessionId, req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err)
		}
		return response.Error(c, http.StatusInternalServerError, "failed to submit quiz", 500)
	}

	return response.Success(c, http.StatusOK, "Quiz Submitted Successfully", nil)
}

func (qs *QuizSessionHandler) GetQuizQuestionByOrderMode(c echo.Context) error {
	quizSessionId, err := uuid.Parse(c.Param("quizSessionId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: "+err.Error(), nil)
	}
	studentId := claims.UserID

	data, err := qs.qsService.GetOrderedQuizQuestions(c.Request().Context(), uuid.MustParse(studentId), quizSessionId)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err)
		}
		return response.Error(c, http.StatusInternalServerError, "failed to get question quiz", 500)
	}

	return response.Success(c, http.StatusOK, "Get Question Quiz Successfully", data)
}

func (qs *QuizSessionHandler) GetQuizSessionStudent(c echo.Context) error {
	data, err := qs.qsService.GetQuizSessionStudentByQuiz(c.Request().Context())
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err)
		}
		return response.Error(c, http.StatusInternalServerError, "failed to get quiz session student", 500)
	}

	return response.Success(c, http.StatusOK, "Get Quiz Session Student Successfully", data)
}

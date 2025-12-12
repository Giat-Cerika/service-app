package questionhandler

import (
	"encoding/json"
	answerrequest "giat-cerika-service/internal/dto/request/answer_request"
	questionrequest "giat-cerika-service/internal/dto/request/question_request"
	questionresponse "giat-cerika-service/internal/dto/response/question_response"
	questionservice "giat-cerika-service/internal/services/question_service"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"giat-cerika-service/pkg/constant/response"
	"giat-cerika-service/pkg/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type QuestionHandler struct {
	questionService questionservice.IQuestionService
}

func NewQuestionHandler(questionService questionservice.IQuestionService) *QuestionHandler {
	return &QuestionHandler{
		questionService: questionService,
	}
}

func (qh *QuestionHandler) CreateQuestion(c echo.Context) error {
	var req questionrequest.CreateQuestionRequest

	req.QuestionText = c.FormValue("question_text")

	// Parse quiz_id
	quizIdStr := c.FormValue("quiz_id")
	quizUUID, err := uuid.Parse(quizIdStr)
	if err != nil {
		return response.Error(c, 400, "Invalid quiz_id", err.Error())
	}
	req.QuizId = quizUUID

	// Parse image
	if questionImage, err := c.FormFile("question_image"); err == nil {
		req.QuestionImage = questionImage
	}

	// Parse answers (HARUS DISET KE req.Answers)
	if answers := c.FormValue("answers"); answers != "" {
		var parsed []answerrequest.CreateAnswerRequest
		if err := json.Unmarshal([]byte(answers), &parsed); err != nil {
			return response.Error(c, 400, "Invalid answers format", err.Error())
		}
		req.Answers = parsed // <<< PERBAIKAN PENTING
	}

	// Call service
	err = qh.questionService.CreateQuestion(c.Request().Context(), req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, 500, "Failed to create question", err.Error())
	}

	return response.Success(c, 200, "Question created successfully", nil)
}

func (qh *QuestionHandler) GetAllQuestion(c echo.Context) error {
	pageInt, limitInt := utils.ParsePaginationParams(c, 10)
	search := c.QueryParam("search")
	questions, total, err := qh.questionService.FindAllQuestions(c.Request().Context(), pageInt, limitInt, search)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, 500, "Failed to get questions", err.Error())
	}
	meta := utils.BuildPaginationMeta(c, pageInt, limitInt, total)
	data := make([]questionresponse.QuestionResponse, len(questions))
	for i, question := range questions {
		data[i] = questionresponse.ToQuestionResponse(*question)
	}

	return response.PaginatedSuccess(c, 200, "Get All Questions Successfully", data, meta)
}

func (qh *QuestionHandler) GetByIdQuestion(c echo.Context) error {
	questionId, err := uuid.Parse(c.Param("questionId"))
	if err != nil {
		return response.Error(c, 400, "Bad request", err.Error())
	}
	question, err := qh.questionService.FindQuestionById(c.Request().Context(), questionId)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, 500, "Failed to get question", err.Error())
	}
	res := questionresponse.ToQuestionResponse(*question)

	return response.Success(c, 200, "Get Question Successfully", res)
}

func (qh *QuestionHandler) UpdateQuestion(c echo.Context) error {
	questionId, err := uuid.Parse(c.Param("questionId"))
	if err != nil {
		return response.Error(c, 400, "Bad request", err.Error())
	}
	var req questionrequest.UpdateQuestionRequest

	req.QuestionText = c.FormValue("question_text")
	// Parse quiz_id
	quizIdStr := c.FormValue("quiz_id")
	if quizIdStr != "" {
		quizUUID, err := uuid.Parse(quizIdStr)
		if err != nil {
			return response.Error(c, 400, "Invalid quiz_id", err.Error())
		}
		req.QuizId = quizUUID
	}
	// Parse image
	if questionImage, err := c.FormFile("question_image"); err == nil {
		req.QuestionImage = questionImage
	}
	// Parse answers (HARUS DISET KE req.Answers)
	if answers := c.FormValue("answers"); answers != "" {
		var parsed []answerrequest.CreateAnswerRequest
		if err := json.Unmarshal([]byte(answers), &parsed); err != nil {
			return response.Error(c, 400, "Invalid answers format", err.Error())
		}
		req.Answers = parsed // <<< PERBAIKAN PENTING
	}
	// Call service
	err = qh.questionService.UpdateQuestion(c.Request().Context(), questionId, req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, 500, "Failed to update question", err.Error())
	}
	return response.Success(c, 200, "Question updated successfully", nil)
}

func (qh *QuestionHandler) DeleteQuestion(c echo.Context) error {
	questionId, err := uuid.Parse(c.Param("questionId"))
	if err != nil {
		return response.Error(c, 400, "Bad request", err.Error())
	}
	if err := qh.questionService.DeleteQuestion(c.Request().Context(), questionId); err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, 500, "Failed to delete question", err.Error())
	}
	return response.Success(c, 200, "Question deleted successfully", nil)
}

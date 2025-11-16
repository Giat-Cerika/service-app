package questionnairehandler

import (
	questionnairerequest "giat-cerika-service/internal/dto/request/questionnaire_request"
	questionnaireresponse "giat-cerika-service/internal/dto/response/questionnaire_response"
	questionnaireservice "giat-cerika-service/internal/services/questionnaire_service"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"giat-cerika-service/pkg/constant/response"
	"giat-cerika-service/pkg/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type QuestionnaireHandler struct {
	questionnaireService questionnaireservice.IQuestionnaireService
}

func NewQuestionnaireHandler(service questionnaireservice.IQuestionnaireService) *QuestionnaireHandler {
	return &QuestionnaireHandler{questionnaireService: service}
}

func (ch *QuestionnaireHandler) CreateQuestionnaire(c echo.Context) error {
	var req questionnairerequest.CreateQuestionnaireRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	err := ch.questionnaireService.CreateQuestionnaire(c.Request().Context(), req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to create questionnaire")
	}

	return response.Success(c, http.StatusOK, "Questionnaire Created Succssfully", nil)
}

func (ch *QuestionnaireHandler) GetAllQuestionnaire(c echo.Context) error {
	pageInt, limitInt := utils.ParsePaginationParams(c, 10)
	search := c.QueryParam("search")

	questionnairees, total, err := ch.questionnaireService.GetAllQuestionnaire(c.Request().Context(), pageInt, limitInt, search)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get questionnairees")
	}

	meta := utils.BuildPaginationMeta(c, pageInt, limitInt, total)
	data := make([]questionnaireresponse.QuestionnaireResponse, len(questionnairees))
	for i, questionnaire := range questionnairees {
		data[i] = questionnaireresponse.ToQuestionnaireResponse(*questionnaire)
	}

	return response.PaginatedSuccess(c, http.StatusOK, "Get All Questionnairees Successfully", data, meta)
}

func (ch *QuestionnaireHandler) GetByIdQuestionnaire(c echo.Context) error {
	questionnaireId, err := uuid.Parse(c.Param("questionnaireId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	questionnaire, err := ch.questionnaireService.GetByIdQuestionnaire(c.Request().Context(), questionnaireId)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get questionnaire")
	}

	res := questionnaireresponse.ToQuestionnaireResponse(*questionnaire)

	return response.Success(c, http.StatusOK, "Get Questionnaire Successfully", res)
}

func (ch *QuestionnaireHandler) UpdateQuestionnaire(c echo.Context) error {
	questionnaireId, err := uuid.Parse(c.Param("questionnaireId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	var req questionnairerequest.UpdateQuestionnaireRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	err = ch.questionnaireService.UpdateQuestionnaire(c.Request().Context(), questionnaireId, req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to update questionnaire")
	}

	return response.Success(c, http.StatusOK, "Questionnaire Updated Successfully", nil)
}

func (ch *QuestionnaireHandler) DeleteQuestionnaire(c echo.Context) error {
	questionnaireId, err := uuid.Parse(c.Param("questionnaireId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	if err := ch.questionnaireService.DeleteQuestionnaire(c.Request().Context(), questionnaireId); err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to delete questionnaire")
	}

	return response.Success(c, http.StatusOK, "Questionnaire Deleted Successfully", nil)
}

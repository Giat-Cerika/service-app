package predictionhandler

import (
	predictionrequest "giat-cerika-service/internal/dto/request/prediction_request"
	predictionresponse "giat-cerika-service/internal/dto/response/prediction_response"
	predictionservice "giat-cerika-service/internal/services/prediction_service"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"giat-cerika-service/pkg/constant/response"
	"giat-cerika-service/pkg/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PredictionHandler struct {
	service predictionservice.IPredictionService
}

func NewPredictionHandler(service predictionservice.IPredictionService) *PredictionHandler {
	return &PredictionHandler{service: service}
}

func (ph *PredictionHandler) CreatePrediction(c echo.Context) error {
	var req predictionrequest.CreatePredictionRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	if err := ph.service.Create(c.Request().Context(), req); err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err)
		}
		return response.Error(c, http.StatusInternalServerError, "failed to create data", 500)
	}

	return response.Success(c, http.StatusCreated, "Save Prediction Successfully", nil)
}

func (ph *PredictionHandler) GetAllPredictions(c echo.Context) error {
	pageInt, limitInt := utils.ParsePaginationParams(c, 10)
	search := c.QueryParam("search")

	items, total, err := ph.service.GetAllPrediction(c.Request().Context(), pageInt, limitInt, search)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err)
		}
		return response.Error(c, http.StatusInternalServerError, "failed to get data", 500)
	}

	meta := utils.BuildPaginationMeta(c, pageInt, limitInt, total)
	data := make([]predictionresponse.PredictionResponse, len(items))
	for i, p := range items {
		data[i] = predictionresponse.ToPredictionResponse(*p)
	}

	return response.PaginatedSuccess(c, http.StatusOK, "Get All Prediction Succesfully", data, meta)
}

func (ph *PredictionHandler) DeletePrediction(c echo.Context) error {
	predictionId, err := uuid.Parse(c.Param("predictionId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	if err := ph.service.DeletePrediction(c.Request().Context(), predictionId); err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err)
		}
		return response.Error(c, http.StatusInternalServerError, "failed to delete data", 500)
	}

	return response.Success(c, http.StatusOK, "Deleted Prediction Successfully", nil)
}

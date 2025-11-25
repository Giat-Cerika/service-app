package materialhandler

import (
	materialrequest "giat-cerika-service/internal/dto/request/material_request"
	materialresponse "giat-cerika-service/internal/dto/response/material_response"
	materialservice "giat-cerika-service/internal/services/material_service"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"giat-cerika-service/pkg/constant/response"
	"giat-cerika-service/pkg/utils"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type MaterialHandler struct {
	materialService materialservice.IMaterialService
}

func NewMaterialHandler(service materialservice.IMaterialService) *MaterialHandler {
	return &MaterialHandler{materialService: service}
}

func (ch *MaterialHandler) CreateMaterial(c echo.Context) error {
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: "+err.Error(), nil)
	}

	adminID := claims.UserID
	var req materialrequest.CreateMaterialRequest
	req.Title = c.FormValue("title")
	req.Description = c.FormValue("description")
	if cover, err := c.FormFile("cover"); err == nil {
		req.Cover = cover
	}

	gallery := []*multipart.FileHeader{}
	formGallery, _ := c.MultipartForm()
	if formGallery != nil {
		if files, ok := formGallery.File["gallery"]; ok {
			gallery = files
		}
	}
	req.Gallery = gallery

	err = ch.materialService.CreateMaterial(c.Request().Context(), uuid.MustParse(adminID), req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to create material")
	}

	return response.Success(c, http.StatusOK, "Material Created Succssfully", nil)
}

func (ch *MaterialHandler) GetAllMaterial(c echo.Context) error {
	pageInt, limitInt := utils.ParsePaginationParams(c, 10)
	search := c.QueryParam("search")

	materiales, total, err := ch.materialService.GetAllMaterial(c.Request().Context(), pageInt, limitInt, search)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get materiales")
	}

	meta := utils.BuildPaginationMeta(c, pageInt, limitInt, total)
	data := make([]materialresponse.MaterialResponse, len(materiales))
	for i, material := range materiales {
		data[i] = materialresponse.ToMaterialResponse(*material)
	}

	return response.PaginatedSuccess(c, http.StatusOK, "Get All Materiales Successfully", data, meta)
}

func (ch *MaterialHandler) GetByIdMaterial(c echo.Context) error {
	materialId, err := uuid.Parse(c.Param("materialId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	material, err := ch.materialService.GetByIdMaterial(c.Request().Context(), materialId)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get material")
	}

	res := materialresponse.ToMaterialResponse(*material)

	return response.Success(c, http.StatusOK, "Get Material Successfully", res)
}

func (ch *MaterialHandler) UpdateMaterial(c echo.Context) error {
	materialId, err := uuid.Parse(c.Param("materialId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	var req materialrequest.UpdateMaterialRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	err = ch.materialService.UpdateMaterial(c.Request().Context(), materialId, req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to update material")
	}

	return response.Success(c, http.StatusOK, "Material Updated Successfully", nil)
}

func (ch *MaterialHandler) DeleteMaterial(c echo.Context) error {
	materialId, err := uuid.Parse(c.Param("materialId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	if err := ch.materialService.DeleteMaterial(c.Request().Context(), materialId); err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to delete material")
	}

	return response.Success(c, http.StatusOK, "Material Deleted Successfully", nil)
}

func (ch *MaterialHandler) GetAllLatestMateriaL(c echo.Context) error {
	materiales, err := ch.materialService.GetAllLatestMaterial(c.Request().Context())
	if err != nil {
		// Tangani Custom Error dari service layer
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		// Tangani Internal Server Error
		return response.Error(c, http.StatusInternalServerError, "Failed to get latest materials", err.Error())
	}

	data := make([]materialresponse.MaterialResponse, len(materiales))
	for i, material := range materiales {
		data[i] = materialresponse.ToMaterialResponse(*material)
	}

	return response.Success(c, http.StatusOK, "Get Latest Materials Successfully", data)
}

func (ch *MaterialHandler) GetAllPublicMaterial(c echo.Context) error {
	pageInt, limitInt := utils.ParsePaginationParams(c, 10)
	search := c.QueryParam("search")

	materiales, total, err := ch.materialService.GetAllPublicMaterial(c.Request().Context(), pageInt, limitInt, search)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get materiales")
	}

	meta := utils.BuildPaginationMeta(c, pageInt, limitInt, total)
	data := make([]materialresponse.MaterialResponse, len(materiales))
	for i, material := range materiales {
		data[i] = materialresponse.ToMaterialResponse(*material)
	}

	return response.PaginatedSuccess(c, http.StatusOK, "Get All Materiales Successfully", data, meta)
}

func (ch *MaterialHandler) GetByIdPublicMaterial(c echo.Context) error {
	materialId, err := uuid.Parse(c.Param("materialId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	material, err := ch.materialService.GetByIdPublicMaterial(c.Request().Context(), materialId)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get material")
	}

	res := materialresponse.ToMaterialResponse(*material)

	return response.Success(c, http.StatusOK, "Get Material Successfully", res)
}

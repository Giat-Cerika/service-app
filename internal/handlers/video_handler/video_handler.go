package videohandler

import (
	videorequest "giat-cerika-service/internal/dto/request/video_request"
	videoresponse "giat-cerika-service/internal/dto/response/video_response"
	videoservice "giat-cerika-service/internal/services/video_service"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"giat-cerika-service/pkg/constant/response"
	"giat-cerika-service/pkg/utils"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type VideoHandler struct {
	videoService videoservice.IVideoService
}

func NewVideoHandler(service videoservice.IVideoService) *VideoHandler {
	return &VideoHandler{videoService: service}
}

func (ch *VideoHandler) CreateVideo(c echo.Context) error {
	var req videorequest.CreateVideoRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "unauthorized", "invalid token type")
	}

	claims, ok := token.Claims.(*utils.JWTClaims)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "unauthorized", "invalid token claims")
	}

	creatorID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "invalid user id in token", err)
	}

	// LEMPAR KE SERVICE + CREATOR-ID
	err = ch.videoService.CreateVideo(c.Request().Context(), req, creatorID)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to create video")
	}

	return response.Success(c, http.StatusOK, "Video created successfully", nil)
}

func (ch *VideoHandler) GetAllVideo(c echo.Context) error {
	pageInt, limitInt := utils.ParsePaginationParams(c, 10)
	search := c.QueryParam("search")

	videoes, total, err := ch.videoService.GetAllVideo(c.Request().Context(), pageInt, limitInt, search)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get videoes")
	}

	meta := utils.BuildPaginationMeta(c, pageInt, limitInt, total)
	data := make([]videoresponse.VideoResponse, len(videoes))
	for i, video := range videoes {
		data[i] = videoresponse.ToVideoResponse(*video)
	}

	return response.PaginatedSuccess(c, http.StatusOK, "Get All Videoes Successfully", data, meta)
}

func (ch *VideoHandler) GetByIdVideo(c echo.Context) error {
	videoId, err := uuid.Parse(c.Param("videoId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	video, err := ch.videoService.GetByIdVideo(c.Request().Context(), videoId)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get video")
	}

	res := videoresponse.ToVideoResponse(*video)

	return response.Success(c, http.StatusOK, "Get Video Successfully", res)
}

func (ch *VideoHandler) UpdateVideo(c echo.Context) error {
	videoId, err := uuid.Parse(c.Param("videoId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	var req videorequest.UpdateVideoRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	err = ch.videoService.UpdateVideo(c.Request().Context(), videoId, req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to update video")
	}

	return response.Success(c, http.StatusOK, "Video Updated Successfully", nil)
}

func (ch *VideoHandler) DeleteVideo(c echo.Context) error {
	videoId, err := uuid.Parse(c.Param("videoId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	if err := ch.videoService.DeleteVideo(c.Request().Context(), videoId); err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to delete video")
	}

	return response.Success(c, http.StatusOK, "Video Deleted Successfully", nil)
}

func (ch *VideoHandler) GetAllLatestVideo(c echo.Context) error {
	videos, err := ch.videoService.GetAllLatestVideo(c.Request().Context())
	if err != nil {
		// Tangani Custom Error dari service layer
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		// Tangani Internal Server Error
		return response.Error(c, http.StatusInternalServerError, "Failed to get latest videos", err.Error())
	}

	data := make([]videoresponse.VideoResponse, len(videos))
	for i, video := range videos {
		data[i] = videoresponse.ToVideoResponse(*video)
	}

	return response.Success(c, http.StatusOK, "Get Latest Videos Successfully", data)

}

func (ch *VideoHandler) GetAllPublicVideo(c echo.Context) error {
	pageInt, limitInt := utils.ParsePaginationParams(c, 10)
	search := c.QueryParam("search")

	videoes, total, err := ch.videoService.GetAllPublicVideo(c.Request().Context(), pageInt, limitInt, search)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get videoes")
	}

	meta := utils.BuildPaginationMeta(c, pageInt, limitInt, total)
	data := make([]videoresponse.VideoResponse, len(videoes))
	for i, video := range videoes {
		data[i] = videoresponse.ToVideoResponse(*video)
	}

	return response.PaginatedSuccess(c, http.StatusOK, "Get All Videoes Successfully", data, meta)
}

func (ch *VideoHandler) GetByIdPublicVideo(c echo.Context) error {
	videoId, err := uuid.Parse(c.Param("videoId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	video, err := ch.videoService.GetByIdPublicVideo(c.Request().Context(), videoId)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get video")
	}

	res := videoresponse.ToVideoResponse(*video)

	return response.Success(c, http.StatusOK, "Get Video Successfully", res)
}

package handlers

import (
	// "fmt"
	"net/http"
	"user-service/api/dto"
	"user-service/api/helper"
	"user-service/config"
	"user-service/services"
	"github.com/gin-gonic/gin"
)


type LocationHandler struct {
	service *services.LocationService
	labelService *services.LabelService
}

func NewLocationHandler(cfg *config.Config) *LocationHandler {
	return &LocationHandler{
		service: services.NewLocationService(cfg),
		labelService: services.NewLabelService(cfg),
	}
}

// CreateLocation godoc
// @Summary Create a Location
// @Description Create a Location
// @Tags Locations
// @Accept json
// @produces json
// @Param Request body dto.CreateLocationRequest true "Create a Location"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.LocationResponse} "Location response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v2/auth/resource/locations/ [post]
// @Security AuthBearer
func (h *LocationHandler) Create(c *gin.Context) {
	req := dto.CreateLocationRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}

	res, err := h.service.Create(c, &req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}

	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(res, true, helper.Success))
}

// UpdateLocation godoc
// @Summary Update a Location
// @Description Update a Location
// @Tags Locations
// @Accept json
// @produces json
// @Param Request body dto.UpdateLocationRequest true "Update a Location"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.LocationResponse} "Location response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v2/auth/resource/locations/ [put]
// @Security AuthBearer
func (h *LocationHandler) Update(c *gin.Context) {	
	req := dto.UpdateLocationRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, 121, err))
		return
	}
	res, err := h.service.Update(c, &req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, 121, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, 0))
}


// GetLocation godoc
// @Summary Get a Location
// @Description Get a Location
// @Tags Locations
// @Accept json
// @produces json
// @Success 200 {object} helper.BaseHttpResponse{result=dto.LocationResponse} "Location response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v2/auth/resource/locations/ [get]
// @Security AuthBearer
func (h *LocationHandler) GetLocation(c *gin.Context) {
	res, err := h.service.GetById(c)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, helper.Success))
}
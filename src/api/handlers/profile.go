package handlers

import (
	// "fmt"
	"net/http"
	"user-service/api/dto"
	"user-service/api/helper"
	"user-service/config"
	"user-service/constant"

	// "user-service/constant"
	"user-service/services"

	"github.com/gin-gonic/gin"
)


type ProfileHandler struct {
	service *services.ProfileService
	labelService *services.LabelService
}

func NewProfileHandler(cfg *config.Config) *ProfileHandler {
	return &ProfileHandler{
		service: services.NewProfileService(cfg),
		labelService: services.NewLabelService(cfg),
	}
}

// CreateProfile godoc
// @Summary Create a Profile
// @Description Create a Profile
// @Tags Profiles
// @Accept json
// @produces json
// @Param Request body dto.CreateProfileRequest true "Create a Profile"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.ProfileResponse} "Profile response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v2/auth/resource/profiles/ [post]
// @Security AuthBearer
func (h *ProfileHandler) Create(c *gin.Context) {
	req := dto.CreateProfileRequest{}
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

	labelReq := &dto.CreateLabelRequest{
		Key:         constant.LabelKeyProfile,
		Value:       constant.LabelValueSubmitted,
		Scope:       constant.LabelScopePrivate,
		Description: constant.LabelDescriptionProfileSubmitted,
	}

	_, err = h.labelService.Create(c, labelReq)

	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(res, true, helper.Success))
}

// UpdateProfile godoc
// @Summary Update a Profile
// @Description Update a Profile
// @Tags Profiles
// @Accept json
// @produces json
// @Param Request body dto.UpdateProfileRequest true "Update a Profile"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.ProfileResponse} "Profile response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v2/auth/resource/profiles/ [put]
// @Security AuthBearer
func (h *ProfileHandler) Update(c *gin.Context) {	
	req := dto.UpdateProfileRequest{}
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


// GetProfile godoc
// @Summary Get a Profile
// @Description Get a Profile
// @Tags Profiles
// @Accept json
// @produces json
// @Success 200 {object} helper.BaseHttpResponse{result=dto.ProfileResponse} "Profile response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v2/auth/resource/profiles/me [get]
// @Security AuthBearer
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	res, err := h.service.GetById(c)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, helper.Success))
}
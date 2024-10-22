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


type PhoneHandler struct {
	service *services.PhoneService
	labelService *services.LabelService
}

func NewPhoneHandler(cfg *config.Config) *PhoneHandler {
	return &PhoneHandler{
		service: services.NewPhoneService(cfg),
		labelService: services.NewLabelService(cfg),
	}
}

// CreatePhone godoc
// @Summary Create a Phone
// @Description Create a Phone
// @Tags Phones
// @Accept json
// @produces json
// @Param Request body dto.CreatePhoneRequest true "Create a Phone"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.PhoneResponse} "Phone response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v2/auth/resource/phones/ [post]
// @Security AuthBearer
func (h *PhoneHandler) Create(c *gin.Context) {
	req := dto.CreatePhoneRequest{}
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
		Key:         constant.LabelKeyPhone,
		Value:       constant.LabelValueSubmitted,
		Scope:       constant.LabelScopePrivate,
		Description: constant.LabelDescriptionPhoneSubmitted,
	}

	_, err = h.labelService.Create(c, labelReq)

	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(res, true, helper.Success))
}

// UpdatePhone godoc
// @Summary Update a Phone
// @Description Update a Phone
// @Tags Phones
// @Accept json
// @produces json
// @Param Request body dto.UpdatePhoneRequest true "Update a Phone"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.PhoneResponse} "Phone response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v2/auth/resource/phones/ [put]
// @Security AuthBearer
func (h *PhoneHandler) Update(c *gin.Context) {	
	req := dto.UpdatePhoneRequest{}
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
package handlers

import (
	"fmt"
	"net/http"
	"user-service/api/dto"
	"user-service/api/helper"
	"user-service/config"
	"user-service/constant"
	"user-service/services"

	"github.com/gin-gonic/gin"
)


type LabelHandler struct {
	service *services.LabelService
}

func NewLabelHandler(cfg *config.Config) *LabelHandler {
	return &LabelHandler{
		service: services.NewLabelService(cfg),
	}
}

// CreateLabel godoc
// @Summary Create a Label
// @Description Create a Label
// @Tags Labels
// @Accept json
// @produces json
// @Param Request body dto.CreateLabelRequest true "Create a Label"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.LabelResponse} "Label response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v2/auth/resource/labels [post]
// @Security AuthBearer
func (h *LabelHandler) Create(c *gin.Context) {
	req := dto.CreateLabelRequest{}
	
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}

	// Set the scope to public
	req.Scope = constant.LabelScopePublic

	res, err := h.service.Create(c, &req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}

	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(res, true, helper.Success))
}


// UpdateLabel godoc
// @Summary Update a Label
// @Description Update a Label. This operation is allowed only for labels with Scope: public.
// @Tags Labels
// @Accept json
// @Produces json
// @Param key path string true "Label Key"
// @Param Request body dto.UpdateLabelRequest true "Update a Label"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.LabelResponse} "Label response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Failure 403 {object} helper.BaseHttpResponse "Forbidden: Label scope is not public"
// @Router /v2/auth/resource/labels/{key} [put]
// @Security AuthBearer
func (h *LabelHandler) Update(c *gin.Context) {
	key := c.Param("key") // Extract the key from the URL path

	req := dto.UpdateLabelRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}

	res, err := h.service.Update(c, key, &req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, helper.Success))
}


// GetLabelByKey godoc
// @Summary Get a Label by Key
// @Description Retrieve a Label by its Key. This operation is allowed for labels with Scope: public or private.
// @Tags Labels
// @Param key path string true "Label Key"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.LabelResponse} "Label response"
// @Failure 404 {object} helper.BaseHttpResponse "Not found: Label not found or already deleted"
// @Failure 403 {object} helper.BaseHttpResponse "Forbidden: Label not found or scope is not accessible"
// @Router /v2/auth/resource/labels/{key} [get]
// @Security AuthBearer
func (h *LabelHandler) GetByKey(c *gin.Context) {
	key := c.Param("key") // Extract the key from the URL path

	// Call the service to get the label
	res, err := h.service.GetByKey(c, key)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, helper.Success))
}


// DeleteLabel godoc
// @Summary Delete a Label
// @Description Delete a Label. This operation is allowed only for labels with Scope: public.
// @Tags Labels
// @Param key path string true "Label Key"
// @Success 200 {object} helper.BaseHttpResponse "Label deleted successfully"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Failure 403 {object} helper.BaseHttpResponse "Forbidden: Label scope is not public"
// @Failure 404 {object} helper.BaseHttpResponse "Not found: Label not found or already deleted"
// @Router /v2/auth/resource/labels/{key} [delete]
// @Security AuthBearer
func (h *LabelHandler) Delete(c *gin.Context) {
	key := c.Param("key") // Extract the key from the URL path

	// Call the service to delete the label
	err := h.service.Delete(c, key)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, true, helper.Success))
}


// ListLabels godoc
// @Summary List all Labels for the current user
// @Description Retrieve all labels associated with the current user.
// @Tags Labels
// @Accept json
// @Produces json
// @Success 200 {object} helper.BaseHttpResponse{result=[]dto.LabelResponse} "List of Labels"
// @Failure 404 {object} helper.BaseHttpResponse "Not found: No labels found"
// @Failure 500 {object} helper.BaseHttpResponse "Internal server error"
// @Router /v2/auth/resource/labels/ [get]
// @Security AuthBearer
func (h *LabelHandler) ListLabels(c *gin.Context) {
	res, err := h.service.ListLabelsForCurrentUser(c)
	fmt.Println("res",res)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, helper.Success))
}

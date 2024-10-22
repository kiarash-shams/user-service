package handlers

import (
	"net/http"
	"time"
	"user-service/api/helper"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck godoc
// @Summary Health Check
// @Description Health Check
// @Tags Public
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Router /v2/auth/public/health/ [get]
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, helper.GenerateBaseResponse("Working!", true, helper.Success))
}

// GetAuthVersion godoc
// @Summary Get Auth Version
// @Description Auth Version
// @Tags Public
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Router /v2/auth/public/version/ [get]
func (h *HealthHandler) Version(c *gin.Context) {
	c.JSON(http.StatusOK, helper.GenerateBaseResponse("version: 0.0.1", true, helper.Success))
}


// GetServerCurrentUnixTimestamp godoc
// @Summary Get Server Current Unix Timestamp
// @Description Auth Server Current Unix Timestamp
// @Tags Public
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Router /v2/auth/public/time/ [get]
func (h *HealthHandler) GetServerCurrentUnixTimestamp (c *gin.Context) {
	currentUnixTimestamp := time.Now().Unix()
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(currentUnixTimestamp, true, helper.Success))
}
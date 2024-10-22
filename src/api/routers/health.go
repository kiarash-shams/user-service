package routers

import (
	"github.com/gin-gonic/gin"
	"user-service/api/handlers"
)

func Health(r *gin.RouterGroup) {
	handler := handlers.NewHealthHandler()

	r.GET("/health", handler.Health)
	r.GET("/version", handler.Version)
	r.GET("/time", handler.GetServerCurrentUnixTimestamp)

}

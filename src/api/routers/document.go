package routers

import (
	"github.com/gin-gonic/gin"
	"user-service/api/handlers"
	"user-service/config"
)


func Document(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewDocumentHandler(cfg)
	
	r.POST("/", h.Create)
	r.GET("/:id", h.GetById)
}

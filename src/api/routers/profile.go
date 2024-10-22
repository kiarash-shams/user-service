package routers

import (
	"github.com/gin-gonic/gin"
	"user-service/api/handlers"
	"user-service/config"
)


func Profile(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewProfileHandler(cfg)
	r.POST("/", h.Create)
	r.PUT("/", h.Update)
	r.GET("/me", h.GetProfile)
}

func Label(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewLabelHandler(cfg)
	r.GET("/", h.ListLabels)
	r.POST("/", h.Create)
	r.PUT("/:key", h.Update)
	r.DELETE("/:key", h.Delete)
	r.GET("/:key", h.GetByKey)
}

func Phone(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewPhoneHandler(cfg)
	r.POST("/", h.Create)
	r.PUT("/", h.Update)
}


func Location(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewLocationHandler(cfg)
	r.POST("/", h.Create)
	r.PUT("/", h.Update)
	r.GET("/me", h.GetLocation)
}
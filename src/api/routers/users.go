package routers

import (
	"github.com/gin-gonic/gin"
	"user-service/api/handlers"
	// "user-service/api/middleware"
	"user-service/config"
)

func User(router *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewUserHandler(cfg)

	router.POST("/login", h.LoginByEmail)
	router.POST("/register", h.RegisterByEmail)
	// router.POST("/login-by-mobile", h.RegisterLoginByMobileNumber)
	// router.POST("/send-otp", middleware.OtpLimiter(cfg), h.SendOtp)
}

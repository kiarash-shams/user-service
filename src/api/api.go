package api

import (
	"fmt"
	"log"
	// "log"
	// "net/http"
	"time"
	"user-service/api/middleware"
	"user-service/api/routers"
	"user-service/api/validation"
	"user-service/config"
	"user-service/docs"
	"user-service/pkg/logging"
	"user-service/pkg/metrics"
	"user-service/services"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus"

	// "github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var logger = logging.NewLogger(config.GetConfig())

func InitServer(cfg *config.Config) {

	gin.SetMode(cfg.Server.RunMode)
	r := gin.New()

	
    rabbitmqService, err := services.NewRabbitMQService("amqp://user:password@localhost:5672/")
    if err != nil {
        log.Fatalf("Failed to initialize RabbitMQ service: %v", err)
    }
    defer rabbitmqService.Close()

    producer := services.NewNotificationProducer(rabbitmqService)

    // Create and start consumers
    consumers := []string{"sms", "email", "webhook"}
    for _, consumerType := range consumers {
        consumer := services.NewNotificationConsumer(rabbitmqService, consumerType)
        go func() {
            if err := consumer.Start(); err != nil {
                log.Printf("Failed to start %s consumer: %v", consumerType, err)
            }
        }()
    }

    // Example: Queue a notification
    notification := &services.Notification{
        Type:      "email",
        Recipient: "+1234567890",
        Message:   "Hello, World!",
        CreatedAt: time.Now(),
        Status:    "pending",
    }

	
	

    if err := producer.QueueNotification(notification); err != nil {
        log.Printf("Failed to queue notification: %v", err)
    }

    // Keep the application running
    select {}
// ---------------------------------------------------------------------------------

	

	// r.POST("/send-notification", func(c *gin.Context) {
    //     var req struct {
    //         Recipient string `json:"recipient" binding:"required"`
    //         Message   string `json:"message" binding:"required"`
    //         Type      string `json:"type" binding:"required"`
    //     }

    //     if err := c.ShouldBindJSON(&req); err != nil {
    //         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    //         return
    //     }

    //     // Get the notifier based on the request type
    //     notifier := services.NewNotifier(req.Type)
    //     if notifier == nil {
    //         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification type"})
    //         return
    //     }

    //     // Create the notification struct with a pending status
    //     notification := &services.Notification{
    //         Type:      req.Type,
    //         Recipient: req.Recipient,
    //         Message:   req.Message,
    //         CreatedAt: time.Now(),
    //         Status:    "pending",
    //     }

    //     // Send the notification and update the status
    //     err := notifier.Send(notification)
    //     if err != nil {
    //         notification.Status = "failed"
    //         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification", "status": notification.Status})
    //         return
    //     }

    //     // Notification sent successfully, return the updated status
    //     c.JSON(http.StatusOK, gin.H{
    //         "status":     notification.Status,
    //         "recipient":  notification.Recipient,
    //         "type":       notification.Type,
    //         "created_at": notification.CreatedAt,
    //     })
    // })



	RegisterValidators()
	// RegisterPrometheus()

	r.Use(middleware.DefaultStructuredLogger(cfg))
	r.Use(middleware.Cors(cfg))
	r.Use(middleware.Prometheus())
	r.Use(gin.Logger(), gin.CustomRecovery(middleware.ErrorHandler), middleware.LimitByRequest())

	RegisterRoutes(r, cfg)
	RegisterSwagger(r, cfg)

	err = r.Run(fmt.Sprintf(":%s", cfg.Server.InternalPort))
	if err != nil {
		logger.Error(logging.General, logging.Startup, err.Error(), nil)
	}

}

func RegisterRoutes(r *gin.Engine, cfg *config.Config) {
	api := r.Group("/api/v2/auth")

	public := api.Group("/public")
	{
		health := public.Group("/")
		routers.Health(health)
	}

	identity := api.Group("/identity")
	{
		users := identity.Group("/users")
		routers.User(users, cfg)
	}

	v1 := api.Group("/resource")
	{
		profiles := v1.Group("/profiles", middleware.Authentication(cfg))
		labels := v1.Group("/labels", middleware.Authentication(cfg))
		phones := v1.Group("/phones", middleware.Authentication(cfg))
		documents := v1.Group("/documents", middleware.Authentication(cfg))
		locations := v1.Group("/locations", middleware.Authentication(cfg))

		routers.Profile(profiles, cfg)
		routers.Label(labels, cfg)
		routers.Phone(phones, cfg)	
		routers.Document(documents, cfg)
		routers.Location(locations, cfg)
	}



	// routers.File(files, cfg)
	r.Static("/static", "./uploads")
	// r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

func RegisterValidators() {
	val, ok := binding.Validator.Engine().(*validator.Validate)

	if ok {
		err := val.RegisterValidation("mobile", validation.IranianMobileNumberValidator)
		if err != nil {
			logger.Error(logging.Validation, logging.Startup, err.Error(), nil)
		}

		err = val.RegisterValidation("password", validation.PasswordValidator)
		if err != nil {
			logger.Error(logging.Validation, logging.Startup, err.Error(), nil)
		}
	}
}

func RegisterSwagger(r *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.Title = "Auth User Microservice"
	docs.SwaggerInfo.Description = "This microservice handles user authentication, including registration, login, password management, and 2FA. It also manages user profiles, KYC (Know Your Customer) verification, and OAuth integration. The microservice provides endpoints for managing user information, verifying identity, and securing user access."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Server.ExternalPort)
	docs.SwaggerInfo.Schemes = []string{"http"}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func RegisterPrometheus() {
	err := prometheus.Register(metrics.DbCall)
	if err != nil {
		logger.Error(logging.Prometheus, logging.Startup, err.Error(), nil)
	}

	err = prometheus.Register(metrics.HttpDuration)
	if err != nil {
		logger.Error(logging.Prometheus, logging.Startup, err.Error(), nil)
	}
}

package main

import (
	// "log"
	"user-service/api"
	"user-service/config"
	"user-service/data/cache"
	"user-service/data/db"
	"user-service/data/db/migration"
	"user-service/services"

	// "user-service/data/mq"
	// "user-service/data/vault"
	"user-service/pkg/logging"
)

// @securityDefinitions.apikey AuthBearer
// @in header
// @name Authorization
func main() {

	// Get Config
	cfg := config.GetConfig()

	logger := logging.NewLogger(cfg)
	// Start Redis
	err := cache.InitRedis(cfg)
	defer cache.CloseRedis()
	if err != nil {
		logger.Fatal(logging.Redis, logging.Startup, err.Error(), nil)
	}

	// Start DB Postgres
	err = db.InitDB(cfg)
	defer db.CloseDb()
	if err != nil {
		logger.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}


	// Initialize RabbitMQ Service
	rabbitmqService, err := services.NewRabbitMQService("amqp://user:password@localhost:5672/",logger)
	if err != nil {
		logger.Fatal(logging.RabbitMQ, logging.Startup, err.Error(), nil)
	}
	defer rabbitmqService.Close()

	producer := services.NewNotificationProducer(rabbitmqService)

	// Create and start consumers
	consumers := []string{"sms", "email", "webhook"}
	for _, consumerType := range consumers {
		consumer := services.NewNotificationConsumer(rabbitmqService, consumerType)
		go func(ct string) {
			if err := consumer.Start(); err != nil {
				logger.Errorf("Failed to start %s consumer: %v", ct, err)
			}
		}(consumerType)
	}


	migration.Up1()

	// Run Server
	api.InitServer(cfg, producer)

}

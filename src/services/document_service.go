package services

import (
	"context"
	"reflect"
	"time"
	"user-service/api/dto"
	_ "user-service/api/helper"
	"user-service/config"
	"user-service/constant"
	"user-service/data/db"
	models "user-service/data/models"
	"user-service/pkg/logging"
	"user-service/pkg/metrics"

	"gorm.io/gorm"
)


type DocumentService struct {
	database *gorm.DB
	logger   logging.Logger
	Preloads []preload
}


func NewDocumentService(cfg *config.Config) *DocumentService {
	return &DocumentService{
		database: db.GetDb(),
		logger:   logging.NewLogger(cfg),
		Preloads: []preload{{string: "File"},},
	}
}


func (s *DocumentService) Create(ctx context.Context, req *dto.CreateDocumentRequest) (*dto.DocumentResponse, error) {
	Document := models.Document{
		DocCategory: req.DocCategory,
		UserId: int(ctx.Value(constant.UserIdKey).(float64)),
	}
	Document.CreatedAt = time.Now().UTC()

	tx := s.database.WithContext(ctx).Begin()
	err := tx.
		Create(&Document).
		Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(Document).String(), "Create", "Failed ").Inc()
		return nil, err
	}
	tx.Commit()
	
	metrics.DbCall.WithLabelValues(reflect.TypeOf(Document).String(), "Create", "Success").Inc()
	
	return s.GetById(ctx,Document.Id)
}


func (s *DocumentService) GetById(ctx context.Context, id int) (*dto.DocumentResponse, error) {
    // Retrieve the Document by ID
    var document models.Document
    err := s.database.WithContext(ctx).
        Model(&models.Document{}).
        Preload("File").
        Where("id = ? AND deleted_by IS NULL", id).
        First(&document).
        Error
    if err != nil {
        s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
        metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Document{}).String(), "GetById", "Failed").Inc()
        return nil, err
    }

    // Map the Document to the DTO response
    dtoResponse := &dto.DocumentResponse{
        Id:          document.Id,
        DocCategory: document.DocCategory,
        Name:        document.Name,
        Directory:   document.Directory,
        Description: document.Description,
        MimeType:    document.MimeType,
    }

    metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Document{}).String(), "GetById", "Success").Inc()
    return dtoResponse, nil
}


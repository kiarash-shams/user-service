package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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


type LabelService struct {
	database *gorm.DB
	logger   logging.Logger
}


func NewLabelService(cfg *config.Config) *LabelService {
	return &LabelService{
		database: db.GetDb(),
		logger:   logging.NewLogger(cfg),
	}
}


func (s *LabelService) Create(ctx context.Context, req *dto.CreateLabelRequest) (*dto.LabelResponse, error) {
	Label := models.Label{
		Key: req.Key,
		Value: req.Value,
		Scope: req.Scope,
		Description: req.Description,
		UserId: int(ctx.Value(constant.UserIdKey).(float64)),
	}
	Label.CreatedAt = time.Now().UTC()
	tx := s.database.WithContext(ctx).Begin()
	err := tx.
		Create(&Label).
		Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(Label).String(), "Create", "Failed ").Inc()
		return nil, err
	}
	tx.Commit()
	dto := &dto.LabelResponse{
		Key: req.Key,
		Value: req.Value,
		Description: req.Description,
	}
	metrics.DbCall.WithLabelValues(reflect.TypeOf(Label).String(), "Create", "Success").Inc()
	return dto, nil
}


func (s *LabelService) Update(ctx context.Context, key string, req *dto.UpdateLabelRequest) (*dto.LabelResponse, error) {
	// Get the user ID from the context
	userId := ctx.Value(constant.UserIdKey).(float64)

	// First, retrieve the label by key, user ID, and check its scope
	var label models.Label
	err := s.database.WithContext(ctx).
		Model(&models.Label{}).
		Where("key = ? AND user_id = ? AND deleted_by IS NULL", key, int64(userId)).
		First(&label).
		Error
	if err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Label{}).String(), "Select", "Failed").Inc()
		return nil, err
	}

	// Check if the scope is public
	if label.Scope != constant.LabelScopePublic {
		return nil, errors.New("update not allowed: label scope is not public")
	}

	// Prepare the update map
	updateMap := map[string]interface{}{
		"Value":       req.Value,
		"Description": req.Description,
		"modified_by": &sql.NullInt64{Int64: int64(userId), Valid: true},
		"modified_at": sql.NullTime{Valid: true, Time: time.Now().UTC()},
	}

	// Begin transaction
	tx := s.database.WithContext(ctx).Begin()
	err = tx.
		Model(&models.Label{}).
		Where("key = ? AND user_id = ? AND deleted_by IS NULL", key, int64(userId)).
		Updates(updateMap).
		Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Label{}).String(), "Update", "Failed").Inc()
		return nil, err
	}

	// Commit the transaction
	tx.Commit()

	dto := &dto.LabelResponse{
		Key:         label.Key,
		Value:       req.Value,
		Description: req.Description,
	}
	metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Label{}).String(), "Update", "Success").Inc()
	return dto, nil
}


func (s *LabelService) Delete(ctx context.Context, key string) error {
	
	// Get the user ID from the context
	userId := ctx.Value(constant.UserIdKey).(float64)

	// First, retrieve the label by key, user ID, and check its scope
	var label models.Label
	err := s.database.WithContext(ctx).
		Model(&models.Label{}).
		Where("key = ? AND user_id = ? AND deleted_by IS NULL", key, int64(userId)).
		First(&label).
		Error
	if err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Label{}).String(), "Select", "Failed").Inc()
		return err
	}

	// Check if the scope is public
	if label.Scope != constant.LabelScopePublic {
		return errors.New("delete not allowed: label scope is not public")
	}

	// Prepare the delete map
	deleteMap := map[string]interface{}{
		"deleted_by": &sql.NullInt64{Int64: int64(userId), Valid: true},
		"deleted_at": sql.NullTime{Valid: true, Time: time.Now().UTC()},
	}

	// Begin transaction
	tx := s.database.WithContext(ctx).Begin()
	err = tx.
		Model(&models.Label{}).
		Where("key = ? AND user_id = ? AND deleted_by IS NULL", key, int64(userId)).
		Updates(deleteMap).
		Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Delete, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Label{}).String(), "Delete", "Failed").Inc()
		return err
	}

	// Commit the transaction
	tx.Commit()
	metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Label{}).String(), "Delete", "Success").Inc()
	return nil
}


func (s *LabelService) GetByKey(ctx context.Context, key string) (*dto.LabelResponse, error) {
	// Get the user ID from the context
	userId := ctx.Value(constant.UserIdKey).(float64)

	// Retrieve the label by key and user ID
	var label models.Label
	err := s.database.WithContext(ctx).
		Model(&models.Label{}).
		Where("key = ? AND user_id = ? AND deleted_by IS NULL", key, int64(userId)).
		First(&label).
		Error
	if err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Label{}).String(), "GetByKey", "Failed").Inc()
		return nil, err
	}

	// Map the label to the DTO response
	dto := &dto.LabelResponse{
		Key:         label.Key,
		Value:       label.Value,
		Description: label.Description,
	}

	metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Label{}).String(), "GetByKey", "Success").Inc()
	return dto, nil
}

// ListLabelsForCurrentUser retrieves all labels for the current user.
func (s *LabelService) ListLabelsForCurrentUser(ctx context.Context) ([]dto.LabelResponse, error) {
	var labels []models.Label

	// Retrieve the User ID from the context
	userID := int64(ctx.Value(constant.UserIdKey).(float64))
	fmt.Println(userID)

	// Query the database for labels belonging to the current user
	err := s.database.WithContext(ctx).
		Where("user_id = ? AND deleted_by IS NULL", userID).
		Find(&labels).
		Error
	if err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Label{}).String(), "ListLabel", "Failed").Inc()
		return nil, err
	}

	// Prepare the response DTOs
	var labelResponses []dto.LabelResponse
	for _, label := range labels {
		labelResponses = append(labelResponses, dto.LabelResponse{
			Key:         label.Key,
			Value:       label.Value,
			Scope:       label.Scope,
			Description: label.Description,
		})
	}

	metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Label{}).String(), "ListLabel", "Success").Inc()
	return labelResponses, nil
}



package services

import (
	"context"
	"database/sql"
	"reflect"
	"fmt"
	"time"
	"user-service/api/dto"
	_ "user-service/api/dto"
	_ "user-service/api/helper"
	"user-service/config"
	"user-service/constant"
	"user-service/data/db"
	models "user-service/data/models"
	"user-service/pkg/logging"
	"user-service/pkg/metrics"

	"gorm.io/gorm"
)


type LocationService struct {
	database *gorm.DB
	logger   logging.Logger
}


func NewLocationService(cfg *config.Config) *LocationService {
	return &LocationService{
		database: db.GetDb(),
		logger:   logging.NewLogger(cfg),
	}
}


func (s *LocationService) Create(ctx context.Context, req *dto.CreateLocationRequest) (*dto.LocationResponse, error) {

	// Check if the PostalCode already exists
	postalCodeExists, err := s.existsByPostalCode(req.PostalCode)
	if err != nil {
		return nil, err
	}
	if postalCodeExists {
		return nil, fmt.Errorf("The PostalCode is already in use.")
	}

	// Then, we check if there is a duplicate UserId
	userId := int(ctx.Value(constant.UserIdKey).(float64))
	userIdExists, err := s.existsByUserId(userId)
	if err != nil {
		return nil, err
	}
	if userIdExists {
		return nil, fmt.Errorf("The user has already registered their Location.")
	}


	Location := models.Location{
		PostalCode: req.PostalCode,
		Address: req.Address,
		City: req.City,
		Country: req.Country,
		StaticPhoneNumber: req.StaticPhoneNumber,
		Metadata: req.Metadata,
		UserId: int(ctx.Value(constant.UserIdKey).(float64)),
	}
	Location.CreatedAt = time.Now().UTC()
	tx := s.database.WithContext(ctx).Begin()
	err = tx.
		Create(&Location).
		Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(Location).String(), "Create", "Failed ").Inc()
		return nil, err
	}
	tx.Commit()
	
	metrics.DbCall.WithLabelValues(reflect.TypeOf(Location).String(), "Create", "Success").Inc()
	return s.GetById(ctx)
	
}


func (s *LocationService) Update(ctx context.Context, req *dto.UpdateLocationRequest) (*dto.LocationResponse, error) {
	userId := int(ctx.Value(constant.UserIdKey).(float64))

	// Check Location verification status
	isVerified, err := s.isLocationVerified(ctx, userId)
	if err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return nil, err
	}
	if isVerified {
		return nil, fmt.Errorf("location is already validated and cannot be updated")
	}

	updateMap := map[string]interface{}{
		"PostalCode":   	req.PostalCode,
		"Address":    		req.Address,
		"City":    			req.City,
		"Country":         	req.Country,
		"StaticPhoneNumber":req.StaticPhoneNumber,
		"Metadata":         req.Metadata,
		"modified_by": 		&sql.NullInt64{Int64: int64(ctx.Value(constant.UserIdKey).(float64)), Valid: true},
		"modified_at": 		sql.NullTime{Valid: true, Time: time.Now().UTC()},
	}
	tx := s.database.WithContext(ctx).Begin()
	err = tx.
		Model(&models.Location{}).
		Where("user_id = ?", ctx.Value(constant.UserIdKey)).
		Updates(updateMap).
		Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Location{}).String(), "Update", "Failed ").Inc()
		return nil, err
	}
	Location := &models.Location{}
	err = tx.
		Model(&models.Location{}).
		Where("user_id = ? AND deleted_by is null", ctx.Value(constant.UserIdKey)).
		First(&Location).
		Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Location{}).String(), "Select", "Failed ").Inc()
		return nil, err
	}
	tx.Commit()
	
	metrics.DbCall.WithLabelValues(reflect.TypeOf(Location).String(), "Update", "Success ").Inc()
	return s.GetById(ctx)
}


func (s *LocationService) GetById(ctx context.Context) (*dto.LocationResponse, error) {
	Location := &models.Location{}
	
	err := s.database.
		Where("user_id = ? AND deleted_by is null", ctx.Value(constant.UserIdKey)).
		First(&Location).
		Error
	if err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(Location).String(), "GetById", "Failed").Inc()
		return nil, err
	}
	dto := &dto.LocationResponse{
		PostalCode: Location.PostalCode,
		Address: Location.Address,
		City: Location.City,
		Country: Location.Country,
		StaticPhoneNumber: Location.StaticPhoneNumber,
		Metadata: Location.Metadata,
	}
	metrics.DbCall.WithLabelValues(reflect.TypeOf(Location).String(), "GetById", "Success").Inc()
	return dto, nil
}


// Method to check if a postal code already exists in the database
func (s *LocationService) existsByPostalCode(postalCode string) (bool, error) {
	var count int64
	err := s.database.Model(&models.Location{}).
		Where("postal_code = ?", postalCode).
		Count(&count).
		Error
	if err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return count > 0, nil
}

func (s *LocationService) existsByUserId(userId int) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.Location{}).
		Select("count(*) > 0").
		Where("user_id = ?", userId).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Location{}).String(), "existsByUserId", "Failed ").Inc()
		return false, err
	}
	metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Location{}).String(), "existsByUserId", "Success").Inc()
	return exists, nil
}

func (s *LocationService) isLocationVerified(ctx context.Context, userId int) (bool, error) {
	Location := &models.Location{}
	err := s.database.WithContext(ctx).
		Model(&models.Location{}).
		Where("user_id = ? AND deleted_by is null", userId).
		First(&Location).
		Error
	if err != nil {
		return false, err
	}
	return Location.Verification, nil
}
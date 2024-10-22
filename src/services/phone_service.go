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


type PhoneService struct {
	database *gorm.DB
	logger   logging.Logger
}


func NewPhoneService(cfg *config.Config) *PhoneService {
	return &PhoneService{
		database: db.GetDb(),
		logger:   logging.NewLogger(cfg),
	}
}


func (s *PhoneService) Create(ctx context.Context, req *dto.CreatePhoneRequest) (*dto.PhoneResponse, error) {
	
	// First, we check if there is a duplicate Phone Number
	nidExists, err := s.existsByPhone(req.MobileNumber)
	if err != nil {
		return nil, err
	}
	if nidExists {
		return nil, fmt.Errorf("PhoneNumber '%s' already exists", req.MobileNumber)
	}

	// Then, we check if there is a duplicate UserId
	userId := int(ctx.Value(constant.UserIdKey).(float64))
	userIdExists, err := s.existsByUserId(userId)
	if err != nil {
		return nil, err
	}
	if userIdExists {
		return nil, fmt.Errorf("The User has already registered their PhoneNumber.")
	}


	Phone := models.Phone{
		Number: req.MobileNumber,
		UserId: int(ctx.Value(constant.UserIdKey).(float64)),
	}
	Phone.CreatedAt = time.Now().UTC()
	tx := s.database.WithContext(ctx).Begin()
	err = tx.
		Create(&Phone).
		Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(Phone).String(), "Create", "Failed ").Inc()
		return nil, err
	}
	tx.Commit()
	
	metrics.DbCall.WithLabelValues(reflect.TypeOf(Phone).String(), "Create", "Success").Inc()
	return s.GetById(ctx)
}

func (s *PhoneService) Update(ctx context.Context, req *dto.UpdatePhoneRequest) (*dto.PhoneResponse, error) {
	// Check validation status
	isValidated, err := s.IsPhoneValidated(ctx)
	if err != nil {
		return nil, err
	}

	// If the phone number is already validated, do not allow an update
	if isValidated {
		return nil, errors.New("phone number has already been validated and cannot be updated")
	}
	
	updateMap := map[string]interface{}{
		"Number":   req.MobileNumber,
		"modified_by": &sql.NullInt64{Int64: int64(ctx.Value(constant.UserIdKey).(float64)), Valid: true},
		"modified_at": sql.NullTime{Valid: true, Time: time.Now().UTC()},
	}
	tx := s.database.WithContext(ctx).Begin()
	err = tx.
		Model(&models.Phone{}).
		Where("user_id = ?", ctx.Value(constant.UserIdKey)).
		Updates(updateMap).
		Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Phone{}).String(), "Update", "Failed ").Inc()
		return nil, err
	}
	Phone := &models.Phone{}
	err = tx.
		Model(&models.Phone{}).
		Where("user_id = ? AND deleted_by is null", ctx.Value(constant.UserIdKey)).
		First(&Phone).
		Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Phone{}).String(), "Select", "Failed ").Inc()
		return nil, err
	}
	tx.Commit()
	metrics.DbCall.WithLabelValues(reflect.TypeOf(Phone).String(), "Update", "Success ").Inc()
	return s.GetById(ctx)
}

func (s *PhoneService) GetById(ctx context.Context) (*dto.PhoneResponse, error) {
	Phone := &models.Phone{}
	
	err := s.database.
		Where("user_id = ? AND deleted_by is null", ctx.Value(constant.UserIdKey)).
		First(&Phone).
		Error
	if err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(Phone).String(), "GetById", "Failed").Inc()
		return nil, err
	}
	dto := &dto.PhoneResponse{
		MobileNumber: Phone.Number,
	}
	metrics.DbCall.WithLabelValues(reflect.TypeOf(Phone).String(), "GetById", "Success").Inc()
	return dto, nil
}

func (s *PhoneService) existsByPhone(phone string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.Phone{}).
		Select("count(*) > 0").
		Where("number = ?", phone).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Phone{}).String(), "existsByPhone", "Failed ").Inc()
		return false, err
	}
	metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Phone{}).String(), "existsByPhone", "Success").Inc()
	return exists, nil
}

func (s *PhoneService) existsByUserId(userId int) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.Phone{}).
		Select("count(*) > 0").
		Where("user_id = ?", userId).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Phone{}).String(), "existsPhoneByUserId", "Failed ").Inc()
		return false, err
	}
	metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Phone{}).String(), "existsPhoneByUserId", "Success").Inc()
	return exists, nil
}

func (s *PhoneService) IsPhoneValidated(ctx context.Context) (bool, error) {
	phone := &models.Phone{}
	err := s.database.
		WithContext(ctx).
		Model(&models.Phone{}).
		Where("user_id = ? AND deleted_by IS NULL", ctx.Value(constant.UserIdKey)).
		First(&phone).
		Error

	if err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}

	// If ValidatedAt is not set, the phone number is not yet validated
	return phone.ValidatedAt.Valid, nil
}

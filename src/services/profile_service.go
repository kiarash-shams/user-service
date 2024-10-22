package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"time"
	"user-service/api/dto"
	_ "user-service/api/dto"
	_ "user-service/api/helper"
	"user-service/config"
	"user-service/constant"
	"user-service/data/db"
	models "user-service/data/models"
	"user-service/data/vault"
	"user-service/pkg/logging"
	"user-service/pkg/metrics"

	"gorm.io/gorm"
)


type ProfileService struct {
	database *gorm.DB
	logger   logging.Logger
	vault    *vault.VaultClient 
}


func NewProfileService(cfg *config.Config) *ProfileService {
	vc, err := vault.NewVaultClient(cfg.Vault.Address, cfg.Vault.Token)
	if err != nil {
		log.Fatalf("Error initializing Vault client: %v", err)
	}

	return &ProfileService{
		database: db.GetDb(),
		logger:   logging.NewLogger(cfg),
		vault:    vc, // ذخیره کلاینت Vault
	}
}

// تابع جدید برای رمزنگاری داده‌ها
func (s *ProfileService) encryptData(plaintext string) (string, error) {
	ciphertext, err := s.vault.Encrypt(plaintext, "my-key")
	if err != nil {
		return "", fmt.Errorf("Error encrypting data: %v", err)
	}
	return ciphertext, nil
}


func (s *ProfileService) decryptData(ciphertext string) (string, error) {
	decryptedText, err := s.vault.Decrypt(ciphertext, "my-key")
	if err != nil {
		return "", fmt.Errorf("Error decrypting data: %v", err)
	}
	return decryptedText, nil
}


func (s *ProfileService) Create(ctx context.Context, req *dto.CreateProfileRequest) (*dto.ProfileResponse, error) {
	
	// First, encrypt the NID
	encryptedNid, err := s.encryptData(req.NidEncrypted)
	if err != nil {
		return nil, fmt.Errorf("Error encrypting NID: %v", err)
	}

	// First, we check if there is a duplicate NID
	nidExists, err := s.existsByNid(encryptedNid)
	if err != nil {
		return nil, err
	}
	if nidExists {
		return nil, fmt.Errorf("NID '%s' already exists", req.NidEncrypted)
	}

	// Then, we check if there is a duplicate UserId
	userId := int(ctx.Value(constant.UserIdKey).(float64))
	userIdExists, err := s.existsByUserId(userId)
	if err != nil {
		return nil, err
	}
	if userIdExists {
		return nil, fmt.Errorf("The user has already registered their profile.")
	}

	// رمزنگاری فیلدهای مربوط به پروفایل
	nidEncrypted, err := s.encryptData(req.NidEncrypted)
	if err != nil {
		return nil, err
	}

	firstNameEncrypted, err := s.encryptData(req.FirstNameEncrypted)
	if err != nil {
		return nil, err
	}

	lastNameEncrypted, err := s.encryptData(req.LastNameEncrypted)
	if err != nil {
		return nil, err
	}


	dobEncrypted, err := s.encryptData(req.DobEncrypted)
	if err != nil {
		return nil, err
	}

	profile := models.Profile{
		FirstNameEncrypted: firstNameEncrypted,
		LastNameEncrypted: lastNameEncrypted,
		FatherName: req.FatherName,
		NidEncrypted:nidEncrypted,
		DobEncrypted: dobEncrypted,
		UserId: int(ctx.Value(constant.UserIdKey).(float64)),
	}
	profile.CreatedAt = time.Now().UTC()
	tx := s.database.WithContext(ctx).Begin()
	err = tx.
		Create(&profile).
		Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(profile).String(), "Create", "Failed ").Inc()
		return nil, err
	}
	tx.Commit()
	
	metrics.DbCall.WithLabelValues(reflect.TypeOf(profile).String(), "Create", "Success").Inc()
	return s.GetById(ctx)
}


func (s *ProfileService) Update(ctx context.Context, req *dto.UpdateProfileRequest) (*dto.ProfileResponse, error) {
    userId := int(ctx.Value(constant.UserIdKey).(float64))

    isVerified, err := s.isProfileVerified(ctx, userId)
    if err != nil {
        s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
        return nil, err
    }
    if isVerified {
        return nil, fmt.Errorf("Profile is already validated and cannot be updated")
    }

    // Prepare update map
    updateMap := make(map[string]interface{})
    
    // Encrypt only the fields that are provided
    if req.FirstNameEncrypted != "" {
        encryptedFirstName, err := s.encryptData(req.FirstNameEncrypted)
        if err != nil {
            return nil, fmt.Errorf("Error encrypting FirstName: %v", err)
        }
        updateMap["FirstNameEncrypted"] = encryptedFirstName
    }
    
    if req.LastNameEncrypted != "" {
        encryptedLastName, err := s.encryptData(req.LastNameEncrypted)
        if err != nil {
            return nil, fmt.Errorf("Error encrypting LastName: %v", err)
        }
        updateMap["LastNameEncrypted"] = encryptedLastName
    }
    
    if req.FatherName != "" {
        updateMap["FatherName"] = req.FatherName
    }
    
    if req.NidEncrypted != "" {
        encryptedNid, err := s.encryptData(req.NidEncrypted)
        if err != nil {
            return nil, fmt.Errorf("Error encrypting NID: %v", err)
        }
        updateMap["NidEncrypted"] = encryptedNid
    }

    if req.DobEncrypted != "" {
        encryptedDob, err := s.encryptData(req.DobEncrypted)
        if err != nil {
            return nil, fmt.Errorf("Error encrypting DOB: %v", err)
        }
        updateMap["DobEncrypted"] = encryptedDob
    }

    updateMap["modified_by"] = &sql.NullInt64{Int64: int64(userId), Valid: true}
    updateMap["modified_at"] = sql.NullTime{Valid: true, Time: time.Now().UTC()}

    tx := s.database.WithContext(ctx).Begin()
    err = tx.Model(&models.Profile{}).
        Where("user_id = ?", userId).
        Updates(updateMap).
        Error
    if err != nil {
        tx.Rollback()
        s.logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
        metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Profile{}).String(), "Update", "Failed").Inc()
        return nil, err
    }
    
    // Fetch the updated profile
    profile := &models.Profile{}
    err = tx.Model(&models.Profile{}).
        Where("user_id = ? AND deleted_by is null", userId).
        First(&profile).
        Error
    if err != nil {
        tx.Rollback()
        s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
        metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Profile{}).String(), "Select", "Failed").Inc()
        return nil, err
    }
    tx.Commit()

    metrics.DbCall.WithLabelValues(reflect.TypeOf(profile).String(), "Update", "Success").Inc()
    return s.GetById(ctx)
}


func (s *ProfileService) GetById(ctx context.Context) (*dto.ProfileResponse, error) {
	profile := &models.Profile{}
	err := s.database.
		Where("user_id = ? AND deleted_by is null", ctx.Value(constant.UserIdKey)).
		First(&profile).
		Error
	if err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(profile).String(), "GetById", "Failed").Inc()
		return nil, err
	}

	firstName, err := s.decryptData(profile.FirstNameEncrypted)
	if err != nil {
		s.logger.Error(logging.Decrypt, "Error decrypting FirstName: %v", err.Error(),nil)
		return nil, err
	}
	fmt.Errorf("firstName: %v", firstName)

	lastName, err := s.decryptData(profile.LastNameEncrypted)
	if err != nil {
		s.logger.Error(logging.Decrypt, "Error decrypting LastName: %v", err.Error(), nil)
		return nil, err
	}
	fmt.Errorf("lastName: %v", lastName)

	nid, err := s.decryptData(profile.NidEncrypted)
	if err != nil {
		s.logger.Error(logging.Decrypt, "Error decrypting NID: %v", err.Error(), nil)
		return nil, err
	}
	fmt.Errorf("nid: %v", nid)

	dob, err := s.decryptData(profile.DobEncrypted)
	if err != nil {
		s.logger.Error(logging.Decrypt, "Error decrypting Dob: %v", err.Error(), nil)
		return nil, err
	}
	fmt.Errorf("dob: %v", dob)

	dto := &dto.ProfileResponse{
		FirstNameEncrypted: firstName,
		LastNameEncrypted: lastName,
		FatherName: profile.FatherName,
		NidEncrypted: nid,
		DobEncrypted: dob,
	}
	metrics.DbCall.WithLabelValues(reflect.TypeOf(profile).String(), "GetById", "Success").Inc()
	return dto, nil
}

func (s *ProfileService) existsByNid(nid string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.Profile{}).
		Select("count(*) > 0").
		Where("nid_encrypted = ?", nid).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Profile{}).String(), "existsByNid", "Failed ").Inc()
		return false, err
	}
	metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Profile{}).String(), "existsByNid", "Success").Inc()
	return exists, nil
}

func (s *ProfileService) existsByUserId(userId int) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.Profile{}).
		Select("count(*) > 0").
		Where("user_id = ?", userId).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Profile{}).String(), "existsByUserId", "Failed ").Inc()
		return false, err
	}
	metrics.DbCall.WithLabelValues(reflect.TypeOf(&models.Profile{}).String(), "existsByUserId", "Success").Inc()
	return exists, nil
}

func (s *ProfileService) isProfileVerified(ctx context.Context, userId int) (bool, error) {
	profile := &models.Profile{}
	err := s.database.WithContext(ctx).
		Model(&models.Profile{}).
		Where("user_id = ? AND deleted_by is null", userId).
		First(&profile).
		Error
	if err != nil {
		return false, err
	}
	return profile.Verification, nil
}
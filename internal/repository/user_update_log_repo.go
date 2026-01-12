package repository

import (
	"bug-bounty-lite/internal/domain"
	"errors"
	"time"

	"gorm.io/gorm"
)

type userUpdateLogRepo struct {
	db *gorm.DB
}

func NewUserUpdateLogRepo(db *gorm.DB) domain.UserUpdateLogRepository {
	return &userUpdateLogRepo{db: db}
}

func (r *userUpdateLogRepo) Create(log *domain.UserUpdateLog) error {
	return r.db.Create(log).Error
}

func (r *userUpdateLogRepo) FindByUserID(userID uint) ([]domain.UserUpdateLog, error) {
	var logs []domain.UserUpdateLog
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&logs).Error
	return logs, err
}

func (r *userUpdateLogRepo) GetLastUpdateAt(userID uint, field string) (*time.Time, error) {
	var log domain.UserUpdateLog
	err := r.db.Where("user_id = ? AND field = ?", userID, field).Order("created_at desc").First(&log).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &log.CreatedAt, nil
}

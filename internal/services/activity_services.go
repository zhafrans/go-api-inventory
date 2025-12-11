package services

import (
	"time"

	"inventory-api/internal/database"
	"inventory-api/internal/models"

	"gorm.io/gorm"
)

type ActivityService struct {
	db *gorm.DB
}

func NewActivityService() *ActivityService {
	return &ActivityService{
		db: database.DB,
	}
}

func (s *ActivityService) LogActivity(activity *models.ActivityLog) error {
	return s.db.Create(activity).Error
}

func (s *ActivityService) GetAllActivities(page, limit int, activityType, itemID, userID string) ([]models.ActivityLog, int64, error) {
	var activities []models.ActivityLog
	var total int64
	
	query := s.db.Model(&models.ActivityLog{})
	
	if activityType != "" {
		query = query.Where("activity_type = ?", activityType)
	}
	if itemID != "" {
		query = query.Where("item_id = ?", itemID)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	offset := (page - 1) * limit
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&activities).Error
	
	return activities, total, err
}

func (s *ActivityService) GetActivitiesByItemID(itemID string, page, limit int) ([]models.ActivityLog, int64, error) {
	var activities []models.ActivityLog
	var total int64
	
	query := s.db.Where("item_id = ?", itemID)
	
	if err := query.Model(&models.ActivityLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	offset := (page - 1) * limit
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&activities).Error
	
	return activities, total, err
}

func (s *ActivityService) GetRecentActivities(limit int) ([]models.ActivityLog, error) {
	var activities []models.ActivityLog
	
	err := s.db.Order("created_at DESC").
		Limit(limit).
		Find(&activities).Error
	
	return activities, err
}

func (s *ActivityService) GetActivitiesByDateRange(startDate, endDate time.Time, page, limit int) ([]models.ActivityLog, int64, error) {
	var activities []models.ActivityLog
	var total int64
	
	query := s.db.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	
	if err := query.Model(&models.ActivityLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	offset := (page - 1) * limit
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&activities).Error
	
	return activities, total, err
}
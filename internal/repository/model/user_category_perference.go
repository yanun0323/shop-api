package model

import "main/internal/domain/entity"

type UserCategoryPreference struct {
	ID         int64                  `gorm:"column:id;primaryKey;autoIncrement"`
	UserID     int64                  `gorm:"column:user_id"`
	CategoryID entity.ProductCategory `gorm:"column:category_id"`
	UpdatedAt  int64                  `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserCategoryPreference) TableName() string {
	return "user_category_preference"
}

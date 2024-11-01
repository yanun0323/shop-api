package model

import (
	"main/internal/domain/entity"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID          int64                  `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string                 `gorm:"column:name"`
	Description string                 `gorm:"column:description"`
	Category    entity.ProductCategory `gorm:"column:category_id"`
	Price       decimal.Decimal        `gorm:"column:price"`
	Rank        int                    `gorm:"column:rank"`
	CreatedAt   int64                  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   int64                  `gorm:"column:updated_at;autoUpdateTime"`
}

func (Product) TableName() string {
	return "product"
}

func (p Product) ToEntity() *entity.Product {
	return &entity.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CategoryID:  p.Category,
		Price:       p.Price,
		Rank:        p.Rank,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func NewProduct(p entity.Product) *Product {
	return &Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Category:    p.CategoryID,
		Price:       p.Price,
		Rank:        p.Rank,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

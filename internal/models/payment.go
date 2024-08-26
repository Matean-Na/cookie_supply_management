package models

import (
	"cookie_supply_management/pkg/base/base_model"
	"cookie_supply_management/utils/types"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Payment struct {
	base_model.Entity
	StoreID uint            `gorm:"not null" json:"store_id"`
	Store   Store           `gorm:"foreignKey:StoreID;references:ID" json:"store"`
	Date    types.DateOnly  `gorm:"not null" json:"date"`
	Price   decimal.Decimal `gorm:"not null" json:"price"`
}

func (s *Payment) AfterCreate(tx *gorm.DB) (err error) {
	if err := tx.First(&s.Store, s.StoreID).Error; err != nil {
		return err
	}
	s.Store.Debt = s.Store.Debt.Sub(s.Price)
	if err := tx.Save(&s.Store).Error; err != nil {
		return err
	}

	return nil
}

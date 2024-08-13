package models

import (
	"cookie_supply_management/pkg/base/base_model"
	"cookie_supply_management/utils/types"
	"fmt"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Sale struct {
	base_model.Entity
	CookieID     uint            `gorm:"not null" json:"cookie_id"`
	Cookie       Cookie          `gorm:"foreignKey:CookieID;references:ID" json:"cookie"`
	StoreID      uint            `gorm:"not null" json:"store_id"`
	Store        Store           `gorm:"foreignKey:StoreID;references:ID" json:"store"`
	Quantity     int             `gorm:"not null" json:"quantity"`
	Date         types.DateOnly  `gorm:"not null" json:"date"`
	PricePerUnit decimal.Decimal `gorm:"not null" json:"price_per_unit"`
	Price        decimal.Decimal `gorm:"not null" json:"price"`
}

func (s *Sale) BeforeCreate(tx *gorm.DB) (err error) {
	if err := tx.First(&s.Cookie, s.CookieID).Error; err != nil {
		return err
	}
	if err := tx.First(&s.Store, s.StoreID).Error; err != nil {
		return err
	}

	if s.Cookie.Quantity < s.Quantity {
		return fmt.Errorf("Кол-во продоваемого больше чем есть в наличии")
	}

	s.PricePerUnit = s.Cookie.Price
	s.Price = s.PricePerUnit.Mul(decimal.NewFromInt(int64(s.Quantity)))

	return nil
}

func (s *Sale) AfterCreate(tx *gorm.DB) (err error) {
	s.Store.Debt = s.Store.Debt.Add(s.Price)
	if err := tx.Save(&s.Store).Error; err != nil {
		return err
	}

	s.Cookie.Quantity = s.Cookie.Quantity - s.Quantity

	if err := tx.Save(&s.Cookie).Error; err != nil {
		return err
	}

	return nil
}

func (s *Sale) BeforeUpdate(tx *gorm.DB) (err error) {
	// Load old values before updating
	var oldSale Sale
	var store Store
	var cookie Cookie

	if err := tx.Preload(clause.Associations).Preload("Cookie.Type").Preload("Store").First(&oldSale, s.ID).Error; err != nil {
		return err
	}
	store = oldSale.Store
	cookie = oldSale.Cookie

	cookie.Quantity = cookie.Quantity + oldSale.Quantity
	if err := tx.Model(&cookie).Update("Quantity", cookie.Quantity).Error; err != nil {
		return err
	}

	store.Debt = store.Debt.Sub(oldSale.Price)
	if err := tx.Model(&store).Update("Debt", store.Debt).Error; err != nil {
		return err
	}

	s.PricePerUnit = oldSale.PricePerUnit
	s.Price = s.PricePerUnit.Mul(decimal.NewFromInt(int64(s.Quantity)))

	return nil
}

func (s *Sale) AfterUpdate(tx *gorm.DB) (err error) {
	if err := tx.First(&s.Store, s.StoreID).Error; err != nil {
		return err
	}
	fmt.Println(s.Store.Debt)
	s.Store.Debt = s.Store.Debt.Add(s.Price)
	fmt.Println(s.Store.Debt)
	fmt.Println(s.Price)
	if err := tx.Save(&s.Store).Error; err != nil {
		return err
	}

	if err := tx.First(&s.Cookie, s.CookieID).Error; err != nil {
		return err
	}
	s.Cookie.Quantity = s.Cookie.Quantity - s.Quantity
	if err := tx.Save(&s.Cookie).Error; err != nil {
		return err
	}

	return nil
}

package models

import (
	"cookie_supply_management/pkg/base/base_model"
	"github.com/shopspring/decimal"
)

type Cookie struct {
	base_model.Entity
	Name     string          `gorm:"unique;not null"` // название
	TypeID   uint            `gorm:"not null"`
	Type     CookieType      `gorm:"foreignKey:TypeID;references:ID"` // тип печенья
	Quantity int             `gorm:"not null"`                        // Количество коробок
	Price    decimal.Decimal `gorm:"not null"`                        // цена
}

package models

import (
	"cookie_supply_management/pkg/base/base_model"
	"github.com/shopspring/decimal"
)

type Cookie struct {
	base_model.Entity
	Name     string          `gorm:"unique;not null" json:"name"` // название
	TypeID   uint            `gorm:"not null" json:"type_id"`
	Type     CookieType      `gorm:"foreignKey:TypeID;references:ID" json:"type"` // тип печенья
	Quantity int             `gorm:"not null" json:"quantity"`                    // Количество коробок
	Price    decimal.Decimal `gorm:"not null" json:"price"`                       // цена
}

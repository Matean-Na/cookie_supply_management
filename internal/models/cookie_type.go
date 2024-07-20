package models

import (
	"cookie_supply_management/pkg/base/base_model"
	"github.com/shopspring/decimal"
)

type CookieType struct {
	base_model.Entity
	Weight decimal.Decimal `gorm:"not null" json:"weight"` // в кг
}

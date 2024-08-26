package models

import (
	"cookie_supply_management/pkg/base/base_model"
	"github.com/shopspring/decimal"
)

type Store struct {
	base_model.Entity
	Name        string          `json:"name"`
	Address     string          `json:"address"`
	Contact     string          `json:"contact"`
	PhoneNumber string          `json:"phone_number"`
	Debt        decimal.Decimal `json:"debt"`
	Payment     []Payment       `json:"payment"`
}

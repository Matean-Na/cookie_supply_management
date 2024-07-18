package models

import "github.com/shopspring/decimal"

type Store struct {
	Entity
	Name        string          `json:"name"`
	Address     string          `json:"address"`
	Contact     string          `json:"role"`
	PhoneNumber string          `json:"phoneNumber"`
	Debt        decimal.Decimal `json:"debt"`
}

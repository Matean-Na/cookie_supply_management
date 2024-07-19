package base_repository

import (
	"cookie_supply_management/pkg/base/base_model"
	"cookie_supply_management/utils/parsers"
	"gorm.io/gorm"

	"strings"
)

type OrderFilter interface {
	processValues(db *gorm.DB) []string
	sort() func(db *gorm.DB) *gorm.DB
}

type Order struct {
	Values []string
	Model  base_model.HasId
}

func NewOrder(values []string, model base_model.HasId) *Order {
	return &Order{
		Values: values,
		Model:  model,
	}
}

func (o *Order) processValues(db *gorm.DB) []string {
	var results []string
	var order string = "asc"

	for _, value := range o.Values {
		attribute := parsers.ToSnakeCase(value)

		if strings.HasPrefix(attribute, "-") {
			order = "desc"
			attribute = strings.TrimPrefix(attribute, "-")
		}

		if db.Migrator().HasColumn(o.Model, attribute) {
			results = append(results, strings.Join([]string{attribute, order}, " "))
		}
	}

	return results
}

func (o *Order) sort() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		values := o.processValues(db)

		if len(values) == 0 {
			return db
		}
		return db.Order(strings.Join(values, ", "))
	}
}

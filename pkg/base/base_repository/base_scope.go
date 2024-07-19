package base_repository

import "gorm.io/gorm"

type Scope func(*gorm.DB) *gorm.DB

func NoScope(db *gorm.DB) *gorm.DB {
	return db
}

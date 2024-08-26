package database

import (
	"cookie_supply_management/core/config"
	"cookie_supply_management/internal/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDsn(config config.Database) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		config.Host, config.Username, config.Password, config.Name, config.Port, config.SslMode)
}

func Recreate(config config.Database, sysDb string) error {
	sysConf := config
	sysConf.Name = sysDb
	dsn := GetDsn(sysConf)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to %s database", sysConf.Name)
	}

	res := db.Exec(fmt.Sprintf("drop database if exists %s", config.Name))
	if res.Error != nil {
		return fmt.Errorf("failed to drop %s database", config.Name)
	}

	res = db.Exec(fmt.Sprintf("create database %s", config.Name))
	if res.Error != nil {
		return fmt.Errorf("failed to create %s database", config.Name)
	}

	sqlDb, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sqldb")
	}

	if err = sqlDb.Close(); err != nil {
		return fmt.Errorf("failed to close sqldb")
	}

	return nil
}

var Models = []interface{}{
	//auth
	&models.User{},
	&models.Store{},
	&models.Sale{},
	&models.Payment{},

	//cookie
	&models.Cookie{},
	&models.CookieType{},
}

func Connect(config config.Database) (*gorm.DB, error) {

	dsn := GetDsn(config)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: ProNamingStrategy{},
		Logger:         logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s database", config.Name)
	}

	err = database.AutoMigrate(Models...)
	if err != nil {
		return nil, fmt.Errorf("failed auto migration error: %s", err.Error())
	}

	return database, nil
}

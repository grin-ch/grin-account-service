package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Provider struct {
	db *gorm.DB
}

func RegistryDatabase(dsn string) (*Provider, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	providerMigrator(db)
	return &Provider{
		db: db,
	}, nil
}

func providerMigrator(db *gorm.DB) {
	db.AutoMigrate(
		User{},
	)
}

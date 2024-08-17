package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"shortlink/config"
)

var DB *gorm.DB

func ConnectToDatabase() {
	dsn := config.Default("DSN", config.DSN.String())
	DB, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %v", err))
	}
	// Setup sharding
	if config.DefaultBool("ENABLE_SHARDING", config.EnableSharding.Bool()) {
		SetupSharding(DB)
	}
}

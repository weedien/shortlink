package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"shortlink/internal/common/config"
)

func ConnectToDatabase() *gorm.DB {
	dsn := config.Default("DSN", config.DSN.String())
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		},
	})
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %v", err))
	}
	// Setup sharding
	if config.DefaultBool("ENABLE_SHARDING", config.EnableSharding.Bool()) {
		setupSharding(db)
	}
	return db
}

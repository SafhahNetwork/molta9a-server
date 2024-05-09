package db

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error

	var dsn = os.Getenv("POSTGRES_DSN")

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	DB.AutoMigrate(&User{}, &Track{}, &Quest{}, &Module{}, &UserTrack{}, &UserModule{})
}

package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB gorm connector
var DB *gorm.DB

// ConnectDB connects to the database
func ConnectDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("gojobs.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database")
	}

	fmt.Println("Connection Opened to Database")

	// Manually create the users table
	DB.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		age INTEGER,
		birthday DATETIME,
		member_number TEXT,
		activated_at DATETIME,
		created_at DATETIME,
		updated_at DATETIME
	)`)

	// Commented out auto-migrate call
	// DB.AutoMigrate(&models.User{})
	fmt.Println("Database Migrated")
}
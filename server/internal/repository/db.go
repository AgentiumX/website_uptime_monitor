package repository

import (
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"uptime-monitor/server/internal/config"
	"uptime-monitor/server/internal/model"
)

// DB is the global database handle used by all repository functions.
var DB *gorm.DB

// Init opens the SQLite database, runs AutoMigrate, and seeds the admin user.
func Init(cfg *config.DatabaseConfig, adminUser, adminPass string) error {
	// Ensure the parent directory exists.
	dir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Agent{},
		&model.Monitor{},
		&model.AlertChannel{},
		&model.MonitorAlertChannel{},
		&model.AlertHistory{},
	); err != nil {
		return err
	}

	DB = db

	if err := seedAdmin(adminUser, adminPass); err != nil {
		return err
	}

	return nil
}

// seedAdmin creates the admin account only when the users table is empty.
func seedAdmin(username, password string) error {
	var count int64
	DB.Model(&model.User{}).Count(&count)
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := model.User{
		Username:     username,
		PasswordHash: string(hash),
	}
	result := DB.Create(&admin)
	if result.Error != nil {
		log.Printf("[seedAdmin] failed to create admin user: %v", result.Error)
	}
	return result.Error
}

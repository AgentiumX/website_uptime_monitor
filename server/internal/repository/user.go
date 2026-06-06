package repository

import (
	"uptime-monitor/server/internal/model"
)

// FindUserByUsername returns the user with the given username, or an error if not found.
func FindUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

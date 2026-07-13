package services

import (
	"errors"
	"fmt"
	"ui_greenmetric/app/facades"
	"ui_greenmetric/app/models"

	"github.com/goravel/framework/contracts/http"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// Login verifies the user credentials and returns a JWT token and user info
func (s *AuthService) Login(ctx http.Context, email, password string) (string, models.User, error) {
	var user models.User
	err := facades.Orm().Query().Where("email = ?", email).First(&user)
	if err != nil {
		return "", user, errors.New("invalid email or password")
	}

	// Verify hashed password
	if !facades.Hash().Check(password, user.Password) {
		return "", user, errors.New("invalid email or password")
	}

	// Load campus relation
	var campus models.Campus
	if err := facades.Orm().Query().Where("id = ?", user.CampusID).First(&campus); err == nil {
		user.Campus = campus
	}

	// Generate JWT token
	token, err := facades.Auth(ctx).Login(&user)
	if err != nil {
		return "", user, fmt.Errorf("failed to generate token: %v", err)
	}

	return token, user, nil
}

// Logout invalidates the user's current JWT token
func (s *AuthService) Logout(ctx http.Context) error {
	return facades.Auth(ctx).Logout()
}

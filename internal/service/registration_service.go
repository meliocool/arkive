package service

import (
	"context"
	"fmt"
	"github.com/meliocool/arkive/internal/helper"
	"github.com/meliocool/arkive/internal/repository/users"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type RegistrationService struct {
	UserRepository users.UserRepository
}

func NewRegistrationService(userRepository users.UserRepository) *RegistrationService {
	return &RegistrationService{UserRepository: userRepository}
}

func (rs *RegistrationService) Register(ctx context.Context, username string, email string, password string) (*users.User, error) {
	code, codeErr := helper.GenerateVerificationCode()
	if codeErr != nil {
		return nil, fmt.Errorf("failed generating verification code: %w", codeErr)
	}

	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		return nil, fmt.Errorf("error hashing password: %w", hashErr)
	}

	userData := users.User{
		Username:         username,
		Email:            email,
		PasswordHash:     string(hashedPassword),
		IsVerified:       false,
		VerificationCode: code,
	}

	user, createErr := rs.UserRepository.CreateUser(ctx, &userData)
	if createErr != nil {
		return nil, fmt.Errorf("failed to create account: %w", createErr)
	}

	log.Print("Email has been sent!")
	log.Print("User ID:", user.ID)

	return user, nil
}

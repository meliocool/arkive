package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/meliocool/arkive/internal/helper"
	"github.com/meliocool/arkive/internal/repository/users"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type RegistrationService struct {
	UserRepository users.UserRepository
	EmailService   *EmailService
	JwtSecret      string
}

func NewRegistrationService(userRepository users.UserRepository, emailService *EmailService, jwtSecret string) *RegistrationService {
	return &RegistrationService{UserRepository: userRepository, EmailService: emailService, JwtSecret: jwtSecret}
}

func (rs *RegistrationService) VerifyUser(ctx context.Context, email string, verificationCode string) (*users.User, string, error) {
	user, findErr := rs.UserRepository.FindByEmail(ctx, email)
	if findErr != nil {
		return nil, "", findErr
	}

	if user.VerificationCode != verificationCode {
		return nil, "", fmt.Errorf("invalid verification code")
	}

	verifErr := rs.UserRepository.UpdateIsVerified(ctx, user.ID, user.IsVerified)
	if verifErr != nil {
		return nil, "", verifErr
	}

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   user.ID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, tokenErr := token.SignedString([]byte(rs.JwtSecret))
	if tokenErr != nil {
		return nil, "", fmt.Errorf("token failed to generated")
	}

	return user, signedToken, nil
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

	emailErr := rs.EmailService.SendVerificationEmail(user.Email, user.Username, user.VerificationCode, user.CreatedAt)

	if emailErr != nil {
		return nil, fmt.Errorf("failed to send email: %w", emailErr)
	}

	log.Print("Email has been sent!")
	log.Print("User ID:", user.ID)

	return user, nil
}

package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/meliocool/arkive/internal/helper"
	"github.com/meliocool/arkive/internal/repository/users"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type LoginService struct {
	UserRepository users.UserRepository
	JwtSecret      string
}

func NewLoginService(userRepository users.UserRepository, jwtSecret string) *LoginService {
	return &LoginService{UserRepository: userRepository, JwtSecret: jwtSecret}
}

func (ls *LoginService) Login(ctx context.Context, email, password string) (string, error) {
	if email == "" || password == "" {
		return "", fmt.Errorf("invalid input!")
	}

	user, findErr := ls.UserRepository.FindByEmail(ctx, email)
	if findErr != nil {
		return "", findErr
	}

	if user.IsVerified == false {
		return "", helper.ErrUnauthorized
	}

	checkErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if checkErr != nil {
		return "", checkErr
	}

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   user.ID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, tokenErr := token.SignedString([]byte(ls.JwtSecret))
	if tokenErr != nil {
		return "", fmt.Errorf("token failed to generated")
	}
	return signedToken, nil
}

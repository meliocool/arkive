package users

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID               uuid.UUID `json:"id,omitempty" db:"id"`
	Username         string    `json:"username,omitempty" db:"username"`
	Email            string    `json:"email,omitempty" db:"email"`
	PasswordHash     string    `json:"password_hash,omitempty" db:"password_hash"`
	IsVerified       bool      `json:"is_verified,omitempty" db:"is_verified"`
	VerificationCode string    `json:"verification_code,omitempty" db:"verification_code"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	ProfileImageCID  string    `json:"profile_image_cid,omitempty" db:"profile_image_cid"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	UpdateIsVerified(ctx context.Context, id uuid.UUID, isVerified bool) error
}

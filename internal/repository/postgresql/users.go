package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meliocool/arkive/internal/repository/users"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: pool}
}

func (u UserRepo) CreateUser(ctx context.Context, user *users.User) (*users.User, error) {
	SQL := `INSERT INTO users (username, email, password_hash, is_verified, verification_code, profile_image_cid)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, username, email, password_hash, is_verified, verification_code, created_at, updated_at, profile_image_cid
			`

	var newUser users.User

	err := u.db.QueryRow(
		ctx,
		SQL,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.IsVerified,
		user.VerificationCode,
		user.ProfileImageCID,
	).Scan(
		&newUser.ID,
		&newUser.Username,
		&newUser.Email,
		&newUser.PasswordHash,
		&newUser.IsVerified,
		&newUser.VerificationCode,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
		&newUser.ProfileImageCID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user in database: %w", err)
	}

	return &newUser, nil
}

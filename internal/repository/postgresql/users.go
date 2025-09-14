package postgresql

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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

func (u UserRepo) FindByEmail(ctx context.Context, email string) (*users.User, error) {
	if email == "" {
		return nil, fmt.Errorf("invalid email")
	}

	SQL := "SELECT * FROM users WHERE email = $1"

	var userFound users.User

	err := u.db.QueryRow(
		ctx,
		SQL,
		email,
	).Scan(
		&userFound.ID,
		&userFound.Username,
		&userFound.Email,
		&userFound.PasswordHash,
		&userFound.IsVerified,
		&userFound.VerificationCode,
		&userFound.CreatedAt,
		&userFound.UpdatedAt,
		&userFound.ProfileImageCID,
	)

	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &userFound, nil
}

func (u *UserRepo) UpdateIsVerified(ctx context.Context, id uuid.UUID, isVerified bool) error {
	if isVerified == true {
		return fmt.Errorf("account already verified")
	}

	if id == uuid.Nil {
		return fmt.Errorf("invalid id")
	}

	SQL := "UPDATE users SET is_verified = TRUE WHERE id = $1"

	exec, dbErr := u.db.Exec(ctx, SQL, id)
	if dbErr != nil {
		return dbErr
	}

	if exec.RowsAffected() > 0 {
		fmt.Println("Row Updated!")
	}
	return nil
}

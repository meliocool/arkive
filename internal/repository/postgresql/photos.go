package postgresql

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meliocool/arkive/internal/repository/photos"
)

type PhotoRepo struct {
	db *pgxpool.Pool
}

func NewPhotoRepo(pool *pgxpool.Pool) *PhotoRepo {
	return &PhotoRepo{db: pool}
}

func (p *PhotoRepo) Create(ctx context.Context, photo *photos.Photo) (*photos.Photo, error) {
	SQL := `INSERT INTO photos (ipfs_cid, filename, user_id)
			VALUES ($1, $2, $3) RETURNING *`

	var newPhoto photos.Photo

	err := p.db.QueryRow(ctx, SQL, photo.IPFSCid, photo.Filename, photo.UserID).Scan(
		&newPhoto.ID,
		&newPhoto.IPFSCid,
		&newPhoto.Filename,
		&newPhoto.CreatedAt,
		&newPhoto.UpdatedAt,
		&newPhoto.UserID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create photo in database: %w", err)
	}
	return &newPhoto, nil
}

func (p *PhotoRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*photos.Photo, error) {
	SQL := `SELECT * FROM photos WHERE user_id = $1`

	rows, queryErr := p.db.Query(ctx, SQL, userID)
	if queryErr != nil {
		return nil, fmt.Errorf("failed to find user data: %w", queryErr)
	}
	defer rows.Close()

	var Photos []*photos.Photo

	for rows.Next() {
		var photo photos.Photo
		scanErr := rows.Scan(&photo.ID, &photo.IPFSCid, &photo.Filename, &photo.CreatedAt, &photo.UpdatedAt, &photo.UserID)
		if scanErr != nil {
			return nil, fmt.Errorf("failed to retrieve rows: %w", scanErr)
		}
		Photos = append(Photos, &photo)
	}

	if rowErr := rows.Err(); rowErr != nil {
		return nil, rowErr
	}
	return Photos, nil
}

func (p *PhotoRepo) Delete(ctx context.Context, photoID uuid.UUID) error {
	SQL := `DELETE FROM photos WHERE id = $1`
	cmd, execErr := p.db.Exec(ctx, SQL, photoID)
	if execErr != nil {
		return execErr
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("photo does not exist")
	}

	return nil
}

func (p *PhotoRepo) FindAll(ctx context.Context) ([]*photos.Photo, error) {
	SQL := `SELECT * FROM photos`

	rows, queryErr := p.db.Query(ctx, SQL)
	if queryErr != nil {
		return nil, fmt.Errorf("failed to find all photos: %w", queryErr)
	}
	defer rows.Close()

	var Photos []*photos.Photo

	for rows.Next() {
		var photo photos.Photo
		scanErr := rows.Scan(&photo.ID, &photo.IPFSCid, &photo.Filename, &photo.CreatedAt, &photo.UpdatedAt, &photo.UserID)
		if scanErr != nil {
			return nil, fmt.Errorf("failed to retrieve rows: %w", scanErr)
		}
		Photos = append(Photos, &photo)
	}

	if rowErr := rows.Err(); rowErr != nil {
		return nil, rowErr
	}
	return Photos, nil
}

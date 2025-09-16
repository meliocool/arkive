package photos

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Photo struct {
	ID        uuid.UUID
	IPFSCid   string
	Filename  string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
}

type PhotoRepository interface {
	Create(ctx context.Context, photo *Photo) (*Photo, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*Photo, error)
	Delete(ctx context.Context, photoID uuid.UUID) error
}

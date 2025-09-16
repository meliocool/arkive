package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/meliocool/arkive/internal/repository/photos"
	"github.com/meliocool/arkive/internal/repository/users"
)

type PublicService struct {
	PhotoRepository photos.PhotoRepository
	UserRepository  users.UserRepository
}

func NewPublicService(photoRepository photos.PhotoRepository, userRepository users.UserRepository) *PublicService {
	return &PublicService{
		PhotoRepository: photoRepository,
		UserRepository:  userRepository,
	}
}

func (ps *PublicService) FindAll(ctx context.Context) ([]*photos.Photo, error) {
	photoList, findErr := ps.PhotoRepository.FindAll(ctx)
	if findErr != nil {
		return nil, findErr
	}
	return photoList, nil
}

func (ps *PublicService) FindUserProfile(ctx context.Context, userId uuid.UUID) (*users.User, []*photos.Photo, error) {
	user, findUserErr := ps.UserRepository.FindByID(ctx, userId)
	if findUserErr != nil {
		return nil, nil, fmt.Errorf("failure in finding user: %w", findUserErr)
	}
	userPhotos, findPhotosErr := ps.PhotoRepository.FindByUserID(ctx, userId)
	if findPhotosErr != nil {
		return nil, nil, fmt.Errorf("failure in finding photos for this user: %w", findPhotosErr)
	}
	return user, userPhotos, nil
}

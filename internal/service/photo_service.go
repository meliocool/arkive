package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/meliocool/arkive/internal/repository/photos"
	"io"
)

type PhotoService struct {
	PhotoRepository photos.PhotoRepository
	IpfsService     IpfsService
}

func NewPhotoService(photoRepository photos.PhotoRepository, ipfsService IpfsService) *PhotoService {
	return &PhotoService{
		PhotoRepository: photoRepository,
		IpfsService:     ipfsService,
	}
}

func (ps *PhotoService) UploadPhoto(ctx context.Context, userID uuid.UUID, filename string, file io.Reader) (*photos.Photo, error) {
	ipfsCid, uploadErr := ps.IpfsService.UploadFile(ctx, filename, file)
	if uploadErr != nil {
		return nil, uploadErr
	}
	photo := photos.Photo{
		IPFSCid:  ipfsCid,
		Filename: filename,
		UserID:   userID,
	}
	savedPhoto, createErr := ps.PhotoRepository.Create(ctx, &photo)
	if createErr != nil {
		return nil, createErr
	}
	return savedPhoto, nil
}

func (ps *PhotoService) ListPhotos(ctx context.Context, userID uuid.UUID) ([]*photos.Photo, error) {
	photoList, findErr := ps.PhotoRepository.FindByUserID(ctx, userID)
	if findErr != nil {
		return nil, findErr
	}
	return photoList, nil
}

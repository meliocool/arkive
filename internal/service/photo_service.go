package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/meliocool/arkive/internal/helper"
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

func (ps *PhotoService) DeletePhoto(ctx context.Context, userID uuid.UUID, photoID uuid.UUID) error {
	allPhotos, getAllPhotosErr := ps.PhotoRepository.FindByUserID(ctx, userID)
	if getAllPhotosErr != nil {
		return fmt.Errorf("could not find all photos owned by this user: %w", getAllPhotosErr)
	}
	var ipfsCID string
	found := false
	for i := range allPhotos {
		if allPhotos[i].ID == photoID {
			ipfsCID = allPhotos[i].IPFSCid
			found = true
			break
		}
	}
	if !found {
		return helper.ErrNotFound
	}

	if unpinErr := ps.IpfsService.UnpinFile(ctx, ipfsCID); unpinErr != nil {
		return fmt.Errorf("unpin ipfs cid %s: %w", ipfsCID, unpinErr)
	}

	if deleteErr := ps.PhotoRepository.Delete(ctx, photoID); deleteErr != nil {
		return fmt.Errorf("delete photo record failed: %w", deleteErr)
	}
	return nil
}

package handler

import (
	"errors"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/meliocool/arkive/internal/helper"
	"github.com/meliocool/arkive/internal/repository/photos"
	"github.com/meliocool/arkive/internal/repository/users"
	"github.com/meliocool/arkive/internal/service"
	"log"
	"net/http"
)

type PublicHandler struct {
	PublicService service.PublicService
}

type ProfileResponse struct {
	UserInfo   *users.User
	UserPhotos []*photos.Photo
}

func NewPublicHandler(publicService *service.PublicService) *PublicHandler {
	return &PublicHandler{PublicService: *publicService}
}

func (ph *PublicHandler) ListAllPublicPhotos(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ctx := request.Context()
	photos, listErr := ph.PublicService.FindAll(ctx)
	if listErr != nil {
		log.Printf("Error listing all photos: %v", listErr)
		helper.WriteErr(writer, helper.ErrInternal)
		return
	}
	helper.WriteToResponseBody(writer, photos)
}

func (ph *PublicHandler) ViewUserProfile(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ctx := request.Context()
	userId := params.ByName("userId")
	if userId == "" {
		helper.WriteErr(writer, helper.ErrInvalidInput)
		return
	}
	userIDUUID, parseUserErr := uuid.Parse(userId)
	if parseUserErr != nil {
		helper.WriteErr(writer, helper.ErrInvalidInput)
		return
	}
	userInfo, userPhotos, getUserProfileErr := ph.PublicService.FindUserProfile(ctx, userIDUUID)
	if getUserProfileErr != nil {
		switch {
		case errors.Is(getUserProfileErr, helper.ErrNotFound):
			helper.WriteErr(writer, helper.ErrNotFound)
			return
		case errors.Is(getUserProfileErr, helper.ErrUnauthorized):
			helper.WriteErr(writer, helper.ErrUnauthorized)
			return
		default:
			log.Printf("ViewUserProfile error userID=%s: %v", userId, getUserProfileErr)
			helper.WriteErr(writer, helper.ErrInternal)
			return
		}
	}
	helper.WriteToResponseBody(writer, &ProfileResponse{
		UserInfo:   userInfo,
		UserPhotos: userPhotos,
	})
}

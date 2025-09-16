package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/meliocool/arkive/internal/helper"
	"github.com/meliocool/arkive/internal/middleware"
	"github.com/meliocool/arkive/internal/service"
	"log"
	"net/http"
)

type PhotoHandler struct {
	PhotoService service.PhotoService
}

func NewPhotoHandler(photoService service.PhotoService) *PhotoHandler {
	return &PhotoHandler{PhotoService: photoService}
}

func (ph *PhotoHandler) UploadPhoto(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ctx := request.Context()
	userID, ok := ctx.Value(middleware.ContextKeyUserID).(string)
	if !ok {
		helper.WriteErr(writer, helper.ErrUnauthorized)
		return
	}

	userIDUUID, parseUserErr := uuid.Parse(userID)
	if parseUserErr != nil {
		helper.WriteErr(writer, helper.ErrUnauthorized)
		return
	}

	maxMemory := int64(10 << 20)
	if parseErr := request.ParseMultipartForm(maxMemory); parseErr != nil {
		helper.WriteErr(writer, helper.ErrInternal)
		return
	}
	file, fileHeader, fileErr := request.FormFile("file")
	if fileErr != nil {
		helper.WriteErr(writer, helper.ErrInvalidInput)
		return
	}
	defer file.Close()

	newPhoto, uploadErr := ph.PhotoService.UploadPhoto(ctx, userIDUUID, fileHeader.Filename, file)
	if uploadErr != nil {
		log.Printf("Error uploading photo: %v", uploadErr)
		helper.WriteErr(writer, helper.ErrInternal)
		return
	}
	helper.WriteToResponseBody(writer, newPhoto)
}

func (ph *PhotoHandler) ListPhotos(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ctx := request.Context()

	userID, ok := ctx.Value(middleware.ContextKeyUserID).(string)
	if !ok {
		helper.WriteErr(writer, helper.ErrUnauthorized)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		helper.WriteErr(writer, helper.ErrUnauthorized)
		return
	}
	photos, listErr := ph.PhotoService.ListPhotos(ctx, userUUID)
	if listErr != nil {
		log.Printf("Error listing photos: %v", listErr)
		helper.WriteErr(writer, helper.ErrInternal)
		return
	}
	helper.WriteToResponseBody(writer, photos)
}

func (ph *PhotoHandler) DeletePhoto(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	photoId := params.ByName("photoId")
	ctx := request.Context()

	userID, ok := ctx.Value(middleware.ContextKeyUserID).(string)
	if !ok {
		helper.WriteErr(writer, helper.ErrUnauthorized)
		return
	}

	userUUID, userUUIDErr := uuid.Parse(userID)
	if userUUIDErr != nil {
		helper.WriteErr(writer, helper.ErrUnauthorized)
		return
	}

	photoUUID, photoIdErr := uuid.Parse(photoId)
	if photoIdErr != nil {
		helper.WriteErr(writer, helper.ErrUnauthorized)
		return
	}

	deleteErr := ph.PhotoService.DeletePhoto(ctx, userUUID, photoUUID)
	if deleteErr != nil {
		helper.WriteErr(writer, helper.ErrInternal)
		return
	}

	writer.Header().Add("Content-Type", "application/json")

	response := helper.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   "Photo Deleted Successfully!",
	}
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		helper.WriteErr(writer, helper.ErrInternal)
		return
	}
}

func (ph *PhotoHandler) SetProfilePicture(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	photoId := params.ByName("photoId")
	ctx := request.Context()

	userID, ok := ctx.Value(middleware.ContextKeyUserID).(string)
	if !ok {
		helper.WriteErr(writer, helper.ErrUnauthorized)
		return
	}

	userUUID, userUUIDErr := uuid.Parse(userID)
	if userUUIDErr != nil {
		helper.WriteErr(writer, helper.ErrUnauthorized)
		return
	}

	photoUUID, photoIdErr := uuid.Parse(photoId)
	if photoIdErr != nil {
		helper.WriteErr(writer, helper.ErrUnauthorized)
		return
	}

	updateErr := ph.PhotoService.SetProfilePictureCID(ctx, userUUID, photoUUID)
	if updateErr != nil {
		helper.WriteErr(writer, helper.ErrInternal)
		return
	}

	writer.Header().Add("Content-Type", "application/json")

	response := helper.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   "Profile Photo Set!",
	}
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		helper.WriteErr(writer, helper.ErrInternal)
		return
	}
}

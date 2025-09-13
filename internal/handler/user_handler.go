package handler

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/meliocool/arkive/internal/helper"
	"github.com/meliocool/arkive/internal/service"
	"net/http"
)

type UserHandler struct {
	RegistrationService service.RegistrationService
}

type RegisterRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func NewUserHandler(registrationService *service.RegistrationService) *UserHandler {
	return &UserHandler{RegistrationService: *registrationService}
}

func (userHandler *UserHandler) RegisterUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	decoder := json.NewDecoder(request.Body)
	reqBody := RegisterRequest{}
	if err := decoder.Decode(&reqBody); err != nil {
		helper.WriteErr(writer, helper.ErrInvalidInput)
		return
	}
	if reqBody.Password == "" || reqBody.ConfirmPassword == "" || reqBody.Password != reqBody.ConfirmPassword {
		helper.WriteErr(writer, helper.ErrInvalidInput)
		return
	}
	user, regErr := userHandler.RegistrationService.Register(context.Background(), reqBody.Username, reqBody.Email, reqBody.Password)
	if regErr != nil {
		helper.WriteErr(writer, helper.ErrInternal)
		return
	}
	helper.WriteToResponseBody(writer, user)
}

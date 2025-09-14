package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/meliocool/arkive/internal/helper"
	"github.com/meliocool/arkive/internal/service"
	"net/http"
	"time"
)

type UserHandler struct {
	RegistrationService service.RegistrationService
	LoginService        service.LoginService
}

type RegisterRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type RegisterResponse struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Message   string    `json:"message"`
}

type VerifyRequest struct {
	Email            string `json:"email"`
	VerificationCode string `json:"verificationCode"`
}

type VerifyResponse struct {
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Token     string    `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func NewUserHandler(registrationService *service.RegistrationService, loginService *service.LoginService) *UserHandler {
	return &UserHandler{RegistrationService: *registrationService, LoginService: *loginService}
}

func (UserHandler *UserHandler) RegisterUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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
	user, regErr := UserHandler.RegistrationService.Register(context.Background(), reqBody.Username, reqBody.Email, reqBody.Password)
	if regErr != nil {
		helper.WriteErr(writer, helper.ErrInternal)
		return
	}
	writer.Header().Add("Content-Type", "application/json")

	response := helper.WebResponse{
		Code:   http.StatusOK,
		Status: "Registration Successful!",
		Data: RegisterResponse{
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			Message:   "Please Continue with the Account Verification",
		},
	}
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		helper.WriteErr(writer, helper.ErrInternal)
		return
	}
}

func (UserHandler *UserHandler) VerifyUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	decoder := json.NewDecoder(request.Body)
	reqBody := VerifyRequest{}
	if err := decoder.Decode(&reqBody); err != nil {
		helper.WriteErr(writer, helper.ErrInvalidInput)
		return
	}
	if reqBody.Email == "" || reqBody.VerificationCode == "" {
		helper.WriteErr(writer, helper.ErrInvalidInput)
		return
	}

	user, token, verifyErr := UserHandler.RegistrationService.VerifyUser(context.Background(), reqBody.Email, reqBody.VerificationCode)
	if verifyErr != nil {
		helper.WriteErr(writer, helper.ErrInternal)
		return
	}

	writer.Header().Add("Content-Type", "application/json")

	response := helper.WebResponse{
		Code:   http.StatusOK,
		Status: "User Verified Successfully!",
		Data: VerifyResponse{
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Token:     token,
		},
	}
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		helper.WriteErr(writer, helper.ErrInternal)
		return
	}
}

func (UserHandler *UserHandler) LoginUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	decoder := json.NewDecoder(request.Body)
	reqBody := LoginRequest{}
	if err := decoder.Decode(&reqBody); err != nil {
		helper.WriteErr(writer, helper.ErrInvalidInput)
		return
	}
	if reqBody.Email == "" || reqBody.Password == "" {
		helper.WriteErr(writer, helper.ErrInvalidInput)
		return
	}

	token, loginErr := UserHandler.LoginService.Login(context.Background(), reqBody.Email, reqBody.Password)
	if loginErr != nil {
		if errors.Is(loginErr, helper.ErrUnauthorized) {
			helper.WriteErr(writer, helper.ErrUnauthorized)
			return
		}
		helper.WriteErr(writer, helper.ErrInternal)
		return
	}

	writer.Header().Add("Content-Type", "application/json")

	response := helper.WebResponse{
		Code:   http.StatusOK,
		Status: "Login Success!",
		Data: LoginResponse{
			Email: reqBody.Email,
			Token: token,
		},
	}
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		helper.WriteErr(writer, helper.ErrInternal)
		return
	}
}

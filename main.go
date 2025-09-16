package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/meliocool/arkive/config"
	"github.com/meliocool/arkive/internal/handler"
	"github.com/meliocool/arkive/internal/middleware"
	"github.com/meliocool/arkive/internal/repository/postgresql"
	"github.com/meliocool/arkive/internal/service"
	"log"
	"net/http"
)

func wrapRouterHandler(routerHandler httprouter.Handle) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routerHandler(w, r, nil)
	})
}

func main() {
	cfg, cfgErr := config.LoadConfig()
	if cfgErr != nil {
		log.Fatal(cfgErr)
		return
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBName)

	db, dbErr := postgresql.NewPostgresDB(connString)

	if dbErr != nil {
		log.Fatal(dbErr)
		return
	}

	defer db.Close()

	userRepository := postgresql.NewUserRepo(db)
	emailService := service.NewEmailService(cfg.ZohoUser, cfg.ZohoPassword, cfg.ZohoHost, cfg.ZohoPort)
	loginService := service.NewLoginService(userRepository, cfg.JwtSecret)
	registrationService := service.NewRegistrationService(userRepository, emailService, cfg.JwtSecret)
	userHandler := handler.NewUserHandler(registrationService, loginService)
	photoRepository := postgresql.NewPhotoRepo(db)
	ipfsService := service.NewIpfsService(cfg.IPFSAPIKey, cfg.IPFSAPISecret)
	photoService := service.NewPhotoService(photoRepository, userRepository, *ipfsService)
	photoHandler := handler.NewPhotoHandler(*photoService)

	router := httprouter.New()
	router.GET("/health", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		fmt.Fprint(writer, "Server is Up and Running!")
	})
	router.POST("/users/register", userHandler.RegisterUser)
	router.POST("/users/verify", userHandler.VerifyUser)
	router.POST("/users/login", userHandler.LoginUser)
	router.POST("/photos", middleware.AuthMiddleware(photoHandler.UploadPhoto, cfg.JwtSecret))
	router.GET("/photos", middleware.AuthMiddleware(photoHandler.ListPhotos, cfg.JwtSecret))
	router.DELETE("/photos/:photoId", middleware.AuthMiddleware(photoHandler.DeletePhoto, cfg.JwtSecret))
	router.POST("/photos/:photoId", middleware.AuthMiddleware(photoHandler.SetProfilePicture, cfg.JwtSecret))

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	server.ListenAndServe()
}

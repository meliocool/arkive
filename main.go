package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/meliocool/arkive/config"
	"github.com/meliocool/arkive/internal/handler"
	"github.com/meliocool/arkive/internal/repository/postgresql"
	"github.com/meliocool/arkive/internal/service"
	"log"
	"net/http"
)

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
	registrationService := service.NewRegistrationService(userRepository)
	userHandler := handler.NewUserHandler(registrationService)

	router := httprouter.New()
	router.GET("/health", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		fmt.Fprint(writer, "Server is Up and Running!")
	})
	router.POST("/users/register", userHandler.RegisterUser)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	server.ListenAndServe()
}

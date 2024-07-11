package main

import (
	"log"
	"net/http"

	"github.com/Daffc/GO-Sales/api/handler"
	"github.com/Daffc/GO-Sales/api/middleware"
	_ "github.com/Daffc/GO-Sales/docs"
	"github.com/Daffc/GO-Sales/internal/config"
	"github.com/Daffc/GO-Sales/internal/database/mariadb"
	"github.com/Daffc/GO-Sales/repository"
	"github.com/Daffc/GO-Sales/usecase"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title			GO Sales API
// @version			1.0
// @description		This is a basic GO API example constructed for study purposes.
// @BasePath		/
// @schemes			http
// @license.nameGNU	GPL
// @license.url		https://www.gnu.org/licenses/lgpl-3.0.html
// @host			localhost:8080
func main() {

	config, err := config.NewConfigParser()
	if err != nil {
		panic(err)
	}

	mariadb, err := mariadb.NewDatabaseConnection(config)
	if err != nil {
		panic(err)
	}

	usersRepository, err := repository.NewMysqlUserRepository(mariadb)
	if err != nil {
		panic(err)
	}

	usersUseCase := usecase.NewUsersUseCase(usersRepository)
	authUseCase := usecase.NewAuthUseCase(usersRepository, config.Server.JwtSigningKey, config.Server.JwtSessionDuration)

	userHandler := handler.NewUsersHandler(usersUseCase)
	authHandler := handler.NewAuthHandler(authUseCase)

	sm := http.NewServeMux()

	sm.HandleFunc("POST /login", authHandler.Login)
	sm.HandleFunc("POST /users", userHandler.CreateUser)
	sm.HandleFunc("/users", userHandler.ListUsers)
	sm.HandleFunc("/users/{userId}", userHandler.FindUserById)
	sm.Handle("POST /users/{userId}/password", middleware.NewJwtAuthenticator(userHandler.UpdateUserPassword, config.Server.JwtSigningKey))

	sm.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Printf("Linstening on %s ...\n", config.Server.Port)
	err = http.ListenAndServe(":"+config.Server.Port, sm)
	if err != nil {
		log.Println(err)
	}
}

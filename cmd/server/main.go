package main

import (
	"fmt"
	"net/http"

	_ "github.com/Daffc/GO-Sales/docs"
	"github.com/Daffc/GO-Sales/internal/config"
	"github.com/Daffc/GO-Sales/internal/database/mariadb"
	"github.com/Daffc/GO-Sales/internal/domain/repository"
	"github.com/Daffc/GO-Sales/internal/domain/usecase"
	"github.com/Daffc/GO-Sales/internal/handler"
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

	uersUseCase := usecase.NewUsersUseCase(usersRepository)

	userHandler := handler.NewUsersHandler(uersUseCase)

	sm := http.NewServeMux()

	sm.HandleFunc("POST /users", userHandler.CreateUser)
	sm.HandleFunc("/users", userHandler.ListUsers)
	sm.HandleFunc("/users/{userId}", userHandler.FindUserById)

	// Documentation.
	// sm.HandleFunc("GET /swagger/*", httpSwagger.Handler(
	// 	httpSwagger.URL(config.Server.Host+":"+config.Server.Port+"/swagger/doc.json"), //The url pointing to API definition
	// ))

	sm.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	fmt.Printf("Linstening on %s ...\n", config.Server.Port)
	err = http.ListenAndServe(":"+config.Server.Port, sm)
	if err != nil {
		panic(err)
	}
}

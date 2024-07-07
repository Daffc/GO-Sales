package main

import (
	"fmt"
	"net/http"

	"github.com/Daffc/GO-Sales/internal/config"
	"github.com/Daffc/GO-Sales/internal/database/mariadb"
	"github.com/Daffc/GO-Sales/internal/domain/repository"
	"github.com/Daffc/GO-Sales/internal/domain/usecase"
	"github.com/Daffc/GO-Sales/internal/handler"
)

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

	sm.HandleFunc("/users/{userId}", userHandler.GetUserById)

	fmt.Printf("Linstening on %s ...\n", config.ServerPort)
	err = http.ListenAndServe(":"+config.ServerPort, sm)
	if err != nil {
		panic(err)
	}
}

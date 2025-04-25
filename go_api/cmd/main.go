package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	db, err := mariadb.NewDatabaseConnection(&config.Database)
	if err != nil {
		panic(err)
	}

	err = mariadb.RunMigrations(db, config.Database.MigrationsFolderPath)
	if err != nil {
		panic(err)
	}

	usersRepository, err := repository.NewMysqlUserRepository(db)
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

	srv := &http.Server{
		Addr:         config.Server.Port,
		WriteTimeout: time.Second * time.Duration(config.Server.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(config.Server.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(config.Server.IdleTimeout),
		Handler:      sm,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	log.Printf("Linstening on %s ...\n", config.Server.Port)

	// Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	shutdownTimeout := time.Duration(2) * time.Nanosecond
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	srv.Shutdown(ctx)
	log.Printf("shutting down server.")
	os.Exit(0)
}

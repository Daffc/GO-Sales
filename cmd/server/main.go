package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Daffc/GO-Sales/internal/users"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	http.HandleFunc("/users/{userId}", users.GetUserById)
	http.HandleFunc("POST /users", users.PostUser)
	fmt.Printf("Linstening on %s ...\n", os.Getenv("APP_PORT"))
	http.ListenAndServe(":"+os.Getenv("APP_PORT"), nil)
}

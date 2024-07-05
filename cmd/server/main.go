package main

import (
	"fmt"
	"net/http"

	"github.com/Daffc/GO-Sales/internal/users"
)

func main() {

	http.HandleFunc("/users/{userId}", users.GetUserById)
	http.HandleFunc("POST /users", users.PostUser)
	fmt.Println("Linsten....")
	http.ListenAndServe(":8080", nil)
}

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Daffc/GO-Sales/internal/domain/usecase"
)

type UsersHandler struct {
	UsersUseCase *usecase.UsersUseCase
}

func NewUsersHandler(usersUseCase *usecase.UsersUseCase) *UsersHandler {
	return &UsersHandler{UsersUseCase: usersUseCase}
}

func (uh *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateUserInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err := uh.UsersUseCase.CreateUser(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
func (uh *UsersHandler) FindUserById(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("userId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	input := usecase.FindUserInputDTO{ID: userId}

	output, err := uh.UsersUseCase.FindUserById(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

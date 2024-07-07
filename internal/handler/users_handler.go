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

func (uh *UsersHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("userId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	input := usecase.GetUserInputDTO{ID: userId}

	output, err := uh.UsersUseCase.GetUserById(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

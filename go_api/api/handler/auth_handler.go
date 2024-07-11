package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Daffc/GO-Sales/usecase"
)

type AuthHandler struct {
	AuthUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{AuthUseCase: authUseCase}
}

// CreateUser 	Logging User.
// @Summary		Logging User.
// @Description	Logging User.
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Param		input	body	usecase.AuthLoginInputDTO	true	"User credentials"
// @Success		200	{object}	usecase.AuthLoginOutputDTO
// @Failure		500	{object}	string
// @Router		/login [post]
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input usecase.AuthLoginInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err := ah.AuthUseCase.Login(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

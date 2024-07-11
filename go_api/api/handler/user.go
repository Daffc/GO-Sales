package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Daffc/GO-Sales/domain"
	"github.com/Daffc/GO-Sales/domain/dto"
	"github.com/Daffc/GO-Sales/usecase"
)

type UsersHandler struct {
	UsersUseCase *usecase.UsersUseCase
}

func NewUsersHandler(usersUseCase *usecase.UsersUseCase) *UsersHandler {
	return &UsersHandler{UsersUseCase: usersUseCase}
}

// CreateUser 	Create a new user.
// @Summary		Create a new user.
// @Description	Create a new user.
// @Tags		Users
// @Accept		json
// @Produce		json
// @Param		input	body	dto.UserInputDTO true	"User input data"
// @Success		200
// @Failure		500	{object}	string
// @Router		/users [post]
func (uh *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input dto.UserInputDTO
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

// ListUsers 	List all non deleted users.
// @Summary		List all non deleted users.
// @Description	List all non deleted users.
// @Tags		Users
// @Accept		json
// @Produce		json
// @Success		200	{object}	[]dto.UserOutputDTO
// @Failure		500	{object}	string
// @Router		/users [get]
func (uh *UsersHandler) ListUsers(w http.ResponseWriter, r *http.Request) {

	output, err := uh.UsersUseCase.ListUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

// FindUserById Recover user by userId.
// @Summary		Recover user by userId.
// @Description	Recover user by userId.
// @Tags		Users
// @Accept		json
// @Produce		json
// @Param		userId	path		int	true	"User ID"
// @Success		200		{object}	dto.UserOutputDTO
// @Failure		500		{object}	string
// @Router		/users/{userId} [get]
func (uh *UsersHandler) FindUserById(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("userId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	input := dto.UserInputDTO{ID: uint(userId)}

	output, err := uh.UsersUseCase.FindUserById(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

// UpdateUserPassword Update user password.
// @Summary		Update user password.
// @Description	Update user password.
// @Tags		Users
// @Accept		json
// @Produce		json
// @Param		userId	path	int	true	"User ID"
// @Param		input	body	dto.UpdateUserPasswordInputDTO	true	"New user Password"
// @Success		200
// @Failure		500		{object}	string
// @Router		/users/{userId}/password	[post]
func (uh *UsersHandler) UpdateUserPassword(w http.ResponseWriter, r *http.Request, authUser *domain.User) {

	var input dto.UpdateUserPasswordInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId, err := strconv.Atoi(r.PathValue("userId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if uint(userId) != authUser.ID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	input.ID = uint(userId)
	err = uh.UsersUseCase.UpdateUserPassword(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

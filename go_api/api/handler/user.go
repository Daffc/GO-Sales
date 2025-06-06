package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Daffc/GO-Sales/domain/dto"
	"github.com/Daffc/GO-Sales/internal/util"
	"github.com/Daffc/GO-Sales/usecase"
)

type UserHandler struct {
	UserUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{UserUseCase: userUseCase}
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
func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input dto.UserInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Println(err)
		util.JSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := uh.UserUseCase.CreateUser(&input)
	if err != nil {
		log.Println(err)
		util.JSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, output, http.StatusOK)
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
func (uh *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {

	output, err := uh.UserUseCase.ListUsers()
	if err != nil {
		log.Println(err)
		util.JSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, output, http.StatusOK)
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
func (uh *UserHandler) FindUserById(w http.ResponseWriter, r *http.Request) {
	// Extract userId directly from the URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		util.JSONResponse(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	userId, err := strconv.ParseUint(pathParts[len(pathParts)-1], 10, 32)
	if err != nil {
		log.Println(err)
		util.JSONResponse(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	output, err := uh.UserUseCase.FindUserById(uint(userId))
	if err != nil {
		log.Println(err)
		util.JSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, output, http.StatusOK)
}

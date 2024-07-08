package usecase

import (
	"time"

	"github.com/Daffc/GO-Sales/internal/domain/model"
	"github.com/Daffc/GO-Sales/internal/domain/repository"
)

type CreateUserInputDTO struct {
	ID       int
	Name     string
	Email    string
	Password string
}
type FindUserInputDTO struct {
	ID int
}

type UserOutputDTO struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UsersUseCase struct {
	repository *repository.UserRepository
}

func NewUsersUseCase(repository *repository.UserRepository) *UsersUseCase {
	return &UsersUseCase{repository: repository}
}
func (uc UsersUseCase) CreateUser(input CreateUserInputDTO) (*UserOutputDTO, error) {
	nu := model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	err := nu.Validate()
	if err != nil {
		return nil, err
	}

	user, err := uc.repository.CreateUser(&nu)
	if err != nil {
		return nil, err
	}

	userDTO := UserOutputDTO{
		ID:        int(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &userDTO, nil
}

func (uc UsersUseCase) FindUserById(input FindUserInputDTO) (*UserOutputDTO, error) {
	user, err := uc.repository.FindUserById(input.ID)
	if err != nil {
		return nil, err
	}

	userDTO := UserOutputDTO{
		ID:        int(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &userDTO, nil
}

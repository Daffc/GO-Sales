package usecase

import (
	"time"

	"github.com/Daffc/GO-Sales/internal/domain/repository"
)

type GetUserInputDTO struct {
	ID int
}

type GetUserOutputDTO struct {
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

func (uc UsersUseCase) GetUserById(input GetUserInputDTO) (*GetUserOutputDTO, error) {
	user, err := uc.repository.GetUserById(input.ID)
	if err != nil {
		return nil, err
	}

	userDTO := GetUserOutputDTO{
		ID:        int(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &userDTO, nil

}

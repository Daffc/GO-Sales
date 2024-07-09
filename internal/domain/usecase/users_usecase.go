package usecase

import (
	"time"

	"github.com/Daffc/GO-Sales/internal/domain/model"
	"github.com/Daffc/GO-Sales/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserInputDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
	u := model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	err := u.Validate()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 0)
	if err != nil {
		return nil, err
	}

	u.Password = string(hashedPassword)

	user, err := uc.repository.CreateUser(&u)
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

func (uc UsersUseCase) ListUsers() ([]*UserOutputDTO, error) {
	us, err := uc.repository.ListUsers()
	if err != nil {
		return nil, err
	}

	usersDTO := make([]*UserOutputDTO, len(us))

	for i, u := range us {
		usersDTO[i] = &UserOutputDTO{
			ID:        int(u.ID),
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}
	}

	return usersDTO, nil
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

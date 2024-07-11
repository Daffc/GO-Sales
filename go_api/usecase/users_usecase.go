package usecase

import (
	"time"

	"github.com/Daffc/GO-Sales/domain"
	"github.com/Daffc/GO-Sales/repository"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserInputDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type FindUserInputDTO struct {
	ID uint
}

type UserOutputDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateUserPasswordInputDTO struct {
	ID       uint   `json:"id"`
	Password string `json:"password"`
}

type UsersUseCase struct {
	repository *repository.UserRepository
}

func NewUsersUseCase(repository *repository.UserRepository) *UsersUseCase {
	return &UsersUseCase{repository: repository}
}
func (uc UsersUseCase) CreateUser(input CreateUserInputDTO) (*UserOutputDTO, error) {
	u := domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	err := u.ValidateAll()
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
		ID:        user.ID,
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
			ID:        u.ID,
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
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &userDTO, nil
}

func (uc UsersUseCase) UpdateUserPassword(input UpdateUserPasswordInputDTO) error {

	u := &domain.User{
		ID:       uint(input.ID),
		Password: input.Password,
	}

	err := u.ValidatePassword()
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 0)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	err = uc.repository.UpdateUserPassword(u)
	if err != nil {
		return err
	}

	return nil
}

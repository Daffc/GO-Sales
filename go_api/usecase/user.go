package usecase

import (
	"github.com/Daffc/GO-Sales/domain"
	"github.com/Daffc/GO-Sales/domain/dto"
	"github.com/Daffc/GO-Sales/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	CreateUser(input dto.UserInputDTO) (*dto.UserOutputDTO, error)
	ListUsers() ([]*dto.UserOutputDTO, error)
	FindUserById(input dto.UserInputDTO) (*dto.UserOutputDTO, error)
	UpdateUserPassword(input dto.UpdateUserPasswordInputDTO) error
}
type userUseCase struct {
	repository repository.UserRepository
}

func NewUserUseCase(repository repository.UserRepository) UserUseCase {
	return &userUseCase{repository: repository}
}

func (uc *userUseCase) CreateUser(input dto.UserInputDTO) (*dto.UserOutputDTO, error) {
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

	userDTO := dto.UserOutputDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &userDTO, nil
}

func (uc *userUseCase) ListUsers() ([]*dto.UserOutputDTO, error) {
	us, err := uc.repository.ListUsers()
	if err != nil {
		return nil, err
	}

	usersDTO := make([]*dto.UserOutputDTO, len(us))

	for i, u := range us {
		usersDTO[i] = &dto.UserOutputDTO{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}
	}

	return usersDTO, nil
}

func (uc *userUseCase) FindUserById(input dto.UserInputDTO) (*dto.UserOutputDTO, error) {
	user, err := uc.repository.FindUserById(input.ID)
	if err != nil {
		return nil, err
	}

	userDTO := dto.UserOutputDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &userDTO, nil
}

func (uc *userUseCase) UpdateUserPassword(input dto.UpdateUserPasswordInputDTO) error {

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

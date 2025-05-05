package repository

import (
	"errors"
	"time"

	"github.com/Daffc/GO-Sales/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(u *domain.User) (*domain.User, error)
	ListUsers() ([]*domain.User, error)
	FindUserById(id uint) (*domain.User, error)
	FindUserByEmail(email string) (*domain.User, error)
	UpdateUserPassword(u *domain.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewMysqlUserRepository(db *gorm.DB) (UserRepository, error) {
	return &userRepository{db: db}, nil
}

func (r *userRepository) CreateUser(u *domain.User) (*domain.User, error) {

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	result := r.db.Create(u)
	if result.Error != nil {
		return nil, result.Error
	}

	return u, nil
}

func (r *userRepository) ListUsers() ([]*domain.User, error) {
	us := []*domain.User{}

	result := r.db.Find(&us)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return us, nil
}

func (r *userRepository) FindUserById(id uint) (*domain.User, error) {
	u := &domain.User{}

	result := r.db.First(&u, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return u, nil
}

func (r *userRepository) FindUserByEmail(email string) (*domain.User, error) {
	u := &domain.User{}

	result := r.db.First(&u, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}

	return u, nil
}

func (r *userRepository) UpdateUserPassword(u *domain.User) error {

	result := r.db.Model(&u).Where("id = ?", u.ID).Update("password", u.Password)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

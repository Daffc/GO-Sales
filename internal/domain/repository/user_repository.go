package repository

import (
	"errors"
	"time"

	"github.com/Daffc/GO-Sales/internal/domain/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewMysqlUserRepository(db *gorm.DB) (*UserRepository, error) {
	return &UserRepository{db: db}, nil
}

func (r *UserRepository) CreateUser(u *model.User) (*model.User, error) {

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	result := r.db.Create(u)
	if result.Error != nil {
		return nil, result.Error
	}

	return u, nil
}

func (r *UserRepository) ListUsers() ([]*model.User, error) {
	us := []*model.User{}

	result := r.db.Find(&us)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return us, nil
}

func (r *UserRepository) FindUserById(id int) (*model.User, error) {
	u := &model.User{}

	result := r.db.First(&u, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return u, nil
}

func (r *UserRepository) FindUserByEmail(email string) (*model.User, error) {
	u := &model.User{}

	result := r.db.First(&u, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}

	return u, nil
}

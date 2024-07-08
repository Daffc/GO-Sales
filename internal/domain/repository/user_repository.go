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

func (r *UserRepository) FindUserById(id int) (*model.User, error) {
	user := &model.User{}

	result := r.db.First(&user, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

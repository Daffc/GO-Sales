package repository

import (
	"errors"

	"github.com/Daffc/GO-Sales/internal/domain/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewMysqlUserRepository(db *gorm.DB) (*UserRepository, error) {
	return &UserRepository{db: db}, nil
}

func (r *UserRepository) GetUserById(id int) (*model.User, error) {
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

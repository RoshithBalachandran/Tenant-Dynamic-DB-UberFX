package repository

import (
	"errors"
	"tenant-Dynamin-DB/internals/models"

	"gorm.io/gorm"
)

type UserRepo struct{}

func NewUserRepository() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) Create(db *gorm.DB, user *models.User) error {
	return db.Create(user).Error
}

func (r *UserRepo) FindByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user models.User

	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, errors.New("user not registered")
	}

	return &user, nil
}

func (r *UserRepo) GetAll(db *gorm.DB) ([]models.User, error) {
	var users []models.User
	err := db.Find(&users).Error
	return users, err
}

func (r *UserRepo) FindByID(db *gorm.DB, id uint) (*models.User, error) {
	var user models.User

	err := db.First(&user, id).Error
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

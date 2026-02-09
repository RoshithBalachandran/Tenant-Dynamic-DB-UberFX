package service

import (
	"errors"
	"tenant-Dynamin-DB/internals/models"
	"tenant-Dynamin-DB/internals/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(rep *repository.UserRepo) *UserService {
	return &UserService{repo: rep}
}

func (s *UserService) Register(db *gorm.DB, u *models.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}
	pass := string(hash)
	user := &models.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: pass,
		Tenant:   u.Tenant,
	}
	return s.repo.Create(db, user)
}

func (s *UserService) Login(db *gorm.DB, email, pass, tenant string) (*models.User, error) {
	u, err := s.repo.FindByEmailTenant(db, email, tenant)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass))
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) GetProfile(db *gorm.DB, id uint) (*models.User, error) {
	return s.repo.FindByID(db, id)
}
func (s *UserService) UpdateProfile(db *gorm.DB, id uint, tenant string, name, email, password string) (*models.User, error) {

	update := map[string]interface{}{}

	if name != "" {
		update["name"] = name
	}

	if email != "" {
		existing, _ := s.repo.FindByEmailTenant(db, email, tenant)
		if existing != nil && existing.ID != id {
			return nil, errors.New("email already used")
		}
		update["email"] = email
	}

	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		update["password"] = string(hash)
	}

	if len(update) == 0 {
		return nil, errors.New("no fields to update")
	}

	if err := s.repo.UpdateFields(db, id, update); err != nil {
		return nil, err
	}

	return s.repo.FindByID(db, id)
}

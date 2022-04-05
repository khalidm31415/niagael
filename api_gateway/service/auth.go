package service

import (
	"api_gateway/entity"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Signup(username, password, displayName string) error
	Authenticate(username, password string) (*entity.User, error)
	GetByID(ID string) (*entity.User, error)
}

type authService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) AuthService {
	return authService{db: db}
}

func (a authService) Signup(username, password, displayName string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		zap.S().Error(err)
		return err
	}

	user := entity.User{
		Username:     username,
		DisplayName:  displayName,
		PasswordHash: string(hashedPassword),
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	user.ID = id.String()

	if err := a.db.Create(&user).Error; err != nil {
		zap.S().Error(err)
		return err
	}
	return nil
}

func (a authService) Authenticate(username, password string) (*entity.User, error) {
	var user entity.User
	if err := a.db.Take(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("unauthorized")
	}

	return &user, nil
}

func (a authService) GetByID(ID string) (*entity.User, error) {
	var user entity.User
	if err := a.db.Take(&user, "id = ?", ID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

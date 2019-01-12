package domainservice

import (
	"go-binar/user/repository"

	"github.com/satori/go.uuid"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo repository.UserRepository
}

func (s UserService) HashPassword(password string) (string, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(pwd), nil
}

func (s UserService) IsEmailAlreadyExists(email string) (bool, error) {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return false, err
	}

	if user == nil {
		return false, nil
	}

	return true, nil
}

func (s UserService) GenerateToken() string {
	// Just generate a simple access token with UUID
	return uuid.NewV4().String()
}

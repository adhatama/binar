package repository

import (
	"go-binar/user/domain"

	uuid "github.com/satori/go.uuid"
)

type UserRepository interface {
	Save(user domain.User) error
	Login(username string, password string) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
}

type AuthRepository interface {
	FindByUserID(userID uuid.UUID) (*domain.Auth, error)
	Save(auth domain.Auth) error
}

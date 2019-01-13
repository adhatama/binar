package domain

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Name     string    `db:"name" json:"name"`
	Email    string    `db:"email" json:"email"`
	Password string    `db:"password" json:"-"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type UserService interface {
	HashPassword(password string) (string, error)
	IsEmailAlreadyExists(email string) (bool, error)
	GenerateToken() string
}

func NewUser(service UserService, name, email, password string) (*User, error) {
	now := time.Now()

	pwd, err := service.HashPassword(password)
	if err != nil {
		return nil, err
	}

	isExists, err := service.IsEmailAlreadyExists(email)
	if isExists {
		return nil, errors.New("Email has already been taken")
	}

	return &User{
		ID:        uuid.NewV4(),
		Name:      name,
		Email:     email,
		Password:  pwd,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) ChangePassword(password string) error {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(pwd)

	return nil
}

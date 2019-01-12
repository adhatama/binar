package sqlite

import (
	"database/sql"
	"go-binar/user/domain"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepositorySqlite(db *sqlx.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r UserRepository) Save(user domain.User) error {
	_, err := r.DB.Exec(`INSERT INTO user (id, name, email, password, created_at, updated_at)
		VALUES (?, ?, ? ,? ,? ,?)`,
		user.ID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) Login(email string, password string) (*domain.User, error) {
	user := domain.User{}

	err := r.DB.Get(&user, `SELECT * FROM user WHERE email = ?`, email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r UserRepository) FindByEmail(email string) (*domain.User, error) {
	user := domain.User{}

	err := r.DB.Get(&user, `SELECT * FROM user WHERE email = ?`, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

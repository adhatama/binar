package sqlite

import (
	"database/sql"
	"go-binar/user/domain"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

type AuthRepository struct {
	DB *sqlx.DB
}

func NewAuthRepositorySqlite(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (r AuthRepository) FindByUserID(userID uuid.UUID) (*domain.Auth, error) {
	auth := domain.Auth{}

	err := r.DB.Get(&auth, `SELECT * FROM user_auth WHERE user_id = ?`, userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &auth, nil
}

func (r AuthRepository) Save(auth domain.Auth) error {
	_, err := r.DB.Exec(`INSERT INTO user_auth (id, user_id, access_token, expired_at, created_at, updated_at)
		VALUES (?, ?, ? ,? ,? ,?)`,
		auth.ID, auth.UserID, auth.AccessToken, auth.ExpiredAt, auth.CreatedAt, auth.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

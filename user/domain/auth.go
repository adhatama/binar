package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Auth struct {
	ID          uuid.UUID `db:"id" json:"id"`
	UserID      uuid.UUID `db:"user_id" json:"user_id"`
	AccessToken string    `db:"access_token" json:"access_token"`
	ExpiredAt   time.Time `db:"expired_at" json:"expired_at"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func NewAuth(userID uuid.UUID, accessToken string, expiredAt time.Time) (*Auth, error) {
	return &Auth{
		ID:          uuid.NewV4(),
		UserID:      userID,
		AccessToken: accessToken,
		ExpiredAt:   expiredAt,
	}, nil
}

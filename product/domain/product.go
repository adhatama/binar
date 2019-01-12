package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Product struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Name     string    `db:"name" json:"name"`
	Price    int       `db:"price" json:"price"`
	ImageURL string    `db:"imageurl" json:"imageurl"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func NewProduct(name string, price int, imageURL string) (*Product, error) {
	now := time.Now()

	return &Product{
		ID:        uuid.NewV4(),
		Name:      name,
		Price:     price,
		ImageURL:  imageURL,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (p *Product) ChangeInfo(name *string, price *int, imageURL *string) error {
	now := time.Now()

	if name != nil {
		p.Name = *name
		p.UpdatedAt = now
	}
	if price != nil {
		p.Price = *price
		p.UpdatedAt = now
	}
	if imageURL != nil {
		p.ImageURL = *imageURL
		p.UpdatedAt = now
	}

	return nil
}

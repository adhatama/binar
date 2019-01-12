package repository

import (
	"go-binar/product/domain"

	uuid "github.com/satori/go.uuid"
)

type ProductRepository interface {
	FindAll() ([]*domain.Product, error)
	FindByID(id uuid.UUID) (*domain.Product, error)
	Save(product domain.Product) error
	Update(product domain.Product) error
	Delete(product domain.Product) error
}

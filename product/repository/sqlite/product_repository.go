package sqlite

import (
	"database/sql"
	"go-binar/product/domain"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

type ProductRepository struct {
	DB *sqlx.DB
}

func NewProductRepositorySqlite(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r ProductRepository) FindAll() ([]*domain.Product, error) {
	products := []*domain.Product{}

	err := r.DB.Select(&products, `SELECT * FROM product`)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r ProductRepository) FindByID(id uuid.UUID) (*domain.Product, error) {
	product := domain.Product{}

	err := r.DB.Get(&product, `SELECT * FROM product WHERE id = ?`, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &product, nil
}

func (r ProductRepository) Save(product domain.Product) error {
	_, err := r.DB.Exec(`INSERT INTO product (id, name, price, imageurl, created_at, updated_at)
		VALUES (?, ?, ? ,? ,? ,?)`,
		product.ID, product.Name, product.Price, product.ImageURL, product.CreatedAt, product.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r ProductRepository) Update(product domain.Product) error {
	_, err := r.DB.Exec(`UPDATE product SET name = ?, price = ?, imageurl = ?, updated_at = ?
		WHERE id = ?`,
		product.Name, product.Price, product.ImageURL, product.UpdatedAt, product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r ProductRepository) Delete(product domain.Product) error {
	_, err := r.DB.Exec(`DELETE FROM product WHERE id = ?`, product.ID)
	if err != nil {
		return err
	}

	return nil
}

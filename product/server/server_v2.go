package server

import (
	"go-binar/product/repository"
	"go-binar/product/repository/sqlite"
	"go-binar/response"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

type ServerV2 struct {
	ProductRepo repository.ProductRepository
}

func NewServerV2(db *sqlx.DB) (*ServerV2, error) {
	s := ServerV2{
		ProductRepo: sqlite.NewProductRepositorySqlite(db),
	}

	return &s, nil
}

func (s *ServerV2) Mount(g *echo.Group) {
	g.GET("/products", s.GetAllProduct)
}

func (s *ServerV2) GetAllProduct(c echo.Context) error {
	resp := response.Response{}
	resp.Result = "Hello there"

	return c.JSON(http.StatusOK, resp)
}

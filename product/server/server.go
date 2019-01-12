package server

import (
	"errors"
	"fmt"
	"go-binar/product/domain"
	"go-binar/product/repository"
	"go-binar/product/repository/sqlite"
	"net/http"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"

	"github.com/jmoiron/sqlx"
)

type Server struct {
	ProductRepo repository.ProductRepository
}

func NewServer(db *sqlx.DB) (*Server, error) {
	s := Server{
		ProductRepo: sqlite.NewProductRepositorySqlite(db),
	}

	return &s, nil
}

func (s *Server) Mount(g *echo.Group) {
	g.GET("/products", s.GetAllProduct)
	g.GET("/products/:id", s.GetProductByID)
	g.POST("/products", s.SaveProduct)
	g.PUT("/products/:id", s.UpdateProduct)
	g.DELETE("/products/:id", s.DeleteProduct)
}

func (s *Server) GetAllProduct(c echo.Context) error {
	resp := map[string]interface{}{}
	resp["result"] = nil

	products, err := s.ProductRepo.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp["result"] = products

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) GetProductByID(c echo.Context) error {
	resp := map[string]interface{}{}
	resp["result"] = nil

	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	product, err := s.ProductRepo.FindByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp["result"] = product

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) SaveProduct(c echo.Context) error {
	req := CreateProductFormRequest{}
	resp := map[string]interface{}{}
	resp["result"] = nil

	if err := c.Bind(&req); err != nil {
		resp["errors"] = err
		return c.JSON(http.StatusBadRequest, err)
	}

	errs := req.Validate()
	if len(errs) != 0 {
		resp["errors"] = errs
		return c.JSON(http.StatusBadRequest, resp)
	}

	product, err := domain.NewProduct(req.Name, req.Price, req.ImageURL)
	if err != nil {
		fmt.Println("err", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	err = s.ProductRepo.Save(*product)
	if err != nil {
		fmt.Println("err", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp["result"] = product

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) UpdateProduct(c echo.Context) error {
	req := UpdateProductFormRequest{}
	resp := map[string]interface{}{}
	resp["result"] = nil

	if err := c.Bind(&req); err != nil {
		resp["errors"] = err
		return c.JSON(http.StatusBadRequest, err)
	}

	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		resp["errors"] = err
		return c.JSON(http.StatusBadRequest, err)
	}

	product, err := s.ProductRepo.FindByID(id)
	if err != nil {
		fmt.Println("err", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	if product == nil {
		resp["errors"] = errors.New("Product not found")
		return c.JSON(http.StatusBadRequest, err)
	}

	err = product.ChangeInfo(req.Name, req.Price, req.ImageURL)
	if err != nil {
		fmt.Println("err", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	err = s.ProductRepo.Update(*product)
	if err != nil {
		fmt.Println("err", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp["result"] = product

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) DeleteProduct(c echo.Context) error {
	resp := map[string]interface{}{}
	resp["result"] = nil

	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		resp["errors"] = err
		return c.JSON(http.StatusBadRequest, err)
	}

	product, err := s.ProductRepo.FindByID(id)
	if err != nil {
		fmt.Println("err", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	if product == nil {
		resp["errors"] = errors.New("Product not found")
		return c.JSON(http.StatusBadRequest, err)
	}

	err = s.ProductRepo.Delete(*product)
	if err != nil {
		fmt.Println("err", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp["result"] = "success"

	return c.JSON(http.StatusOK, resp)
}

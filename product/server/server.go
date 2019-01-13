package server

import (
	"go-binar/product/domain"
	"go-binar/product/repository"
	"go-binar/product/repository/sqlite"
	"go-binar/response"
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
	resp := response.Response{}

	products, err := s.ProductRepo.FindAll()
	if err != nil {
		resp.Errors = map[string]interface{}{
			"message": err.Error(),
		}
		return response.JSON(c, http.StatusInternalServerError, resp)
	}

	resp.Result = products

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) GetProductByID(c echo.Context) error {
	resp := response.Response{}

	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		resp.Errors = map[string]interface{}{
			"message": err.Error(),
		}
		return response.JSON(c, http.StatusBadRequest, resp)
	}

	product, err := s.ProductRepo.FindByID(id)
	if err != nil {
		resp.Errors = map[string]interface{}{
			"message": err.Error(),
		}
		return response.JSON(c, http.StatusInternalServerError, resp)
	}

	if product == nil {
		resp.Errors = map[string]interface{}{
			"message": "Product not found",
		}
		return response.JSON(c, http.StatusBadRequest, resp)
	}

	resp.Result = product

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) SaveProduct(c echo.Context) error {
	req := CreateProductFormRequest{}
	resp := response.Response{}

	if err := c.Bind(&req); err != nil {
		resp.Errors = map[string]interface{}{
			"message": err.Error(),
		}
		return response.JSON(c, http.StatusBadRequest, resp)
	}

	errs := req.Validate()
	if len(errs) != 0 {
		resp.Errors = errs
		return response.JSON(c, http.StatusInternalServerError, resp)
	}

	product, err := domain.NewProduct(req.Name, req.Price, req.ImageURL)
	if err != nil {
		resp.Errors = map[string]interface{}{
			"message": err.Error(),
		}
		return response.JSON(c, http.StatusBadRequest, resp)
	}

	err = s.ProductRepo.Save(*product)
	if err != nil {
		resp.Errors = map[string]interface{}{
			"message": err.Error(),
		}
		return response.JSON(c, http.StatusInternalServerError, resp)
	}

	resp.Result = product

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) UpdateProduct(c echo.Context) error {
	req := UpdateProductFormRequest{}
	resp := response.Response{}

	if err := c.Bind(&req); err != nil {
		resp.Errors = map[string]interface{}{
			"message": err.Error(),
		}
		return response.JSON(c, http.StatusBadRequest, resp)
	}

	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		resp.Errors = map[string]interface{}{
			"message": err.Error(),
		}
		return response.JSON(c, http.StatusBadRequest, resp)
	}

	product, err := s.ProductRepo.FindByID(id)
	if err != nil {
		resp.Errors = map[string]interface{}{
			"message": err.Error(),
		}
		return response.JSON(c, http.StatusInternalServerError, resp)
	}

	if product == nil {
		resp.Errors = map[string]interface{}{
			"message": "Product not found",
		}
		return response.JSON(c, http.StatusBadRequest, resp)
	}

	err = product.ChangeInfo(req.Name, req.Price, req.ImageURL)
	if err != nil {
		resp.Errors = map[string]interface{}{
			"message": err.Error(),
		}
		return response.JSON(c, http.StatusBadRequest, resp)
	}

	err = s.ProductRepo.Update(*product)
	if err != nil {
		resp.Errors = map[string]interface{}{
			"message": err.Error(),
		}
		return response.JSON(c, http.StatusInternalServerError, resp)
	}

	resp.Result = product

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) DeleteProduct(c echo.Context) error {
	resp := response.Response{}

	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		resp.Errors = map[string]interface{}{
			"message": err.Error(),
		}
		return response.JSON(c, http.StatusBadRequest, resp)
	}

	product, err := s.ProductRepo.FindByID(id)
	if err != nil {
		resp.Errors = map[string]interface{}{
			"message": err.Error(),
		}
		return response.JSON(c, http.StatusInternalServerError, resp)
	}

	if product == nil {
		resp.Errors = map[string]interface{}{
			"message": "Product not found",
		}
		return response.JSON(c, http.StatusBadRequest, resp)
	}

	err = s.ProductRepo.Delete(*product)
	if err != nil {
		resp.Errors = map[string]interface{}{
			"message": err.Error(),
		}
		return response.JSON(c, http.StatusInternalServerError, resp)
	}

	resp.Result = "success"

	return c.JSON(http.StatusOK, resp)
}

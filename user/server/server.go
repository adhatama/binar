package server

import (
	"errors"
	"fmt"
	"go-binar/user/domain"
	"go-binar/user/domainservice"
	"go-binar/user/repository"
	"go-binar/user/repository/sqlite"
	"net/http"
	"time"

	"github.com/labstack/echo"

	"github.com/jmoiron/sqlx"
)

type Server struct {
	UserRepo    repository.UserRepository
	UserService domainservice.UserService
	AuthRepo    repository.AuthRepository
}

func NewServer(db *sqlx.DB) (*Server, error) {
	userRepo := sqlite.NewUserRepositorySqlite(db)
	authRepo := sqlite.NewAuthRepositorySqlite(db)

	userService := domainservice.UserService{
		UserRepo: userRepo,
	}

	return &Server{
		UserRepo:    userRepo,
		UserService: userService,
		AuthRepo:    authRepo,
	}, nil
}

func (s *Server) Mount(g *echo.Group) {
	g.POST("/auth/signup", s.Signup)
	g.POST("/auth/login", s.Login)
}

func (s *Server) Signup(c echo.Context) error {
	req := SignupFormRequest{}
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

	user, err := domain.NewUser(s.UserService, req.Name, req.Email, req.Password)
	if err != nil {
		resp["errors"] = err
		return c.JSON(http.StatusBadRequest, resp)
	}

	err = s.UserRepo.Save(*user)
	if err != nil {
		resp["errors"] = errs
		return c.JSON(http.StatusInternalServerError, resp)
	}

	resp["result"] = user

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) Login(c echo.Context) error {
	req := LoginFormRequest{}
	resp := map[string]interface{}{}
	resp["result"] = nil

	if err := c.Bind(&req); err != nil {
		fmt.Println(err)
		resp["errors"] = err
		return c.JSON(http.StatusBadRequest, err)
	}

	errs := req.Validate()
	if len(errs) != 0 {
		fmt.Println(errs)
		resp["errors"] = errs
		return c.JSON(http.StatusBadRequest, resp)
	}

	user, err := s.UserRepo.Login(req.Email, req.Password)
	if err != nil {
		fmt.Println(err)
		resp["errors"] = err
		return c.JSON(http.StatusInternalServerError, resp)
	}

	if user == nil {
		fmt.Println("error user")
		resp["errors"] = errors.New("User not found")
		return c.JSON(http.StatusBadRequest, resp)
	}

	accessToken := s.UserService.GenerateToken()
	expiredAt := time.Now().Add(24 * 30 * time.Hour) // Expired in 30 days

	auth, err := domain.NewAuth(user.ID, accessToken, expiredAt)
	if err != nil {
		fmt.Println(err)
		resp["errors"] = err
		return c.JSON(http.StatusBadRequest, resp)
	}

	err = s.AuthRepo.Save(*auth)
	if err != nil {
		fmt.Println(err)
		resp["errors"] = err
		return c.JSON(http.StatusBadRequest, resp)
	}

	resp["result"] = map[string]string{
		"access_token": accessToken,
	}

	return c.JSON(http.StatusOK, resp)
}

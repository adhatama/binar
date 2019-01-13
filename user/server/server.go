package server

import (
	"go-binar/response"
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
		return response.JSON(c, http.StatusBadRequest, resp)
	}

	user, err := domain.NewUser(s.UserService, req.Name, req.Email, req.Password)
	if err != nil {
		resp.Errors = map[string]interface{}{
			"message": err.Error(),
		}
		return response.JSON(c, http.StatusBadRequest, resp)
	}

	err = s.UserRepo.Save(*user)
	if err != nil {
		resp.Errors = map[string]interface{}{
			"message": err.Error(),
		}
		return response.JSON(c, http.StatusInternalServerError, resp)
	}

	resp.Result = user

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) Login(c echo.Context) error {
	req := LoginFormRequest{}
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
		return response.JSON(c, http.StatusBadRequest, resp)
	}

	user, err := s.UserRepo.Login(req.Email, req.Password)
	if err != nil {
		resp.Errors = map[string]interface{}{
			"message": err.Error(),
		}
		return response.JSON(c, http.StatusInternalServerError, resp)
	}

	if user == nil {
		resp.Errors = map[string]interface{}{
			"message": "Invalid email or password",
		}
		return response.JSON(c, http.StatusBadRequest, resp)
	}

	authQueryResult, err := s.AuthRepo.FindByUserID(user.ID)
	if err != nil {
		resp.Errors = map[string]interface{}{
			"message": err.Error(),
		}
		return response.JSON(c, http.StatusInternalServerError, resp)
	}

	accessToken := s.UserService.GenerateToken()
	expiredAt := time.Now().Add(24 * 30 * time.Hour) // Expired in 30 days

	if authQueryResult != nil {
		accessToken = authQueryResult.AccessToken
		expiredAt = authQueryResult.ExpiredAt
	} else {
		auth, err := domain.NewAuth(user.ID, accessToken, expiredAt)
		if err != nil {
			resp.Errors = map[string]interface{}{
				"message": err.Error(),
			}
			return response.JSON(c, http.StatusInternalServerError, resp)
		}

		err = s.AuthRepo.Save(*auth)
		if err != nil {
			resp.Errors = map[string]interface{}{
				"message": err.Error(),
			}
			return response.JSON(c, http.StatusInternalServerError, resp)
		}
	}

	resp.Result = map[string]string{
		"access_token": accessToken,
	}

	return c.JSON(http.StatusOK, resp)
}

package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	productserver "go-binar/product/server"
	"go-binar/response"
	userserver "go-binar/user/server"

	log "github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
)

func main() {
	db := initSqlite()

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	RootApiMiddlewares := []echo.MiddlewareFunc{
		requestLogging(),
	}
	APIMiddlewares := []echo.MiddlewareFunc{
		tokenValidation(db),
		requestLogging(),
	}

	root := e.Group("", RootApiMiddlewares...)

	userServer, err := userserver.NewServer(db)
	if err != nil {
		panic(err)
	}
	userServer.Mount(root)

	api := e.Group("/api/v1", APIMiddlewares...)

	productServer, err := productserver.NewServer(db)
	if err != nil {
		panic(err)
	}
	productServer.Mount(api)

	apiV2 := e.Group("/api/v2", APIMiddlewares...)

	productServerV2, err := productserver.NewServerV2(db)
	if err != nil {
		panic(err)
	}
	productServerV2.Mount(apiV2)

	e.Logger.Fatal(e.Start(":1323"))
}

func initSqlite() *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", "./database.db?_fk=true")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func tokenValidation(db *sqlx.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			resp := response.Response{}

			authorization := c.Request().Header.Get("Authorization")

			if authorization == "" {
				resp.Errors = map[string]interface{}{
					"message": "Unauthorized",
				}
				return response.JSON(c, http.StatusUnauthorized, resp)
			}

			id := ""
			err := db.Get(&id, `SELECT user_id
				FROM user_auth WHERE access_token = ?`, authorization)
			if err != nil {
				resp.Errors = map[string]interface{}{
					"message": "Unauthorized",
				}
				return response.JSON(c, http.StatusUnauthorized, resp)
			}

			userID, err := uuid.FromString(id)
			if err != nil {
				resp.Errors = map[string]interface{}{
					"message": err.Error(),
				}
				return response.JSON(c, http.StatusInternalServerError, resp)
			}

			c.Set("USER_ID", userID)

			return next(c)
		}
	}
}

func requestLogging() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			start := time.Now()

			if err := next(c); err != nil {
				c.Error(err)
			}

			res := c.Response()
			stop := time.Now()

			fields := map[string]interface{}{
				"request_id":      res.Header().Get(echo.HeaderXRequestID),
				"ip":              c.RealIP(),
				"host":            req.Host,
				"uri":             req.RequestURI,
				"method":          req.Method,
				"user_agent":      req.UserAgent(),
				"status":          res.Status,
				"roundtrip":       strconv.FormatInt(stop.Sub(start).Nanoseconds()/1000, 10),
				"roundtrip_human": stop.Sub(start).String(),
			}

			// We will add log Form Values and Query String if...
			if res.Status == http.StatusInternalServerError {
				if !strings.HasPrefix(req.Header.Get(echo.HeaderContentType), echo.MIMEMultipartForm) {
					qs := c.QueryString()

					forms, err := c.FormParams()
					if err != nil {
						c.Error(err)
					}

					fields["query_string"] = qs
					fields["form_values"] = forms
				}
			}

			log.WithFields(fields).Info()

			return nil
		}
	}
}

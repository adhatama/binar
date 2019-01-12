package main

import (
	"log"
	"net/http"

	productserver "go-binar/product/server"
	userserver "go-binar/user/server"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
)

func main() {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	db := initSqlite()

	APIMiddlewares := []echo.MiddlewareFunc{
		tokenValidation(db),
	}

	root := e.Group("")

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
			authorization := c.Request().Header.Get("Authorization")

			if authorization == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"result": "Unauthorized"})
			}

			id := ""
			err := db.Get(&id, `SELECT user_id
				FROM user_auth WHERE access_token = ?`, authorization)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"result": "Unauthorized"})
			}

			userID, err := uuid.FromString(id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]error{"result": err})
			}

			c.Set("USER_ID", userID)

			return next(c)
		}
	}
}

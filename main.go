package main

import (
	"log"

	productserver "go-binar/product/server"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	db := initSqlite()

	api := e.Group("/api/v1")

	productServer, err := productserver.NewServer(db)
	if err != nil {
		panic(err)
	}
	productServer.Mount(api)

	apiV2 := e.Group("/api/v2")

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

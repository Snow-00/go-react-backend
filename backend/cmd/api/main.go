package main

import (
	"flag"
	"log"

	"github.com/Snow-00/go-react-movies-backend/internal/repository"
	"github.com/Snow-00/go-react-movies-backend/internal/repository/dbrepo"
)

const PORT = "8080"

type Application struct {
	Domain string
	DSN    string // data source name
	DB repository.DatabaseRepo
}

func main() {
	// set app config
	var app Application

	// read from command line, think still better .env
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string") // timeout 5 secs
	flag.Parse()

	// connect to db
	conn, err := app.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}

	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()
	
	app.Domain = "example.com"

	// start web server
	app.Routes().Run(":" + PORT)
}

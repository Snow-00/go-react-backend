package main

import (
	"flag"
	"log"
	"time"

	"github.com/Snow-00/go-react-movies-backend/internal/repository"
	"github.com/Snow-00/go-react-movies-backend/internal/repository/dbrepo"
)

const PORT = "8080"

type Application struct {
	Domain       string
	DSN          string // data source name
	DB           repository.DatabaseRepo
	Auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
}

func main() {
	// set app config
	var app Application

	// read from command line, think still better .env
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string") // timeout 5 secs
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "Signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "Signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "Signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "Cookie domain")
	flag.StringVar(&app.Domain, "domain", "example.com", "Domain")
	flag.Parse()

	// connect to db
	conn, err := app.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}

	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	app.Auth = Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "refreshToken",
		CookieDomain:  app.CookieDomain,
	}

	// start web server
	app.Routes().Run(":" + PORT)
}

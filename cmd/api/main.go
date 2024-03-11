package main

import (
	"flag"
	"log"
	"time"
	"os"

	"github.com/Snow-00/go-react-movies-backend/internal/repository"
	"github.com/Snow-00/go-react-movies-backend/internal/repository/dbrepo"
)

const PORT = os.Getenv("PORT")

type Application struct {
	Domain       string
	DSN          string // data source name
	DB           repository.DatabaseRepo
	Auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	APIKey       string
}

func main() {
	// set app config
	var app Application

	// read from command line, think still better .env
	flag.StringVar(&app.DSN, "dsn", "host=viaduct.proxy.rlwy.net port=34159 user=postgres password=mNwgjZHHkhwJKAZDTkfMoPClLyFACHtJ dbname=railway sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string") // timeout 5 secs
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "Signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "go-react-frontend.vercel.app", "Signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "go-react-frontend.vercel.app", "Signing audience")
	// flag.StringVar(&app.CookieDomain, "cookie-domain", "127.0.0.1", "Cookie domain")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "", "Cookie domain")
	flag.StringVar(&app.Domain, "domain", "go-react-backend-production.up.railway.app", "Domain")
	flag.StringVar(&app.APIKey, "api-key", "8561321331c4cb04c25726ac2e38328a", "api key")
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
		TokenExpiry:   time.Minute * 10,
		RefreshExpiry: time.Minute * 30,
		CookiePath:    "/",
		CookieName:    "refreshToken",
		CookieDomain:  app.CookieDomain,
	}

	// start web server
	app.Routes().Run("0.0.0.0:" + PORT)
}

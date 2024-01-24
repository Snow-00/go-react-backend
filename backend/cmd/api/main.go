package main

import "flag"

const PORT = "8080"

type Application struct {
	Domain string
	DSN    string // data source name
}

func main() {
	// set app config
	var app Application

	// read from command line, think still better .env
	flag.StringVar(&app.DSN, "dsn", "HOST=localhost PORT=5432 USER=postgres PASSWORD=postgres DBNAME=movies SSLMODE=disable TIMEZONE=UTC CONNECT_TIMEOUT=5", "Postgres connection string") // timeout 5 secs
	flag.Parse()

	// connect to db
	app.Domain = "example.com"

	// start web server
	app.Routes().Run(":" + PORT)
}

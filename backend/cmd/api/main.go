package main

const PORT = "8000"

type Application struct {
	Domain string
}

func main() {
	// set app config
	var app Application

	// read from command line

	// connect to db
	app.Domain = "example.com"

	// start web server
	app.Routes().Run(":" + PORT)
}

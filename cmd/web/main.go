package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/jumbhood/bookings/pkg/config"
	"github.com/jumbhood/bookings/pkg/handlers"
	"github.com/jumbhood/bookings/pkg/render"
)

// port number of the HTTP application
const portNumber = ":8080"

// application config
var app config.AppConfig

// session of the request
var session *scs.SessionManager

func main() {
	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	// on production will need to be true
	session.Cookie.Secure = app.InProduction
	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create teamplate cache")
	}
	log.Println("Templates created in cache")

	app.TemplateCache = tc

	// using debug mode can turn off the this
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	log.Println("Starting the application on port", portNumber)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"net/http"

	"github.com/mclyashko/everlock/internal/config"
	"github.com/mclyashko/everlock/internal/db"
	"github.com/mclyashko/everlock/internal/web"
)

func main() {
	config := config.LoadConfig()

	db.LoadDbPool(&config.Db)

	web.ConfigureRouter()

	log.Printf("Starting Everlock on http://localhost:%s\n", config.Web.Port)
	log.Fatal(http.ListenAndServe(":"+config.Web.Port, nil))
}

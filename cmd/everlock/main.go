package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mclyashko/everlock/internal/config"
	"github.com/mclyashko/everlock/internal/db"
)

func main() {
	config := config.LoadConfig()

	db.LoadDbPool(&config.Db)

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, config)
	})

	log.Printf("Starting Everlock on http://localhost:%s\n", config.Web.Port)
	log.Fatal(http.ListenAndServe(":"+config.Web.Port, nil))
}

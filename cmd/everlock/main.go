package main

import (
	"fmt"
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to Everlock!")
}

func main() {
	http.HandleFunc("/", homeHandler)

	port := "8080"
	log.Printf("Starting Everlock on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

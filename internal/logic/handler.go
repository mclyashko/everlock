package logic

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/hoisie/mustache"
)

// предоставляет доступ к шаблону главной страницы
func MainPageHandler(w http.ResponseWriter, _ *http.Request) {
	template, err := mustache.ParseFile(filepath.Join("..", "..", "internal", "template", "index.html"))
	if err != nil {
		log.Printf("Error loading main page template, error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	renderedTemplate := template.Render()

	w.Header().Set("Content-Type", "text/html")

	_, err = w.Write([]byte(renderedTemplate))
	if err != nil {
		log.Printf("Error writing main page response, error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

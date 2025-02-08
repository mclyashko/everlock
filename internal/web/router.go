package web

import (
	"net/http"

	"github.com/mclyashko/everlock/internal/logic"
)

// настраивает HTTP роутер, задавая пути и их обработчики
func ConfigureRouter() {
	http.HandleFunc("/", logic.MainPageHandler)
}

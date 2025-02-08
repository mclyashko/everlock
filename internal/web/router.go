package web

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mclyashko/everlock/internal/config"
	"github.com/mclyashko/everlock/internal/logic"
)

// настраивает HTTP роутер, задавая пути и их обработчики
func ConfigureRouter(_ *config.App, p *pgxpool.Pool) {
	http.HandleFunc("/", logic.MainPageHandler)
	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		logic.SubmitMessageHandler(p, w, r)
	})
}

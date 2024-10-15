package Handlers

import (
	"net/http"
)

func HandleLbFrontend(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./src/web/leaderboards.html")
}

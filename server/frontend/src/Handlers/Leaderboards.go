package Handlers

import (
	"net/http"
)

func HandleLbFrontend(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/src/web/leaderboards.html")
}

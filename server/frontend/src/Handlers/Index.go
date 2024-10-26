package Handlers

import (
	"net/http"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/src/web/index.html")
}

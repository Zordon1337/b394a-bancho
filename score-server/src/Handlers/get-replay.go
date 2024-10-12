package Handlers

import (
	"net/http"
)

func HandleReplay(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("c")
	filePath := "./replays/" + id + ".osr"
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}

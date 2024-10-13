package Handlers

import (
	"net/http"
)

func HandleAvatar(w http.ResponseWriter, r *http.Request) {
	filePath := "./avatar.png"
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}

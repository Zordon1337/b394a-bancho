package web

import (
	"net/http"
	"os"
)

func HandleAvatar(w http.ResponseWriter, r *http.Request) {
	filePath := "./Avatars/" + r.URL.Query().Get("avatar") + ".png"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		filePath = "./avatar.png"
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}

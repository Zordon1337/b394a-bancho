package web

import (
	"fmt"
	"net/http"
	"os"
)

func HandleAvatar(w http.ResponseWriter, r *http.Request) {
	filePath := "./frontend/Avatars/" + r.URL.Query().Get("avatar") + ".png"
	fmt.Println(filePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		filePath = "./frontend/avatar.png"
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}

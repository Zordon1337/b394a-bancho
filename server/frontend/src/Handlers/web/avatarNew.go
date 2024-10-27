package web

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func HandleAvatarNew(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	filePath := "./frontend/Avatars/" + id + ".png"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		filePath = "./frontend/avatar.png"
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}

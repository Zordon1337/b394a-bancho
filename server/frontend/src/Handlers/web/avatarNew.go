package web

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func HandleAvatarNew(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("Before replacement:", id)

	id = strings.Replace(id, "_000.png", "", -1)
	id = strings.Replace(id, "_000.jpg", "", -1)

	fmt.Println("After replacement:", id)

	filePath := "./frontend/Avatars/" + id + ".png"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		filePath = "./frontend/avatar.png"
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}

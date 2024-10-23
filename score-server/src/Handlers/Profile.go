package Handlers

import (
	"net/http"
)

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./src/web/profile.html")
}

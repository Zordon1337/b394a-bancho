package Handlers

import (
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/src/web/login.html")
}

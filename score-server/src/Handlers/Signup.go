package Handlers

import (
	"net/http"
)

func HandleSignup(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./src/web/register.html")
}

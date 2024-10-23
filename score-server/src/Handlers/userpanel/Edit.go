package userpanel

import (
	"net/http"
)

func HandleEdit(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./src/web/editprofile.html")
}

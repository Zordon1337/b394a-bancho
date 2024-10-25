package api

import (
	"fmt"
	"net/http"
	"score-server/src/Utils"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	ses, err := Utils.Sessions.Get(r, "usersession")
	if err != nil {
		fmt.Fprintf(w, "No")
		return
	}
	ses.Options.MaxAge = -1
	err = ses.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to destroy session", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

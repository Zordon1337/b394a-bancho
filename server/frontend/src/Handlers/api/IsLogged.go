package api

import (
	"fmt"
	"net/http"
	"retsu/Utils"
	"retsu/shared/db"
)

func HandleIsLogged(w http.ResponseWriter, r *http.Request) {
	ses, err := Utils.Sessions.Get(r, "usersession")
	if err != nil {
		fmt.Fprintf(w, "No")
		return
	}
	if username, ok := ses.Values["username"].(string); ok && username != "" {
		if password, ok := ses.Values["password"].(string); ok && password != "" {
			if db.IsCorrectCred(username, password) {
				fmt.Fprintf(w, "Yes")
			} else {
				fmt.Fprintf(w, "No")
				ses.Options.MaxAge = -1
				err := ses.Save(r, w)
				if err != nil {
					http.Error(w, "Failed to destroy session", http.StatusInternalServerError)
					return
				}
			}
			return
		} else {
			fmt.Fprintf(w, "No")
			return
		}
	} else {

		fmt.Fprintf(w, "No")
		return
	}
}

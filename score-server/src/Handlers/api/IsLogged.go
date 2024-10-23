package api

import (
	"fmt"
	"net/http"
	"score-server/src/Utils"
	"score-server/src/db"
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

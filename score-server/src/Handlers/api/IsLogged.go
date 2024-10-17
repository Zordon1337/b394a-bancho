package api

import (
	"fmt"
	"net/http"
	"score-server/src/Utils"
)

func HandleIsLogged(w http.ResponseWriter, r *http.Request) {
	ses, err := Utils.Sessions.Get(r, "usersession")
	if err != nil {
		fmt.Fprintf(w, "No")
		return
	}
	if username, ok := ses.Values["username"].(string); ok && username != "" {
		if password, ok := ses.Values["password"].(string); ok && password != "" {
			fmt.Fprintf(w, "Yes")
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

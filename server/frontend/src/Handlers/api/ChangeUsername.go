package api

import (
	"fmt"
	"net/http"
	"retsu/Utils"
	"retsu/shared/db"
)

func ChangeUsername(w http.ResponseWriter, r *http.Request) {
	ses, err := Utils.Sessions.Get(r, "usersession")
	if err != nil {
		return
	}
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	newusername := r.FormValue("newusername")
	if newusername == "" {
		http.Error(w, "No username provided", http.StatusBadRequest)
	}
	if username, ok := ses.Values["username"].(string); ok && username != "" {
		if password, ok := ses.Values["password"].(string); ok && password != "" {
			if db.IsCorrectCred(username, password) {
				result := db.SetUsername(username, newusername)
				if result {
					ses.Values["username"] = newusername
					ses.Options.MaxAge = -1
					err := ses.Save(r, w)
					if err != nil {
						http.Error(w, "Failed to destroy session", http.StatusInternalServerError)
						return
					}
					fmt.Fprintf(w, "Changed username!")
					return
				} else {

					http.Error(w, "Failed to change username, already taken?", http.StatusBadRequest)
					return
				}
			} else {
				ses.Options.MaxAge = -1
				err := ses.Save(r, w)
				if err != nil {
					http.Error(w, "Failed to destroy session", http.StatusInternalServerError)
					return
				}

			}
			return
		} else {
			return
		}
	} else {
		return
	}
}

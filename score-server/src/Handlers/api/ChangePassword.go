package api

import (
	"fmt"
	"net/http"
	"score-server/src/Utils"
	"score-server/src/db"
)

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	ses, err := Utils.Sessions.Get(r, "usersession")
	if err != nil {
		return
	}
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	oldpass := r.FormValue("oldpass")
	newpass := Utils.HashMD5(r.FormValue("newpass"))
	if newpass == "" || oldpass == "" {
		http.Error(w, "No password provided", http.StatusBadRequest)
	}
	if username, ok := ses.Values["username"].(string); ok && username != "" {
		if password, ok := ses.Values["password"].(string); ok && password != "" {
			if db.IsCorrectCred(username, password) {
				result := db.SetPassword(username, password, newpass)
				if result && oldpass == password {
					ses.Options.MaxAge = -1
					ses.Values = make(map[interface{}]interface{})
					err := ses.Save(r, w)
					if err != nil {
						http.Error(w, "Failed to destroy session", http.StatusInternalServerError)
						return
					}
					fmt.Fprintf(w, "Succesfully changed password")
					return
				} else {

					http.Error(w, "Failed to change password", http.StatusBadRequest)
					return
				}
			} else {

			}
			return
		} else {
			return
		}
	} else {
		return
	}
}

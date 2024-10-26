package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"retsu/Utils"
	"retsu/shared/db"
)

type ResponseJSON struct {
	LoggedIn bool   `json:"loggedIn"`
	Username string `json:"username"`
}

func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	ses, err := Utils.Sessions.Get(r, "usersession")
	if err != nil {
		fmt.Fprintf(w, "No")
		return
	}

	if username, ok := ses.Values["username"].(string); ok && username != "" {
		if password, ok := ses.Values["password"].(string); ok && password != "" {

			user := ResponseJSON{
				LoggedIn: true,
				Username: username,
			}
			jsonData, err := json.MarshalIndent(user, "", "  ")
			if err != nil {
				fmt.Println("Error marshaling to JSON:", err)
				return
			}
			if !db.IsCorrectCred(username, password) {
				user.LoggedIn = false
			}
			fmt.Fprint(w, string(jsonData))
			return
		} else {
			user := ResponseJSON{
				LoggedIn: false,
				Username: "",
			}
			jsonData, err := json.MarshalIndent(user, "", "  ")
			if err != nil {
				fmt.Println("Error marshaling to JSON:", err)
				return
			}
			fmt.Fprint(w, string(jsonData))
			return
		}
	} else {

		fmt.Fprintf(w, "No")
		return
	}
}

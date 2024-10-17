package api

import (
	"fmt"
	"net/http"
	"score-server/src/Utils"
	"score-server/src/db"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	ses, err := Utils.Sessions.Get(r, "usersession")
	if err != nil {
		fmt.Fprintf(w, "-1")
		return
	}
	if r.URL.Query().Get("u") != "" && r.URL.Query().Get("p") != "" {
		user := r.URL.Query().Get("u")
		pass := r.URL.Query().Get("p")
		res := db.IsCorrectCred(user, Utils.HashMD5(pass))
		if res {
			ses.Values["username"] = user
			ses.Values["password"] = pass
		} else {

			delete(ses.Values, "username")
			delete(ses.Values, "password")
			fmt.Fprintf(w, "ERR\nWrong username or password")
			return
		}
		ses.Save(r, w)
	}
	fmt.Fprintf(w, "SUCCESS")
}

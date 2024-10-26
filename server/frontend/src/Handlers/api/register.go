package api

import (
	"fmt"
	"net/http"
	"retsu/Utils"
	"retsu/shared/db"
)

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	ses, err := Utils.Sessions.Get(r, "usersession")
	if err != nil {
		fmt.Fprintf(w, "-1")
		return
	}
	if r.URL.Query().Get("u") != "" && r.URL.Query().Get("p") != "" {
		user := r.URL.Query().Get("u")
		pass := r.URL.Query().Get("p")
		err := db.RegisterUser(user, pass)
		if err != nil {
			fmt.Fprintf(w, "ERR\n%s", err.Error())
			return
		}
		ses.Values["username"] = user
		ses.Values["password"] = Utils.HashMD5(pass)
		ses.Save(r, w)
	}
	fmt.Fprintf(w, "SUCCESS")
}

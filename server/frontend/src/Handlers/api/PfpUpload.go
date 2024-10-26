package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"retsu/Utils"
	"retsu/shared/db"
	"strconv"
)

func PfpUpload(w http.ResponseWriter, r *http.Request) {
	ses, err := Utils.Sessions.Get(r, "usersession")
	if err != nil {
		return
	}
	if username, ok := ses.Values["username"].(string); ok && username != "" {
		if password, ok := ses.Values["password"].(string); ok && password != "" {
			if db.IsCorrectCred(username, password) {
				err := r.ParseMultipartForm(10 << 20)
				if err != nil {
					fmt.Print("Unable to parse form", http.StatusBadRequest)
					return
				}
				file, _, err := r.FormFile("avatar")
				if err != nil {
					Utils.LogErr("Unable to get avatar from form")
					return
				}
				defer file.Close()
				destFile, err := os.Create("./frontend/Avatars/" + strconv.Itoa(int(db.GetUserIdByUsername(username))) + ".png")
				if err != nil {
					Utils.LogErr("Unable to create destination file")
					return
				}
				defer destFile.Close()

				_, err = io.Copy(destFile, file)
				if err != nil {
					Utils.LogErr("Failed to save avatar")
					return
				}
			}
		}
	}

}

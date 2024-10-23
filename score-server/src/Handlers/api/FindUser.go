package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"score-server/src/db"

	"github.com/gorilla/mux"
)

func HandleFindUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["user"]
	if id == "" {
		fmt.Fprint(w, "No user provided")
		return
	}
	ser := db.GetUserFromDatabase(id)
	jsonData, err := json.MarshalIndent(ser, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}
	fmt.Fprint(w, string(jsonData))
}

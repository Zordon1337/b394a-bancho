package Handlers

import (
	"fmt"
	"net/http"
	"score-server/src/db"
)

func HandleScores(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, db.GetScores(r.URL.Query().Get("c")))
}

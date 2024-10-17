package web

import (
	"fmt"
	"net/http"
	"score-server/src/db"
)

func HandleScores(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, db.GetMapStatus(r.URL.Query().Get("c")))
	if db.IsRanked(r.URL.Query().Get("c")) {
		fmt.Fprintf(w, db.GetScores(r.URL.Query().Get("c")))
	}
}

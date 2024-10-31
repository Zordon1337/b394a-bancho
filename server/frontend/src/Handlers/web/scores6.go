package web

import (
	"fmt"
	"net/http"
	"retsu/shared/db"
)

func HandleScores6(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, db.GetMapStatus(r.URL.Query().Get("c")))
	if db.IsRanked(r.URL.Query().Get("c")) {
		fmt.Fprintf(w, db.GetScores6(r.URL.Query().Get("c"), r.URL.Query().Get("u"), r.URL.Query().Get("m")))
	}
}

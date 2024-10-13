package main

import (
	"net/http"
	"score-server/src/Handlers"
	"score-server/src/Utils"

	"github.com/gorilla/mux"
)

func main() {
	Utils.LogInfo("Starting score server on port 80")
	r := mux.NewRouter()
	r.HandleFunc("/web/osu-submit.php", Handlers.HandleScore).Methods("GET", "POST")
	r.HandleFunc("/web/osu-getreplay.php", Handlers.HandleReplay).Methods("GET", "POST")
	r.HandleFunc("/web/osu-getscores3.php", Handlers.HandleScores).Methods("GET", "POST")
	r.HandleFunc("/forum/download.php", Handlers.HandleAvatar).Methods("GET", "POST")
	http.ListenAndServe(":8080", r)
}

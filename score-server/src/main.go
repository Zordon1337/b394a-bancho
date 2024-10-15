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
	r.HandleFunc("/api/v1/gettop", Handlers.HandleLbs).Methods("GET", "POST")
	r.HandleFunc("/leaderboard", Handlers.HandleLbFrontend).Methods("GET", "POST")
	r.HandleFunc("/", Handlers.HandleIndex).Methods("GET", "POST")
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./src/web/css/"))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./src/web/js/"))))
	r.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("./src/web/img/"))))
	http.ListenAndServe(":8080", r)
}

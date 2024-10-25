package main

import (
	"net/http"
	"score-server/src/Handlers"
	"score-server/src/Handlers/api"
	"score-server/src/Handlers/userpanel"
	"score-server/src/Handlers/web"
	"score-server/src/Utils"
	"score-server/src/db"

	"github.com/gorilla/mux"
)

func main() {
	Utils.LogInfo("Starting score server on port 80")
	db.InitDatabase()
	r := mux.NewRouter()
	r.HandleFunc("/web/osu-submit.php", web.HandleScore).Methods("GET", "POST")
	r.HandleFunc("/web/osu-getreplay.php", web.HandleReplay).Methods("GET", "POST")
	r.HandleFunc("/web/osu-getscores3.php", web.HandleScores).Methods("GET", "POST")
	r.HandleFunc("/forum/download.php", web.HandleAvatar).Methods("GET", "POST")
	r.HandleFunc("/api/v1/gettop", api.HandleLbs).Methods("GET", "POST")
	r.HandleFunc("/api/v1/GetUser", api.HandleGetUser).Methods("GET", "POST")
	r.HandleFunc("/api/v1/isloggedin", api.HandleIsLogged).Methods("GET", "POST")
	r.HandleFunc("/api/v1/register", api.HandleRegister).Methods("GET", "POST")
	r.HandleFunc("/api/v1/login", api.HandleLogin).Methods("GET", "POST")
	r.HandleFunc("/api/v1/findplayer/{user}", api.HandleFindUser).Methods("GET", "POST")
	r.HandleFunc("/leaderboard", Handlers.HandleLbFrontend).Methods("GET", "POST")
	r.HandleFunc("/", Handlers.HandleIndex).Methods("GET", "POST")
	r.HandleFunc("/register", Handlers.HandleSignup).Methods("GET", "POST")
	r.HandleFunc("/login", Handlers.HandleLogin).Methods("GET", "POST")
	r.HandleFunc("/profile/{user}", Handlers.HandleProfile).Methods("GET", "POST")
	r.HandleFunc("/userpanel/edit", userpanel.HandleEdit).Methods("GET", "POST")
	r.HandleFunc("/api/v1/userpanel/UpdatePassword", api.ChangePassword).Methods("GET", "POST")
	r.HandleFunc("/api/v1/userpanel/UpdateUsername", api.ChangeUsername).Methods("GET", "POST")
	r.HandleFunc("/api/v1/userpanel/SetAvatar", api.PfpUpload).Methods("GET", "POST")
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./src/web/css/"))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./src/web/js/"))))
	r.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("./src/web/img/"))))
	http.ListenAndServe(":8080", r)
}

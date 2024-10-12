package Handlers

import (
	"fmt"
	"net/http"
)

func HandleScore(w http.ResponseWriter, r *http.Request) {
	score := r.URL.Query().Get("score")
	password := r.URL.Query().Get("pass")
	fmt.Println("got score submit: %s %s", score, password)
}

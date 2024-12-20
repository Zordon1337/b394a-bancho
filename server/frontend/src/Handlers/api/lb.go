package api

import (
	"encoding/json"
	"net/http"
	"retsu/shared/db"
)

func HandleLbs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	topUsers, err := db.GetTopUsers()
	if err != nil {
		http.Error(w, "Failed to get lb", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(topUsers)
}

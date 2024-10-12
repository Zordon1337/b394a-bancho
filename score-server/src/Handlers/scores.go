package Handlers

import (
	"fmt"
	"net/http"
)

func HandleScores(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "\n1|ZRD|1|1|1|1|1|1|1|1|False|0|1|1")
}

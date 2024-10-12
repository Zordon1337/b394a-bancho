package Handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"score-server/src/Utils"
)

func HandleScore(w http.ResponseWriter, r *http.Request) {
	score := r.URL.Query().Get("score")
	//password := r.URL.Query().Get("pass")
	score2 := Utils.FormattedToScore(score)
	fmt.Println(Utils.CalculateAccuracy(score2))
	fmt.Printf("got score submit from %s with acc %f", score2.Username, Utils.CalculateAccuracy(score2))
	if true {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			fmt.Print("Unable to parse form", http.StatusBadRequest)
			return
		}
		file, _, err := r.FormFile("score")
		if err != nil {
			fmt.Print("Unable to get file from form", http.StatusBadRequest)
			return
		}
		defer file.Close()
		destFile, err := os.Create("./test.osr")
		if err != nil {
			fmt.Print("Unable to create destination file", http.StatusInternalServerError)
			return
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, file)
		if err != nil {
			fmt.Print("Failed to save file", http.StatusInternalServerError)
			return
		}

	}
}

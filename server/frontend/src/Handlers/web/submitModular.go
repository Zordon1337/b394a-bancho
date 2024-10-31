package web

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"retsu/Utils"
	"retsu/shared/db"
	"strconv"
	"strings"
)

func HandleModularScore(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Print("Unable to parse form", http.StatusBadRequest)
		return
	}
	score := r.FormValue("score")
	password := r.FormValue("pass")
	if score == "" || password == "" {
		return
	}
	score2 := Utils.FormattedToScore(score, true)

	if db.IsCorrectCred(score2.Username, password) {
		if score2.Pass != "True" {
			return
		}
		if strings.Split(score, ":")[15] != "0" {
			fmt.Fprintf(w, "error: disabled")
		}
		if db.IsRestricted(db.GetUserIdByUsername(score2.Username)) {
			fmt.Fprintf(w, "error: ban")
		}
		if !db.IsRanked(score2.FileChecksum) {
			fmt.Fprintf(w, "error: beatmap")
		}
		Utils.LogInfo("%s has submitted score", score2.Username)
		scoreid := db.GetNewScoreId()

		file, _, err := r.FormFile("score")
		if err != nil {
			Utils.LogErr("Unable to get file from form")
			return
		}
		defer file.Close()
		destFile, err := os.Create("./frontend/replays/" + strconv.Itoa(int(scoreid)) + ".osr")
		if err != nil {
			Utils.LogErr("Unable to create destination file")
			return
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, file)
		if err != nil {
			Utils.LogErr("Failed to save replay")
			return
		}
		db.InsertScore(score2, scoreid)
		db.UpdateRankedScore(score2.Username)
		db.UpdatePlaycount(score2.Username)
		db.UpdateTotalScore(score2.Username)
		db.UpdateAccuracy(score2.Username)
		db.UpdateRank(score2.Username)
	} else {

		fmt.Fprintf(w, "error: pass")
	}
}

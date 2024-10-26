package web

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"retsu/Utils"
	"retsu/shared/db"
	"strconv"
)

func HandleScore(w http.ResponseWriter, r *http.Request) {
	score := r.URL.Query().Get("score")
	password := r.URL.Query().Get("pass")
	score2 := Utils.FormattedToScore(score)
	Utils.LogInfo("%s has submitted score", score2.Username)

	if db.IsCorrectCred(score2.Username, password) && !db.IsRestricted(db.GetUserIdByUsername(score2.Username)) && score2.Pass == "True" && db.IsRanked(score2.FileChecksum) {
		scoreid := db.GetNewScoreId()
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			fmt.Print("Unable to parse form", http.StatusBadRequest)
			return
		}
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
	}
}

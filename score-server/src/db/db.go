package db

import (
	"database/sql"
	"score-server/src/Utils"

	_ "github.com/go-sql-driver/mysql"
)

var connectionstring string = "test:test@tcp(127.0.0.1:3306)/osu!"

func IsCorrectCred(username string, password string) bool {
	db, err := sql.Open("mysql", connectionstring)
	if err != nil {
		Utils.LogErr(err.Error())
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM users WHERE username = ? AND password = ?", username, password)
	if err != nil {
		Utils.LogErr(err.Error())
	}
	rowsreturned := 0
	defer rows.Close()
	for rows.Next() {
		rowsreturned++
	}
	return rowsreturned > 0
}
func GetUserIdByUsername(username string) int32 {
	db, err := sql.Open("mysql", connectionstring)
	if err != nil {
		Utils.LogErr(err.Error())
	}
	defer db.Close()
	rows, err := db.Query("SELECT userid FROM users WHERE username = ?", username)
	if err != nil {
		Utils.LogErr(err.Error())
	}

	defer rows.Close()
	for rows.Next() {
		var userid int32
		err := rows.Scan(userid)
		if err != nil {
			return -1
		}
		return userid
	}
	return -1
}
func IsRestricted(userid int32) bool {
	db, err := sql.Open("mysql", connectionstring)
	if err != nil {
		Utils.LogErr(err.Error())
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM restricts WHERE userid = ?", userid)
	if err != nil {
		Utils.LogErr(err.Error())
	}
	rowsreturned := 0
	defer rows.Close()
	for rows.Next() {
		rowsreturned++
	}
	return rowsreturned > 0
}
func GetNewScoreId() int32 {
	db, err := sql.Open("mysql", connectionstring)
	if err != nil {
		Utils.LogErr(err.Error())
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM scores")
	if err != nil {
		Utils.LogErr(err.Error())
	}
	rowsreturned := 0
	defer rows.Close()
	for rows.Next() {
		rowsreturned++
	}
	return int32(rowsreturned) + 1
}
func InsertScore(score Utils.Score, scoreid int32) {
	db, err := sql.Open("mysql", connectionstring)
	if err != nil {
		Utils.LogErr(err.Error())
		return
	}
	defer db.Close()
	rows, err := db.Query("INSERT INTO `scores`(`scoreid`, `mapchecksum`,`username`, `OnlineScoreChecksum`, `Count300`, `Count100`, `Count50`, `CountGeki`, `CountKatu`, `CountMiss`, `TotalScore`, `MaxCombo`, `Perfect`, `Ranking`, `EnabledMods`, `Pass`, `Accuracy`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		scoreid, score.FileChecksum, score.Username, score.OnlineScoreChecksum, score.Count300, score.Count100, score.Count50, score.CountGeki, score.CountKatu, score.CountMiss, score.TotalScore, score.MaxCombo, score.Perfect, score.Ranking, score.EnabledMods, score.EnabledMods, Utils.CalculateAccuracy(score))
	defer rows.Close()
}

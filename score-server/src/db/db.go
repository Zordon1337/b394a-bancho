package db

import (
	"database/sql"
	"log"
	"score-server/src/Utils"

	"fmt"

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

func GetScores(mapchecksum string) string {
	var response string
	db, err := sql.Open("mysql", connectionstring) // Assuming `connectionstring` is defined elsewhere
	if err != nil {
		log.Printf("Error connecting to the database: %s", err)
		return ""
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT 
			s.scoreid,
			s.mapchecksum,
			s.Username,
			s.OnlineScoreChecksum,
			s.Count300,
			s.Count100,
			s.Count50,
			s.CountGeki,
			s.CountKatu,
			s.CountMiss,
			s.TotalScore,
			s.MaxCombo,
			s.Perfect,
			s.Ranking,
			s.EnabledMods,
			s.Pass
		FROM scores AS s
		INNER JOIN (
			SELECT Username, MAX(TotalScore) AS MaxScore
			FROM scores
			WHERE mapchecksum = ?
			GROUP BY Username
		) AS maxScores ON s.Username = maxScores.Username AND s.TotalScore = maxScores.MaxScore
		WHERE s.mapchecksum = ?
		ORDER BY s.TotalScore DESC;`, mapchecksum, mapchecksum)

	if err != nil {
		log.Printf("Error executing query: %s", err)
		return ""
	}
	defer rows.Close()
	rowam := 1
	for rows.Next() {
		var score Utils.Score
		var scoreid int32
		err := rows.Scan(
			&scoreid,
			&score.FileChecksum,
			&score.Username,
			&score.OnlineScoreChecksum,
			&score.Count300,
			&score.Count100,
			&score.Count50,
			&score.CountGeki,
			&score.CountKatu,
			&score.CountMiss,
			&score.TotalScore,
			&score.MaxCombo,
			&score.Perfect,
			&score.Ranking,
			&score.EnabledMods,
			&score.Pass,
		)
		if err != nil {
			Utils.LogErr("Error scanning row: %s", err)
			continue
		}

		response += fmt.Sprintf(
			"\n%d|%s|%d|%d|%d|%d|%d|%d|%d|%d|%t|%s|%d|test",
			scoreid,
			score.Username,
			score.TotalScore,
			score.MaxCombo,
			score.Count50,
			score.Count100,
			score.Count300,
			score.CountMiss,
			score.CountKatu,
			score.CountGeki,
			score.Perfect,
			score.EnabledMods,
			rowam,
		)
		rowam++
	}
	return response
}
func UpdateRankedScore(user string) {
	db, err := sql.Open("mysql", connectionstring) // Assuming `connectionstring` is defined elsewhere
	if err != nil {
		log.Printf("Error connecting to the database: %s", err)
		return
	}
	defer db.Close()

	rows, err := db.Query(`
    SELECT 
        s.mapchecksum,
        MAX(s.TotalScore) AS MaxScore
    FROM scores AS s
    WHERE s.Username = ?
    GROUP BY s.mapchecksum
    ORDER BY MaxScore DESC;`, user)

	if err != nil {
		log.Printf("Error executing query: %s", err)
		return
	}
	defer rows.Close()
	rowam := 1
	rankedscore := 0
	for rows.Next() {
		var score Utils.Score
		err := rows.Scan(
			&score.FileChecksum,
			&score.TotalScore,
		)
		if err != nil {
			Utils.LogErr("Error scanning row: %s", err)
			continue
		}
		rankedscore += int(score.TotalScore)
		rowam++
	}
	db.Query("UPDATE `users` SET `ranked_score`= ? WHERE username = ?", rankedscore, user)
}

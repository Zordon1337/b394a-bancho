package db

import (
	"database/sql"
	"log"
	"score-server/src/Utils"
	"time"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var connectionstring string = "test:test@tcp(127.0.0.1:3306)/osu!"
var db *sql.DB

func InitDatabase() {
	var err error
	db, err = sql.Open("mysql", connectionstring)
	if err != nil {
		Utils.LogErr("Failed to create database connection!")
	}
}
func RegisterUser(username string, password string) error {

	if IsNameTaken(username) {
		return fmt.Errorf("User already taken!")
	}
	hashedPassword := Utils.HashMD5(password)
	stmt, err := db.Prepare(`INSERT INTO users (userid, username, password, ranked_score, accuracy, playcount, total_score, rank, lastonline, joindate)
                             VALUES (?,?, ?, 0, 0.0, 0, 0, 0, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err = stmt.Exec(GetNewUserId(), username, hashedPassword, now, now)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}

	return nil
}
func IsNameTaken(username string) bool {

	rows, err := db.Query("SELECT * FROM users WHERE username = ?", username)
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
func IsCorrectCred(username string, password string) bool {

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

	rows, err := db.Query("SELECT userid FROM users WHERE username = ?", username)
	if err != nil {
		Utils.LogErr(err.Error())
	}

	defer rows.Close()
	for rows.Next() {
		var userid int32
		err := rows.Scan(&userid)
		if err != nil {
			return -1
		}
		return userid
	}
	return -1
}
func IsRestricted(userid int32) bool {

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
func GetNewUserId() int32 {

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		Utils.LogErr(err.Error())
	}
	rowsreturned := 1
	defer rows.Close()
	for rows.Next() {
		rowsreturned++
	}
	Utils.LogWarning("%d", rowsreturned+1)
	return int32(rowsreturned) + 1
}
func GetNewScoreId() int32 {

	rows, err := db.Query("SELECT * FROM scores")
	if err != nil {
		Utils.LogErr(err.Error())
	}
	rowsreturned := 0
	defer rows.Close()
	for rows.Next() {
		rowsreturned++
	}
	return int32(rowsreturned)
}
func InsertScore(score Utils.Score, scoreid int32) {

	rows, err := db.Query("INSERT INTO `scores`(`scoreid`, `mapchecksum`,`username`, `OnlineScoreChecksum`, `Count300`, `Count100`, `Count50`, `CountGeki`, `CountKatu`, `CountMiss`, `TotalScore`, `MaxCombo`, `Perfect`, `Ranking`, `EnabledMods`, `Pass`, `Accuracy`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		scoreid, score.FileChecksum, score.Username, score.OnlineScoreChecksum, score.Count300, score.Count100, score.Count50, score.CountGeki, score.CountKatu, score.CountMiss, score.TotalScore, score.MaxCombo, score.Perfect, score.Ranking, score.EnabledMods, score.EnabledMods, Utils.CalculateAccuracy(score))
	defer rows.Close()
	if err != nil {
	}
}

func GetScores(mapchecksum string) string {
	var response string

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
func UpdateTotalScore(user string) {

	rows, err := db.Query(`
    SELECT 
        mapchecksum,
        TotalScore
    FROM scores
    WHERE Username = ?`, user)

	if err != nil {
		log.Printf("Error executing query: %s", err)
		return
	}
	defer rows.Close()
	rowam := 1
	totalscore := 0
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
		totalscore += int(score.TotalScore)
		rowam++
	}
	db.Query("UPDATE `users` SET `total_score`= ? WHERE username = ?", totalscore, user)
}
func UpdatePlaycount(user string) {

	rows, err := db.Query(`
    SELECT 
        mapchecksum,
        TotalScore
    FROM scores
    WHERE Username = ?`, user)

	if err != nil {
		log.Printf("Error executing query: %s", err)
		return
	}
	defer rows.Close()
	rowam := 0
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
		rowam++
	}
	db.Query("UPDATE `users` SET `playcount`= ? WHERE username = ?", rowam, user)
}

func IsAdmin(userid int32) bool {

	rows, err := db.Query("SELECT * FROM admins WHERE userid = ?", userid)
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

/*
-1 = NotSubmitted
0 = Pending
1 = UpdateAvailable
2 = Ranked?
3 = Approved

we like only approved and ranked
*/
func GetMapStatus(checksum string) string {

	rows, err := db.Query("SELECT status from beatmaps WHERE checksum = ?", checksum)
	if err != nil {
		return "-1"
	}
	for rows.Next() {
		var status string
		rows.Scan(&status)
		return status
	}
	return "-1"
}
func IsRanked(checksum string) bool {
	return GetMapStatus(checksum) == "2" || GetMapStatus(checksum) == "3"
}
func GetTopUsers() ([]map[string]interface{}, error) {

	rows, err := db.Query(`
        SELECT username, ranked_score, total_score, accuracy 
        FROM users 
        ORDER BY ranked_score DESC LIMIT 25`)
	if err != nil {
		Utils.LogErr(err.Error())
		return nil, err
	}
	defer rows.Close()

	var topUsers []map[string]interface{}
	rank := 1

	for rows.Next() {
		var username string
		var rankedScore, totalScore int64
		var accuracy float64

		err := rows.Scan(&username, &rankedScore, &totalScore, &accuracy)
		if err != nil {
			Utils.LogErr("Error scanning row: %s", err)
			continue
		}

		user := map[string]interface{}{
			"rank":        fmt.Sprintf("%d", rank),
			"Username":    username,
			"RankedScore": fmt.Sprintf("%d", rankedScore),
			"TotalScore":  fmt.Sprintf("%d", totalScore),
			"Accuracy":    fmt.Sprintf("%.2f", accuracy),
		}
		topUsers = append(topUsers, user)
		rank++
	}

	return topUsers, nil
}

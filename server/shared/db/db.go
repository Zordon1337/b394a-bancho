package db

import (
	"database/sql"
	"log"
	"retsu/Utils"
	"retsu/cho/Structs"
	"time"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var connectionstring string = "test:test@tcp(127.0.0.1:3306)/osu!"
var db *sql.DB

type User struct {
	Username    string  `json:"Username"`
	UserId      int     `json:"UserId"`
	RankedScore int64   `json:"RankedScore"`
	Accuracy    float32 `json:"Accuracy"`
	PlayCount   int32   `json:"PlayCount"`
	TotalScore  int64   `json:"TotalScore"`
	Rank        int32   `json:"Rank"`
	JoinDate    string  `json:"JoinDate"`
	LastOnline  string  `json:"LastOnline"`
}

func InitDatabase() {
	var err error
	db, err = sql.Open("mysql", connectionstring)
	if err != nil {
		Utils.LogErr("Failed to create database connection!")
	}
}
func GetJoinDate(username string) string {

	rows, err := db.Query("SELECT joindate FROM users WHERE username = ?", username)
	if err != nil {
		Utils.LogErr(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var joindate string
		rows.Scan(&joindate)
		return joindate
	}
	return "unknown"
}
func GetUserFromDatabasePassword(username string, password string) Structs.UserStats {
	user := new(Structs.UserStats)

	rows, err := db.Query("SELECT userid, ranked_score, accuracy, playcount, total_score, rank FROM users WHERE username = ? AND password = ?", username, password)
	if err != nil {
		Utils.LogErr(err.Error())
	}

	defer rows.Close()
	for rows.Next() {
		rows.Scan(&user.UserID, &user.RankedScore, &user.Accuracy, &user.PlayCount, &user.TotalScore, &user.Rank)
		return *user
	}
	user.UserID = -1
	return *user
}
func GetUserFromDatabase(username string) User {
	user := new(User)

	rows, err := db.Query("SELECT userid, username, ranked_score, accuracy, playcount, total_score, rank, lastonline, joindate FROM users WHERE username = ?", username)
	if err != nil {
		Utils.LogErr(err.Error())
	}

	defer rows.Close()
	for rows.Next() {
		rows.Scan(&user.UserId, &user.Username, &user.RankedScore, &user.Accuracy, &user.PlayCount, &user.TotalScore, &user.Rank, &user.LastOnline, &user.JoinDate)
		return *user
	}
	user.UserId = -1
	return *user
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
func SetUsername(username string, newusername string) bool {
	if IsNameTaken(newusername) {
		return false
	} else {
		db.Query("UPDATE `users` SET `username`= ? WHERE username = ?", newusername, username)
		return true
	}
}
func SetPassword(username string, currentpass string, newpass string) bool {
	if IsCorrectCred(username, currentpass) {
		db.Query("UPDATE `users` SET `password`= ? WHERE username = ?", newpass, username)
		return true
	} else {
		return false
	}
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
func UpdateAccuracy(user string) {

	rows, err := db.Query(`
SELECT 
    SUM(Accuracy * (Count300 + Count100 + Count50 + CountMiss)) / 
    SUM(Count300 + Count100 + Count50 + CountMiss) AS total_accuracy
FROM scores
WHERE username = ?;
	`, user)

	if err != nil {
		log.Printf("Error executing query: %s", err)
		return
	}
	defer rows.Close()
	rowam := 1
	var accuracy float32
	for rows.Next() {
		err := rows.Scan(
			&accuracy,
		)
		if err != nil {
			Utils.LogErr("Error scanning row: %s", err)
			continue
		}
		rowam++
	}
	db.Query("UPDATE `users` SET `accuracy`= ? WHERE username = ?", accuracy, user)
}
func UpdateRank(username string) {

	user := new(User)

	rows, err := db.Query("SELECT userid, username, ranked_score, accuracy, playcount, total_score, rank, lastonline, joindate FROM users ORDER BY ranked_score DESC")
	if err != nil {
		Utils.LogErr(err.Error())
	}
	rownum := 1
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&user.UserId, &user.Username, &user.RankedScore, &user.Accuracy, &user.PlayCount, &user.TotalScore, &user.Rank, &user.LastOnline, &user.JoinDate)
		if user.Username == username {
			db.Query("UPDATE `users` SET `rank`= ? WHERE username = ?", rownum, username)
			return
		}
		rownum++
	}
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
        ORDER BY ranked_score DESC LIMIT 100`)
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
func UpdateLastOnline(username string) {

	rows, err := db.Query("UPDATE `users` SET `lastonline`= ? WHERE username = ?", time.Now().Format("2006-01-02 15:04:05"), username)
	if err != nil {
		Utils.LogErr(err.Error())
	}
	defer rows.Close()
}
func DoesExist(username string) bool {

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
func UnRestrictUser(username string) {
	id := GetUserIdByUsername(username)
	rows, err := db.Query("DELETE from restricts WHERE userid = ?", id)
	if err != nil {
	}
	defer rows.Close()
}
func RestrictUser(username string, admin string, reason string) {

	id := GetUserIdByUsername(username)
	rows, err := db.Query("INSERT INTO `restricts`(`bandate`, `userid`, `bannedby`, `reason`) VALUES (?,?,?,?)", time.Now().Format("2006-01-02 15:04:05"), id, admin, reason)
	if err != nil {
	}
	defer rows.Close()
}

func DoesMapExistInDB(checksum string) bool {

	rows, err := db.Query("SELECT status from beatmaps WHERE checksum = ?", checksum)
	if err != nil {
		return false
	}
	for rows.Next() {
		var status string
		rows.Scan(&status)
		return true
	}
	return false
}
func SetStatus(checksum string, newstatus string) string {
	curstatus := GetMapStatus(checksum)
	if newstatus == curstatus {
		return "Map already has this status!"
	}

	if DoesMapExistInDB(checksum) {

		db.Query("UPDATE `beatmaps` SET `status`= ? WHERE checksum = ?", newstatus, checksum)
	} else {
		db.Query("INSERT INTO `beatmaps`(`checksum`, `status`) VALUES (?,?)", checksum, newstatus)
	}
	return "Successfully updated beatmap"
}

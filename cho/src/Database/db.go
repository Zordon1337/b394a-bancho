package db

import (
	"database/sql"
	"socket-server/src/Structs"

	"socket-server/src/Utils"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var connectionstring string = "test:test@tcp(127.0.0.1:3306)/osu!"

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
		err := rows.Scan(&userid)
		if err != nil {
			return -1
		}
		return userid
	}
	return -1
}
func GetJoinDate(username string) string {
	db, err := sql.Open("mysql", connectionstring)
	if err != nil {
		Utils.LogErr(err.Error())
	}
	defer db.Close()
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
func GetUserFromDatabase(username string, password string) Structs.UserStats {
	user := new(Structs.UserStats)
	db, err := sql.Open("mysql", connectionstring)
	if err != nil {
		Utils.LogErr(err.Error())
	}
	defer db.Close()
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

func UpdateLastOnline(username string) {
	db, err := sql.Open("mysql", connectionstring)
	if err != nil {
		Utils.LogErr(err.Error())
	}
	defer db.Close()
	rows, err := db.Query("UPDATE `users` SET `lastonline`= ? WHERE username = ?", time.Now().Format("2006-01-02 15:04:05"), username)
	if err != nil {
		Utils.LogErr(err.Error())
	}
	defer rows.Close()
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
func DoesExist(username string) bool {
	db, err := sql.Open("mysql", connectionstring)
	if err != nil {
		Utils.LogErr(err.Error())
	}
	defer db.Close()
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
	db, err := sql.Open("mysql", connectionstring)
	if err != nil {
		Utils.LogErr(err.Error())
		return
	}
	defer db.Close()
	id := GetUserIdByUsername(username)
	rows, err := db.Query("DELETE from restricts WHERE userid = ?", id)
	defer rows.Close()
}
func RestrictUser(username string, admin string, reason string) {
	db, err := sql.Open("mysql", connectionstring)
	if err != nil {
		Utils.LogErr(err.Error())
		return
	}
	defer db.Close()
	id := GetUserIdByUsername(username)
	rows, err := db.Query("INSERT INTO `restricts`(`bandate`, `userid`, `bannedby`, `reason`) VALUES (?,?,?,?)", time.Now().Format("2006-01-02 15:04:05"), id, admin, reason)
	defer rows.Close()
}
func IsAdmin(userid int32) bool {
	db, err := sql.Open("mysql", connectionstring)
	if err != nil {
		Utils.LogErr(err.Error())
	}
	defer db.Close()
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
	db, err := sql.Open("mysql", connectionstring)
	if err != nil {
		return "-1"
	}
	defer db.Close()
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

func DoesMapExistInDB(checksum string) bool {
	db, err := sql.Open("mysql", connectionstring)
	if err != nil {
		return false
	}
	defer db.Close()
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
func IsRanked(checksum string) bool {
	return GetMapStatus(checksum) == "2" || GetMapStatus(checksum) == "3"
}
func SetStatus(checksum string, newstatus string) string {
	curstatus := GetMapStatus(checksum)
	if newstatus == curstatus {
		return "Map already has this status!"
	}
	db, err := sql.Open("mysql", connectionstring)
	if err != nil {
		return "Failed to update map"
	}
	defer db.Close()
	if DoesMapExistInDB(checksum) {

		db.Query("UPDATE `beatmaps` SET `status`= ? WHERE checksum = ?", newstatus, checksum)
	} else {
		db.Query("INSERT INTO `beatmaps`(`checksum`, `status`) VALUES (?,?)", checksum, newstatus)
	}
	return "Successfully updated beatmap"
}

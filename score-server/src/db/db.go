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

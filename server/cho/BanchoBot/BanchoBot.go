package BanchoBot

import (
	"fmt"
	"retsu/Utils"
	"retsu/cho/Structs"
	db "retsu/shared/db"
	"strconv"
	"strings"
)

func GenerateProfile() *Structs.Player {
	stats := Structs.UserStats{
		UserID:      1,
		RankedScore: 1337,
		Accuracy:    0.1337,
		PlayCount:   1337,
		TotalScore:  1337,
		Rank:        0,
	}
	status := Structs.Status{
		Status:          0,
		BeatmapUpdate:   false,
		StatusText:      "",
		BeatmapChecksum: "",
		CurrentMods:     0,
		BeatmapId:       0,
		PlayMode:        0,
	}
	player := Structs.Player{
		Username:  "BanchoBot",
		Conn:      nil, // bot doesn't actually connect but just fakes its presence
		Stats:     stats,
		Status:    status,
		IsInLobby: false,
		Timezone:  byte(24),
		Country:   "Satelite",
	}
	return &player
}
func HandleMsg(sender *Structs.Player, msg string, target string) string {
	if strings.HasPrefix(msg, "!ping") {
		return "Pong!"
	}
	if strings.HasPrefix(msg, "!whoami") {
		id := db.GetUserIdByUsername(sender.Username)
		return fmt.Sprintf("You are %s\nUserId: %d\nIsAdmin: %t\nIsRestricted: %t\nJoin Date: %s", sender, id, db.IsAdmin(id), db.IsRestricted(id), db.GetJoinDate(sender.Username))
	}
	if strings.HasPrefix(msg, "!unrestrict") {
		id := db.GetUserIdByUsername(sender.Username)
		if db.IsAdmin(id) {
			args := strings.Split(msg, " ")
			if len(args) < 2 {
				return "Missing arguments! correct command: !restrict <username>"
			}
			username := args[1]
			if db.DoesExist(username) {
				if !db.IsRestricted(db.GetUserIdByUsername(username)) {
					return "User is not restricted!"
				} else {
					db.UnRestrictUser(username)
					return "Succesfully removed restriction from" + username
				}
			} else {
				return "User does not exist!"
			}
		} else {
			return "You are not an admin!"
		}
	}
	if strings.HasPrefix(msg, "!restrict") {
		id := db.GetUserIdByUsername(sender.Username)
		if db.IsAdmin(id) {
			args := strings.Split(msg, " ")
			if len(args) < 3 {
				return "Missing arguments! correct command: !restrict <username> <reason(without spaces)>"
			}
			username := args[1]
			reason := args[2]
			if db.DoesExist(username) {
				if db.IsRestricted(db.GetUserIdByUsername(username)) {
					return "User is already restricted!"
				} else {
					db.RestrictUser(username, sender.Username, reason)
					return "Succesfully restricted " + username
				}
			} else {
				return "User does not exist!"
			}
		} else {
			return "You are not an admin!"
		}
	}
	fmt.Println(msg)
	if strings.HasPrefix(msg, "!updatebeatmapstatus") {
		id := db.GetUserIdByUsername(sender.Username)
		if db.IsAdmin(id) {
			args := strings.Split(msg, " ")
			if len(args) < 3 {
				return "Missing arguments! correct command: !updatebeatmapstatus <beatmapmd5> <newstatus>"
			}
			md5 := args[1]
			status := args[2]
			return db.SetStatus(md5, status)
		} else {
			return "You are not an admin!"
		}
	}
	if strings.HasPrefix(msg, "!updatecurrentmap") {
		uid := db.GetUserIdByUsername(sender.Username)
		if db.IsAdmin(uid) {
			id := sender.LastNp
			if id == "" {
				return "/np first"
			}
			if sender.LastNpIsSet {
				maps := Utils.GetBeatmapsBySetId(id)
				for _, map1 := range maps {
					db.SetStatus(map1.FileMD5, strconv.Itoa(map1.Approved))
				}
				return "OK(Set)"
			} else {
				maps := Utils.GetBeatmapsById(id)
				for _, map1 := range maps {
					db.SetStatus(map1.FileMD5, strconv.Itoa(map1.Approved))
				}
				return "OK(Single)"
			}
		} else {
			return "You are not an admin!"
		}
	}

	return ""
}

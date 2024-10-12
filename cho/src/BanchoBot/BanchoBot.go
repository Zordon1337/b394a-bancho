package BanchoBot

import (
	"socket-server/src/Structs"
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
	}
	player := Structs.Player{
		Username:  "BanchoBot",
		Conn:      nil, // bot doesn't actually connect but just fakes its presence
		Stats:     stats,
		Status:    status,
		IsInLobby: false,
		Timezone:  byte(24),
	}
	return &player
}
func HandleMsg(sender string, msg string, target string) string {
	if strings.HasPrefix(msg, "!ping") {
		return "Pong!"
	}
	return ""
}

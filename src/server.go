package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"socket-server/src/Packets"
	"socket-server/src/Structs"
	"socket-server/src/Utils"
	"strings"
	"sync"
	"time"
)

var (
	addr        string     = ":13381"
	players                = make(map[string]*Structs.Player) // Map to hold logged-in players
	playersMu   sync.Mutex                                    // Mutex to synchronize access to the player list
	MpMatches   = make(map[byte]*Structs.Match)
	MpMatchesMu sync.Mutex
)

func main() {
	var ln, err = net.Listen("tcp", addr)
	if err != nil {
		Utils.LogErr(err.Error())
		return
	}
	Utils.LogInfo("Socket listening for clients on " + addr)
	Utils.LogInfo("Lets play osu!")
	for {
		client, err := ln.Accept()
		if err != nil {
			Utils.LogErr(err.Error())
			ln.Close() // a bit brutal but okay
		}
		if client != nil {
			Utils.LogInfo("Client connected, Welcome")
			go handleClient(client)
		}
	}
}
func handleClient(client net.Conn) {
	defer client.Close()
	initialMessage := make([]byte, 1024)
	n, err := client.Read(initialMessage)
	if err != nil {
		Utils.LogErr("Error reading login info", err)
		return
	}
	lines := strings.Split(string(initialMessage[:n]), "\n")
	username := strings.TrimSpace(lines[0])
	//md5Hash := lines[1]
	build := strings.Split(lines[2], "|")[0]
	Utils.LogInfo(username + " logged in on build " + build)
	Packets.WriteLoginReply(client, int32(len(players))+1)
	Packets.WriteChannelJoinSucess(client, "#osu")
	stats := Structs.UserStats{
		UserID:      int32(len(players)) + 1,
		RankedScore: 1337,
		Accuracy:    1,
		PlayCount:   0,
		TotalScore:  1337,
		Rank:        1,
	}
	status := Structs.Status{
		Status:          0,
		BeatmapUpdate:   false,
		StatusText:      "",
		BeatmapChecksum: "",
		CurrentMods:     0,
	}
	player := Structs.Player{
		Username:  (username),
		Conn:      client,
		Stats:     stats,
		Status:    status,
		IsInLobby: false,
	}

	addPlayer(&player)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	Packets.WriteAnnounce(player.Conn, "Welcome to bancho, Mr. "+player.Username)
	for _, player1 := range players {
		Packets.WriteUserStats(player1.Conn, player, 2)
	}
	for _, player1 := range players {
		Packets.WriteUserStats(player.Conn, *player1, 2)
	}
	go func() {
		for {
			select {
			case <-ticker.C:
				Packets.WritePing(client)
				if err != nil {
					return
				}
			}
		}
	}()

	for {

		header := make([]byte, 7) // int16 + bool + int32
		_, err := io.ReadFull(client, header)
		if err != nil {
			Utils.LogErr("Error reading from client:", err)
			removePlayer(player.Username, player.Stats.UserID)
			return
		}
		var packetType int16
		var compression bool
		var dataLength int32
		buf := bytes.NewReader(header)

		if err := binary.Read(buf, binary.LittleEndian, &packetType); err != nil {
			Utils.LogErr("Failed to read packet type:", err)
			removePlayer(player.Username, player.Stats.UserID)
			return
		}
		if err := binary.Read(buf, binary.LittleEndian, &compression); err != nil {
			Utils.LogErr("Failed to read bool compression:", err)
			removePlayer(player.Username, player.Stats.UserID)
			return
		}
		if err := binary.Read(buf, binary.LittleEndian, &dataLength); err != nil {
			Utils.LogErr("Failed to read data length:", err)
			removePlayer(player.Username, player.Stats.UserID)
			return
		}

		if dataLength < 0 || dataLength > 4096 {
			fmt.Printf("Invalid data length: %d\n", dataLength)
		}

		data := make([]byte, dataLength)
		if _, err := io.ReadFull(client, data); err != nil {
			Utils.LogErr("Failed to read packet data:", err)
			return
		}
		switch packetType {
		case 0: // status change
			{
				handleStatus(player, data)
			}
		case 1: // irc msg
			{
				handleMsg(player, data)
				break
			}
		case 2: // osu quit
			{
				removePlayer(player.Username, player.Stats.UserID)
				return
			}
		case 4: // pong
			{
				break
			}
		case 3: // user stats request
			{
				Packets.WriteUserStats(player.Conn, player, 2)
				break
			}

		case 32: // Create Match
			{
				buf := bytes.NewReader(data)
				match := new(Structs.Match)
				err := binary.Read(buf, binary.LittleEndian, &match.MatchId)
				if err != nil {
					Utils.LogErr("Error occurred on MatchId creation: ", err.Error())
					return
				}
				err = binary.Read(buf, binary.LittleEndian, &match.InProgress)
				if err != nil {
					Utils.LogErr("Error occurred on InProgress creation: ", err.Error())
					return
				}
				err = binary.Read(buf, binary.LittleEndian, &match.MatchType)
				if err != nil {
					Utils.LogErr("Error occurred on MatchType creation: ", err.Error())
					return
				}
				err = binary.Read(buf, binary.LittleEndian, &match.ActiveMods)
				if err != nil {
					Utils.LogErr("Error occurred on ActiveMods creation: ", err.Error())
					return
				}
				match.GameName, err = Utils.ReadOsuString(buf)
				if err != nil {
					Utils.LogErr("Error occurred on GameName creation: ", err.Error())
					return
				}
				match.BeatmapName, err = Utils.ReadOsuString(buf)
				if err != nil {
					Utils.LogErr("Error occurred on BeatmapName creation: ", err.Error())
					return
				}
				err = binary.Read(buf, binary.LittleEndian, &match.BeatmapId)
				if err != nil {
					Utils.LogErr("Error occurred on BeatmapId creation: ", err.Error())
					return
				}
				match.BeatmapChecksum, err = Utils.ReadOsuString(buf)
				if err != nil {
					Utils.LogErr("Error occurred on BeatmapChecksum creation: ", err.Error())
					return
				}
				match.SlotStatus = [8]byte{4, 1, 1, 1, 1, 1, 1, 1}
				match.SlotId = [8]int32{player.Stats.UserID, -1, -1, -1, -1, -1, -1, -1}
				Utils.LogErr("match created", match.MatchId, match.InProgress, match.MatchType, match.ActiveMods, match.GameName, match.BeatmapName, match.BeatmapId, match.BeatmapChecksum)
				Packets.WriteMatchJoinSuccess(player.Conn, *match)
				player.CurrentMatch = match
				AddMatch(match)
				break
			}
		case 33: // Join Match
			{
				buf := bytes.NewReader(data)
				var matchid byte
				err := binary.Read(buf, binary.LittleEndian, &matchid)
				if err != nil {
					Utils.LogErr("Error occurred on match creation: ", err.Error())
					return
				}
				Utils.LogInfo("Player %s joined match %b", player.Username, matchid)
				match := FindMatchById(matchid)
				player.CurrentMatch = match
				if JoinMatch(&player, match) {
					player.IsInLobby = false
				}
				break
			}
		case 31:
			{ // Join Lobby
				player.IsInLobby = true
				for _, match := range MpMatches {
					Packets.WriteMatchUpdate(player.Conn, *match)
				}
				break
			}
		case 30:
			{
				player.IsInLobby = false
				break
			}
		case 34: // Part Match aka quit match
			{
				PartMatch(&player, player.CurrentMatch)
				break
			}
		case 56: // unready
			{
				SetSlotStatusById(player.Stats.UserID, player.CurrentMatch, 4)
				break
			}
		case 45: // Start match
			{
				m := player.CurrentMatch
				m.InProgress = true
				m.LoadingPeople = int32(CalculatePlayerInLobby(m))
				for i := 0; i < 8; i++ {
					if m.SlotId[i] != -1 {
						m.SlotStatus[i] = 32 // playing
					}
				}
				for i := 0; i < 8; i++ {
					if m.SlotId[i] != -1 {
						plr := GetPlayerById(m.SlotId[i])
						Packets.WriteMatchUpdate(plr.Conn, *m)
						Packets.WriteMatchStart(plr.Conn, *m)
					}
				}

				break
			}
		case 48:
			{ // score update

				score := Structs.ReadScoreFrameFromStream(data)
				score.SlotId = byte(FindUserSlotInMatchById(player.Stats.UserID, player.CurrentMatch))
				m := player.CurrentMatch
				for i := 0; i < 8; i++ {
					if m.SlotId[i] != -1 && m.SlotStatus[i] == 32 { // is player and is playing
						plr := GetPlayerById(m.SlotId[i])
						Packets.WriteMatchScoreUpdate(plr.Conn, Structs.WriteScoreFrameToBytes(score))
					}
				}
				break
			}
		case 50: // Match complete
			{

				m := player.CurrentMatch
				for i := 0; i < 8; i++ {
					if m.SlotId[i] != -1 && m.SlotStatus[i] == 32 { // is player and is playing
						plr := GetPlayerById(m.SlotId[i])
						m.SlotStatus[i] = 4
						m.InProgress = false
						Packets.WriteMatchUpdate(plr.Conn, *m)
						Packets.WriteMatchComplete(plr.Conn)
					}
				}
				break
			}
		case 53: // loaded into game
			{
				m := player.CurrentMatch
				m.LoadingPeople--
				if m.LoadingPeople < 1 { // everybody loaded and we are happy
					for i := 0; i < 8; i++ {
						if m.SlotId[i] != -1 && m.SlotStatus[i] == 32 { // is player and is playing
							plr := GetPlayerById(m.SlotId[i])
							Packets.MatchAllPlayersLoaded(plr.Conn)
						}
					}
				}

				break
			}
		case 42: // Beatmap change
			{
				buf := bytes.NewReader(data)
				match := player.CurrentMatch
				err := binary.Read(buf, binary.LittleEndian, &match.MatchId)
				if err != nil {
					Utils.LogErr("Error occurred on MatchId creation: ", err.Error())
					return
				}
				err = binary.Read(buf, binary.LittleEndian, &match.InProgress)
				if err != nil {
					Utils.LogErr("Error occurred on InProgress creation: ", err.Error())
					return
				}
				err = binary.Read(buf, binary.LittleEndian, &match.MatchType)
				if err != nil {
					Utils.LogErr("Error occurred on MatchType creation: ", err.Error())
					return
				}
				err = binary.Read(buf, binary.LittleEndian, &match.ActiveMods)
				if err != nil {
					Utils.LogErr("Error occurred on ActiveMods creation: ", err.Error())
					return
				}
				match.GameName, err = Utils.ReadOsuString(buf)
				if err != nil {
					Utils.LogErr("Error occurred on GameName creation: ", err.Error())
					return
				}
				match.BeatmapName, err = Utils.ReadOsuString(buf)
				if err != nil {
					Utils.LogErr("Error occurred on BeatmapName creation: ", err.Error())
					return
				}
				err = binary.Read(buf, binary.LittleEndian, &match.BeatmapId)
				if err != nil {
					Utils.LogErr("Error occurred on BeatmapId creation: ", err.Error())
					return
				}
				match.BeatmapChecksum, err = Utils.ReadOsuString(buf)
				if err != nil {
					Utils.LogErr("Error occurred on BeatmapChecksum creation: ", err.Error())
					return
				}
				for i := 0; i < 8; i++ {
					if match.SlotId[i] != -1 {
						plr := GetPlayerById(match.SlotId[i])
						Packets.WriteMatchUpdate(plr.Conn, *match)
					}
				}
				break
			}
		case 52: // Mods
			{
				m := player.CurrentMatch
				buf := bytes.NewReader(data)
				var mods int16
				err := binary.Read(buf, binary.LittleEndian, &mods)
				if err != nil {
					return
				}
				m.ActiveMods = mods
				for i := 0; i < 8; i++ {
					if m.SlotId[i] != -1 {
						plr := GetPlayerById(m.SlotId[i])
						Packets.WriteMatchUpdate(plr.Conn, *m)
					}
				}
				break
			}
		case 21:
			{
				buf := bytes.NewReader(data)
				dat, err := Utils.ReadOsuString(buf)
				if err != nil {
					Utils.LogErr(err.Error())
					break
				}
				Utils.LogInfo("Osu! reported an error: ", dat)
				break
			}
		case 40: // Ready
			{
				match := player.CurrentMatch
				if match != nil { // we don't want to ready when we are not in a match
					slot := FindUserSlotInMatchById(player.Stats.UserID, match)
					if slot != -1 { // making sure we actually exist in that match
						match.SlotStatus[slot] = 8
					}
				}
				for i := 0; i < 8; i++ {
					if match.SlotId[i] != -1 {
						plr := GetPlayerById(match.SlotId[i])
						Packets.WriteMatchUpdate(plr.Conn, *match)
					}
				}
				break
			}
		default:
			{
				Utils.LogWarning("Received unhandled packet", packetType)
				break
			}
		}
	}
}
func handleMsg(player Structs.Player, data []byte) {
	buf := bytes.NewReader(data)
	Utils.ReadOsuString(buf)
	sender := player.Username // for some reason sender is empty???
	msg, err := Utils.ReadOsuString(buf)
	if err != nil {
		Utils.LogErr("Failed to read msg:", err)
		return
	}

	target, err := Utils.ReadOsuString(buf)
	if err != nil {
		Utils.LogErr("Failed to read target:", err)
		return
	}
	for _, player1 := range players {
		if player1.Username != player.Username {
			Packets.WriteMessage(player1.Conn, sender, msg, target)
		}
	}
	Utils.LogInfo(sender + "->" + target + ": " + msg)
}
func handleStatus(player Structs.Player, data []byte) {
	var status byte
	var bmapUpdate bool
	var mods uint16
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &status)
	if err != nil {
		Utils.LogErr("Error occurred while reading Status from " + player.Username + " " + err.Error())
	}
	err = binary.Read(buf, binary.LittleEndian, &bmapUpdate)
	if err != nil {
		Utils.LogErr("Error occurred while reading BeatmapUpdate from " + player.Username + " " + err.Error())
	}
	player.Status.Status = status
	player.Status.BeatmapUpdate = bmapUpdate
	if bmapUpdate {
		statusText, err := Utils.ReadOsuString(buf)
		if err != nil {
			Utils.LogErr("Failed to read statusText:", err)
			return
		}
		beatmapMd5, err := Utils.ReadOsuString(buf)
		if err != nil {
			Utils.LogErr("Failed to read beatmapMd5:", err)
			return
		}
		err = binary.Read(buf, binary.LittleEndian, &mods)
		if err != nil {
			Utils.LogErr("Error occurred while reading mods from " + player.Username)
		}
		player.Status.StatusText = statusText
		player.Status.BeatmapChecksum = beatmapMd5
		player.Status.CurrentMods = mods
	}
	for _, player1 := range players {
		Packets.WriteUserStats(player1.Conn, player, 0)
	}
}
func addPlayer(player *Structs.Player) {
	playersMu.Lock()
	defer playersMu.Unlock()
	players[player.Username] = player
}
func removePlayer(username string, id int32) {
	playersMu.Lock()
	defer playersMu.Unlock()
	delete(players, username)
	Utils.LogInfo("Player disconnected: %s\n", username)
	for _, player1 := range players {
		Packets.WriteIrcQuit(player1.Conn, username)
		Packets.WriteUserQuit(player1.Conn, id)
	}
}
func AddMatch(match *Structs.Match) {
	MpMatchesMu.Lock()
	defer MpMatchesMu.Unlock()
	MpMatches[match.MatchId] = match
	for _, player1 := range players {
		if player1.IsInLobby {
			Packets.WriteMatchUpdate(player1.Conn, *match)
		}
	}
}
func RemoveMatch(id byte) {
	MpMatchesMu.Lock()
	defer MpMatchesMu.Unlock()
	delete(MpMatches, id)
}
func FindMatchById(id byte) *Structs.Match {
	MpMatchesMu.Lock()
	defer MpMatchesMu.Unlock()
	if MpMatches[id] != nil {
		return MpMatches[id]
	}
	return nil
}
func FindUserSlotInMatchById(userid int32, match *Structs.Match) int {
	if match == nil {
		return 0
	}
	for i := 0; i < 8; i++ {
		if match.SlotId[i] == userid {
			return i
		}
	}
	return -1
}
func GetPlayerById(userid int32) *Structs.Player {
	for _, player := range players {
		if player.Stats.UserID == userid {
			return player
		}
	}
	return nil
}
func FindSlotForPlayer(match *Structs.Match) int {
	if match == nil {
		return 0
	}
	for i := 0; i < 8; i++ {
		if match.SlotId[i] == -1 {
			return i
		}
	}
	return -1
}
func SetSlotStatusById(userid int32, match *Structs.Match, newstatus byte) {
	if match == nil {
		return
	}
	slot := FindUserSlotInMatchById(userid, match)
	if slot != -1 {
		match.SlotStatus[slot] = newstatus
	}
	for i := 0; i < 8; i++ {
		if match.SlotId[i] != -1 {
			plr := GetPlayerById(match.SlotId[i])
			Packets.WriteMatchUpdate(plr.Conn, *match)
		}
	}
}
func JoinMatch(player *Structs.Player, match *Structs.Match) bool {
	if match == nil {
		return false
	}
	newslot := FindSlotForPlayer(match)
	if newslot == -1 {
		// match full
		Packets.WriteMatchJoinFail(player.Conn)
		return false
	} else {
		match.SlotId[newslot] = player.Stats.UserID
		match.SlotStatus[newslot] = 4
		for i := 0; i < 8; i++ {
			if match.SlotId[i] != -1 && match.SlotId[i] != player.Stats.UserID {
				plr := GetPlayerById(match.SlotId[i])
				Packets.WriteMatchUpdate(plr.Conn, *match)
			}
		}
		Packets.WriteMatchJoinSuccess(player.Conn, *match)
		return true
	}
}
func CalculatePlayerInLobby(match *Structs.Match) int {
	people := 0
	for i := 0; i < 8; i++ {
		if match.SlotId[i] != -1 {
			people++
		}
	}
	return people
}
func PartMatch(player *Structs.Player, match *Structs.Match) {
	if match == nil {
		return
	}
	slot := FindUserSlotInMatchById(player.Stats.UserID, match)
	match.SlotId[slot] = -1
	match.SlotStatus[slot] = 1
	for i := 0; i < 8; i++ {
		if match.SlotId[i] != -1 && match.SlotId[i] != player.Stats.UserID {
			plr := GetPlayerById(match.SlotId[i])
			Packets.WriteMatchUpdate(plr.Conn, *match)
		}
	}
	peopleinlobby := 0
	for i := 0; i < 8; i++ {
		if match.SlotId[i] != -1 {
			peopleinlobby++
		}
	}
	if peopleinlobby == 0 {
		// lobby empty, delete it
		RemoveMatch(match.MatchId)
		for _, player1 := range players {
			if player1.IsInLobby {
				Packets.WriteDisbandMatch(player1.Conn, int(match.MatchId))
			}
		}
	}
}

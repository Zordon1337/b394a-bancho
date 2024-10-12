package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"socket-server/src/Packets"
	"socket-server/src/Structs"
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
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Socket listening for clients on " + addr)
	fmt.Println("Lets play osu!")
	for {
		client, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			ln.Close() // a bit brutal but okay
		}
		if client != nil {
			fmt.Println("Client connected, Welcome")
			go handleClient(client)
		}
	}
}
func handleClient(client net.Conn) {
	defer client.Close()
	initialMessage := make([]byte, 1024)
	n, err := client.Read(initialMessage)
	if err != nil {
		fmt.Println("Error reading login info", err)
		return
	}
	lines := strings.Split(string(initialMessage[:n]), "\n")
	username := strings.TrimSpace(lines[0])
	//md5Hash := lines[1]
	build := strings.Split(lines[2], "|")[0]
	fmt.Println(username + " logged in on build " + build)
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
			fmt.Println("Error reading from client:", err)
			removePlayer(player.Username, player.Stats.UserID)
			return
		}
		var packetType int16
		var compression bool
		var dataLength int32
		buf := bytes.NewReader(header)

		if err := binary.Read(buf, binary.LittleEndian, &packetType); err != nil {
			fmt.Println("Failed to read packet type:", err)
			removePlayer(player.Username, player.Stats.UserID)
			return
		}
		if err := binary.Read(buf, binary.LittleEndian, &compression); err != nil {
			fmt.Println("Failed to read bool compression:", err)
			removePlayer(player.Username, player.Stats.UserID)
			return
		}
		if err := binary.Read(buf, binary.LittleEndian, &dataLength); err != nil {
			fmt.Println("Failed to read data length:", err)
			removePlayer(player.Username, player.Stats.UserID)
			return
		}

		if dataLength < 0 || dataLength > 4096 {
			fmt.Printf("Invalid data length: %d\n", dataLength)
		}

		data := make([]byte, dataLength)
		if _, err := io.ReadFull(client, data); err != nil {
			fmt.Println("Failed to read packet data:", err)
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
					fmt.Println("Error occurred on MatchId creation: ", err.Error())
					return
				}
				err = binary.Read(buf, binary.LittleEndian, &match.InProgress)
				if err != nil {
					fmt.Println("Error occurred on InProgress creation: ", err.Error())
					return
				}
				err = binary.Read(buf, binary.LittleEndian, &match.MatchType)
				if err != nil {
					fmt.Println("Error occurred on MatchType creation: ", err.Error())
					return
				}
				err = binary.Read(buf, binary.LittleEndian, &match.ActiveMods)
				if err != nil {
					fmt.Println("Error occurred on ActiveMods creation: ", err.Error())
					return
				}
				match.GameName, err = Packets.ReadOsuString(buf)
				if err != nil {
					fmt.Println("Error occurred on GameName creation: ", err.Error())
					return
				}
				match.BeatmapName, err = Packets.ReadOsuString(buf)
				if err != nil {
					fmt.Println("Error occurred on BeatmapName creation: ", err.Error())
					return
				}
				err = binary.Read(buf, binary.LittleEndian, &match.BeatmapId)
				if err != nil {
					fmt.Println("Error occurred on BeatmapId creation: ", err.Error())
					return
				}
				match.BeatmapChecksum, err = Packets.ReadOsuString(buf)
				if err != nil {
					fmt.Println("Error occurred on BeatmapChecksum creation: ", err.Error())
					return
				}
				match.SlotStatus = [8]byte{4, 1, 1, 1, 1, 1, 1, 1}
				match.SlotId = [8]int32{player.Stats.UserID, -1, -1, -1, -1, -1, -1, -1}
				fmt.Println("match created", match.MatchId, match.InProgress, match.MatchType, match.ActiveMods, match.GameName, match.BeatmapName, match.BeatmapId, match.BeatmapChecksum)
				Packets.WriteMatchJoinSuccess(player.Conn, *match)
				player.CurrentMatch = match
				AddMatch(match)
			}
		case 33: // Join Match
			{
				buf := bytes.NewReader(data)
				var matchid byte
				err := binary.Read(buf, binary.LittleEndian, &matchid)
				if err != nil {
					fmt.Println("Error occurred on match creation: ", err.Error())
					return
				}
				fmt.Printf("Player %s joined match %b", player.Username, matchid)
				match := FindMatchById(matchid)
				player.CurrentMatch = match
				if JoinMatch(&player, match) {
					player.IsInLobby = false
				}
			}
		case 31:
			{ // Join Lobby
				player.IsInLobby = true
				for _, match := range MpMatches {
					Packets.WriteMatchUpdate(player.Conn, *match)
				}
			}
		case 30:
			{
				player.IsInLobby = false
			}
		case 34:
			{
				PartMatch(&player, player.CurrentMatch)
			}
		case 40:
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
			}
		default:
			{
				fmt.Println("Received unhandled packet", packetType)
				break
			}
		}
	}
}
func handleMsg(player Structs.Player, data []byte) {
	buf := bytes.NewReader(data)
	Packets.ReadOsuString(buf)
	sender := player.Username // for some reason sender is empty???
	msg, err := Packets.ReadOsuString(buf)
	if err != nil {
		fmt.Println("Failed to read msg:", err)
		return
	}

	target, err := Packets.ReadOsuString(buf)
	if err != nil {
		fmt.Println("Failed to read target:", err)
		return
	}
	for _, player1 := range players {
		if player1.Username != player.Username {
			Packets.WriteMessage(player1.Conn, sender, msg, target)
		}
	}
	fmt.Println(sender + "->" + target + ": " + msg)
}
func handleStatus(player Structs.Player, data []byte) {
	var status byte
	var bmapUpdate bool
	var mods uint16
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &status)
	if err != nil {
		fmt.Println("Error occurred while reading Status from " + player.Username + " " + err.Error())
	}
	err = binary.Read(buf, binary.LittleEndian, &bmapUpdate)
	if err != nil {
		fmt.Println("Error occurred while reading BeatmapUpdate from " + player.Username + " " + err.Error())
	}
	player.Status.Status = status
	player.Status.BeatmapUpdate = bmapUpdate
	if bmapUpdate {
		statusText, err := Packets.ReadOsuString(buf)
		if err != nil {
			fmt.Println("Failed to read statusText:", err)
			return
		}
		beatmapMd5, err := Packets.ReadOsuString(buf)
		if err != nil {
			fmt.Println("Failed to read beatmapMd5:", err)
			return
		}
		err = binary.Read(buf, binary.LittleEndian, &mods)
		if err != nil {
			fmt.Println("Error occurred while reading mods from " + player.Username)
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
	fmt.Printf("Player disconnected: %s\n", username)
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
	for i := 0; i < 8; i++ {
		if match.SlotId[i] == -1 {
			return i
		}
	}
	return -1
}

func JoinMatch(player *Structs.Player, match *Structs.Match) bool {
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
func PartMatch(player *Structs.Player, match *Structs.Match) {
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

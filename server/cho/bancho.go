package cho

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"regexp"
	"retsu/Utils"
	"retsu/cho/BanchoBot"
	"retsu/cho/Packets"
	"retsu/cho/Structs"
	db "retsu/shared/db"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	addr        string = "0.0.0.0:13381"
	players            = make(map[string]*Structs.Player)
	playersMu   sync.Mutex
	MpMatches   = make(map[byte]*Structs.Match)
	MpMatchesMu sync.Mutex
)

func Bancho() {
	var ln, err = net.Listen("tcp", addr)
	if err != nil {
		Utils.LogErr(err.Error())
		return
	}
	db.InitDatabase()
	Utils.LogInfo("Socket listening for clients on " + addr)
	Utils.LogInfo("Lets play osu!")
	addPlayer(BanchoBot.GenerateProfile())
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
	if len(lines) < 3 {
		return // make client reconnect
	}
	username := strings.TrimSpace(lines[0])
	password := strings.TrimSpace(lines[1])
	//md5Hash := lines[1]
	build := strings.Split(lines[2], "|")[0]
	timezone, err := strconv.Atoi(strings.Split(lines[2], "|")[1])
	ip := client.RemoteAddr().(*net.TCPAddr).IP.String()
	country, err := Utils.GetCountryFromIP(ip)
	if err != nil {
		Utils.LogErr("Error getting country for %s: %s", username, err.Error())
		return
	}
	stats := db.GetUserFromDatabasePassword(username, password)
	status := Structs.Status{
		Status:          0,
		BeatmapUpdate:   false,
		StatusText:      "",
		BeatmapChecksum: "",
		CurrentMods:     0,
	}
	re := regexp.MustCompile(`\D`)
	cleaned := re.ReplaceAllString(build, "")
	build2, err := strconv.ParseInt(cleaned, 0, 64)
	fmt.Println(cleaned)
	player := Structs.Player{
		Username:  (username),
		Conn:      client,
		Stats:     stats,
		Status:    status,
		IsInLobby: false,
		Timezone:  byte(24 + timezone),
		Country:   country,
		Build:     build2,
	}
	if stats.UserID < 1 {
		Utils.LogInfo(username + " attempted to login with invalid password")
		Packets.WriteLoginReply(client, -1, int(player.Build))
		return
	}
	Utils.LogInfo(username + " logged in on build " + build)
	Packets.WriteLoginReply(client, player.Stats.UserID, int(player.Build))
	if build2 > 394 {
		Packets.WriteProtcolNegotiation(client, 5, int(player.Build))
	}
	Packets.WriteChannelJoinSuccess(client, "#osu", int(player.Build))
	if db.IsRestricted(stats.UserID) {
		Packets.WriteMessage(client, "BanchoBot", "You are restricted, because of that you can't submit any new scores", player.Username, int(player.Build))
	}
	db.UpdateLastOnline(player.Username)
	addPlayer(&player)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	Packets.WriteAnnounce(player.Conn, "Welcome to bancho, Mr. "+player.Username, int(player.Build))
	go func() {
		time.Sleep(1000)
		playersMu.Lock()
		for _, player1 := range players {
			if player1.Build > 1717 {
				Packets.WriteUserStats(*player1, player, 2, int(player1.Build))
				Packets.WriteUserPresense(*player1, player, 2, int(player1.Build))
			} else {
				Packets.WriteUserStats(*player1, player, 2, int(player1.Build))
			}
		}
		for _, player1 := range players {
			if player.Build > 1717 && player.Stats.UserID != player1.Stats.UserID {
				Packets.WriteUserPresense(player, *player1, 2, int(player.Build))
			}
			Packets.WriteUserStats(player, *player1, 2, int(player.Build))

		}
		playersMu.Unlock()
	}()
	go func() {

		for {
			select {
			case <-ticker.C:
				stats2 := db.GetUserFromDatabasePassword(username, password)
				if stats2.Accuracy != player.Stats.Accuracy || stats2.PlayCount != player.Stats.PlayCount || stats2.RankedScore != player.Stats.RankedScore || stats2.Rank != player.Stats.Rank || stats2.TotalScore != player.Stats.TotalScore {
					player.Stats = stats2
					playersMu.Lock()
					for _, player1 := range players {
						Packets.WriteUserStats(*player1, player, 2, int(player1.Build))
					}
					playersMu.Unlock()
				}
				Packets.WritePing(client, int(player.Build))
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

				handleStatus(&player, data)

				break
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

				break
			}

		case Utils.CalculatePacketOffset(int(player.Build), 32): // Create Match
			{
				buf := bytes.NewReader(data)
				match := new(Structs.Match)
				var matchidtemp int16
				if player.Build > 1717 {
					err := binary.Read(buf, binary.LittleEndian, &matchidtemp)
					if err != nil {
						Utils.LogErr("Error occurred on MatchId creation: ", err.Error())
						return
					}
					match.MatchId = byte(matchidtemp)
				} else {
					err := binary.Read(buf, binary.LittleEndian, &match.MatchId)
					if err != nil {
						Utils.LogErr("Error occurred on MatchId creation: ", err.Error())
						return
					}
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
				if player.Build > 590 {
					match.Password, err = Utils.ReadOsuString(buf)
					if err != nil {
						Utils.LogErr("Error occurred on Password creation: ", err.Error())
						return
					}
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
				if build2 > 394 {
					for i := 0; i < 8; i++ {
						err = binary.Read(buf, binary.LittleEndian, &match.SlotStatus[i])
					}
					if build2 > 553 {
						for i := 0; i < 8; i++ {
							err = binary.Read(buf, binary.LittleEndian, &match.SlotTeam[i])
						}
					}
					for i := 0; i < 8; i++ {
						if match.SlotStatus[i]&0x07C > 0 {
							err = binary.Read(buf, binary.LittleEndian, &match.SlotId[i])
						}
					}
				}
				err = binary.Read(buf, binary.LittleEndian, &match.HostId)
				if err != nil {
					Utils.LogErr("Error occurred on HostId creation: ", err.Error())
					return
				}
				if player.Build > 483 {

					err = binary.Read(buf, binary.LittleEndian, &match.Mode)
					if err != nil {
						Utils.LogErr("Error occurred on Mode creation: ", err.Error())
						return
					}

				}
				if player.Build > 535 {
					err = binary.Read(buf, binary.LittleEndian, &match.ScoringType)
					if err != nil {
						Utils.LogErr("Error occurred on ScoringType creation: ", err.Error())
						return
					}
					err = binary.Read(buf, binary.LittleEndian, &match.TeamType)
					if err != nil {
						Utils.LogErr("Error occurred on TeamType creation: ", err.Error())
						return
					}
				}
				match.SlotStatus = [8]byte{4, 1, 1, 1, 1, 1, 1, 1}
				match.SlotId = [8]int32{player.Stats.UserID, -1, -1, -1, -1, -1, -1, -1}
				Utils.LogInfo("%s created an match with id %s", player.Username, match.MatchId)
				Packets.WriteMatchJoinSuccess(player.Conn, *match, int32(player.Build))
				player.CurrentMatch = match
				AddMatch(match)
				break
			}

		case Utils.CalculatePacketOffset(int(player.Build), 17): // StartSpec
			{
				buf := bytes.NewReader(data)
				var userid int32
				binary.Read(buf, binary.LittleEndian, &userid)
				for _, player1 := range players {
					if player1.Stats.UserID == userid {
						if player.Build > 1700 && player1.Build < 1700 {
							Packets.WriteAnnounce(player.Conn, "Player that you want spectate, has too old client", int(player.Build))
							break
						} else if player.Build < 1700 && player1.Build > 1700 {
							Packets.WriteAnnounce(player.Conn, "Player that you want spectate, has too new client", int(player.Build))
							break
						}
						fmt.Printf("%s started spectating %s\n", player.Username, player1.Username)
						Packets.WriteSpecJoined(player1.Conn, player.Stats.UserID, int(player1.Build))
						if player1.Spectators == nil {
							player1.Spectators = make(map[int32]*Structs.Player)
						}
						player1.Spectators[player.Stats.UserID] = &player
						player.CurrentlySpectating = player1
					}
				}
				break
			}
		case Utils.CalculatePacketOffset(int(player.Build), 18): // StopSpec
			{

				if player.CurrentlySpectating != nil {
					player1 := player.CurrentlySpectating
					if player.Build > 1700 && player1.Build < 1700 {
						Packets.WriteAnnounce(player.Conn, "Player that you want spectate, has too old client", int(player.Build))
						return
					} else if player.Build < 1700 && player1.Build > 1700 {
						Packets.WriteAnnounce(player.Conn, "Player that you want spectate, has too new client", int(player.Build))
						return
					}
					Packets.WriteSpecQuited(player1.Conn, player.Stats.UserID, int(player1.Build))
					if player1.Spectators != nil {
						delete(player1.Spectators, player.Stats.UserID)
					}
					player.CurrentlySpectating = nil
					fmt.Printf("%s stopped spectating %s\n", player.Username, player1.Username)

				}
				break
			}
		case Utils.CalculatePacketOffset(int(player.Build), 19): // SendFrames
			{
				for _, player1 := range player.Spectators {
					packet, err := Utils.SerializePacket(15, data)
					if err != nil {
						break
					}
					fmt.Printf("%s sent replay frames to %s\n", player.Username, player1.Username)
					player1.Conn.Write(packet)
				}
				break
			}
		case Utils.CalculatePacketOffset(int(player.Build), 33): // Join Match
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
		case Utils.CalculatePacketOffset(int(player.Build), 31):
			{ // Join Lobby
				player.IsInLobby = true
				for _, match := range MpMatches {
					Packets.WriteMatchUpdate(player, *match, int(player.Build))
				}
				break
			}
		case Utils.CalculatePacketOffset(int(player.Build), 30):
			{
				player.IsInLobby = false
				break
			}
		case Utils.CalculatePacketOffset(int(player.Build), 34): // Part Match aka quit match
			{
				PartMatch(&player, player.CurrentMatch)
				break
			}
		case Utils.CalculatePacketOffset(int(player.Build), 56): // unready
			{
				SetSlotStatusById(player.Stats.UserID, player.CurrentMatch, 4)
				break
			}
		case Utils.CalculatePacketOffset(int(player.Build), 45): // Start match
			{
				m := player.CurrentMatch
				m.InProgress = true
				m.LoadingPeople = int32(CalculatePlayerInLobby(m))
				m.SkippingNeededToSkip = int32(CalculatePlayerInLobby(m))
				for i := 0; i < 8; i++ {
					if m.SlotId[i] != -1 {
						m.SlotStatus[i] = 32 // playing
					}
				}
				for i := 0; i < 8; i++ {
					if m.SlotId[i] != -1 {
						plr := GetPlayerById(m.SlotId[i])
						Packets.WriteMatchUpdate(*plr, *m, int(player.Build))
						Packets.WriteMatchStart(*plr, *m, int(player.Build))
					}
				}

				break
			}
		case Utils.CalculatePacketOffset(int(player.Build), 61):
			{
				m := player.CurrentMatch
				m.SkippingNeededToSkip--
				if m.SkippingNeededToSkip < 1 { // everybody loaded and we are happy
					for i := 0; i < 8; i++ {
						if m.SlotId[i] != -1 && m.SlotStatus[i] == 32 { // is player and is playing
							plr := GetPlayerById(m.SlotId[i])
							Packets.WriteMatchSkip(plr.Conn, int(player.Build))
						}
					}
				}
				break
			}
		case Utils.CalculatePacketOffset(int(player.Build), 48):
			{ // score update

				score := Structs.ReadScoreFrameFromStream(data)
				score.SlotId = byte(FindUserSlotInMatchById(player.Stats.UserID, player.CurrentMatch))
				m := player.CurrentMatch
				for i := 0; i < 8; i++ {
					if m.SlotId[i] != -1 && m.SlotStatus[i] == 32 { // is player and is playing
						plr := GetPlayerById(m.SlotId[i])
						Packets.WriteMatchScoreUpdate(plr.Conn, Structs.WriteScoreFrameToBytes(score), int(player.Build))
					}
				}
				break
			}
		case Utils.CalculatePacketOffset(int(player.Build), 50): // Match complete
			{

				m := player.CurrentMatch
				for i := 0; i < 8; i++ {
					if m.SlotId[i] != -1 && m.SlotStatus[i] == 32 { // is player and is playing
						plr := GetPlayerById(m.SlotId[i])
						m.SlotStatus[i] = 4
						m.InProgress = false
						Packets.WriteMatchUpdate(*plr, *m, int(player.Build))
						Packets.WriteMatchComplete(plr.Conn, int(player.Build))
					}
				}
				break
			}
		case Utils.CalculatePacketOffset(int(player.Build), 53): // loaded into game
			{
				m := player.CurrentMatch
				m.LoadingPeople--
				if m.LoadingPeople < 1 { // everybody loaded and we are happy
					for i := 0; i < 8; i++ {
						if m.SlotId[i] != -1 && m.SlotStatus[i] == 32 { // is player and is playing
							plr := GetPlayerById(m.SlotId[i])
							Packets.MatchAllPlayersLoaded(plr.Conn, int(player.Build))
						}
					}
				}

				break
			}
		case Utils.CalculatePacketOffset(int(player.Build), 42): // Beatmap change
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
				if player.Build > 590 {
					match.Password, err = Utils.ReadOsuString(buf)
					if err != nil {
						Utils.LogErr("Error occurred on Password creation: ", err.Error())
						return
					}
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
				if build2 > 394 {
					for i := 0; i < 8; i++ {
						err = binary.Read(buf, binary.LittleEndian, &match.SlotStatus[i])
					}
					if build2 > 553 {
						for i := 0; i < 8; i++ {
							err = binary.Read(buf, binary.LittleEndian, &match.SlotTeam[i])
						}
					}
					for i := 0; i < 8; i++ {
						if match.SlotStatus[i]&0x07C > 0 {
							err = binary.Read(buf, binary.LittleEndian, &match.SlotId[i])
						}
					}
				}
				err = binary.Read(buf, binary.LittleEndian, &match.HostId)
				if err != nil {
					Utils.LogErr("Error occurred on HostId creation: ", err.Error())
					return
				}
				if player.Build > 483 {

					err = binary.Read(buf, binary.LittleEndian, &match.Mode)
					if err != nil {
						Utils.LogErr("Error occurred on Mode creation: ", err.Error())
						return
					}

				}
				if player.Build > 535 {
					err = binary.Read(buf, binary.LittleEndian, &match.ScoringType)
					if err != nil {
						Utils.LogErr("Error occurred on ScoringType creation: ", err.Error())
						return
					}
					err = binary.Read(buf, binary.LittleEndian, &match.TeamType)
					if err != nil {
						Utils.LogErr("Error occurred on TeamType creation: ", err.Error())
						return
					}
				}
				for i := 0; i < 8; i++ {
					if match.SlotId[i] != -1 {
						plr := GetPlayerById(match.SlotId[i])
						Packets.WriteMatchUpdate(*plr, *match, int(player.Build))
					}
				}
				break
			}
		case Utils.CalculatePacketOffset(int(player.Build), 52): // Mods
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
						Packets.WriteMatchUpdate(*plr, *m, int(player.Build))
					}
				}
				break
			}
		case Utils.CalculatePacketOffset(int(player.Build), 21):
			{
				buf := bytes.NewReader(data)
				dat, err := Utils.ReadOsuString(buf)
				if err != nil {
					Utils.LogErr(err.Error())
					break
				}
				Utils.LogInfo("Osu! reported an error: %s", dat)
				break
			}
		case Utils.CalculatePacketOffset(int(player.Build), 40): // Ready
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
						Packets.WriteMatchUpdate(*plr, *match, int(player.Build))
					}
				}
				break
			}
		case Utils.CalculatePacketOffset(int(player.Build), 86):
			{
				buf := bytes.NewReader(data)
				var leng int16
				err := binary.Read(buf, binary.LittleEndian, &leng)
				if err != nil {
					Utils.LogErr("Failed to read infotype: ", err.Error())
					break
				}

				for _, player1 := range players {
					if player.Build > 1717 && player.Stats.UserID != player1.Stats.UserID {
						Packets.WriteUserPresense(player, *player1, 2, int(player.Build))
					}
					Packets.WriteUserStats(player, *player1, 2, int(player.Build))

				}
				break
			}
		case Utils.CalculatePacketOffset(int(player.Build), 80):
			{
				buf := bytes.NewReader(data)
				var infotype int32
				err := binary.Read(buf, binary.LittleEndian, &infotype)
				if err != nil {
					Utils.LogErr("Failed to read infotype: ", err.Error())
					break
				}

				for _, player1 := range players {
					if player.Build > 1717 && player.Stats.UserID != player1.Stats.UserID {
						Packets.WriteUserPresense(player, *player1, 2, int(player.Build))
					}
					Packets.WriteUserStats(player, *player1, 2, int(player.Build))

				}

				break
			}
		/*case 71: // TransferHost
		{

			match := player.CurrentMatch
			buf := bytes.NewReader(data)
			if match.HostId == player.Stats.UserID {
				var newhost int32
				var oldhostslot int
				binary.Read(buf, binary.LittleEndian, &newhost)
				for i := 0; i < 8; i++ {
					if match.SlotId[i] == match.HostId {
						oldhostslot = i
					}
				}
				match.SlotId[newhost] = match.SlotId[oldhostslot]
				match.SlotStatus[newhost] = match.SlotStatus[oldhostslot]
				match.SlotId[oldhostslot] = match.SlotId[newhost]
				match.SlotStatus[oldhostslot] = match.SlotStatus[newhost]
				match.HostId = newhost
				fmt.Println("newhostslot: ", newhost)
				fmt.Println("oldhostslot: ", oldhostslot)
				fmt.Println("newhost: ", newhost)

			}

			for i := 0; i < 8; i++ {
				if match.SlotId[i] != -1 {
					plr := GetPlayerById(match.SlotId[i])
					Packets.WriteMatchUpdate(*plr, *match)
				}
			}
			break
		}*/
		default:
			{
				Utils.LogWarning("Received unhandled packet %d", packetType)
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
	playersMu.Lock()
	for _, player1 := range players {
		if player1.Username != player.Username {
			Packets.WriteMessage(player1.Conn, sender, msg, target, int(player.Build))
		}
	}
	playersMu.Unlock()
	banchobotthoughts := BanchoBot.HandleMsg(sender, msg, target)
	if banchobotthoughts != "" && target == "#osu" {
		playersMu.Lock()
		for _, player1 := range players {
			lines := strings.Split(banchobotthoughts, "\n")
			for i := 0; i < len(lines); i++ {
				Packets.WriteMessage(player1.Conn, "BanchoBot", lines[i], "#osu", int(player.Build))
			}
		}
		playersMu.Unlock()
	}
	Utils.LogInfo(sender + "->" + target + ": " + msg)
}
func handleStatus(player *Structs.Player, data []byte) {
	var status byte
	var bmapUpdate bool
	var mods uint16
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &status)
	if err != nil {
		Utils.LogErr("Error occurred while reading Status from " + player.Username + " " + err.Error())
	}
	if player.Build < 1717 {
		err = binary.Read(buf, binary.LittleEndian, &bmapUpdate)
		if err != nil {
			Utils.LogErr("Error occurred while reading BeatmapUpdate from " + player.Username + " " + err.Error())
		}
		player.Status.BeatmapUpdate = bmapUpdate
	}

	player.Status.Status = status
	if bmapUpdate || player.Build > 1717 {
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
		if player.Build > 483 {
			err = binary.Read(buf, binary.LittleEndian, &player.Status.PlayMode)
			if err != nil {
				Utils.LogErr("Error occurred while reading PlayMode from " + player.Username)
			}
			err = binary.Read(buf, binary.LittleEndian, &player.Status.BeatmapId)
			if err != nil {
				Utils.LogErr("Error occurred while reading beatmapid from " + player.Username)
			}
		}

		player.Status.StatusText = statusText
		player.Status.BeatmapChecksum = beatmapMd5
		player.Status.CurrentMods = mods
	}
	playersMu.Lock()
	for _, player1 := range players {
		if player1.Build > 1717 {
			Packets.WriteUserStats(*player1, *player, 2, int(player1.Build))
			Packets.WriteUserPresense(*player1, *player, 2, int(player1.Build))
		} else {
			Packets.WriteUserStats(*player1, *player, 2, int(player1.Build))
		}
	}
	playersMu.Unlock()
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
		Packets.WriteIrcQuit(player1.Conn, username, int(player1.Build))
		Packets.WriteUserQuit(player1.Conn, id, int(player1.Build))
	}
}
func AddMatch(match *Structs.Match) {
	MpMatchesMu.Lock()
	defer MpMatchesMu.Unlock()
	MpMatches[match.MatchId] = match
	playersMu.Lock()
	for _, player1 := range players {
		if player1.IsInLobby {
			Packets.WriteMatchUpdate(*player1, *match, int(player1.Build))
		}
	}
	playersMu.Unlock()
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
			Packets.WriteMatchUpdate(*plr, *match, int(plr.Build))
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
		Packets.WriteMatchJoinFail(player.Conn, int(player.Build))
		return false
	} else {
		match.SlotId[newslot] = player.Stats.UserID
		match.SlotStatus[newslot] = 4
		for i := 0; i < 8; i++ {
			if match.SlotId[i] != -1 && match.SlotId[i] != player.Stats.UserID {
				plr := GetPlayerById(match.SlotId[i])
				Packets.WriteMatchUpdate(*plr, *match, int(plr.Build))
			}
		}
		Packets.WriteMatchJoinSuccess(player.Conn, *match, int32(player.Build))
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
			Packets.WriteMatchUpdate(*plr, *match, int(plr.Build))
		}
	}
	setnewleader := false
	peopleinlobby := 0
	for i := 0; i < 8; i++ {
		if match.SlotId[i] != -1 {
			if !setnewleader {
				match.HostId = match.SlotId[i]
				setnewleader = true
			}
			peopleinlobby++
		}
	}
	if peopleinlobby == 0 {
		// lobby empty, delete it
		RemoveMatch(match.MatchId)
		playersMu.Lock()
		for _, player1 := range players {
			if player1.IsInLobby {
				Packets.WriteDisbandMatch(player1.Conn, int(match.MatchId), int(player1.Build))
			}
		}
		playersMu.Unlock()
	}
}

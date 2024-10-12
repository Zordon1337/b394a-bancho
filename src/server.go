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
	addr      string     = ":13381"
	players              = make(map[string]*Structs.Player) // Map to hold logged-in players
	playersMu sync.Mutex                                    // Mutex to synchronize access to the player list
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
	username := lines[0]
	//md5Hash := lines[1]
	//clientInfo := lines[2]
	fmt.Printf("Login %s\n", string(initialMessage[:n]))
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
		Username: (username),
		Conn:     client,
		Stats:    stats,
		Status:   status,
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
		var flag bool
		var dataLength int32
		buf := bytes.NewReader(header)

		if err := binary.Read(buf, binary.LittleEndian, &packetType); err != nil {
			fmt.Println("Failed to read packet type:", err)
			removePlayer(player.Username, player.Stats.UserID)
			return
		}
		if err := binary.Read(buf, binary.LittleEndian, &flag); err != nil {
			fmt.Println("Failed to read bool flag:", err)
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
		case 1:
			{
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
				break
			}
		case 2:
			{
				removePlayer(player.Username, player.Stats.UserID)
				return
			}
		case 4:
			{
				break
			}
		case 3:
			{
				Packets.WriteUserStats(player.Conn, player, 2)
				break
			}

		default:
			{
				fmt.Printf("Packet Type: %d, Bool Flag: %v, Data Length: %d\n", packetType, flag, dataLength)
				break
			}
		}
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

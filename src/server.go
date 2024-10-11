package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"socket-server/src/Packets"
	"time"
)

var addr string = ":13381"

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
			handleClient(client)
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
	fmt.Printf("Login %s\n", string(initialMessage[:n]))
	Packets.WriteLoginReply(client)
	Packets.WriteUserStats(client)
	Packets.WriteChannelJoinSucess(client, "#osu")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				Packets.WritePing(client)
				if err != nil {
					fmt.Println("Error sending ping to client:", err)
					return
				}
				fmt.Println("Pinged client")
			}
		}
	}()
	for {

		header := make([]byte, 7) // int16 + bool + int32
		_, err := io.ReadFull(client, header)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return
		}
		var packetType int16
		var flag bool
		var dataLength int32
		buf := bytes.NewReader(header)

		if err := binary.Read(buf, binary.LittleEndian, &packetType); err != nil {
			fmt.Println("Failed to read packet type:", err)
			return
		}
		if err := binary.Read(buf, binary.LittleEndian, &flag); err != nil {
			fmt.Println("Failed to read bool flag:", err)
			return
		}
		if err := binary.Read(buf, binary.LittleEndian, &dataLength); err != nil {
			fmt.Println("Failed to read data length:", err)
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
		case 4:
			{
				fmt.Println("Client sent pong")
				break
			}
		default:
			{
				fmt.Printf("Packet Type: %d, Bool Flag: %v, Data Length: %d\n", packetType, flag, dataLength)
			}
		}
	}
}

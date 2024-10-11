package main

import (
	"fmt"
	"net"
	"socket-server/src/Packets"
)

var addr string = ":13381"

func main() {
	var ln, err = net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Socket listening for clients on port 6969")
	fmt.Println("Lets play osu!")
	for {
		client, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			ln.Close() // a bit brutal but okay
		}
		if client != nil {
			fmt.Println("Client connected, Welcome")
			Packets.WriteLoginReply(client)
			Packets.WriteUserStats(client)
			Packets.WriteChannelJoinSucess(client, "#osu")
		}
	}
}

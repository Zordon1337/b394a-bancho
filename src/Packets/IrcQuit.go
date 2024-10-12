package Packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

func WriteIrcQuit(client net.Conn, username string) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, WriteOsuString(username))
	if err != nil {
		fmt.Println("An error occurred while sending IrcQuit from " + username)
		return
	}
	resp, err := SerializePacket(10, buf.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}

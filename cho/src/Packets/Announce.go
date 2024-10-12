package Packets

import (
	"bytes"
	"encoding/binary"
	"net"
	"socket-server/src/Utils"
)

func WriteAnnounce(client net.Conn, msg string) {
	if client == nil {
		return
	}
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, Utils.WriteOsuString(msg))
	if err != nil {
		return
	}
	resp, err := Utils.SerializePacket(25, buffer.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}

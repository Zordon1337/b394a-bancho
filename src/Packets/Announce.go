package Packets

import (
	"bytes"
	"encoding/binary"
	"net"
)

func WriteAnnounce(client net.Conn, msg string) {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, WriteOsuString(msg))
	if err != nil {
		return
	}
	resp, err := SerializePacket(25, buffer.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}

package Packets

import (
	"bytes"
	"encoding/binary"
	"net"
)

func WriteLoginReply(client net.Conn, response int32) {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, int32(response))
	if err != nil {
		return
	}
	resp, err := SerializePacket(5, buffer.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}

package Packets

import (
	"bytes"
	"encoding/binary"
	"net"
)

func WriteChannelJoinSucess(client net.Conn, channel string) {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, []byte(channel))
	if err != nil {
		return
	}
	resp, err := SerializePacket(65, buffer.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}

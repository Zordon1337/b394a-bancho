package Packets

import (
	"bytes"
	"encoding/binary"
	"net"
)

func WriteUserQuit(client net.Conn, userid int32) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, userid)
	if err != nil {
		return
	}
	resp, err := SerializePacket(13, buf.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}

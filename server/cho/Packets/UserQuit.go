package Packets

import (
	"bytes"
	"encoding/binary"
	"net"
	"retsu/Utils"
)

func WriteUserQuit(client net.Conn, userid int32) {
	if client == nil {
		return
	}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, userid)
	if err != nil {
		return
	}
	resp, err := Utils.SerializePacket(13, buf.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}

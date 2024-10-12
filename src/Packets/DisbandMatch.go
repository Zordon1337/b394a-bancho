package Packets

import (
	"bytes"
	"encoding/binary"
	"net"
)

func WriteDisbandMatch(client net.Conn, matchid int) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, int32(matchid))
	if err != nil {
		return
	}
	resp, err := SerializePacket(29, buf.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}

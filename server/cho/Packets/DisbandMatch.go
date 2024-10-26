package Packets

import (
	"bytes"
	"encoding/binary"
	"net"
	"retsu/Utils"
)

func WriteDisbandMatch(client net.Conn, matchid int) {
	if client == nil {
		return
	}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, int32(matchid))
	if err != nil {
		return
	}
	resp, err := Utils.SerializePacket(29, buf.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}

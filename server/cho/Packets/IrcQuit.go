package Packets

import (
	"bytes"
	"encoding/binary"
	"net"
	"retsu/Utils"
)

func WriteIrcQuit(client net.Conn, username string, build int) {
	if client == nil {
		return
	}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, Utils.WriteOsuString(username))
	if err != nil {
		Utils.LogErr("An error occurred while sending IrcQuit from " + username)
		return
	}
	resp, err := Utils.SerializePacket(int16(Utils.CalculatePacketOffset(build, int(10))), buf.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}

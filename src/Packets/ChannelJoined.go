package Packets

import (
	"bytes"
	"encoding/binary"
	"net"
	"socket-server/src/Utils"
)

func WriteChannelJoinSucess(client net.Conn, channel string) {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, Utils.WriteOsuString(channel))
	if err != nil {
		return
	}
	resp, err := Utils.SerializePacket(65, buffer.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}

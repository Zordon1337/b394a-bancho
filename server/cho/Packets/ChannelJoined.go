package Packets

import (
	"bytes"
	"encoding/binary"
	"net"
	"retsu/Utils"
)

func WriteChannelJoinSuccess(client net.Conn, channel string) {
	if client == nil {
		return
	}
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

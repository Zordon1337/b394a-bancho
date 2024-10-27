package Packets

import (
	"bytes"
	"encoding/binary"
	"net"
	"retsu/Utils"
)

func WriteMessage(client net.Conn, sender string, message string, target string, build int) {
	if client == nil {
		return
	}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, Utils.WriteOsuString(sender))
	if err != nil {
		return
	}
	err = binary.Write(buf, binary.LittleEndian, Utils.WriteOsuString(message))
	if err != nil {
		return
	}
	err = binary.Write(buf, binary.LittleEndian, Utils.WriteOsuString(target))
	if err != nil {
		return
	}
	resp, err := Utils.SerializePacket(int16(Utils.CalculatePacketOffset(int(build), int(7))), buf.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}

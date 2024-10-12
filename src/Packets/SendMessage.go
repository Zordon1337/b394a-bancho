package Packets

import (
	"bytes"
	"encoding/binary"
	"net"
)

func WriteMessage(client net.Conn, sender string, message string, target string) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, WriteOsuString(sender))
	if err != nil {
		return
	}
	err = binary.Write(buf, binary.LittleEndian, WriteOsuString(message))
	if err != nil {
		return
	}
	err = binary.Write(buf, binary.LittleEndian, WriteOsuString(target))
	if err != nil {
		return
	}
	resp, err := SerializePacket(7, buf.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}

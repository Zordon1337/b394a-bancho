package Packets

import (
	"bytes"
	"encoding/binary"
	"net"
	"retsu/Utils"
)

func WriteProtcolNegotiation(client net.Conn, response int32) {
	if client == nil {
		return
	}
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, int32(response))
	if err != nil {
		return
	}
	resp, err := Utils.SerializePacket(76, buffer.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}

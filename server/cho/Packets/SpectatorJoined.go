package Packets

import (
	"bytes"
	"encoding/binary"
	"net"
	"retsu/Utils"
)

func WriteFellowSpecJoined(client net.Conn, response int32, build int) {
	if client == nil {
		return
	}
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, int32(response))
	if err != nil {
		return
	}
	resp, err := Utils.SerializePacket(int16(Utils.CalculatePacketOffset(build, int(43))), buffer.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}
func WriteSpecJoined(client net.Conn, response int32, build int) {
	if client == nil {
		return
	}
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, int32(response))
	if err != nil {
		return
	}
	resp, err := Utils.SerializePacket(int16(Utils.CalculatePacketOffset(build, int(14))), buffer.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}

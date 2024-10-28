package Packets

import (
	"bytes"
	"encoding/binary"
	"net"
	"retsu/Utils"
)

func WriteFellowSpecQuited(client net.Conn, response int32, build int) {
	if client == nil {
		return
	}
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, int32(response))
	if err != nil {
		return
	}
	resp, err := Utils.SerializePacket(int16(Utils.CalculatePacketOffset(build, int(44))), buffer.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}
func WriteSpecQuited(client net.Conn, response int32, build int) {
	if client == nil {
		return
	}
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, int32(response))
	if err != nil {
		return
	}
	resp, err := Utils.SerializePacket(int16(Utils.CalculatePacketOffset(build, int(15))), buffer.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}

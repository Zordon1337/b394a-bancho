package Packets

import (
	"net"
	"retsu/Utils"
)

func WriteMatchJoinFail(client net.Conn, build int) {
	if client == nil {
		return
	}

	resp, err := Utils.SerializePacket(int16(Utils.CalculatePacketOffset(build, int(38))), []byte{})
	if err != nil {
		return
	}
	client.Write(resp)
}

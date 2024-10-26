package Packets

import (
	"net"
	"retsu/Utils"
)

func WriteMatchJoinFail(client net.Conn) {
	if client == nil {
		return
	}

	resp, err := Utils.SerializePacket(38, []byte{})
	if err != nil {
		return
	}
	client.Write(resp)
}

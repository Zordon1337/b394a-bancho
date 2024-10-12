package Packets

import (
	"net"
	"socket-server/src/Utils"
)

func WriteMatchJoinFail(client net.Conn) {

	resp, err := Utils.SerializePacket(38, []byte{})
	if err != nil {
		return
	}
	client.Write(resp)
}

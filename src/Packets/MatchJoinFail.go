package Packets

import (
	"net"
)

func WriteMatchJoinFail(client net.Conn) {

	resp, err := SerializePacket(38, []byte{})
	if err != nil {
		return
	}
	client.Write(resp)
}

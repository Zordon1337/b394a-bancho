package Packets

import (
	"bytes"
	"encoding/binary"
	"net"
	"socket-server/src/Structs"
)

func WriteMatchUpdate(client net.Conn, match Structs.Match) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, match.MatchId)
	if err != nil {
		return
	}
	err = binary.Write(buf, binary.LittleEndian, match.InProgress)
	if err != nil {
		return
	}
	err = binary.Write(buf, binary.LittleEndian, match.MatchType)
	if err != nil {
		return
	}
	err = binary.Write(buf, binary.LittleEndian, match.ActiveMods)
	if err != nil {
		return
	}
	err = binary.Write(buf, binary.LittleEndian, WriteOsuString(match.GameName))
	if err != nil {
		return
	}
	err = binary.Write(buf, binary.LittleEndian, WriteOsuString(match.BeatmapName))
	if err != nil {
		return
	}
	err = binary.Write(buf, binary.LittleEndian, match.BeatmapId)
	if err != nil {
		return
	}
	err = binary.Write(buf, binary.LittleEndian, WriteOsuString(match.BeatmapChecksum))
	if err != nil {
		return
	}
	for i := 0; i < 8; i++ {
		err = binary.Write(buf, binary.LittleEndian, match.SlotStatus[i])
		if err != nil {
			return
		}
	}
	for i := 0; i < 8; i++ {
		err = binary.Write(buf, binary.LittleEndian, match.SlotId[i])
		if err != nil {
			return
		}
	}
	resp, err := SerializePacket(27, buf.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}

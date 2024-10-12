package Structs

import (
	"bytes"
	"encoding/binary"
	"socket-server/src/Utils"
)

/*
SlotStatus:
Open - 1
Locked - 2
NotReady - 4
Ready - 8
NoMap - 16
Playing - 32
Complete - 64
CompHasPlayer - 124
*/
type Match struct {
	MatchId         byte
	InProgress      bool
	MatchType       byte
	ActiveMods      int16
	GameName        string
	BeatmapName     string
	BeatmapId       int32
	BeatmapChecksum string
	SlotStatus      [8]byte
	SlotId          [8]int32
	LoadingPeople   int32
}

func GetBytesFromMatch(match *Match) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, match.MatchId)
	if err != nil {
		return nil
	}
	err = binary.Write(buf, binary.LittleEndian, match.InProgress)
	if err != nil {
		return nil
	}
	err = binary.Write(buf, binary.LittleEndian, match.MatchType)
	if err != nil {
		return nil
	}
	err = binary.Write(buf, binary.LittleEndian, match.ActiveMods)
	if err != nil {
		return nil
	}
	err = binary.Write(buf, binary.LittleEndian, Utils.WriteOsuString(match.GameName))
	if err != nil {
		return nil
	}
	err = binary.Write(buf, binary.LittleEndian, Utils.WriteOsuString(match.BeatmapName))
	if err != nil {
		return nil
	}
	err = binary.Write(buf, binary.LittleEndian, match.BeatmapId)
	if err != nil {
		return nil
	}
	err = binary.Write(buf, binary.LittleEndian, Utils.WriteOsuString(match.BeatmapChecksum))
	if err != nil {
		return nil
	}
	for i := 0; i < 8; i++ {
		err = binary.Write(buf, binary.LittleEndian, match.SlotStatus[i])
		if err != nil {
			return nil
		}
	}
	for i := 0; i < 8; i++ {
		err = binary.Write(buf, binary.LittleEndian, match.SlotId[i])
		if err != nil {
			return nil
		}
	}
	return buf.Bytes()
}

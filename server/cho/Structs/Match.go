package Structs

import (
	"bytes"
	"encoding/binary"
	"retsu/Utils"
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
	MatchId              byte
	InProgress           bool
	MatchType            byte
	ActiveMods           int16
	GameName             string
	Password             string
	BeatmapName          string
	BeatmapId            int32
	BeatmapChecksum      string
	SlotStatus           [8]byte
	SlotTeam             [8]byte
	SlotId               [8]int32
	HostId               int32
	Mode                 byte
	ScoringType          byte
	TeamType             byte
	LoadingPeople        int32
	SkippingNeededToSkip int32
}

func GetBytesFromMatch(match *Match, build int32) []byte {
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
	if build > 590 {

		err = binary.Write(buf, binary.LittleEndian, Utils.WriteOsuString(match.Password))
		if err != nil {
			return nil
		}
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
	if build > 553 {
		for i := 0; i < 8; i++ {
			err = binary.Write(buf, binary.LittleEndian, match.SlotTeam[i])
			if err != nil {
				return nil
			}
		}
	}
	for i := 0; i < 8; i++ {
		if match.SlotStatus[i]&0x07C > 0 {
			err = binary.Write(buf, binary.LittleEndian, match.SlotId[i])
			if err != nil {
				return nil
			}
		}

	}
	if build > 399 {
		err = binary.Write(buf, binary.LittleEndian, match.HostId)
		if err != nil {
			return nil
		}
	}
	if build > 483 {
		err = binary.Write(buf, binary.LittleEndian, byte(match.Mode))
		if err != nil {
			return nil
		}
	}
	if build > 535 {
		err = binary.Write(buf, binary.LittleEndian, byte(match.ScoringType))
		if err != nil {
			return nil
		}
		err = binary.Write(buf, binary.LittleEndian, byte(match.TeamType))
		if err != nil {
			return nil
		}
	}
	return buf.Bytes()
}

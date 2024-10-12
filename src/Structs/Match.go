package Structs

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
}

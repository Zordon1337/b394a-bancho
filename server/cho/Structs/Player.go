package Structs

import (
	"net"
)

type Player struct {
	Username            string
	Status              Status
	Stats               UserStats
	Conn                net.Conn
	CurrentMatch        *Match
	IsInLobby           bool
	Timezone            byte
	Country             string
	Build               int64
	Spectators          map[int32]*Player
	CurrentlySpectating *Player
}

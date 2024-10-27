package Packets

import (
	"bytes"
	"encoding/binary"
	"retsu/Utils"
	"retsu/cho/Structs"
)

func WriteUserPresense(receiver Structs.Player, user Structs.Player, completeness byte, build int) {
	if receiver.Conn == nil {
		return
	}
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.LittleEndian, int32(user.Stats.UserID)) // UserId
	if err != nil {
		return
	}
	output := Utils.WriteOsuString(user.Username)
	err = binary.Write(buffer, binary.LittleEndian, output) // Username
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, byte(1)) // Avatar Extension
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, byte(user.Timezone+24)) // Timezone
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, rune('P')) // CountryId
	if err != nil {
		return
	}
	output = Utils.WriteOsuString("")
	err = binary.Write(buffer, binary.LittleEndian, output) // City
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, byte(0)) // Perms
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, float32(0)) // Longitude
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, float32(0)) // Latitude
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, int32(user.Stats.Rank)) // Rank
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, byte(user.Status.PlayMode)) // PlayMode
	if err != nil {
		return
	}
	resp, err := Utils.SerializePacket(int16(Utils.CalculatePacketOffset(int(build), int(84))), buffer.Bytes())
	if err != nil {
		return
	}
	receiver.Conn.Write(resp)

}

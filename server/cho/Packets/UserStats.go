package Packets

import (
	"bytes"
	"encoding/binary"
	"retsu/Utils"
	"retsu/cho/Structs"
	"strconv"
)

func GetStatusUpdate(receiver Structs.Player, user Structs.Player) ([]byte, error) {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, byte(user.Status.Status)) // status, Idle for testing
	if err != nil {
		return nil, err
	}
	err = binary.Write(buffer, binary.LittleEndian, user.Status.BeatmapUpdate)
	if err != nil {
		return nil, err
	}
	if user.Status.BeatmapUpdate {
		err = binary.Write(buffer, binary.LittleEndian, Utils.WriteOsuString(user.Status.StatusText))
		if err != nil {
			return nil, err
		}
		err = binary.Write(buffer, binary.LittleEndian, Utils.WriteOsuString(user.Status.BeatmapChecksum))
		if err != nil {
			return nil, err
		}
		err = binary.Write(buffer, binary.LittleEndian, uint16(user.Status.CurrentMods))
		if err != nil {
			return nil, err
		}
		if receiver.Build > 483 {
			err = binary.Write(buffer, binary.LittleEndian, user.Status.PlayMode)
			if err != nil {
				return nil, err
			}
			err = binary.Write(buffer, binary.LittleEndian, int32(user.Status.BeatmapId))
			if err != nil {
				return nil, err
			}
		}
	}
	return buffer.Bytes(), nil
}
func WriteUserStats(receiver Structs.Player, user Structs.Player, completeness byte) {
	if receiver.Conn == nil {
		return
	}
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.LittleEndian, int32(user.Stats.UserID)) // UserId
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, completeness) // Completeness
	if err != nil {
		return
	}
	statusupdate, err := GetStatusUpdate(receiver, user)
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, statusupdate) // Status
	if err != nil {
		return
	}
	if completeness > 0 {
		err = binary.Write(buffer, binary.LittleEndian, int64(user.Stats.RankedScore)) // Ranked Score
		if err != nil {
			return
		}
		err = binary.Write(buffer, binary.LittleEndian, float32(user.Stats.Accuracy)) // Accuracy
		if err != nil {
			return
		}
		err = binary.Write(buffer, binary.LittleEndian, int32(user.Stats.PlayCount)) // Play count
		if err != nil {
			return
		}
		err = binary.Write(buffer, binary.LittleEndian, int64(user.Stats.TotalScore)) // Total Score
		if err != nil {
			return
		}
		if receiver.Build > 504 {
			err = binary.Write(buffer, binary.LittleEndian, int32(user.Stats.Rank)) // Rank
			if err != nil {
				return
			}
		} else {
			err = binary.Write(buffer, binary.LittleEndian, uint16(user.Stats.Rank)) // Rank
			if err != nil {
				return
			}
		}
	}
	// full
	if completeness == 2 {
		err = binary.Write(buffer, binary.LittleEndian, Utils.WriteOsuString(user.Username)) // username
		if err != nil {
			return
		}
		err = binary.Write(buffer, binary.LittleEndian, Utils.WriteOsuString(strconv.Itoa(int(user.Stats.UserID)))) // pfp file name
		if err != nil {
			return
		}
		err = binary.Write(buffer, binary.LittleEndian, byte(user.Timezone)) // Timezone
		if err != nil {
			return
		}
		err = binary.Write(buffer, binary.LittleEndian, Utils.WriteOsuString(user.Country)) // Country
		if err != nil {
			return
		}
		if receiver.Build > 394 {
			err = binary.Write(buffer, binary.LittleEndian, byte(0)) // Perms
			if err != nil {
				return
			}
		}
		if receiver.Build > 504 {

			err = binary.Write(buffer, binary.LittleEndian, float32(0)) // Longitude
			if err != nil {
				return
			}
			err = binary.Write(buffer, binary.LittleEndian, float32(0)) // Latitude
			if err != nil {
				return
			}
		}
	}
	resp, err := Utils.SerializePacket(12, buffer.Bytes())
	if err != nil {
		return
	}
	receiver.Conn.Write(resp)
}

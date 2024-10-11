package Packets

import (
	"bytes"
	"encoding/binary"
	"net"
)

func GetStatusUpdate() ([]byte, error) {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, byte(0)) // status, Idle for testing
	if err != nil {
		return nil, err
	}
	err = binary.Write(buffer, binary.LittleEndian, false) // status, Idle for testing
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
func WriteUserStats(client net.Conn) {
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.LittleEndian, int32(1)) // UserId
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, byte(1)) // Completeness
	if err != nil {
		return
	}
	statusupdate, err := GetStatusUpdate()
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, statusupdate) // Status
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, int64(1337420)) // Ranked Score
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, float32(1)) // Accuracy
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, int32(0)) // Play count
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, int64(1337420)) // Total Score
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, int32(0)) // Rank
	if err != nil {
		return
	}
	resp, err := SerializePacket(12, buffer.Bytes())
	if err != nil {
		return
	}
	client.Write(resp)
}

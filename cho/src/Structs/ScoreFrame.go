package Structs

import (
	"bytes"
	"encoding/binary"
	"socket-server/src/Utils"
)

type ScoreFrame struct {
	Time         int32
	SlotId       byte
	Count300     uint16
	Count100     uint16
	Count50      uint16
	CountGeki    uint16
	CountKatu    uint16
	CountMiss    uint16
	TotalScore   int32
	MaxCombo     uint16
	CurrentCombo uint16
	Perfect      bool
	CurrentHP    byte
	TagByte      byte
}

func ReadScoreFrameFromStream(data []byte) ScoreFrame {
	frame := new(ScoreFrame)
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &frame.Time)
	if err != nil {
		Utils.LogErr("Error occurred on Time creation: ", err.Error())
	}
	err = binary.Read(buf, binary.LittleEndian, &frame.SlotId)
	if err != nil {
		Utils.LogErr("Error occurred on SlotId creation: ", err.Error())
	}
	err = binary.Read(buf, binary.LittleEndian, &frame.Count300)
	if err != nil {
		Utils.LogErr("Error occurred on Count300 creation: ", err.Error())
	}
	err = binary.Read(buf, binary.LittleEndian, &frame.Count100)
	if err != nil {
		Utils.LogErr("Error occurred on Count100 creation: ", err.Error())
	}
	err = binary.Read(buf, binary.LittleEndian, &frame.Count50)
	if err != nil {
		Utils.LogErr("Error occurred on Count50 creation: ", err.Error())
	}
	err = binary.Read(buf, binary.LittleEndian, &frame.CountGeki)
	if err != nil {
		Utils.LogErr("Error occurred on CountGeki creation: ", err.Error())
	}
	err = binary.Read(buf, binary.LittleEndian, &frame.CountKatu)
	if err != nil {
		Utils.LogErr("Error occurred on CountKatu creation: ", err.Error())
	}
	err = binary.Read(buf, binary.LittleEndian, &frame.CountMiss)
	if err != nil {
		Utils.LogErr("Error occurred on CountMiss creation: ", err.Error())
	}
	err = binary.Read(buf, binary.LittleEndian, &frame.TotalScore)
	if err != nil {
		Utils.LogErr("Error occurred on TotalScore creation: ", err.Error())
	}
	err = binary.Read(buf, binary.LittleEndian, &frame.MaxCombo)
	if err != nil {
		Utils.LogErr("Error occurred on MaxCombo creation: ", err.Error())
	}
	err = binary.Read(buf, binary.LittleEndian, &frame.CurrentCombo)
	if err != nil {
		Utils.LogErr("Error occurred on CurrentCombo creation: ", err.Error())
	}
	err = binary.Read(buf, binary.LittleEndian, &frame.Perfect)
	if err != nil {
		Utils.LogErr("Error occurred on Perfect creation: ", err.Error())
	}
	err = binary.Read(buf, binary.LittleEndian, &frame.CurrentHP)
	if err != nil {
		Utils.LogErr("Error occurred on CurrentHP creation: ", err.Error())
	}
	return *frame
}
func WriteScoreFrameToBytes(frame ScoreFrame) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &frame.Time)
	if err != nil {
		Utils.LogErr("Error occurred on Time creation: ", err.Error())
	}
	err = binary.Write(buf, binary.LittleEndian, &frame.SlotId)
	if err != nil {
		Utils.LogErr("Error occurred on Time creation: ", err.Error())
	}
	err = binary.Write(buf, binary.LittleEndian, &frame.Count300)
	if err != nil {
		Utils.LogErr("Error occurred on Time creation: ", err.Error())
	}
	err = binary.Write(buf, binary.LittleEndian, &frame.Count100)
	if err != nil {
		Utils.LogErr("Error occurred on Time creation: ", err.Error())
	}
	err = binary.Write(buf, binary.LittleEndian, &frame.Count50)
	if err != nil {
		Utils.LogErr("Error occurred on Time creation: ", err.Error())
	}
	err = binary.Write(buf, binary.LittleEndian, &frame.CountGeki)
	if err != nil {
		Utils.LogErr("Error occurred on Time creation: ", err.Error())
	}
	err = binary.Write(buf, binary.LittleEndian, &frame.CountKatu)
	if err != nil {
		Utils.LogErr("Error occurred on Time creation: ", err.Error())
	}
	err = binary.Write(buf, binary.LittleEndian, &frame.CountMiss)
	if err != nil {
		Utils.LogErr("Error occurred on Time creation: ", err.Error())
	}
	err = binary.Write(buf, binary.LittleEndian, &frame.TotalScore)
	if err != nil {
		Utils.LogErr("Error occurred on Time creation: ", err.Error())
	}
	err = binary.Write(buf, binary.LittleEndian, &frame.MaxCombo)
	if err != nil {
		Utils.LogErr("Error occurred on Time creation: ", err.Error())
	}
	err = binary.Write(buf, binary.LittleEndian, &frame.CurrentCombo)
	if err != nil {
		Utils.LogErr("Error occurred on Time creation: ", err.Error())
	}
	err = binary.Write(buf, binary.LittleEndian, &frame.Perfect)
	if err != nil {
		Utils.LogErr("Error occurred on Time creation: ", err.Error())
	}
	err = binary.Write(buf, binary.LittleEndian, &frame.CurrentHP)
	if err != nil {
		Utils.LogErr("Error occurred on Time creation: ", err.Error())
	}
	err = binary.Write(buf, binary.LittleEndian, &frame.TagByte)
	if err != nil {
		Utils.LogErr("Error occurred on Time creation: ", err.Error())
	}
	return buf.Bytes()
}

package Utils

import (
	"strconv"
	"strings"
)

type Score struct {
	FileChecksum        string
	Username            string
	OnlineScoreChecksum string
	Count300            int32
	Count100            int32
	Count50             int32
	CountGeki           int32
	CountKatu           int32
	CountMiss           int32
	TotalScore          int64
	MaxCombo            int32
	Perfect             bool
	Ranking             string
	EnabledMods         string
	Pass                string
}

func getInt(val string) int32 {
	val1, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return int32(val1)
}
func getInt64(val string) int64 {
	val1, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return int64(val1)
}
func getBool(bo string) bool {
	bo1, err := strconv.ParseBool(bo)
	if err != nil {
		return false
	}
	return bo1
}
func FormattedToScore(formatted string) Score {
	values := strings.Split(formatted, ":")
	score := Score{
		FileChecksum:        values[0],
		Username:            values[1],
		OnlineScoreChecksum: values[2],
		Count300:            getInt(values[3]),
		Count100:            getInt(values[4]),
		Count50:             getInt(values[5]),
		CountGeki:           getInt(values[6]),
		CountKatu:           getInt(values[7]),
		CountMiss:           getInt(values[8]),
		TotalScore:          getInt64(values[9]),
		MaxCombo:            getInt(values[10]),
		Perfect:             getBool(values[11]),
		Ranking:             values[12],
		EnabledMods:         values[13],
		Pass:                values[14],
	}
	return score
}
func CalculateAccuracy(score Score) float32 {
	totalScore := float32(score.Count50*50 + score.Count100*100 + score.Count300*300 + score.CountGeki*300 + score.CountKatu*100)
	totalHits := float32(score.Count50 + score.Count100 + score.Count300 + score.CountGeki + score.CountKatu + score.CountMiss)

	if totalHits > 0 {
		return totalScore / (totalHits * 300)
	} else {
		return 0
	}
}

package Utils

import (
	"strconv"
	"strings"
	"time"
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
	Playmode            string
	Date                string
}

func GetInt(val string) int32 {
	val1, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return int32(val1)
}
func GetInt64(val string) int64 {
	val1, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return int64(val1)
}
func GetBool(bo string) bool {
	bo1, err := strconv.ParseBool(bo)
	if err != nil {
		return false
	}
	return bo1
}
func FormattedToScore(formatted string, isnewscore bool) Score {
	values := strings.Split(formatted, ":")
	score := Score{
		FileChecksum:        values[0],
		Username:            values[1],
		OnlineScoreChecksum: values[2],
		Count300:            GetInt(values[3]),
		Count100:            GetInt(values[4]),
		Count50:             GetInt(values[5]),
		CountGeki:           GetInt(values[6]),
		CountKatu:           GetInt(values[7]),
		CountMiss:           GetInt(values[8]),
		TotalScore:          GetInt64(values[9]),
		MaxCombo:            GetInt(values[10]),
		Perfect:             GetBool(values[11]),
		Ranking:             values[12],
		EnabledMods:         values[13],
		Pass:                values[14],
	}
	if isnewscore {
		score.Playmode = values[15]
		score.Date = values[16]
	} else {
		score.Playmode = "0"
		score.Date = time.Now().Format("2006-01-02 15:04:05")
	}
	return score
}
func CalculateAccuracy(score Score, isCtb bool) float32 {
	if !isCtb {
		totalScore := float32(score.Count50*50 + score.Count100*100 + score.Count300*300 + score.CountGeki*300 + score.CountKatu*100)
		totalHits := float32(score.Count50 + score.Count100 + score.Count300 + score.CountGeki + score.CountKatu + score.CountMiss)

		if totalHits > 0 {
			return totalScore / (totalHits * 300)
		} else {
			return 0
		}
	} else {

		totalScore := float32(score.Count50 + score.Count100 + score.Count300 + score.CountGeki + score.CountKatu)
		totalHits := float32(score.Count50 + score.Count100 + score.Count300 + score.CountGeki + score.CountKatu)

		if totalHits > 0 {
			return (totalScore / totalHits * 300)
		} else {
			return 0
		}
	}
	return 0
}

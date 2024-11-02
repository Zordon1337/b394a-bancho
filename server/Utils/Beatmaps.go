package Utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Beatmap struct {
	BeatmapsetID        int       `json:"beatmapset_id"`
	BeatmapID           int       `json:"beatmap_id"`
	Approved            int       `json:"approved"`
	TotalLength         int       `json:"total_length"`
	HitLength           int       `json:"hit_length"`
	Version             string    `json:"version"`
	FileMD5             string    `json:"file_md5"`
	DiffSize            float64   `json:"diff_size"`
	DiffOverall         float64   `json:"diff_overall"`
	DiffApproach        float64   `json:"diff_approach"`
	DiffDrain           float64   `json:"diff_drain"`
	Mode                int       `json:"mode"`
	ApprovedDate        time.Time `json:"approved_date"`
	LastUpdate          time.Time `json:"last_update"`
	Artist              string    `json:"artist"`
	ArtistUnicode       string    `json:"artist_unicode"`
	Title               string    `json:"title"`
	TitleUnicode        string    `json:"title_unicode"`
	Creator             string    `json:"creator"`
	CreatorID           int       `json:"creator_id"`
	BPM                 float64   `json:"bpm"`
	Source              string    `json:"source"`
	Tags                string    `json:"tags"`
	GenreID             int       `json:"genre_id"`
	LanguageID          int       `json:"language_id"`
	FavouriteCount      int       `json:"favourite_count"`
	Storyboard          int       `json:"storyboard"`
	Video               int       `json:"video"`
	DownloadUnavailable int       `json:"download_unavailable"`
	Playcount           int       `json:"playcount"`
	Passcount           int       `json:"passcount"`
	Packs               []string  `json:"packs"`
	MaxCombo            int       `json:"max_combo"`
	DifficultyRating    float64   `json:"difficultyrating"`
}

func GetBeatmap(md5 string) Beatmap {
	url := "https://osu.direct/api/get_beatmaps?h=" + md5
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: status code %d\n", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	var beatmaps []Beatmap
	err = json.Unmarshal([]byte(body), &beatmaps)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}
	var result *Beatmap
	for _, beatmap := range beatmaps {
		if beatmap.FileMD5 == md5 {
			result = &beatmap
			break
		}
	}
	return *result
}
func GetBeatmapsBySetId(id string) []Beatmap {
	url := "https://osu.direct/api/get_beatmaps?s=" + id
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: status code %d\n", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	var beatmaps []Beatmap
	err = json.Unmarshal([]byte(body), &beatmaps)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}
	return *&beatmaps
}
func GetBeatmapById(id string) Beatmap {
	url := "https://osu.direct/api/get_beatmaps?s=" + id
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: status code %d\n", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	var beatmaps []Beatmap
	err = json.Unmarshal([]byte(body), &beatmaps)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}
	integ, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Failed converting string to int in BeatmapBySetId")
	}
	var result *Beatmap
	for _, beatmap := range beatmaps {
		if beatmap.BeatmapID == integ {
			result = &beatmap
			break
		}
	}
	return *result
}
func GetBeatmapsById(id string) []Beatmap {
	url := "https://osu.direct/api/get_beatmaps?b=" + id
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: status code %d\n", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	var beatmaps []Beatmap
	err = json.Unmarshal([]byte(body), &beatmaps)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}
	return *&beatmaps
}

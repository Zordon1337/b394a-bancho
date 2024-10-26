package Utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type GeoLocation struct {
	Country string `json:"country"`
}

func GetCountryFromIP(ip string) (string, error) {
	if ip == "127.0.0.1" || ip == "::1" || strings.Split(ip, ".")[0] == "192" {
		return "behind you", nil // it still might pop up as US when using radmin/other lan vpn
	}
	apiURL := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var location GeoLocation
	if err := json.Unmarshal(body, &location); err != nil {
		return "", err
	}

	return location.Country, nil
}

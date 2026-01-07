package bilibili

import (
	"encoding/json"
	"fmt"
	"time"
)

type BuvidResponse struct {
	Code int `json:"code"`
	Data struct {
		Buvid string `json:"buvid"`
	} `json:"data"`
}

var ExpiredTime = 24 * time.Hour

func GetBuvid3() (string, error) {
	return GetCookieValue("https://bilibili.com", "buvid3"), nil
}

func GetGuestBuvid3() (string, error) {
	resp, err := AuthClient.Get("https://api.bilibili.com/x/web-frontend/getbuvid")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var buvidRes BuvidResponse
	if err := json.NewDecoder(resp.Body).Decode(&buvidRes); err != nil {
		return "", fmt.Errorf("decode buvid json error: %w", err)
	}

	if buvidRes.Code != 0 {
		return "", fmt.Errorf("buvid api code: %d", buvidRes.Code)
	}

	return buvidRes.Data.Buvid, nil
}

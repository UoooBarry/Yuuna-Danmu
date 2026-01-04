package wbi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type BuvidResponse struct {
	Code int `json:"code"`
	Data struct {
		Buvid string `json:"buvid"`
	} `json:"data"`
}

var ExpiredTime = 24 * time.Hour

func EnsureBuvid() error {
	if GetCookieValue("https://bilibili.com", "buvid3") == "" {
		return RefreshBuvid3()
	}
	return nil
}

func RefreshBuvid3() error {
	resp, err := BiliClient.Get("https://api.bilibili.com/x/web-frontend/getbuvid")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var buvidRes BuvidResponse
	if err := json.NewDecoder(resp.Body).Decode(&buvidRes); err != nil {
		return fmt.Errorf("decode buvid json error: %w", err)
	}

	if buvidRes.Code != 0 {
		return fmt.Errorf("buvid api code: %d", buvidRes.Code)
	}

	u, _ := url.Parse("https://bilibili.com")
	cookies := []*http.Cookie{
		{
			Name:     "buvid3",
			Value:    buvidRes.Data.Buvid,
			Domain:   ".bilibili.com",
			Path:     "/",
			HttpOnly: false,
		},
	}
	BiliClient.Jar.SetCookies(u, cookies)

	return nil
}

func GetBuvid3() (string, error) {
	return GetCookieValue("https://bilibili.com", "buvid3"), nil
}

package bilibili

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const publicKeyPEM = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDLgd2OAkcGVtoE3ThUREbio0Eg
Uc/prcajMKXvkCKFCWhJYJcLkcM2DKKcSeFpD/j6Boy538YXnR6VhcuUJOhH2x71
nzPjfdTcqMz7djHum0qSZA0AyCBDABUqCrfNgCiJ00Ra7GmRj+YCK1NJEuewlb40
JNrRuoEUXpabUzGB8QIDAQAB
-----END PUBLIC KEY-----
`

func CheckAndRefreshCookie(refreshToken string) (string, string, error) {
	if refreshToken == "" {
		return "", "", nil
	}

	needed, err := needRefresh()
	if err != nil {
		return "", "", fmt.Errorf("check need refresh failed: %w", err)
	}
	if !needed {
		return "", "", nil
	}

	ts := time.Now().UnixMilli()
	correspondPath, err := getCorrespondPath(ts)
	if err != nil {
		return "", "", fmt.Errorf("generate correspond path failed: %w", err)
	}

	refreshCSRF, err := getRefreshCSRF(correspondPath)
	if err != nil {
		return "", "", fmt.Errorf("get refresh csrf failed: %w", err)
	}

	csrf := GetCookieValue("https://www.bilibili.com", "bili_jct")
	if csrf == "" {
		return "", "", fmt.Errorf("csrf token (bili_jct) not found in cookies")
	}

	newRefreshToken, err := refreshCookie(csrf, refreshCSRF, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("refresh cookie failed: %w", err)
	}

	newCsrf := GetCookieValue("https://www.bilibili.com", "bili_jct")
	if err := confirmRefresh(newCsrf, refreshToken); err != nil {
		return "", "", fmt.Errorf("confirm refresh failed: %w", err)
	}

	newCookie := getFullCookieString()

	return newCookie, newRefreshToken, nil
}

func needRefresh() (bool, error) {
	resp, err := AuthClient.Get("https://passport.bilibili.com/x/passport-login/web/cookie/info")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result struct {
		Code int `json:"code"`
		Data struct {
			Refresh bool `json:"refresh"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	if result.Code != 0 {
		return false, fmt.Errorf("api error code: %d", result.Code)
	}

	return result.Data.Refresh, nil
}

func getCorrespondPath(ts int64) (string, error) {
	pubKeyBlock, _ := pem.Decode([]byte(publicKeyPEM))
	hash := sha256.New()
	random := rand.Reader
	msg := []byte(fmt.Sprintf("refresh_%d", ts))

	pubInterface, err := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if err != nil {
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)
	encryptedData, err := rsa.EncryptOAEP(hash, random, pub, msg, nil)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(encryptedData), nil
}

func getRefreshCSRF(correspondPath string) (string, error) {
	urlStr := fmt.Sprintf("https://www.bilibili.com/correspond/1/%s", correspondPath)
	resp, err := AuthClient.Get(urlStr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Regex to extract content of <div id="1-name">...</div>
	re := regexp.MustCompile(`<div id="1-name">([^<]+)</div>`)
	matches := re.FindSubmatch(body)
	if len(matches) < 2 {
		return "", fmt.Errorf("refresh_csrf not found in response")
	}
	return string(matches[1]), nil
}

func refreshCookie(csrf, refreshCSRF, refreshToken string) (string, error) {
	params := url.Values{}
	params.Set("csrf", csrf)
	params.Set("refresh_csrf", refreshCSRF)
	params.Set("source", "main_web")
	params.Set("refresh_token", refreshToken)

	resp, err := AuthClient.PostForm("https://passport.bilibili.com/x/passport-login/web/cookie/refresh", params)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			RefreshToken string `json:"refresh_token"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Code != 0 {
		return "", fmt.Errorf("code: %d, message: %s", result.Code, result.Message)
	}

	return result.Data.RefreshToken, nil
}

func confirmRefresh(csrf, refreshToken string) error {
	params := url.Values{}
	params.Set("csrf", csrf)
	params.Set("refresh_token", refreshToken)

	resp, err := AuthClient.PostForm("https://passport.bilibili.com/x/passport-login/web/confirm/refresh", params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Code != 0 {
		return fmt.Errorf("code: %d, message: %s", result.Code, result.Message)
	}
	return nil
}

func getFullCookieString() string {
	u, _ := url.Parse("https://bilibili.com")
	cookies := AuthClient.Jar.Cookies(u)
	var parts []string
	for _, c := range cookies {
		parts = append(parts, fmt.Sprintf("%s=%s", c.Name, c.Value))
	}
	return strings.Join(parts, "; ")
}

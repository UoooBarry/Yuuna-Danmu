package wbi

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/publicsuffix"
)

var (
	BiliClient *http.Client
	once       sync.Once
)

type Transport struct {
	Base http.RoundTripper
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Referer", "https://www.bilibili.com/")
	req.Header.Set("Origin", "https://www.bilibili.com")

	return t.Base.RoundTrip(req)
}

func init() {
	once.Do(func() {
		jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		BiliClient = &http.Client{
			Jar:     jar,
			Timeout: 10 * time.Second,
			Transport: &Transport{
				Base: http.DefaultTransport,
			},
		}
	})
}

func GetCookieValue(rawURL, name string) string {
	u, _ := url.Parse(rawURL)
	cookies := BiliClient.Jar.Cookies(u)
	for _, c := range cookies {
		if c.Name == name {
			return c.Value
		}
	}
	return ""
}

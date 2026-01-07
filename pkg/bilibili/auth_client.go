package bilibili

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/publicsuffix"
)

var (
	AuthClient *http.Client
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
		AuthClient = &http.Client{
			Timeout: 10 * time.Second,
			Transport: &Transport{
				Base: http.DefaultTransport,
			},
			Jar: jar,
		}
	})
}

func SetCookie(cookie string) {
	jar, ok := AuthClient.Jar.(*cookiejar.Jar)
	if !ok {
		return
	}

	initJar(jar, cookie)
}

func GetCookieValue(rawURL, name string) string {
	u, _ := url.Parse(rawURL)
	cookies := AuthClient.Jar.Cookies(u)
	for _, c := range cookies {
		if c.Name == name {
			return c.Value
		}
	}
	return ""
}

func manualSetCookie(jar *cookiejar.Jar, rawCookie string, domain string) {
	u, _ := url.Parse("https://" + domain)
	segments := strings.Split(rawCookie, ";")
	var cookies []*http.Cookie

	for _, seg := range segments {
		parts := strings.SplitN(strings.TrimSpace(seg), "=", 2)
		if len(parts) != 2 {
			continue
		}
		cookies = append(cookies, &http.Cookie{
			Name:   parts[0],
			Value:  parts[1],
			Domain: domain,
			Path:   "/",
		})
	}
	jar.SetCookies(u, cookies)
}

func initJar(jar *cookiejar.Jar, cookie string) {
	domains := []string{"bilibili.com", "live.bilibili.com", ".bilibili.com"}

	for _, domain := range domains {
		manualSetCookie(jar, cookie, domain)
	}
}

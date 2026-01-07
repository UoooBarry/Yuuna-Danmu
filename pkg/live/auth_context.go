package live

import (
	"errors"
	"net/http"
	"strconv"

	"uooobarry/yuuna-danmu/pkg/bilibili"
)

type AuthContext struct {
	UID    int
	Buvid3 string
	Token  string
	Cookie string
}

func (session *Session) prepareAuth(cookie string) error {
	buvid3, err := getBuvidFromCookieOrGuest(cookie)
	if err != nil {
		return err
	}
	session.authContext = &AuthContext{
		UID:    getUIDFromCookie(cookie),
		Buvid3: buvid3,
		Cookie: cookie,
	}
	bilibili.SetCookie(cookie)
	return nil
}

func getUIDFromCookie(cookie string) int {
	if cookie == "" {
		return 0
	}
	val := getCookieValue(cookie, "DedeUserID")
	uid, _ := strconv.Atoi(val)
	return uid
}

func getBuvidFromCookieOrGuest(cookie string) (string, error) {
	if cookie == "" {
		return bilibili.GetGuestBuvid3()
	}
	buvid3 := getCookieValue(cookie, "buvid3")
	if buvid3 == "" {
		return buvid3, errors.New("buvid3 cookie not found")
	}
	return buvid3, nil
}

func getCookieValue(cookie, name string) string {
	header := http.Header{}
	header.Add("Cookie", cookie)

	req := http.Request{Header: header}

	for _, cookie := range req.Cookies() {
		if cookie.Name == name {
			return cookie.Value
		}
	}
	return ""
}

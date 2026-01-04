package live

type Event struct {
	Type      string
	Data      any
	Raw       []byte
	Timestamp int64
}

// DanmuData 简单的弹幕结构
type DanmuData struct {
	UID      int64
	Nickname string
	Content  string
}

type BaseMsg struct {
	Cmd string `json:"cmd"`
}

type DanmuMsg struct {
	Content    string
	UserID     int64
	Nickname   string
	MedalName  string
	MedalLevel int
}

type GiftData struct {
	UID        int64  `json:"uid"`
	Nickname   string `json:"uname"`
	Action     string `json:"action"`
	GiftName   string `json:"giftName"`
	GiftID     int    `json:"giftId"`
	Num        int    `json:"num"`
	TotalCoin  int    `json:"total_coin"`
	CoinType   string `json:"coin_type"`
	ComboCount int    `json:"combo_count"`
	Face       string `json:"face"`
}

type PopularityMsg struct {
	Popularity int
}

var (
	DanmuEvent      = "DANMU_MSG"
	PopularityEvent = "POPULARITY"
	GiftEvent       = "SEND_GIFT"
)

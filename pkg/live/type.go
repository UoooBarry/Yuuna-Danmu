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
	UID            int64     `json:"uid"`
	Uname          string    `json:"uname"`
	Face           string    `json:"face"`
	GiftName       string    `json:"gift_name"`
	GiftNum        int       `json:"gift_num"`
	Price          int       `json:"price"`
	ComboTotalCoin int       `json:"combo_total_coin"`
	TotalCoin      int       `json:"total_coin"`
	CoinType       string    `json:"coin_type"`
	Action         string    `json:"action"`
	GiftInfo       GiftInfo  `json:"gift_info"`
	MedalInfo      MedalInfo `json:"medal_info"`
	ComboSend      ComboSend `json:"combo_send"`
}

type GiftInfo struct {
	ImgBasic string `json:"img_basic"`
	Gif      string `json:"gif"`
}

type MedalInfo struct {
	MedalName  string `json:"medal_name"`
	MedalLevel int    `json:"medal_level"`
}

type ComboSend struct {
	ComboID  string `json:"combo_id"`
	ComboNum int    `json:"combo_num"`
}

type PopularityMsg struct {
	Popularity int
}

var (
	DanmuEvent      = "DANMU_MSG"
	PopularityEvent = "POPULARITY"
	GiftEvent       = "SEND_GIFT"
	SysMsgEvent     = "SYS_MSG"
	ErrorEvent      = "SYS_ERROR"
)

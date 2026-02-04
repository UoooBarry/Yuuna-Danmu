package live

import "encoding/json"

type Event struct {
	Type      string
	Data      any
	Raw       []byte
	Timestamp int64
}

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
	Price          float64   `json:"price"`
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

type SuperChatMsgData struct {
	MedalInfo MedalInfo `json:"medal_info"`
	Message   string    `json:"message"`
	FontColor string    `json:"message_font_color"`
	Price     int       `json:"price"`
	UserInfo  UserInfo  `json:"user_info"`
	StartTime int64     `json:"start_time"`
	EndTime   int64     `json:"end_time"`
}

type UserInfo struct {
	Face  string `json:"face"`
	UName string `json:"uname"`
}

var (
	DanmuEvent           = "DANMU_MSG"
	PopularityEvent      = "POPULARITY"
	GiftEvent            = "SEND_GIFT"
	SysMsgEvent          = "SYS_MSG"
	ErrorEvent           = "SYS_ERROR"
	SuperChatEvent       = "SUPER_CHAT_MESSAGE"
	InteractionEvent     = "DM_INTERACTION"
	UserToastEvent       = "USER_TOAST_MSG"
	GiftStarProcessEvent = "GIFT_STAR_PROCESS"
)

type InteractMsg struct {
	ID     int64           `json:"id"`
	Status int             `json:"status"`
	Type   int             `json:"type"` // 101-106
	Data   json.RawMessage `json:"data"`
}

type InteractData102 struct {
	Combo []InteractCombo `json:"combo"`
}

type InteractCombo struct {
	ID      int64  `json:"id"`
	Status  int    `json:"status"`
	Content string `json:"content"`
	Cnt     int    `json:"cnt"`
	Guide   string `json:"guide"`
}

type InteractDataNotice struct {
	Cnt        int    `json:"cnt"`
	SuffixText string `json:"suffix_text"`
	GiftID     int    `json:"gift_id"`
}

type GuardLevel int

const (
	Governor GuardLevel = 1
	Admiral  GuardLevel = 2
	Captain  GuardLevel = 3
)

type ToastMsgData struct {
	GuardLevel GuardLevel `json:"guard_level"`
	Username   string     `json:"username"`
	Price      int        `json:"price"`
	UID        int64      `json:"uid"`
	Num        int        `json:"num"`
	Unit       string     `json:"unit"`
	RoleName   string     `json:"role_name"`
}

type GiftStarProcessData struct {
	Message string `json:"message"`
}

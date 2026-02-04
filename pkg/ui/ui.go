package ui

import (
	"uooobarry/yuuna-danmu/pkg/config"
	"uooobarry/yuuna-danmu/pkg/live"
)

type UI interface {
	AppendDanmu(medalName string, medalLevel int, nickname, content string)
	AppendGift(gift *live.GiftData)
	AppendError(err error)
	AppendSysMsg(msg string)
	Start() error
	Stop()
	SetOnConfigChange(onConfigChange OnConfigChange)
	AppendSuperChat(superchat *live.SuperChatMsgData)
	AppendInteraction(interaction *live.InteractMsg)
	UpdatePopularity(popularity int)
	AppendGiftStarProcess(data *live.GiftStarProcessData)
}

type ConfigPayload struct {
	RoomID       int                     `json:"room_id"`
	Cookie       string                  `json:"cookie"`
	Servers      []config.ServerSettings `json:"servers"`
	RefreshToken string                  `json:"refresh_token"`
	Transparent  bool                    `json:"transparent"`
}

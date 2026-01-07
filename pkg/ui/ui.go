package ui

import "uooobarry/yuuna-danmu/pkg/live"

type UI interface {
	AppendDanmu(medalName string, medalLevel int, nickname, content string)
	AppendGift(gift *live.GiftData)
	AppendError(err error)
	AppendSysMsg(msg string)
	Start() error
	SetOnConfigChange(onConfigChange OnConfigChange)
	AppendSuperChat(superchat *live.SuperChatMsgData)
}

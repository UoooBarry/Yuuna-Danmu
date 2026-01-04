package ui

type UI interface {
	AppendDanmu(medalName string, medalLevel int, nickname, content string)
	AppendGift(nickname, giftName string, count int)
	AppendError(err error)
	AppendSysMsg(msg string)
	Start() error
}

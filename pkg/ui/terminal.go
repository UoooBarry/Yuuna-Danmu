package ui

import (
	"fmt"
	"time"
)

type TerminalUI struct {
	onConfigChange OnConfigChange
}

func NewTerminalUI() *TerminalUI {
	return &TerminalUI{}
}

func (t *TerminalUI) Start() error {
	fmt.Println(">>> Yuuna Danmu 终端模式已启动")
	return nil
}

func (t *TerminalUI) AppendDanmu(medalName string, medalLevel int, nickname, content string) {
	fmt.Printf("[%s] [弹幕] [%s|%d]%s: %s\n", time.Now().Format(time.TimeOnly), medalName, medalLevel, nickname, content)
}

func (t *TerminalUI) AppendGift(nickname, giftName string, count int) {
	fmt.Printf("[%s] [礼物] [%s] 送出 %s x %d\n", time.Now().Format(time.TimeOnly), nickname, giftName, count)
}

func (t *TerminalUI) AppendError(err error) {
	fmt.Printf("[%s] [错误] %v\n", time.Now().Format(time.TimeOnly), err)
}

func (t *TerminalUI) AppendSysMsg(msg string) {
	fmt.Printf("[%s] [系统] %s\n", time.Now().Format(time.TimeOnly), msg)
}

func (t *TerminalUI) SaveConfig(roomID int, cookie string) string {
	if t.onConfigChange != nil {
		err := t.onConfigChange(roomID, cookie)
		if err != nil {
			return "更新失败: " + err.Error()
		}
	}
	return "保存成功"
}

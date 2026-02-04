package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"uooobarry/yuuna-danmu/pkg/config"
	"uooobarry/yuuna-danmu/pkg/live"
)

type TerminalUI struct {
	onConfigChange OnConfigChange
	quit           chan struct{}
}

func NewTerminalUI() *TerminalUI {
	return &TerminalUI{}
}

func (t *TerminalUI) Start() error {
	t.quit = make(chan struct{})
	fmt.Println(">>> Yuuna Danmu 终端模式已启动")

	<-t.quit
	return nil
}

func (t *TerminalUI) Stop() {
	if t.quit != nil {
		fmt.Println(">>> Yuuna Danmu 终端模式已关闭")
		close(t.quit)
	}
	return
}

func (t *TerminalUI) AppendDanmu(medalName string, medalLevel int, nickname, content string) {
	fmt.Printf("[%s] [弹幕] [%s|%d]%s: %s\n", time.Now().Format(time.TimeOnly), medalName, medalLevel, nickname, content)

	danmuEvent := map[string]interface{}{
		"event":      "DanmuEvent",
		"nickname":   nickname,
		"content":    content,
		"medalName":  medalName,
		"medalLevel": medalLevel,
		"timestamp":  time.Now().Format(time.RFC3339),
	}
	jsonBytes, err := json.Marshal(danmuEvent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling DanmuEvent: %v\n", err)
		return
	}
	fmt.Println(">>> JSON_EVENT:", string(jsonBytes))
}

func (t *TerminalUI) AppendGift(gift *live.GiftData) {
	fmt.Printf("[%s] [礼物] [%s] 送出 %s x %d\n", time.Now().Format(time.TimeOnly), gift.Uname, gift.GiftName, gift.GiftNum)

	giftEvent := map[string]interface{}{
		"event":     "GiftEvent",
		"giftData":  gift,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	jsonBytes, err := json.Marshal(giftEvent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling GiftEvent: %v\n", err)
		return
	}
	fmt.Println(">>> JSON_EVENT:", string(jsonBytes))
}

func (t *TerminalUI) AppendError(err error) {
	fmt.Printf("[%s] [错误] %v\n", time.Now().Format(time.TimeOnly), err)

	errorEvent := map[string]interface{}{
		"event":     "ErrorEvent",
		"error":     err.Error(),
		"timestamp": time.Now().Format(time.RFC3339),
	}
	jsonBytes, jsonErr := json.Marshal(errorEvent)
	if jsonErr != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling ErrorEvent: %v\n", jsonErr)
		return
	}
	fmt.Println(">>> JSON_EVENT:", string(jsonBytes))
}

func (t *TerminalUI) AppendSysMsg(msg string) {
	fmt.Printf("[%s] [系统] %s\n", time.Now().Format(time.TimeOnly), msg)

	sysMsgEvent := map[string]interface{}{
		"event":     "SysMsgEvent",
		"message":   msg,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	jsonBytes, err := json.Marshal(sysMsgEvent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling SysMsgEvent: %v\n", err)
		return
	}
	fmt.Println(">>> JSON_EVENT:", string(jsonBytes))
}

func (t *TerminalUI) SaveConfig(payload ConfigPayload) error {
	if t.onConfigChange != nil {
		err := t.onConfigChange(payload)
		if err != nil {
			t.AppendError(err)
			return err
		}
	}
	t.AppendSysMsg("保存成功")
	return nil
}

func (t *TerminalUI) SetOnConfigChange(onConfigChange OnConfigChange) {
	t.onConfigChange = onConfigChange
}

func (t *TerminalUI) AppendSuperChat(superchat *live.SuperChatMsgData) {
	fmt.Printf("[%s] [超级弹幕] %s: %s\n", time.Now().Format(time.TimeOnly), superchat.UserInfo.UName, superchat.Message)

	superChatEvent := map[string]interface{}{
		"event":     "SuperChatEvent",
		"superchat": superchat,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	jsonBytes, err := json.Marshal(superChatEvent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling SuperChatEvent: %v\n", err)
		return
	}
	fmt.Println(">>> JSON_EVENT:", string(jsonBytes))
}

func (t *TerminalUI) AppendInteraction(interaction *live.InteractMsg) {
	switch interaction.Type {
	case 102: // Combo Danmu
		var data live.InteractData102
		if err := json.Unmarshal(interaction.Data, &data); err == nil {
			for _, combo := range data.Combo {
				fmt.Printf("[%s] [弹幕合并] %s x%d: %s\n", time.Now().Format(time.TimeOnly), combo.Guide, combo.Cnt, combo.Content)
			}

			interactionEvent := map[string]interface{}{
				"event":       "InteractionEvent",
				"interaction": interaction,
				"parsedData":  data,
				"timestamp":   time.Now().Format(time.RFC3339),
			}
			jsonBytes, err := json.Marshal(interactionEvent)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error marshalling InteractionEvent (Combo Danmu): %v\n", err)
				return
			}
			fmt.Println(">>> JSON_EVENT:", string(jsonBytes))
		}
	case 103, 104, 105, 106: // Notice
		var jsonStr string
		if err := json.Unmarshal(interaction.Data, &jsonStr); err == nil {
			var data live.InteractDataNotice
			if err := json.Unmarshal([]byte(jsonStr), &data); err == nil {
				msg := fmt.Sprintf("%d %s", data.Cnt, data.SuffixText)
				if interaction.Type == 104 { // Gift
					msg += fmt.Sprintf(" (GiftID: %d)", data.GiftID)
				}
				fmt.Printf("[%s] [交互] %s\n", time.Now().Format(time.TimeOnly), msg)

				interactionEvent := map[string]interface{}{
					"event":       "InteractionEvent",
					"interaction": interaction,
					"parsedData":  data,
					"message":     msg,
					"timestamp":   time.Now().Format(time.RFC3339),
				}
				jsonBytes, err := json.Marshal(interactionEvent)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error marshalling InteractionEvent (Notice): %v\n", err)
					return
				}
				fmt.Println(">>> JSON_EVENT:", string(jsonBytes))
			}
		}
	case 101:
		fmt.Printf("[%s] [交互] 投票活动\n", time.Now().Format(time.TimeOnly))

		interactionEvent := map[string]interface{}{
			"event":       "InteractionEvent",
			"interaction": interaction,
			"message":     "投票活动",
			"timestamp":   time.Now().Format(time.RFC3339),
		}
		jsonBytes, err := json.Marshal(interactionEvent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshalling InteractionEvent (Vote): %v\n", err)
			return
		}
		fmt.Println(">>> JSON_EVENT:", string(jsonBytes))
	}
}

func (t *TerminalUI) UpdatePopularity(popularity int) {
	fmt.Printf("[%s] [人气] %d\n", time.Now().Format(time.TimeOnly), popularity)

	popularityEvent := map[string]interface{}{
		"event":      "PopularityEvent",
		"popularity": popularity,
		"timestamp":  time.Now().Format(time.RFC3339),
	}
	jsonBytes, err := json.Marshal(popularityEvent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling PopularityEvent: %v\n", err)
		return
	}
	fmt.Println(">>> JSON_EVENT:", string(jsonBytes))
}

func (t *TerminalUI) LoadConfig() *config.AppConfig {
	return config.Load()
}

func (t *TerminalUI) AppendGiftStarProcess(data *live.GiftStarProcessData) {
	fmt.Printf("[%s] [礼物星程] %s\n", time.Now().Format(time.TimeOnly), data.Message)

	giftStarProcessEvent := map[string]interface{}{
		"event":     "GiftStarProcessEvent",
		"data":      data,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	jsonBytes, err := json.Marshal(giftStarProcessEvent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling GiftStarProcessEvent: %v\n", err)
		return
	}
	fmt.Println(">>> JSON_EVENT:", string(jsonBytes))
}

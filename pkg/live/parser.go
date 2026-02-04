package live

import (
	"encoding/json"
	"log"
	"time"
)

func (c *WsClient) routeOperation(header *Header, body []byte) {
	switch header.Operation {
	case OpHeartbeatReply:
		// Heartbeat Reply no longer contains popularity data

	case OpSendMsgReply:
		c.dispatch(body)

	case OpAuthReply:
		c.eventCh <- Event{
			Type:      SysMsgEvent,
			Data:      "连接上了...",
			Timestamp: time.Now().UnixNano(),
		}

	default:
		log.Printf("[Yuuna-Danmu] Unknown operation: %d", header.Operation)
	}
}

func (c *WsClient) dispatch(body []byte) {
	var base BaseMsg
	if err := json.Unmarshal(body, &base); err != nil {
		return
	}

	switch base.Cmd {
	case DanmuEvent:
		data := parseDanmu(body)
		c.eventCh <- Event{
			Type:      DanmuEvent,
			Data:      data,
			Timestamp: time.Now().UnixNano(),
		}
	case GiftEvent:
		data := parseGift(body)
		log.Println(data)
		c.eventCh <- Event{
			Type:      GiftEvent,
			Data:      data,
			Timestamp: time.Now().UnixNano(),
		}
	case ComboSendEvent:
		data := parseComboSend(body)
		log.Println(data)
		c.eventCh <- Event{
			Type:      ComboSendEvent,
			Data:      data,
			Timestamp: time.Now().UnixNano(),
		}
	case SuperChatEvent:
		data := parseSuperChat(body)
		c.eventCh <- Event{
			Type:      SuperChatEvent,
			Data:      data,
			Timestamp: time.Now().UnixNano(),
		}
	case InteractionEvent:
		data := parseInteraction(body)
		c.eventCh <- Event{
			Type:      InteractionEvent,
			Data:      data,
			Timestamp: time.Now().UnixNano(),
		}
	case UserToastEvent:
		data := parseToast(body)
		c.eventCh <- Event{
			Type:      UserToastEvent,
			Data:      data,
			Timestamp: time.Now().UnixNano(),
		}
	case GiftStarProcessEvent:
		data := parseGiftStarProcess(body)
		c.eventCh <- Event{
			Type:      GiftStarProcessEvent,
			Data:      data,
			Timestamp: time.Now().UnixNano(),
		}
	case OnlineRankCountEvent:
		data := parseOnlineRankCount(body)
		c.eventCh <- Event{
			Type:      OnlineRankCountEvent,
			Data:      data,
			Timestamp: time.Now().UnixNano(),
		}
	default:
	}
}

func parseInteraction(body []byte) *InteractMsg {
	return parseJSONData[InteractMsg](body)
}

func parseToast(body []byte) *ToastMsgData {
	var raw struct {
		Data ToastMsgData `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil
	}
	return &raw.Data
}

func parseDanmu(body []byte) *DanmuMsg {
	var raw struct {
		Info []interface{} `json:"info"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil
	}

	if len(raw.Info) < 3 {
		return nil
	}

	content, _ := raw.Info[1].(string)

	userSlice, _ := raw.Info[2].([]interface{})
	nickname, _ := userSlice[1].(string)

	var medalName string
	var medalLevel int
	if medalSlice, ok := raw.Info[3].([]interface{}); ok && len(medalSlice) > 0 {
		medalLevel = int(medalSlice[0].(float64))
		medalName, _ = medalSlice[1].(string)
	}
	danMu := &DanmuMsg{
		Content:    content,
		UserID:     int64(userSlice[0].(float64)),
		Nickname:   nickname,
		MedalName:  medalName,
		MedalLevel: medalLevel,
	}

	return danMu
}

func parseJSONData[T any](body []byte) *T {
	var raw struct {
		Data T `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil
	}
	return &raw.Data
}

func parseGift(body []byte) *GiftData {
	return parseJSONData[GiftData](body)
}

func parseSuperChat(body []byte) *SuperChatMsgData {
	return parseJSONData[SuperChatMsgData](body)
}

func parseGiftStarProcess(body []byte) *GiftStarProcessData {
	return parseJSONData[GiftStarProcessData](body)
}

func parseComboSend(body []byte) *ComboSendData {
	return parseJSONData[ComboSendData](body)
}

func parseOnlineRankCount(body []byte) *OnlineRankCountData {
	return parseJSONData[OnlineRankCountData](body)
}

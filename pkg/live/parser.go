package live

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"time"
)

func (c *WsClient) routeOperation(header *Header, body []byte) {
	switch header.Operation {
	case OpHeartbeatReply:
		if len(body) >= 4 {
			popularity := binary.BigEndian.Uint32(body[:4])
			c.eventCh <- Event{
				Type:      PopularityEvent,
				Data:      PopularityMsg{Popularity: int(popularity)},
				Timestamp: time.Now().UnixNano(),
			}
		}

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
		data := c.parseDanmu(body)
		c.eventCh <- Event{
			Type:      DanmuEvent,
			Data:      data,
			Timestamp: time.Now().UnixNano(),
		}
	case GiftEvent:
		data := c.parseGift(body)
		c.eventCh <- Event{
			Type:      GiftEvent,
			Data:      data,
			Timestamp: time.Now().UnixNano(),
		}
	case SuperChatEvent:
		data := c.parseSuperChat(body)
		c.eventCh <- Event{
			Type:      SuperChatEvent,
			Data:      data,
			Timestamp: time.Now().UnixNano(),
		}
	case "INTERACT_WORD":
	default:
	}
}

func (c *WsClient) parseDanmu(body []byte) *DanmuMsg {
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

func (c *WsClient) parseGift(body []byte) *GiftData {
	var raw struct {
		Data GiftData `json:"data"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil
	}
	return &raw.Data
}

func (c *WsClient) parseSuperChat(body []byte) *SuperChatMsgData {
	var raw struct {
		Data SuperChatMsgData `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil
	}
	return &raw.Data
}

package app

import (
	"context"
	"math/rand"
	"time"

	"uooobarry/yuuna-danmu/pkg/live"
)

func startMockDriver(ctx context.Context, ch chan live.Event) {
	danmuTicker := time.NewTicker(2 * time.Second)
	defer danmuTicker.Stop()

	giftTicker := time.NewTicker(5 * time.Second)
	defer giftTicker.Stop()

	superChatTicker := time.NewTicker(30 * time.Second)
	defer superChatTicker.Stop()

	mockModel := &live.MedalInfo{
		MedalName:  "开发者",
		MedalLevel: 20,
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-danmuTicker.C:
			mockDanmu := &live.DanmuMsg{
				Nickname:   "测试用户_" + time.Now().Format("05"),
				Content:    "这是一条模拟弹幕 " + time.Now().Format(time.TimeOnly),
				MedalName:  "开发者",
				MedalLevel: 20,
			}

			ch <- live.Event{
				Type: live.DanmuEvent,
				Data: mockDanmu,
			}

		case <-giftTicker.C:
			mockGift := &live.GiftData{
				UID:            1,
				Uname:          "花花花花人",
				Face:           "http://i1.hdslb.com/bfs/face/8b9a772ff6414bf9a83b57f6fcc22b00821aeb03.jpg",
				GiftName:       "粉丝团灯牌",
				GiftNum:        1,
				Price:          100,
				TotalCoin:      100,
				ComboTotalCoin: 100,
				CoinType:       "gold",
				Action:         "投喂",
				GiftInfo: live.GiftInfo{
					ImgBasic: "https://i0.hdslb.com/bfs/live/816f8b7aa2132888fce928cdfb17b9cf21cc0823.gif",
					Gif:      "https://s1.hdslb.com/bfs/live/e051dfd4557678f8edcac4993ed00a0935cbd9cc.png",
				},
				MedalInfo: *mockModel,
				ComboSend: live.ComboSend{
					ComboID:  "gift:combo_id:33313931353735d41d8cd98f00b204e9800998ecf8427e:1593304774:34001:1767675372.9443",
					ComboNum: 1,
				},
			}
			ch <- live.Event{
				Type: live.GiftEvent,
				Data: mockGift,
			}
		case <-superChatTicker.C:
			mockSuperChat := &live.SuperChatMsgData{
				MedalInfo: *mockModel,
				Message:   "这是一条模拟超级留言 " + time.Now().Format(time.TimeOnly),
				FontColor: "#FF0000",
				Price:     rand.Intn(100),
				UserInfo: live.UserInfo{
					UName: "花花花花人",
					Face:  "http://i1.hdslb.com/bfs/face/8b9a772ff6414bf9a83b57f6fcc22b00821aeb03.jpg",
				},
				StartTime: time.Now().Unix(),
				EndTime:   time.Now().Add(time.Minute).Unix(),
			}
			ch <- live.Event{
				Type: live.SuperChatEvent,
				Data: mockSuperChat,
			}
		}
	}
}

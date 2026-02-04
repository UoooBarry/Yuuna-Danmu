package app

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"uooobarry/yuuna-danmu/pkg/live"
)

func startMockDriver(ctx context.Context, ch chan live.Event) {
	danmuTicker := time.NewTicker(2 * time.Second)
	defer danmuTicker.Stop()

	giftTicker := time.NewTicker(5 * time.Second)
	defer giftTicker.Stop()

	comboSendTicker := time.NewTicker(8 * time.Second)
	defer comboSendTicker.Stop()

	superChatTicker := time.NewTicker(30 * time.Second)
	defer superChatTicker.Stop()

	interactTicker := time.NewTicker(10 * time.Second)
	defer interactTicker.Stop()

	giftStarProcessTicker := time.NewTicker(15 * time.Second)
	defer giftStarProcessTicker.Stop()

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
		case <-comboSendTicker.C:
			mockComboSend := &live.ComboSendData{
				Action:         "投喂",
				BatchComboID:   "batch:gift:combo_id:510149209:36047134:31036:1673622464.8445",
				BatchComboNum:  3,
				ComboID:        "gift:combo_id:510149209:36047134:31036:1673622464.8434",
				ComboNum:       3,
				ComboTotalCoin: 300,
				Dmscore:        112,
				GiftID:         31036,
				GiftName:       "小花花",
				GiftNum:        0,
				IsJoinReceiver: false,
				IsNaming:       false,
				IsShow:         1,
				MedalInfo:      *mockModel,
				NameColor:      "",
				RUname:         "测试主播",
				ReceiveUserInfo: struct {
					UID   int64  `json:"uid"`
					Uname string `json:"uname"`
				}{
					UID:   36047134,
					Uname: "测试主播",
				},
				Ruid:       36047134,
				SendMaster: nil,
				TotalNum:   3,
				UID:        510149209,
				Uname:      "花花花花人",
			}
			ch <- live.Event{
				Type: live.ComboSendEvent,
				Data: mockComboSend,
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
		case <-interactTicker.C:
			if rand.Intn(2) == 0 {
				data := live.InteractData102{
					Combo: []live.InteractCombo{
						{
							Content: "666666",
							Cnt:     rand.Intn(10) + 1,
							Guide:   "他们都在说:",
						},
					},
				}
				dataBytes, _ := json.Marshal(data)
				ch <- live.Event{
					Type: live.InteractionEvent,
					Data: &live.InteractMsg{
						Type: 102,
						Data: dataBytes,
					},
				}
			} else {
				types := []int{103, 106}
				selectedType := types[rand.Intn(len(types))]
				suffix := "人关注了主播"
				if selectedType == 106 {
					suffix = "人正在点赞"
				}

				noticeData := live.InteractDataNotice{
					Cnt:        rand.Intn(100) + 1,
					SuffixText: suffix,
				}
				noticeBytes, _ := json.Marshal(noticeData)
				noticeStr, _ := json.Marshal(string(noticeBytes))

				ch <- live.Event{
					Type: live.InteractionEvent,
					Data: &live.InteractMsg{
						Type: selectedType,
						Data: noticeStr,
					},
				}
			}
		case <-giftStarProcessTicker.C:
			messages := []string{
				"用户 花花花花人 的礼物星程已更新",
				"礼物星程进度更新",
				"用户 测试用户_" + time.Now().Format("05") + " 完成了礼物星程任务",
				"星程等级提升",
			}
			selectedMessage := messages[rand.Intn(len(messages))]
			mockGiftStarProcess := &live.GiftStarProcessData{
				Message: selectedMessage,
			}
			ch <- live.Event{
				Type: live.GiftStarProcessEvent,
				Data: mockGiftStarProcess,
			}
		}
	}
}

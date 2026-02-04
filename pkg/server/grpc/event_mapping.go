package grpc

import (
	"log"

	"uooobarry/yuuna-danmu/api/grpc/pb"
	"uooobarry/yuuna-danmu/pkg/live"
)

func (s *GRPCServer) mapToProto(event live.Event) *pb.LiveEvent {
	switch event.Type {

	case live.DanmuEvent:
		e, ok := event.Data.(*live.DanmuMsg)
		if !ok {
			return nil
		}

		return &pb.LiveEvent{
			Payload: &pb.LiveEvent_Danmu{
				Danmu: &pb.DanmuMsg{
					Content:    e.Content,
					UserId:     e.UserID,
					Nickname:   e.Nickname,
					MedalName:  e.MedalName,
					MedalLevel: int32(e.MedalLevel),
				},
			},
		}

	case live.PopularityEvent:
		e, ok := event.Data.(*live.PopularityMsg)
		if !ok {
			return nil
		}

		return &pb.LiveEvent{
			Payload: &pb.LiveEvent_Popularity{
				Popularity: &pb.PopularityMsg{
					Popularity: int32(e.Popularity),
				},
			},
		}

	case live.GiftEvent:
		e, ok := event.Data.(*live.GiftData)
		if !ok {
			return nil
		}

		return &pb.LiveEvent{
			Payload: &pb.LiveEvent_Gift{
				Gift: &pb.GiftData{
					Uid:            e.UID,
					Uname:          e.Uname,
					Face:           e.Face,
					GiftName:       e.GiftName,
					GiftNum:        int32(e.GiftNum),
					Price:          e.Price,
					ComboTotalCoin: int32(e.ComboTotalCoin),
					TotalCoin:      int32(e.TotalCoin),
					CoinType:       e.CoinType,
					Action:         e.Action,
					GiftInfo: &pb.GiftInfo{
						ImgBasic: e.GiftInfo.ImgBasic,
						Gif:      e.GiftInfo.Gif,
					},
					MedalInfo: &pb.MedalInfo{
						MedalName:  e.MedalInfo.MedalName,
						MedalLevel: int32(e.MedalInfo.MedalLevel),
					},
					ComboSend: &pb.ComboSend{
						ComboId:  e.ComboSend.ComboID,
						ComboNum: int32(e.ComboSend.ComboNum),
					},
				},
			},
		}

	case live.SuperChatEvent:
		e, ok := event.Data.(*live.SuperChatMsgData)
		if !ok {
			return nil
		}

		return &pb.LiveEvent{
			Payload: &pb.LiveEvent_SuperChat{
				SuperChat: &pb.SuperChatMsg{
					Message:   e.Message,
					FontColor: e.FontColor,
					Price:     int32(e.Price),
					StartTime: e.StartTime,
					EndTime:   e.EndTime,
					MedalInfo: &pb.MedalInfo{
						MedalName:  e.MedalInfo.MedalName,
						MedalLevel: int32(e.MedalInfo.MedalLevel),
					},
					UserInfo: &pb.UserInfo{
						Face:  e.UserInfo.Face,
						Uname: e.UserInfo.UName,
					},
				},
			},
		}

	case live.InteractionEvent:
		e, ok := event.Data.(*live.InteractMsg)
		if !ok {
			return nil
		}

		return &pb.LiveEvent{
			Payload: &pb.LiveEvent_Interaction{
				Interaction: &pb.InteractMsg{
					Id:     e.ID,
					Status: int32(e.Status),
					Type:   int32(e.Type),
					Data:   string(e.Data),
				},
			},
		}

	case live.SysMsgEvent:
		e, ok := event.Data.(string)
		if !ok {
			return nil
		}

		return &pb.LiveEvent{
			Payload: &pb.LiveEvent_SysMsg{
				SysMsg: e,
			},
		}

	case live.ErrorEvent:
		e, ok := event.Data.(string)
		if !ok {
			return nil
		}

		return &pb.LiveEvent{
			Payload: &pb.LiveEvent_Error{
				Error: e,
			},
		}

	case live.OnlineRankCountEvent:
		e, ok := event.Data.(*live.OnlineRankCountData)
		if !ok {
			return nil
		}

		return &pb.LiveEvent{
			Payload: &pb.LiveEvent_OnlineRankCount{
				OnlineRankCount: &pb.OnlineRankCountMsg{
					Count:           int32(e.Count),
					CountText:       e.CountText,
					OnlineCount:     int32(e.OnlineCount),
					OnlineCountText: e.OnlineCountText,
				},
			},
		}

	case live.UserToastEvent:
		e, ok := event.Data.(*live.ToastMsgData)
		if !ok {
			return nil
		}

		return &pb.LiveEvent{
			Payload: &pb.LiveEvent_Toast{
				Toast: &pb.ToastMsg{
					GuardLevel: int32(e.GuardLevel),
					Username:   e.Username,
					Price:      int32(e.Price),
					Uid:        e.UID,
					Num:        int32(e.Num),
					Unit:       e.Unit,
					RoleName:   e.RoleName,
				},
			},
		}

	case live.GiftStarProcessEvent:
		e, ok := event.Data.(*live.GiftStarProcessData)
		if !ok {
			return nil
		}

		return &pb.LiveEvent{
			Payload: &pb.LiveEvent_GiftStarProcess{
				GiftStarProcess: &pb.GiftStarProcessMsg{
					Message: e.Message,
				},
			},
		}

	case live.ComboSendEvent:
		e, ok := event.Data.(*live.ComboSendData)
		if !ok {
			return nil
		}

		return &pb.LiveEvent{
			Payload: &pb.LiveEvent_ComboSend{
				ComboSend: &pb.ComboSendData{
					Action:         e.Action,
					BatchComboId:   e.BatchComboID,
					BatchComboNum:  int32(e.BatchComboNum),
					ComboId:        e.ComboID,
					ComboNum:       int32(e.ComboNum),
					ComboTotalCoin: int32(e.ComboTotalCoin),
					Dmscore:        int32(e.Dmscore),
					GiftId:         int32(e.GiftID),
					GiftName:       e.GiftName,
					GiftNum:        int32(e.GiftNum),
					IsJoinReceiver: e.IsJoinReceiver,
					IsNaming:       e.IsNaming,
					IsShow:         int32(e.IsShow),
					MedalInfo: &pb.MedalInfo{
						MedalName:  e.MedalInfo.MedalName,
						MedalLevel: int32(e.MedalInfo.MedalLevel),
					},
					NameColor: e.NameColor,
					RUname:    e.RUname,
					ReceiveUserInfo: &pb.UserInfo{
						Face:  "",
						Uname: e.ReceiveUserInfo.Uname,
					},
					Ruid:     e.Ruid,
					TotalNum: int32(e.TotalNum),
					Uid:      e.UID,
					Uname:    e.Uname,
				},
			},
		}

	default:
		log.Printf("Unknown event type: %s", event.Type)
		return nil
	}
}

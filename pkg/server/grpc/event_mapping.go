package grpc

import (
	"log"

	"uooobarry/yuuna-danmu/pkg/live"
	"uooobarry/yuuna-danmu/pkg/server/grpc/pb"
)

func (s *GRPCServer) mapToProto(event any) *pb.LiveEvent {
	switch e := event.(type) {

	case *live.DanmuMsg:
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

	case int:
		return &pb.LiveEvent{
			Payload: &pb.LiveEvent_Popularity{
				Popularity: &pb.PopularityMsg{
					Popularity: int32(e),
				},
			},
		}

	case *live.GiftData:
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

	case *live.SuperChatMsgData:
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

	case string:
		return &pb.LiveEvent{
			Payload: &pb.LiveEvent_SysMsg{
				SysMsg: e,
			},
		}

	default:
		log.Printf("Unknown event type: %T", event)
		return nil
	}
}


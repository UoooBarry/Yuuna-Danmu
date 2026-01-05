package live

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"uooobarry/yuuna-danmu/pkg/errors"
	"uooobarry/yuuna-danmu/pkg/wbi"
)

type RoomInitResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		RoomID      int  `json:"room_id"`  // 真实房间ID
		ShortID     int  `json:"short_id"` // 短号
		Uid         int  `json:"uid"`      // 主播UID
		NeedP2p     int  `json:"need_p2p"`
		IsHidden    bool `json:"is_hidden"`
		IsLocked    bool `json:"is_locked"`
		IsEncrypted bool `json:"is_encrypted"`
	} `json:"data"`
}

type DanmuInfoResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Token    string `json:"token"`
		HostList []struct {
			Host    string `json:"host"`
			Port    int    `json:"port"`
			WssPort int    `json:"wss_port"`
			WsPort  int    `json:"ws_port"`
		} `json:"host_list"`
	} `json:"data"`
}

func GetRealRoomID(roomID int) (int, error) {
	api := fmt.Sprintf("https://api.live.bilibili.com/room/v1/Room/room_init?id=%d", roomID)

	resp, err := wbi.BiliClient.Get(api)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result RoomInitResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	if result.Code != 0 {
		return 0, &errors.ApiError{Code: result.Code, Message: result.Message, RawErr: err}
	}

	return result.Data.RoomID, nil
}

func GetDanmuConfig(roomID int) (*DanmuInfoResp, error) {
	api := fmt.Sprintf("https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=%d&type=0&wts=%d", roomID, getWts())
	apiUrl, err := url.Parse(api)
	if err != nil {
		return nil, err
	}
	if err = wbi.Sign(apiUrl); err != nil {
		return nil, err
	}

	resp, err := wbi.BiliClient.Get(apiUrl.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result DanmuInfoResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, &errors.ApiError{Code: result.Code, Message: result.Message, RawErr: err}
	}
	return &result, nil
}

func getWts() int64 {
	return time.Now().Unix()
}

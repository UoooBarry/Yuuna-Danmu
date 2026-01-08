package ui

import (
	"context"
	"fmt"
	"log"
	"time"

	"uooobarry/yuuna-danmu/pkg/config"
	"uooobarry/yuuna-danmu/pkg/live"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type WailsUI struct {
	ctx            context.Context
	onConfigChange OnConfigChange
	emitter        func(ctx context.Context, name string, data ...any)
	assetOpts      *assetserver.Options
}

type OnConfigChange func(roomID int, cookie string) error

func NewWailsUI(assetOts *assetserver.Options, opts ...func(*WailsUI)) *WailsUI {
	ui := &WailsUI{
		assetOpts: assetOts,
		emitter:   runtime.EventsEmit,
	}
	return ui
}

func (w *WailsUI) Stop() error {
	if w.ctx != nil {
		runtime.Quit(w.ctx)
	}
	return nil
}

func (w *WailsUI) SetContext(ctx context.Context) {
	w.ctx = ctx
}

func (w *WailsUI) OnStartup(ctx context.Context) {
	w.ctx = ctx
}

func (w *WailsUI) SetOnConfigChange(onConfigChange OnConfigChange) {
	w.onConfigChange = onConfigChange
}

func (w *WailsUI) Start() error {
	return wails.Run(&options.App{
		Title:            "Yuuna Danmu",
		Width:            320,
		Height:           520,
		Frameless:        true,
		AlwaysOnTop:      true,
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		AssetServer:      w.assetOpts,
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			BackdropType:         windows.Mica,
		},
		OnStartup: w.OnStartup,
		Bind: []interface{}{
			w,
		},
	})
}

func (w *WailsUI) AppendDanmu(medalName string, medalLevel int, nickname, content string) {
	if w.ctx == nil {
		return
	}

	w.emitter(w.ctx, live.DanmuEvent, map[string]string{
		"nickname":   nickname,
		"content":    content,
		"medalName":  medalName,
		"medalLevel": fmt.Sprintf("%d", medalLevel),
		"timestamp":  time.Now().Format(time.TimeOnly),
	})
}

func (w *WailsUI) AppendGift(gift *live.GiftData) {
	if w.ctx == nil {
		return
	}

	w.emitter(w.ctx, live.GiftEvent, gift)
}

func (w *WailsUI) AppendSuperChat(superchat *live.SuperChatMsgData) {
	if w.ctx == nil {
		return
	}

	log.Println(superchat)
	w.emitter(w.ctx, live.SuperChatEvent, superchat)
}

func (w *WailsUI) AppendError(err error) {
	if w.ctx == nil {
		return
	}

	w.emitter(w.ctx, live.ErrorEvent, err.Error())
}

func (w *WailsUI) AppendSysMsg(msg string) {
	if w.ctx == nil {
		return
	}

	w.emitter(w.ctx, live.SysMsgEvent, msg)
}

type configPayload struct {
	RoomID int    `json:"room_id"`
	Cookie string `json:"cookie"`
}

func (w *WailsUI) SaveConfig(payload configPayload) {
	if w.onConfigChange != nil {
		err := w.onConfigChange(payload.RoomID, payload.Cookie)
		if err != nil {
			w.AppendError(err)
		}
	}
	w.AppendSysMsg("保存成功")
}

func (w *WailsUI) LoadConfig() *config.AppConfig {
	return config.Load()
}

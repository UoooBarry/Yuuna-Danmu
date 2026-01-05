package ui

import (
	"context"
	"embed"
	"fmt"
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
	AssetFS        embed.FS
	onConfigChange OnConfigChange
}

type OnConfigChange func(cfg config.AppConfig) error

func NewWailsUI(assetFS embed.FS) *WailsUI {
	return &WailsUI{
		AssetFS: assetFS,
	}
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
		AssetServer: &assetserver.Options{
			Assets: w.AssetFS,
		},
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

	runtime.EventsEmit(w.ctx, live.DanmuEvent, map[string]string{
		"nickname":   nickname,
		"content":    content,
		"medalName":  medalName,
		"medalLevel": fmt.Sprintf("%d", medalLevel),
		"timestamp":  time.Now().Format(time.TimeOnly),
	})
}

func (t *WailsUI) AppendGift(nickname, giftName string, count int) {
	fmt.Printf("[%s] [礼物] [%s] 送出 %s x %d\n", time.Now().Format(time.TimeOnly), nickname, giftName, count)
}

func (w *WailsUI) AppendError(err error) {
	if w.ctx == nil {
		return
	}

	runtime.EventsEmit(w.ctx, live.ErrorEvent, err.Error())
}

func (w *WailsUI) AppendSysMsg(msg string) {
	if w.ctx == nil {
		return
	}

	runtime.EventsEmit(w.ctx, live.SysMsgEvent, msg)
}

func (w *WailsUI) SaveConfig(cfg config.AppConfig) string {
	if w.onConfigChange != nil {
		err := w.onConfigChange(cfg)
		if err != nil {
			return "更新失败: " + err.Error()
		}
	}
	return "保存成功"
}

func (w *WailsUI) LoadConfig() *config.AppConfig {
	return config.Load()
}

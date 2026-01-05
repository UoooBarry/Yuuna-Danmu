package app

import (
	"log"
	"os"

	"uooobarry/yuuna-danmu/pkg/config"
	"uooobarry/yuuna-danmu/pkg/errors"
	"uooobarry/yuuna-danmu/pkg/live"
	"uooobarry/yuuna-danmu/pkg/ui"
)

type App struct {
	RoomID  int
	Cookie  string
	session *live.Session
	ui      ui.UI
}

type Option func(*App)

func WithUI(u ui.UI) Option {
	return func(app *App) {
		app.ui = u
		app.ui.SetOnConfigChange(func(cfg config.AppConfig) error {
			return app.RestartSession(cfg.RoomID, cfg.Cookie)
		})
	}
}

func WithFileLog(filename string) Option {
	return func(a *App) {
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			log.Printf("Unable to open log file: %v", err)
			return
		}
		log.SetOutput(f)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}
}

func NewApp(opts ...Option) *App {
	cfg := config.Load()
	app := &App{
		RoomID:  cfg.RoomID,
		session: live.NewSession(cfg.RoomID, cfg.Cookie),
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

func (app *App) Run() error {
	go func() {
		for event := range app.session.EventCh {
			switch event.Type {
			case live.DanmuEvent:
				if d, ok := event.Data.(*live.DanmuMsg); ok {
					app.ui.AppendDanmu(d.MedalName, d.MedalLevel, d.Nickname, d.Content)
				}
			case live.GiftEvent:
				if g, ok := event.Data.(*live.GiftData); ok {
					app.ui.AppendGift(g.Nickname, g.GiftName, g.Num)
				}
			}
		}
	}()

	go func() {
		if err := app.session.Start(); err != nil {
			if apiErr, ok := err.(*errors.ApiError); ok {
				app.ui.AppendError(apiErr)
			}
		} else {
			app.ui.AppendSysMsg("服务已启动，正在监听数据...")
		}
	}()

	return app.ui.Start()
}

func (app *App) RestartSession(newRoomID int, newCookie string) error {
	config := &config.AppConfig{
		RoomID: newRoomID,
		Cookie: newCookie,
	}
	config.Save()
	if app.session != nil {
		app.session.Stop()
	}
	app.session.UpdateConfig(newRoomID, newCookie)
	log.Printf("[Yuuna-Danmu] RoomID: %d, Cookie: %s", newRoomID, newCookie)
	go func() {
		if err := app.session.Start(); err != nil {
			if apiErr, ok := err.(*errors.ApiError); ok {
				app.ui.AppendError(apiErr)
			}
		}
	}()

	return nil
}

func (app *App) Stop() {
	app.ui.AppendSysMsg("[Yuuna-Danmu]Stopping session...")
	app.session.Stop()
}

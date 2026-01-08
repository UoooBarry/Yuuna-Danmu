package app

import (
	"context"
	"errors"
	"log"
	"os"

	"uooobarry/yuuna-danmu/pkg/config"
	"uooobarry/yuuna-danmu/pkg/live"
	"uooobarry/yuuna-danmu/pkg/server"
	"uooobarry/yuuna-danmu/pkg/ui"
)

type App struct {
	AppConfig  config.AppConfig
	session    *live.Session
	ui         ui.UI
	servers    []*ServerConfig
	cancelMock context.CancelFunc
}

type ServerConfig struct {
	Target server.Server
	Port   int
}

type Option func(*App)

func WithUI(u ui.UI) Option {
	return func(app *App) {
		app.ui = u
		app.ui.SetOnConfigChange(func(roomID int, cookie string) error {
			return app.RestartSession(roomID, cookie)
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

func WithServer(s server.Server, port int) Option {
	return func(a *App) {
		a.servers = append(a.servers, &ServerConfig{Target: s, Port: port})
	}
}

func NewApp(opts ...Option) (*App, error) {
	cfg := config.Load()

	session, err := live.NewSession(cfg.RoomID, cfg.Cookie)
	if err != nil {
		return nil, err
	}

	app := &App{
		AppConfig: *cfg,
		session:   session,
	}

	for _, opt := range opts {
		opt(app)
	}

	if app.ui == nil {
		WithUI(ui.NewTerminalUI())(app)
	}
	return app, nil
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
					app.ui.AppendGift(g)
				}
			case live.ErrorEvent:
				if e, ok := event.Data.(string); ok {
					app.ui.AppendError(errors.New(e))
				}
			case live.SysMsgEvent:
				if msg, ok := event.Data.(string); ok {
					app.ui.AppendSysMsg(msg)
				}
			case live.SuperChatEvent:
				if sc, ok := event.Data.(*live.SuperChatMsgData); ok {
					app.ui.AppendSuperChat(sc)
				}
			}

			for _, s := range app.servers {
				s.Target.Dispatch(event.Data)
			}
		}
	}()

	go app.runSession()

	for _, s := range app.servers {
		s.Target.Start(s.Port)
	}

	return app.ui.Start()
}

func (app *App) RestartSession(newRoomID int, newCookie string) error {
	config := &config.AppConfig{
		RoomID: newRoomID,
		Cookie: newCookie,
		Debug:  app.AppConfig.Debug,
	}
	config.Save()
	if app.session != nil {
		app.session.Stop()
	}

	app.AppConfig = *config
	app.session.UpdateConfig(newRoomID, newCookie)
	app.ui.AppendSysMsg("[Yuuna-Danmu]Session restarting...")
	if !app.AppConfig.Debug {
		go app.runSession()
	}

	return nil
}

func (app *App) runSession() {
	if app.AppConfig.Debug {
		ctx, cancel := context.WithCancel(context.Background())
		app.cancelMock = cancel
		app.ui.AppendSysMsg("[Debug Mode] Starting mock driver...")
		app.startMockDriver(ctx)
	} else {
		if err := app.session.Start(); err != nil {
			app.ui.AppendError(err)
		}
	}
}

func (app *App) Stop() {
	app.ui.AppendSysMsg("[Yuuna-Danmu]Stopping session...")
	if app.cancelMock != nil {
		app.cancelMock()
	}
	for _, s := range app.servers {
		s.Target.Stop()
	}
	if app.session != nil {
		app.session.Stop()
	}
	app.ui.Stop()
	log.Printf("[Yuuna-Danmu]App stopped")
}

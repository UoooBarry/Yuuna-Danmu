package app

import (
	"context"
	"errors"
	"log"
	"os"

	"uooobarry/yuuna-danmu/pkg/config"
	"uooobarry/yuuna-danmu/pkg/live"
	"uooobarry/yuuna-danmu/pkg/server"
	"uooobarry/yuuna-danmu/pkg/server/grpc"
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
		app.ui.SetOnConfigChange(func(payload ui.ConfigPayload) error {
			return app.reload(payload)
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
	go app.consumeEvents()

	app.initServers()

	go app.runSession()

	return app.ui.Start()
}

func (app *App) reload(payload ui.ConfigPayload) error {
	if app.session != nil {
		app.session.Stop()
	}

	app.AppConfig.RoomID = payload.RoomID
	app.AppConfig.Cookie = payload.Cookie
	app.AppConfig.Servers = payload.Servers
	app.AppConfig.Save()

	app.session.UpdateConfig(payload.RoomID, payload.Cookie)

	app.ui.AppendSysMsg("[Yuuna-Danmu]Session restarting...")

	app.initServers()
	go app.runSession()

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
	app.servers = []*ServerConfig{}
	log.Printf("[Yuuna-Danmu]App stopped")
}

func (app *App) consumeEvents() {
	for event := range app.session.EventCh {
		app.handleEvent(event)

		for _, s := range app.servers {
			s.Target.Dispatch(event.Data)
		}
	}
}

func (app *App) handleEvent(event live.Event) {
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
}

func (app *App) initServers() {
	for _, s := range app.servers {
		s.Target.Stop()
	}
	app.servers = nil
	for _, s := range app.AppConfig.Servers {
		if s.Type == config.GRPC && s.Enabled {
			server := grpc.New()
			app.servers = append(app.servers, &ServerConfig{Target: server, Port: s.Port})
			go server.Start(s.Port)
		}
	}
}

package app

import (
	"log"
	"os"

	"uooobarry/yuuna-danmu/pkg/live"
	"uooobarry/yuuna-danmu/pkg/ui"
)

type App struct {
	roomID  int
	session *live.Session
	ui      ui.UI
}

type Option func(*App)

func WithUI(u ui.UI) Option {
	return func(a *App) {
		a.ui = u
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

func NewApp(roomID int, opts ...Option) *App {
	app := &App{
		roomID:  roomID,
		session: live.NewSession(roomID),
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

func (a *App) Run() error {
	if err := a.ui.Start(); err != nil {
		return err
	}

	if err := a.session.Start(); err != nil {
		a.ui.AppendError(err)
		return err
	}

	a.ui.AppendSysMsg("服务已启动，正在监听数据...")

	for event := range a.session.EventCh {
		switch event.Type {
		case live.DanmuEvent:
			a.ui.AppendSysMsg("test")
			if d, ok := event.Data.(*live.DanmuMsg); ok {
				a.ui.AppendDanmu(d.MedalName, d.MedalLevel, d.Nickname, d.Content)
			}
		case live.GiftEvent:
			if g, ok := event.Data.(*live.GiftData); ok {
				a.ui.AppendGift(g.Nickname, g.GiftName, g.Num)
			}
		}
	}

	return nil
}

func (a *App) Stop() {
	a.ui.AppendSysMsg("Stopping session...")
	a.session.Stop()
}

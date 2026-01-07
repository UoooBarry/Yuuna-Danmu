package app

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	"uooobarry/yuuna-danmu/pkg/config"
	"uooobarry/yuuna-danmu/pkg/live"
	"uooobarry/yuuna-danmu/pkg/ui"
)

type App struct {
	AppConfig config.AppConfig
	session   *live.Session
	ui        ui.UI
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
				if g, ok := event.Data.(string); ok {
					app.ui.AppendError(errors.New(g))
				}
			case live.SysMsgEvent:
				if g, ok := event.Data.(string); ok {
					app.ui.AppendSysMsg(g)
				}
			case live.SuperChatEvent:
				if g, ok := event.Data.(*live.SuperChatMsgData); ok {
					app.ui.AppendSuperChat(g)
				}
			}
		}
	}()

	go app.runSession()

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
		app.ui.AppendSysMsg("[Debug Mode] Starting mock driver...")
		app.startMockDriver()
	} else {
		if err := app.session.Start(); err != nil {
			app.ui.AppendError(err)
		}
	}
}

func (app *App) Stop() {
	app.ui.AppendSysMsg("[Yuuna-Danmu]Stopping session...")
	app.session.Stop()
}

func (app *App) startMockDriver() {
	danmuTicker := time.NewTicker(2 * time.Second)
	defer danmuTicker.Stop()

	giftTicker := time.NewTicker(5 * time.Second)
	defer giftTicker.Stop()

	superChatTicker := time.NewTicker(30 * time.Second)
	defer superChatTicker.Stop()

	mockModel := &live.MedalInfo{
		MedalName:  "开发者",
		MedalLevel: 20,
	}
	for {
		select {
		case <-danmuTicker.C:
			mockDanmu := &live.DanmuMsg{
				Nickname:   "测试用户_" + time.Now().Format("05"),
				Content:    "这是一条模拟弹幕 " + time.Now().Format(time.TimeOnly),
				MedalName:  "开发者",
				MedalLevel: 20,
			}

			app.session.EventCh <- live.Event{
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
			app.session.EventCh <- live.Event{
				Type: live.GiftEvent,
				Data: mockGift,
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
			app.session.EventCh <- live.Event{
				Type: live.SuperChatEvent,
				Data: mockSuperChat,
			}
		}
	}
}

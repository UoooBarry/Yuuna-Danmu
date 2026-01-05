package live

import (
	"context"
	"log"
	"sync"
)

type Session struct {
	RoomID     int
	EventCh    chan Event
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	mu         sync.Mutex
	Cookie     string
	RealRoomID int
}

func NewSession(roomID int, cookie string) *Session {
	return &Session{
		RoomID:  roomID,
		EventCh: make(chan Event, 100),
		Cookie:  cookie,
	}
}

func (session *Session) UpdateConfig(roomID int, cookie string) {
	session.mu.Lock()
	defer session.mu.Unlock()
	session.RoomID = roomID
	session.Cookie = cookie
}

func (session *Session) Start() error {
	session.mu.Lock()
	defer session.mu.Unlock()

	if session.cancel != nil {
		session.cancel()
	}

	session.ctx, session.cancel = context.WithCancel(context.Background())

	realRoomID, err := GetRealRoomID(session.RoomID)
	if err != nil {
		return err
	}
	session.RealRoomID = realRoomID

	config, err := GetDanmuConfig(session.RealRoomID)
	if err != nil {
		return err
	}

	session.wg.Add(1)
	go session.connectLoop(config.Data.HostList[0].Host, config.Data.Token)

	return nil
}

func (session *Session) Stop() {
	session.mu.Lock()
	if session.cancel != nil {
		session.cancel()
	}
	session.mu.Unlock()

	session.wg.Wait()
}

func (session *Session) connectLoop(host string, token string) {
	defer session.wg.Done()

	for {
		select {
		case <-session.ctx.Done():
			log.Printf("[Yuuna-Danmu] Session stopped by done")
			return
		default:
			c := NewClient(session, host, token)

			if err := c.Run(session.ctx); err != nil {
				if session.ctx.Err() != nil {
					log.Println(session.ctx.Err())
					return
				}
				log.Printf("[Yuuna-Danmu] Disconnecting: %v", err)
			}
		}
	}
}

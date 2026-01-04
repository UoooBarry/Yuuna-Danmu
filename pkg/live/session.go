package live

import (
	"context"
	"log"
	"sync"
)

type Session struct {
	RoomID  int
	EventCh chan Event
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func NewSession(roomID int) *Session {
	ctx, cancel := context.WithCancel(context.Background())
	return &Session{
		RoomID:  roomID,
		EventCh: make(chan Event, 100),
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (session *Session) Start() error {
	config, err := GetDanmuConfig(session.RoomID)
	if err != nil {
		return err
	}

	session.wg.Add(1)
	go session.connectLoop(config.Data.HostList[0].Host, config.Data.Token)

	return nil
}

func (session *Session) Stop() {
	session.cancel()
	session.wg.Wait()
	close(session.EventCh)
}

func (cg *Session) connectLoop(host string, token string) {
	defer cg.wg.Done()

	for {
		select {
		case <-cg.ctx.Done():
			return
		default:
			c := NewClient(cg, host, token)

			if err := c.Run(cg.ctx); err != nil {
				log.Printf("[Yuuna-Danmu] Disconnecting: %v", err)
			}
		}
	}
}

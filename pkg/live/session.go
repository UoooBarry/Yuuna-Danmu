package live

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
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
	go session.connectLoop(config.Data.HostList, config.Data.Token)

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

func (session *Session) connectLoop(hosts []HostInfo, token string) {
	defer session.wg.Done()

	for i, hostInfo := range hosts {
		host := hostInfo.Host

		for {
			select {
			case <-session.ctx.Done():
				log.Printf("[Yuuna-Danmu] Session stopped by context")
				return
			default:
				session.sendSysMsg(fmt.Sprintf("[Yuuna-Danmu] Connecting to host [%d/%d]: %s:%d", i+1, len(hosts), host, hostInfo.WssPort))
				c := NewClient(session, host, hostInfo.WssPort, token)

				err := c.Run(session.ctx)

				if session.ctx.Err() != nil {
					log.Println("[Yuuna-Danmu] Context cancelled, exiting loop")
					log.Println(session.ctx.Err())
					return
				}

				if err != nil {
					session.sendErrorEvent(fmt.Sprintf("[Yuuna-Danmu] Disconnecting from %s:%d: %v", host, hostInfo.WssPort, err))

					session.sendSysMsg("[Yuuna-Danmu] Attempting to switch to next host...")
					goto nextHost
				}
			}
		}
	nextHost:
		time.Sleep(3 * time.Second)
	}

	session.sendErrorEvent("[Yuuna-Danmu] All hosts tried. Stopping session.")
	go session.Stop()
}

func (session *Session) sendSysMsg(msg string) {
	log.Println(msg)
	session.EventCh <- Event{
		Type:      SysMsgEvent,
		Data:      msg,
		Timestamp: time.Now().Unix(),
	}
}

func (session *Session) sendErrorEvent(errStr string) {
	log.Println(errStr)
	session.EventCh <- Event{
		Type:      ErrorEvent,
		Data:      errStr,
		Timestamp: time.Now().Unix(),
	}
}

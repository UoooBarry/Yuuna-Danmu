package live

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"uooobarry/yuuna-danmu/pkg/wbi"

	"github.com/andybalholm/brotli"
	"github.com/gorilla/websocket"
)

type WsClient struct {
	Conn    *websocket.Conn
	RoomID  int
	Host    string
	Token   string
	EventCh chan Event
	Cookie  string
	WssPort int
}

func NewClient(session *Session, host string, wssPort int, token string) *WsClient {
	return &WsClient{
		RoomID:  session.RoomID,
		Host:    host,
		Token:   token,
		EventCh: session.EventCh,
		WssPort: wssPort,
		Cookie:  session.Cookie,
	}
}

func (c *WsClient) Run(ctx context.Context) error {
	address := fmt.Sprintf("wss://%s:%d/sub", c.Host, c.WssPort)
	header := http.Header{}
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	if c.Cookie != "" {
		header.Set("Cookie", c.Cookie)
	}

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, address, header)
	if err != nil {
		return fmt.Errorf("[Yuuna-Danmu]dial error: %w", err)
	}
	c.Conn = conn
	defer c.Conn.Close()

	if err := c.sendAuth(); err != nil {
		return fmt.Errorf("[Yuuna-Danmu]auth error: %w", err)
	}

	heartbeatCtx, cancelHeartbeat := context.WithCancel(ctx)
	defer cancelHeartbeat()
	go c.heartbeatLoop(heartbeatCtx)

	return c.readLoop(ctx)
}

func (c *WsClient) sendAuth() error {
	buvid, err := wbi.GetBuvid3()
	if err != nil {
		return err
	}
	payload := AuthPayload{
		UID:      0,
		RoomID:   c.RoomID,
		ProtoVer: 3, // compress level 3
		Type:     2,
		Key:      c.Token,
		Buvid:    buvid,
	}

	body, _ := json.Marshal(payload)
	return c.Conn.WriteMessage(websocket.BinaryMessage, Pack(OpAuth, body))
}

func (c *WsClient) heartbeatLoop(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			hb := Pack(OpHeartbeat, []byte("[object Object]"))
			if err := c.Conn.WriteMessage(websocket.BinaryMessage, hb); err != nil {
				return
			}
		}
	}
}

func (c *WsClient) readLoop(ctx context.Context) error {
	errChan := make(chan error, 1)
	defer c.Conn.Close()

	go func() {
		for {
			messageType, data, err := c.Conn.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}
			if messageType == websocket.BinaryMessage {
				c.handleMessage(data)
			}
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		if err == io.EOF {
			return nil
		}
		return err
	}
}

func (c *WsClient) handleMessage(data []byte) {
	header, body, err := Unpack(data)
	if err != nil {
		log.Printf("Unpack error: %v", err)
		return
	}

	switch header.ProtoVer {
	case ProtoRaw, ProtoInt:
		c.routeOperation(header, body)

	case ProtoBrotli:
		reader := brotli.NewReader(bytes.NewReader(body))
		decompressed, _ := io.ReadAll(reader)
		c.handleMultiplePackets(decompressed)

	case ProtoZlib:
		log.Println("zlib not supported yet")
	}
}

func (c *WsClient) handleMultiplePackets(data []byte) {
	for len(data) > 0 {
		header, body, err := Unpack(data)
		if err != nil {
			log.Printf("Unpack error: %v", err)
			break
		}
		c.dispatch(body)
		data = data[header.PacketLen:]
	}
}

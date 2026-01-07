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

	"github.com/andybalholm/brotli"
	"github.com/gorilla/websocket"
)

type WsClient struct {
	Conn        *websocket.Conn
	RoomID      int
	Host        string
	Token       string
	WssPort     int
	authContext *AuthContext
	eventCh     chan Event
}

func init() {
	websocket.DefaultDialer.HandshakeTimeout = 10 * time.Second
}

func (s *Session) NewClient(host string, wssPort int, token string) *WsClient {
	return &WsClient{
		RoomID:      s.RealRoomID,
		Host:        host,
		Token:       token,
		WssPort:     wssPort,
		authContext: s.authContext,
		eventCh:     s.EventCh,
	}
}

func (c *WsClient) getHeader() *http.Header {
	header := http.Header{}
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	if c.authContext.Cookie != "" {
		header.Set("Cookie", c.authContext.Cookie)
	}
	return &header
}

func (c *WsClient) Run(ctx context.Context) error {
	address := fmt.Sprintf("wss://%s:%d/sub", c.Host, c.WssPort)
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, address, *c.getHeader())
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
	buvid := c.authContext.Buvid3

	payload := AuthPayload{
		UID:      c.authContext.UID,
		RoomID:   c.RoomID,
		ProtoVer: 3, // compress level 3
		Type:     2,
		Key:      c.Token,
		Buvid:    buvid,
	}
	log.Println(payload)

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

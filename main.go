package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"uooobarry/yuuna-danmu/pkg/live"
	"uooobarry/yuuna-danmu/pkg/wbi"
)

func main() {
	wbi.EnsureBuvid()
	realRoomID, err := live.GetRealRoomID(50819)
	log.Printf("[Yuuna-Danmu] Connecting to room: %d", realRoomID)
	if err != nil {
		panic(err)
	}
	session := live.NewSession(realRoomID)
	session.Start()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Printf("[Yuuna-Danmu] Received signal: %s", sig)
		session.Stop()
		log.Printf("[Yuuna-Danmu] Session stopped")
	}()

	for event := range session.EventCh {
		switch event.Type {
		case live.DanmuEvent:
			if data, ok := event.Data.(*live.DanmuMsg); ok {
				fmt.Printf("[%s|%d]%s: %s\n", data.MedalName, data.MedalLevel, data.Nickname, data.Content)
			}
		case live.PopularityEvent:
			if data, ok := event.Data.(*live.PopularityMsg); ok {
				fmt.Printf("[Yuuna-Danmu] Popularity: %d\n", data.Popularity)
			}
		case live.GiftEvent:
			if data, ok := event.Data.(*live.GiftData); ok {
				fmt.Printf("[%s] Gift: %s x%d\n", data.Nickname, data.GiftName, data.Num)
			}
		default:
		}
	}
}

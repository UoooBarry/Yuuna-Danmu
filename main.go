package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"uooobarry/yuuna-danmu/pkg/app"
	"uooobarry/yuuna-danmu/pkg/ui"
	"uooobarry/yuuna-danmu/pkg/wbi"
)

func main() {
	wbi.EnsureBuvid()
	app := app.NewApp(
		app.WithUI(ui.NewWailsUI(assets)),
		app.WithFileLog("log/yuuna-danmu.log"),
	)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		app.Stop()
		log.Printf("[Yuuna-Danmu] Session stopped")
	}()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

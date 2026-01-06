package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"uooobarry/yuuna-danmu/pkg/app"
	"uooobarry/yuuna-danmu/pkg/ui"
	"uooobarry/yuuna-danmu/pkg/wbi"

	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func main() {
	wbi.EnsureBuvid()
	app := app.NewApp(
		app.WithUI(ui.NewWailsUI(&assetserver.Options{
			Assets:     assets,
			Middleware: proxyHandler,
		})),
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

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"uooobarry/yuuna-danmu/pkg/app"
	"uooobarry/yuuna-danmu/pkg/ui"

	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func main() {
	app, err := app.NewApp(
		app.WithUI(ui.NewWailsUI(&assetserver.Options{
			Assets:     assets,
			Middleware: proxyHandler,
		})),
		app.WithFileLog("log/yuuna-danmu.log"),
	)
	if err != nil {
		log.Fatal(err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		go func() {
			time.Sleep(1 * time.Second)
			log.Println("[Yuuna-Danmu] Shutdown timed out, force exiting...")
			os.Exit(1)
		}()

		app.Stop()
		os.Exit(0)
	}()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"fmt"
	"github.com/reactivejson/vwap-engine/internal/app"
	"github.com/reactivejson/vwap-engine/internal/storage"
	"github.com/reactivejson/vwap-engine/internal/tunnel"
	"log"
	"os"
	"os/signal"
	"syscall"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2022
 */

//The core entry point into the app. will setup the config, and run the App
func main() {
	ctx, cancelCtxFn := context.WithCancel(context.Background())

	cfg := app.SetupEnvConfig()

	fmt.Printf("wsURL %s\n", cfg.WebsocketUrl)
	fmt.Printf("Trading pairs %s\n", cfg.TradingPairs)

	// Intercepting shutdown signals.
	go waitForSignal(ctx, cancelCtxFn)

	ws, err := tunnel.NewReceiver(cfg.WebsocketUrl)
	if err != nil {
		log.Fatal(err)
	}

	queue, err := storage.NewVwapQueue([]storage.DataPoint{}, cfg.WindowSize)
	if err != nil {
		log.Fatal(err)
	}

	svc := app.NewContext(ws, queue, cfg)

	err = svc.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func waitForSignal(ctx context.Context, cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-signals:
		log.Printf("received signal: %s, exiting gracefully", s)
		cancel()
	case <-ctx.Done():
		log.Printf("Service context done, serving remaining requests and exiting.")
	}
}

package app

import (
	"github.com/reactivejson/vwap-engine/internal/storage"
	"github.com/reactivejson/vwap-engine/internal/tunnel"
	"time"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2022
 */

type envConfig struct {
	Port         uint          `envconfig:"PORT"               required:"false" default:"8080"`
	HTTPTimeout  time.Duration `envconfig:"HTTP_TIMEOUT"       required:"false" default:"1800s"`
	WebsocketUrl string        `envconfig:"WEBSOCKET_URL"      required:"false" default:"wss://ws-feed.pro.coinbase.com"`
	TradingPairs []string      `envconfig:"TRADING_PAIRS"      required:"false" default:"BTC-USD,ETH-USD,ETH-BTC"`
	WindowSize   uint          `envconfig:"WINDOW_SIZE"        required:"false" default:"200"`
}

// Context is application's content
type Context struct {
	cfg        *envConfig
	wsReceiver tunnel.Tunnel
	queue      storage.Vwap
}

// NewContext instantiates new rte context object.
func NewContext(tunnel tunnel.Tunnel, queue storage.Vwap, cfg *envConfig) *Context {
	return &Context{
		cfg:        cfg,
		wsReceiver: tunnel,
		queue:      queue,
	}
}

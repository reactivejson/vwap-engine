package tunnel

import (
	"context"
	"github.com/reactivejson/vwap-engine/api/models"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2022
 */

//Tunnel is a coinbase websocket stream client to receive data from coinbase websocket server
type Tunnel interface {
	//Subscribe sends a subscribe request to the coinbase channel's websocket, using trading pairs (productIDs).
	Subscribe(tradingPairs []string) error

	// Read receives data points from the coinbase and passes them to the receiver channel.
	Read(ctx context.Context, receiver chan *models.CoinbaseResponse)

	// Close closes the websocket connection.
	Close()
}

package tunnel

import (
	"context"
	"github.com/reactivejson/vwap-engine/api/models"
	"testing"

	"github.com/stretchr/testify/require"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2022
 */

var (
	wsURL = "wss://ws-feed.pro.coinbase.com"
)

func TestNewReceiver(t *testing.T) {
	t.Parallel()

	_, err := NewReceiver(wsURL)
	require.NoError(t, err)
}

func TestTunnelSubscribe_AndRead_ShouldSucceed(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	receiver := make(chan *models.CoinbaseResponse)

	ws, err := NewReceiver(wsURL)
	require.NoError(t, err)

	err = ws.Subscribe([]string{"BTC-USD"})
	ws.Read(ctx, receiver)
	require.NoError(t, err)
	defer ws.Close()
	var limit int
	// Check the first couple of responses.
	for response := range receiver {
		if limit >= 3 {
			break
		}

		//check for a proper response type, since the first one is a subscription msg
		//&{[{matches [BTC-USD]}]     subscriptions  0 }
		//Received: &{[{matches [BTC-USD]}]     subscriptions  0 }
		//Received: &{[]  29303.34 BTC-USD 0.00324023 last_match 2022-05-21T09:12:02.350239Z 341498073 buy}
		//Received: &{[]  29303.35 BTC-USD 0.0000299 match 2022-05-21T09:12:04.862866Z 341498074 sell}
		if response.Type == "match" {
			require.Equal(t, "BTC-USD", response.ProductID)
		}
		limit++
	}
}

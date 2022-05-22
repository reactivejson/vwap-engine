package tunnel

import (
	"context"
	"github.com/reactivejson/vwap-engine/api/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2022
 */

var (
	url = "wss://ws-feed.pro.coinbase.com"
)

type ReceiverSuite struct {
	suite.Suite
	receiver chan *models.CoinbaseResponse
	ws       Tunnel
}

func TestReceiverSuite(t *testing.T) {
	suite.Run(t, new(ReceiverSuite))
}

func (suite *ReceiverSuite) SetupTest() {

	suite.ws, _ = NewReceiver(url)
	suite.receiver = make(chan *models.CoinbaseResponse)
}

func (suite *ReceiverSuite) TearDownTest() {
}

func (suite *ReceiverSuite) TestNewReceiver_Subscribe() {
	err := suite.ws.Subscribe([]string{"BTC-USD"})
	require.NoError(suite.T(), err)
}

func (suite *ReceiverSuite) TestTunnelSubscribe_AndRead() {

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)

	suite.receiver = make(chan *models.CoinbaseResponse)

	err := suite.ws.Subscribe([]string{"BTC-USD"})
	suite.ws.Read(ctx, suite.receiver)
	require.NoError(suite.T(), err)
	var limit int

	for {
		select {
		case <-ctx.Done():
			assert.True(suite.T(), false)
		default:
			// Check the first couple of responses.
			for response := range suite.receiver {
				if limit >= 3 {
					break
				}

				//check for a proper response type, since the first one is a subscription msg
				//&{[{matches [BTC-USD]}]     subscriptions  0 }
				//Received: &{[{matches [BTC-USD]}]     subscriptions  0 }
				//Received: &{[]  29303.34 BTC-USD 0.00324023 last_match 2022-05-21T09:12:02.350239Z 341498073 buy}
				//Received: &{[]  29303.35 BTC-USD 0.0000299 match 2022-05-21T09:12:04.862866Z 341498074 sell}
				if response.Type == "match" {
					require.Equal(suite.T(), "BTC-USD", response.ProductID)
				}
				limit++
			}
			<-suite.receiver
			assert.True(suite.T(), true)
			cancel()
			return
		}
	}
}

func (suite *ReceiverSuite) TestTunnel_WithBadUrl_ShouldFail() {
	_, err := NewReceiver("ws://InvalidUrl")
	require.Error(suite.T(), err)
}

package models

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2022
 */

// Channel is a Model to subscribe to the Coinbase Websocket stream channels
type Channel struct {
	Name       string   `json:"name"`
	ProductIDs []string `json:"product_ids" mapstructure:"product_ids"`
}

// CoinbaseRequest JSON Payload request to be sent to subscribe to the Coinbase Websocket stream channels
/**
Sample:
{
    "type": "subscribe",
    "product_ids": [
        "ETH-USD",
        "BTC-USD"
    ],
    "channels": ["ticker_batch"]
}
*/
type CoinbaseRequest struct {
	Type       string    `json:"type"`
	ProductIDs []string  `json:"product_ids"`
	Channels   []Channel `json:"channels"`
}

// CoinbaseResponse is the coinbase response payload.
type CoinbaseResponse struct {
	Channels []Channel `json:"channels"`
	Message  string    `json:"message,omitempty"`
	Price    string    `json:"price"`
	//ProductID  is the coinbase for TradingPair
	ProductID string `json:"product_id"`
	Size      string `json:"size"`
	Type      string `json:"type"`
	Time      string `json:"time"`
	TradeID   int    `json:"trade_id"`
	Side      string `json:"side"`
}

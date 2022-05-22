package tunnel

import (
	"context"
	"fmt"
	ws "github.com/gorilla/websocket"
	"github.com/reactivejson/vwap-engine/api/models"
	"log"
	"net/http"
	"time"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2022
 */

const (
	TunnelSubscribe string = "subscribe"
)

// Receiver connexion
type Receiver struct {
	conn *ws.Conn
	done chan struct{} // the Receiver will close done once it cannot read from the websocket anymore
}

// NewReceiver initializes a new coinbase Tunnel object and dials the coinbase websocket. It takes a coinbase ws urr,
// If a connection cannot be reached it returns an error.
// NewReceiver returns a new websocket client Tunnel.
func NewReceiver(websocketUrl string) (Tunnel, error) {
	dialer := ws.Dialer{HandshakeTimeout: 5 * time.Second}
	conn, _, err := dialer.Dial(websocketUrl, http.Header{})
	if err != nil {
		return nil, fmt.Errorf("error while creating websocket receiver: %v", err)
	}

	done := make(chan struct{})

	log.Printf("Successfully connected to: %s", websocketUrl)

	return &Receiver{
		conn: conn,
		done: done,
	}, nil
}

// NewReceiverWithconn returns a new websocket client.
func NewReceiverWithconn(websocketUrl string, conn *ws.Conn) (Tunnel, error) {
	done := make(chan struct{})

	return &Receiver{
		conn: conn,
		done: done,
	}, nil
}

//Subscribe sends a subscribe request to the coinbase channel's websocket, using trading pairs (productIDs).
func (r *Receiver) Subscribe(tradingPairs []string) error {
	sbPayload := models.CoinbaseRequest{
		Type:       TunnelSubscribe,
		ProductIDs: tradingPairs,
		Channels: []models.Channel{
			{Name: "matches"},
		},
	}

	if err := r.conn.WriteJSON(sbPayload); err != nil {
		return fmt.Errorf("error writing JSON to subscribe: %v", err)
	}

	return nil

}

// Read receives data points from the coinbase and passes them to the receiver channel.
func (r *Receiver) Read(ctx context.Context, receiver chan *models.CoinbaseResponse) {
	//defer close(r.done)

	go func() {
		defer close(receiver)
		for {
			select {
			case <-r.done:
				return
			case <-ctx.Done():
				// Close the connection completely by sending a close message and then waiting (with timeout) for the coinbase server to do so.
				err := r.conn.WriteMessage(ws.CloseMessage, ws.FormatCloseMessage(ws.CloseNormalClosure, ""))
				if err != nil {
					log.Printf("error writing close message %v", err)
					return
				}
				select {
				case <-r.done:
				case <-time.After(2 * time.Second):
					log.Print("timeout waiting for close message response")
					r.Close()
				}
				return

			default:
				response := &models.CoinbaseResponse{}
				err := r.conn.ReadJSON(response)
				if err != nil {
					exceptionHandler(err)
					return
				}
				receiver <- response
			}
		}
	}()
}

// Close shuts down the websocket connection tunnel and logs any close error.
func (r *Receiver) Close() {
	err := r.conn.Close()
	if err != nil {
		exceptionHandler(err)
	} else {
		log.Printf("Tunnel connection closed successfuly")
	}
}

// exceptionHandler handles and logs reading errors.
func exceptionHandler(err error) {
	if closeErr, ok := err.(*ws.CloseError); ok {
		if closeErr.Code == ws.CloseNormalClosure {
			log.Printf("normal connection close for Tunnel %v", closeErr)
		} else {
			log.Printf("close error for Tunnel  %v", closeErr)
		}
	} else {
		log.Printf("read error for Tunnel %v", err)
	}
}

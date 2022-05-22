package tunnel

import (
	"bufio"
	"fmt"
	ws "github.com/gorilla/websocket"
	"github.com/reactivejson/vwap-engine/api/models"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"os"
	"testing"
)

func verifyWsSub(t *testing.T, done chan struct{}, wantErr bool) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		var upgrader = ws.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		assert.NoError(t, err)
		defer conn.Close()

		if !wantErr {
			subMsg := models.CoinbaseRequest{}
			err = conn.ReadJSON(&subMsg)
			assert.NoError(t, err)
			assert.Equal(t, "subscribe", subMsg.Type)
			assert.Equal(t, []string{"BTC-USD"}, subMsg.ProductIDs)
		}
		done <- struct{}{}

	}
}

func wsDial(t *testing.T, wantErr bool) func(w http.ResponseWriter, r *http.Request) {
	if wantErr {
		return func(w http.ResponseWriter, r *http.Request) {
			// do nothing to force dial error
		}
	}
	return wsSuccess(t)
}

func wsSuccess(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var upgrader = ws.Upgrader{}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				break
			}
			err = c.WriteMessage(mt, message)
			if err != nil {
				break
			}
		}

		assert.NoError(t, err)
	}
}

func scannerHelper(t *testing.T) (*bufio.Scanner, *os.File, *os.File) {
	reader, writer, err := os.Pipe()
	if err != nil {
		assert.Fail(t, "couldn't get os Pipe: %v", err)
	}
	log.SetOutput(writer)

	return bufio.NewScanner(reader), reader, writer
}

func resetScanner(reader *os.File, writer *os.File) {
	err := reader.Close()
	if err != nil {
		fmt.Println("error closing reader was ", err)
	}
	if err = writer.Close(); err != nil {
		fmt.Println("error closing writer was ", err)
	}
	log.SetOutput(os.Stderr)
}

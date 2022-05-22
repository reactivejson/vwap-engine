package tunnel

import (
	ws "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2022
 */

var (
	webSocketURL string
)

func setUpWSServer(handlerFunc http.HandlerFunc) *httptest.Server {
	server := httptest.NewServer(handlerFunc)
	webSocketURL = "ws" + strings.TrimPrefix(server.URL, "http")
	return server
}

func TestNewTunnel(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"Valid handshake", false},
		{"InValid handshake", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := setUpWSServer(wsDial(t, tt.wantErr))
			defer server.Close()

			t.Parallel()

			tunnel, err := NewReceiver(webSocketURL)

			if !tt.wantErr {
				defer tunnel.Close()
			}

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestTunnelr_Close(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"Close", false},
		{"Invalid close", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner, reader, writer := scannerHelper(t)
			defer resetScanner(reader, writer)

			server := setUpWSServer(wsSuccess(t))
			defer server.Close()

			dialer := ws.Dialer{HandshakeTimeout: 5 * time.Second}
			conn, _, err := dialer.Dial(webSocketURL, http.Header{})

			tunnel, err := NewReceiverWithconn(webSocketURL, conn)
			require.NoError(t, err)
			defer tunnel.Close()

			if tt.wantErr {
				err := conn.Close()
				assert.NoError(t, err)
			}
			tunnel.Close()

			if tt.wantErr {
				scanner.Scan()
				assert.Contains(t, scanner.Text(), "error for Tunnel close")
			} else {
				assert.Empty(t, scanner.Text())
			}
		})
	}
}

func TestTunnel_Subscribe(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"test valid subscribe", false},
		{"test error from writeJSON", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, reader, writer := scannerHelper(t)
			defer resetScanner(reader, writer)

			done := make(chan struct{})
			server := setUpWSServer(verifyWsSub(t, done, tt.wantErr))
			defer server.Close()

			dialer := ws.Dialer{HandshakeTimeout: 5 * time.Second}
			conn, _, err := dialer.Dial(webSocketURL, http.Header{})

			tunnel, err := NewReceiverWithconn(webSocketURL, conn)
			require.NoError(t, err)
			if tt.wantErr {
				conn.Close()
			} else {
				defer conn.Close()
			}

			err = tunnel.Subscribe([]string{"BTC-USD"})
			if tt.wantErr {
				assert.Error(t, err)
				conn.Close()
			} else {
				conn.Close()
				assert.NoError(t, err)
			}
		})
	}
}

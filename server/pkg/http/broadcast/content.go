package broadcast

import (
	"golang.org/x/net/websocket"
)

func handleContentMessage(wsc *websocket.Conn, message []byte) error {
	documentContent = message
	return nil
}

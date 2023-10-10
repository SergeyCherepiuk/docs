package broadcast

import (
	"encoding/json"

	"golang.org/x/net/websocket"
)

func handleUserMessage(wsc *websocket.Conn, message []byte) error {
	var user User
	if err := json.Unmarshal(message, &user); err != nil {
		return err
	}

	connections[wsc] = user
	return nil
}

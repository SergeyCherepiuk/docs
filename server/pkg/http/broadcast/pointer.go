package broadcast

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/websocket"
)

func handlePointerMessage(wsc *websocket.Conn, message []byte) error {
	var u User
	if err := json.Unmarshal(message, &u); err != nil {
		return err
	}

	user, ok := connections[wsc]
	if !ok {
		return fmt.Errorf("connection wasn't found")
	}

	user.Pointer = u.Pointer
	connections[wsc] = user
	return nil
}

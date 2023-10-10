package broadcast

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/websocket"
)

func handleSelectionMessage(wsc *websocket.Conn, message []byte) error {
	var u User
	if err := json.Unmarshal(message, &u); err != nil {
		return err
	}

	user, ok := connections[wsc]
	if !ok {
		return fmt.Errorf("connection wasn't found")
	}

	user.Selection = u.Selection
	connections[wsc] = user
	return nil
}

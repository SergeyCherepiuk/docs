package broadcast

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

var (
	messageHandlers = map[string]func(*websocket.Conn, []byte) error{
		"pointer":   handlePointerMessage,
		"content":   handleContentMessage,
		"selection": handleSelectionMessage,
		"user":      handleUserMessage,
	}
	connections     = make(map[*websocket.Conn]User)
	documentContent = []byte("")
)

type Message struct {
	MessageType string          `json:"messageType"`
	RawMessage  json.RawMessage `json:"rawMessage"`
}

type User struct {
	ID        string    `json:"id"`
	Pointer   Pointer   `json:"pointer"`
	Selection Selection `json:"selection"`
}

type Pointer struct {
	Position struct {
		x int
		y int
	} `json:"position"`
	Scroll struct {
		x int
		y int
	} `json:"scroll"`
}

type Selection struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

func Connect(c echo.Context) error {
	websocket.Server{Handler: func(wsc *websocket.Conn) {
		defer wsc.Close()
		defer delete(connections, wsc)
		// TODO: Send disconnect message

		// TODO: Send current UI state back

		for {
			var jsonMessage []byte
			if err := websocket.Message.Receive(wsc, &jsonMessage); err != nil {
				break
			}

			var message Message
			if err := json.Unmarshal(jsonMessage, &message); err != nil {
				break
			}

			if handler, ok := messageHandlers[message.MessageType]; ok {
				handler(wsc, []byte(message.RawMessage)) // NOTE: Error is ignored
			}

			if message.MessageType == "user" {
				continue
			}

			for conn := range connections {
				if conn != wsc {
					websocket.Message.Send(conn, string(jsonMessage))
				}
			}
		}
	}}.ServeHTTP(c.Response(), c.Request())
	return c.NoContent(http.StatusSwitchingProtocols)
}

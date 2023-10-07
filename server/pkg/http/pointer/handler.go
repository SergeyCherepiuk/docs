package pointer

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

var conns = make(map[*websocket.Conn]struct{})

func Broadcast(c echo.Context) error {
	websocket.Server{Handler: websocket.Handler(func(ws *websocket.Conn) {
		conns[ws] = struct{}{}
		defer delete(conns, ws)
		defer ws.Close()

		for {
			var message string
			if err := websocket.Message.Receive(ws, &message); err != nil {
				break
			}

			for conn := range conns {
				if conn != ws {
					websocket.Message.Send(conn, message)
				}
			}
		}
	})}.ServeHTTP(c.Response(), c.Request())
	return c.NoContent(http.StatusSwitchingProtocols)
}

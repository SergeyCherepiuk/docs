package broadcast

import (
	"net/http"

	"github.com/SergeyCherepiuk/docs/pkg/http/ws"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

// TODO: Consider sharing the room for pointers and selections (same functionality)
//  to reduce number of concurrent websocket connections
var selectionRoom = ws.NewRoom()

func Selection(c echo.Context) error {
	websocket.Server{Handler: func(wsc *websocket.Conn) {
		selectionRoom.EnterAndListen(wsc, ws.Listener{
			OnEnter: func() {},
			OnMessage: func(message any) any {
				return message
			},
			OnExit: func() {},
		})
	}}.ServeHTTP(c.Response(), c.Request())
	return c.NoContent(http.StatusSwitchingProtocols)
}